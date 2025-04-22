import os
import argparse
import logging
from pathlib import Path
from typing import List, Optional, Dict, Any
from tqdm import tqdm
from concurrent.futures import ThreadPoolExecutor, as_completed
import json
from datetime import datetime

# Import our custom modules
from file_processor import FileProcessor
from Types.document import Document
from Types.metadata import Metadata
from elasticsearch_handler import ElasticsearchHandler

# Import constants
from constants import (
    MAX_FILE_SIZE_DEFAULT,
    ES_DEFAULT_HOST,
    ES_DEFAULT_PORT,
    ES_DEFAULT_INDEX,
    DEFAULT_MAX_WORKERS,
    DEFAULT_BATCH_SIZE,
    EXTENSION_CATEGORIES,
    PATH_CATEGORIES
)

# Configure logging
logging.basicConfig(
    level=logging.INFO,
    format='%(asctime)s - %(name)s - %(levelname)s - %(message)s',
    handlers=[
        logging.FileHandler("document_processing.log"),
        logging.StreamHandler()
    ]
)
logger = logging.getLogger("document_processor")

class DocumentProcessor:
    """Main class for processing and indexing documents"""
    
    def __init__(
        self, 
        es_host: str = ES_DEFAULT_HOST, 
        es_port: int = ES_DEFAULT_PORT,
        es_index: str = ES_DEFAULT_INDEX,
        es_username: Optional[str] = None,
        es_password: Optional[str] = None,
        es_use_ssl: bool = False,
        max_workers: int = DEFAULT_MAX_WORKERS,
        batch_size: int = DEFAULT_BATCH_SIZE,
        max_file_size: int = MAX_FILE_SIZE_DEFAULT,
        skip_existing: bool = True
    ):
        """Initialize the document processor"""
        self.file_processor = FileProcessor(max_file_size=max_file_size)
        self.es_handler = ElasticsearchHandler(
            host=es_host, 
            port=es_port, 
            index_name=es_index,
            username=es_username,
            password=es_password,
            use_ssl=es_use_ssl
        )
        self.max_workers = max_workers
        self.batch_size = batch_size
        self.skip_existing = skip_existing
        
        # Ensure the index exists
        self.es_handler.create_index()
        
        # Stats tracking
        self.stats = {
            "processed": 0,
            "indexed": 0,
            "failed": 0,
            "skipped": 0,
            "total": 0
        }
    
    def get_file_category(self, file_path: str) -> str:
        """Determine document category based on file path or extension"""
        path = Path(file_path)
        ext = path.suffix.lower()
        
        # Check for specific document types in path
        path_str = str(path).lower()
        for key_term, category in PATH_CATEGORIES.items():
            if key_term in path_str:
                return category
            
        # Fall back to extension-based categorization
        return EXTENSION_CATEGORIES.get(ext, 'Document')
    
    def process_file(self, file_path: str) -> Optional[Document]:
        """Process a single file into a Document object"""
        try:
            # Generate file hash first to check if we can skip
            file_hash = self.file_processor.get_file_hash(file_path)
            
            if not file_hash:
                logger.error(f"Could not generate hash for {file_path}")
                return None
                
            if self.skip_existing and self.es_handler.document_exists(file_hash):
                logger.info(f"Skipping existing document: {file_path}")
                self.stats["skipped"] += 1
                return None
            
            # Extract text and metadata
            text, metadata_dict = self.file_processor.extract_text(file_path)
            
            if not text:
                logger.warning(f"No text extracted from {file_path}")
                self.stats["failed"] += 1
                return None
            
            # Create metadata object
            metadata_obj = Metadata(
                document_name=metadata_dict.get('document_name', Path(file_path).name),
                subject=metadata_dict.get('subject', 'Unknown'),
                status=metadata_dict.get('status', 'processed'),
                timestamp=metadata_dict.get('timestamp', datetime.now()),
                case_name=metadata_dict.get('case_name'),
                author=metadata_dict.get('author'),
                judge=metadata_dict.get('judge'),
                court=metadata_dict.get('court')
            )
            
            # Determine document category
            category = self.get_file_category(file_path)
            
            # Create document object
            document = Document(
                file_path=file_path,
                text=text,
                metadata=metadata_obj,
                doc_type="document",
                category=category,
                hash_value=file_hash,
                created_at=datetime.now()
            )
            
            self.stats["processed"] += 1
            return document
            
        except Exception as e:
            logger.error(f"Error processing file {file_path}: {e}")
            self.stats["failed"] += 1
            return None
    
    def process_directory(self, directory_path: str, file_extensions: Optional[List[str]] = None):
        """Process all files in a directory and its subdirectories"""
        logger.info(f"Processing directory: {directory_path}")
        
        # Get list of files to process
        file_paths = []
        for root, _, files in os.walk(directory_path):
            for file in files:
                file_path = os.path.join(root, file)
                
                # Skip files with unwanted extensions
                if file_extensions:
                    ext = os.path.splitext(file_path)[1].lower()
                    if ext not in file_extensions:
                        continue
                        
                # Skip files that can't be processed
                if not self.file_processor.can_process(file_path):
                    continue
                    
                file_paths.append(file_path)
        
        self.stats["total"] = len(file_paths)
        logger.info(f"Found {self.stats['total']} files to process")
        
        # Process files in parallel
        documents_batch = []
        with tqdm(total=len(file_paths), desc="Processing Files") as pbar:
            with ThreadPoolExecutor(max_workers=self.max_workers) as executor:
                future_to_file = {executor.submit(self.process_file, file_path): file_path for file_path in file_paths}
                
                for future in as_completed(future_to_file):
                    file_path = future_to_file[future]
                    try:
                        document = future.result()
                        if document:
                            documents_batch.append(document)
                            
                            # Index in batches
                            if len(documents_batch) >= self.batch_size:
                                success, errors = self.es_handler.bulk_index_documents(documents_batch)
                                self.stats["indexed"] += success
                                self.stats["failed"] += errors
                                documents_batch = []
                                
                    except Exception as e:
                        logger.error(f"Error processing {file_path}: {e}")
                        self.stats["failed"] += 1
                        
                    pbar.update(1)
        
        # Index any remaining documents
        if documents_batch:
            success, errors = self.es_handler.bulk_index_documents(documents_batch)
            self.stats["indexed"] += success
            self.stats["failed"] += errors
            
        logger.info(f"Processing complete. Stats: {json.dumps(self.stats)}")
        return self.stats


if __name__ == "__main__":
    parser = argparse.ArgumentParser(description="Process documents and index them in Elasticsearch")
    parser.add_argument("directory", help="Directory containing documents to process")
    parser.add_argument("--es-host", default=ES_DEFAULT_HOST, help="Elasticsearch host")
    parser.add_argument("--es-port", type=int, default=ES_DEFAULT_PORT, help="Elasticsearch port")
    parser.add_argument("--es-index", default=ES_DEFAULT_INDEX, help="Elasticsearch index name")
    parser.add_argument("--es-username", help="Elasticsearch username")
    parser.add_argument("--es-password", help="Elasticsearch password")
    parser.add_argument("--es-ssl", action="store_true", help="Use SSL for Elasticsearch connection")
    parser.add_argument("--workers", type=int, default=DEFAULT_MAX_WORKERS, help="Number of worker threads")
    parser.add_argument("--batch-size", type=int, default=DEFAULT_BATCH_SIZE, help="Batch size for indexing")
    parser.add_argument("--max-file-size", type=int, default=int(MAX_FILE_SIZE_DEFAULT/(1024*1024)), help="Maximum file size in MB")
    parser.add_argument("--extensions", nargs="+", help="List of file extensions to process")
    parser.add_argument("--process-all", action="store_true", help="Process all files regardless of extension")
    parser.add_argument("--force", action="store_true", help="Process files even if they already exist in the index")
    
    args = parser.parse_args()
    
    # Convert extensions to proper format
    extensions = None
    if not args.process_all and args.extensions:
        extensions = [ext if ext.startswith('.') else f'.{ext}' for ext in args.extensions]
    
    processor = DocumentProcessor(
        es_host=args.es_host,
        es_port=args.es_port,
        es_index=args.es_index,
        es_username=args.es_username,
        es_password=args.es_password,
        es_use_ssl=args.es_ssl,
        max_workers=args.workers,
        batch_size=args.batch_size,
        max_file_size=args.max_file_size * 1024 * 1024,  # Convert to bytes
        skip_existing=not args.force
    )
    
    stats = processor.process_directory(args.directory, file_extensions=extensions)
    
    print("\nProcessing Summary:")
    print(f"Total files found: {stats['total']}")
    print(f"Files processed: {stats['processed']}")
    print(f"Documents indexed: {stats['indexed']}")
    print(f"Files skipped: {stats['skipped']}")
    print(f"Failed files: {stats['failed']}")