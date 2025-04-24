"""
Document processor for processing and indexing documents.
"""
import os
import uuid
import logging
from pathlib import Path
from typing import List, Optional, Dict
from tqdm import tqdm
from concurrent.futures import ThreadPoolExecutor, as_completed
from datetime import datetime
import openai

from src.models.document import Document
from src.models.metadata import Metadata
from src.handlers.file_processor import FileProcessor
from src.handlers.elasticsearch_handler import ElasticsearchHandler
from src.handlers.s3_handler import S3Handler
from src.handlers.document_classifier import process_document_llm
from src.utils.constants import (
    MAX_FILE_SIZE_DEFAULT,
    ES_DEFAULT_HOST,
    ES_DEFAULT_PORT,
    ES_DEFAULT_INDEX,
    DEFAULT_MAX_WORKERS,
    DEFAULT_BATCH_SIZE,
    EXTENSION_CATEGORIES
)

# Configure logging
logger = logging.getLogger("document_processor")


class DocumentProcessor:
    """
    Main class for processing and indexing documents.
    Handles document text extraction, classification, S3 upload, and Elasticsearch indexing.
    """
    
    def __init__(
        self, 
        # Elasticsearch parameters
        es_host: str = ES_DEFAULT_HOST, 
        es_port: int = ES_DEFAULT_PORT,
        es_index: str = ES_DEFAULT_INDEX,
        es_username: Optional[str] = None,
        es_password: Optional[str] = None,
        es_api_key: Optional[str] = None,
        es_cloud_id: Optional[str] = None,
        es_use_ssl: bool = True,
        # S3 parameters
        s3_bucket: Optional[str] = None,
        s3_region: Optional[str] = None,
        s3_access_key: Optional[str] = None,
        s3_secret_key: Optional[str] = None,
        # Processing parameters
        max_workers: int = DEFAULT_MAX_WORKERS,
        batch_size: int = DEFAULT_BATCH_SIZE,
        max_file_size: int = MAX_FILE_SIZE_DEFAULT,
        skip_existing: bool = True,
        use_llm_classification: bool = True,
        openai_api_key: Optional[str] = None
    ):
        """
        Initialize the document processor.
        
        Args:
            es_host: Elasticsearch host
            es_port: Elasticsearch port
            es_index: Elasticsearch index name
            es_username: Elasticsearch username for authentication
            es_password: Elasticsearch password for authentication
            es_use_ssl: Whether to use SSL for Elasticsearch connection
            s3_bucket: S3 bucket name for document storage
            s3_region: AWS region for S3
            s3_access_key: AWS access key ID
            s3_secret_key: AWS secret access key
            max_workers: Maximum number of worker threads
            batch_size: Batch size for indexing
            max_file_size: Maximum file size in bytes
            skip_existing: Whether to skip existing documents
            use_llm_classification: Whether to use LLM for document classification
            openai_api_key: OpenAI API key for document classification
        """
        # Initialize file processor
        self.file_processor = FileProcessor(max_file_size=max_file_size)
        
        # Initialize Elasticsearch handler
        self.es_handler = ElasticsearchHandler(
            host=es_host, 
            port=es_port, 
            index_name=es_index,
            username=es_username,
            password=es_password,
            api_key=es_api_key,
            cloud_id=es_cloud_id,
            use_ssl=es_use_ssl
        )
        
        # Initialize S3 handler if bucket provided
        self.s3_handler = None
        if s3_bucket:
            self.s3_handler = S3Handler(
                bucket_name=s3_bucket,
                region_name=s3_region,
                aws_access_key_id=s3_access_key or os.environ.get('AWS_ACCESS_KEY_ID'),
                aws_secret_access_key=s3_secret_key or os.environ.get('AWS_SECRET_ACCESS_KEY')
            )
        
        # Set up OpenAI for document classification if enabled
        self.use_llm_classification = use_llm_classification
        if use_llm_classification:
            api_key = openai_api_key or os.environ.get('OPENAI_API_KEY')
            if not api_key:
                logger.warning("OpenAI API key not provided, LLM classification will be disabled")
                self.use_llm_classification = False
            else:
                openai.api_key = api_key
                self.openai_client = openai
        
        # Other parameters
        self.max_workers = max_workers
        self.batch_size = batch_size
        self.skip_existing = skip_existing
        
        # Ensure the index exists
        self.es_handler.create_index()
        
        # Initialize statistics dictionary
        self.stats = {
            "total": 0,
            "processed": 0,
            "indexed": 0,
            "uploaded": 0,
            "skipped": 0,
            "failed": 0,
            "unsupported_format": 0,
            "too_large": 0
        }
    
    def get_file_category(self, file_path: str) -> str:
        """
        Determine document category based on file extension.
        
        Args:
            file_path: Path to the file
            
        Returns:
            Category of the document
        """
        path = Path(file_path)
        ext = path.suffix.lower()
        
        # Fall back to extension-based categorization
        return EXTENSION_CATEGORIES.get(ext, 'Document')
    
    def process_file(self, file_path: str) -> Optional[Document]:
        """
        Process a single file into a Document object.
        Extracts text, classifies document, uploads to S3, and prepares for indexing.
        
        Args:
            file_path: Path to the file
            
        Returns:
            Document object if successful, None otherwise
        """
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
            text = self.file_processor.extract_text(file_path)
            
            if not text:
                logger.warning(f"No text extracted from {file_path}")
                self.stats["failed"] += 1
                return None
            
            # Get filename for metadata
            filename = Path(file_path).name
            
            # Classify document and extract metadata using LLM if enabled
            metadata_dict = {}
            doc_type = "unknown"
            
            if self.use_llm_classification and hasattr(self, 'openai_client'):
                try:
                    # Use LLM to classify document and extract metadata
                    logger.info(f"Classifying document: {filename}")
                    llm_result = process_document_llm(filename, text, self.openai_client)
                    
                    # Extract metadata from LLM result
                    doc_type = llm_result.get('document_type', 'Unknown')
                    
                    # Create metadata object from LLM results
                    metadata_obj = Metadata(
                        document_name=filename,
                        subject=llm_result.get('subject', 'Unknown'),
                        status=llm_result.get('status'),
                        case_name=llm_result.get('case_name'),
                        case_number=llm_result.get('case_number'),
                        author=llm_result.get('author'),
                        judge=llm_result.get('judge'),
                        court=llm_result.get('court')
                    )
                    
                    logger.info(f"Document classified as: {doc_type}")
                except Exception as e:
                    logger.error(f"Error during LLM classification: {e}")
                    # Fall back to basic metadata
                    metadata_obj = Metadata(document_name=filename, subject="Unknown", status=None)
            else:
                # Create basic metadata if LLM classification is disabled
                metadata_obj = Metadata(document_name=filename, subject="Unknown", status=None)
            
            # Determine document category based on file extension
            category = self.get_file_category(file_path)
            
            # Upload to S3 if handler is configured
            s3_uri = None
            if self.s3_handler:
                try:
                    # Create a folder structure in S3 based on document type and date
                    today = datetime.now().strftime('%Y/%m/%d')
                    doc_type_folder = doc_type.lower().replace(' ', '_')
                    
                    # Generate a unique filename to avoid collisions
                    unique_id = str(uuid.uuid4())[:8]
                    original_ext = Path(file_path).suffix
                    s3_filename = f"{Path(file_path).stem}_{unique_id}{original_ext}"
                    
                    # Construct the S3 key (path)
                    s3_key = f"{doc_type_folder}/{today}/{s3_filename}"
                    
                    # Upload the file to S3
                    s3_uri = self.s3_handler.upload_file(file_path, s3_key)
                    
                    if s3_uri:
                        logger.info(f"Uploaded to S3: {s3_uri}")
                        self.stats["uploaded"] += 1
                    else:
                        logger.warning(f"Failed to upload to S3: {file_path}")
                except Exception as e:
                    logger.error(f"Error uploading to S3: {e}")
                    # Continue processing even if S3 upload fails
            
            # Create document object with additional S3 info if available
            document_properties = {
                'file_path': s3_uri if s3_uri else file_path,  # Use S3 URI if available
                'text': text,
                'metadata': metadata_obj,
                'doc_type': doc_type,
                'category': category,
                'hash_value': file_hash,
                'created_at': datetime.now()
            }
            
            # Add S3 URI as additional metadata if available
            if s3_uri:
                document_properties['s3_uri'] = s3_uri
                
            document = Document(**document_properties)
            
            self.stats["processed"] += 1
            return document
            
        except Exception as e:
            logger.error(f"Error processing file {file_path}: {e}")
            self.stats["failed"] += 1
            return None
    
    def process_directory(self, directory_path: str, file_extensions: Optional[List[str]] = None) -> Dict[str, int]:
        """
        Process all files in a directory and its subdirectories.
        
        Args:
            directory_path: Path to the directory
            file_extensions: Optional list of file extensions to process
            
        Returns:
            Dictionary with processing statistics
        """
        logger.info(f"Processing directory: {directory_path}")
        
        # Reset stats for this processing run
        self.stats = {
            "total": 0,
            "processed": 0,
            "indexed": 0,
            "uploaded": 0,
            "skipped": 0,
            "failed": 0,
            "unsupported_format": 0,
            "too_large": 0
        }
        
        # Find all files in directory and subdirectories
        file_paths = []
        total_files_found = 0
        
        for root, _, files in os.walk(directory_path):
            for file in files:
                file_path = os.path.join(root, file)
                total_files_found += 1
                
                # Get file extension
                ext = os.path.splitext(file)[1].lower()
                
                # Skip files that don't match the extensions filter
                if file_extensions and ext not in file_extensions:
                    logger.debug(f"Skipping file with non-matching extension: {file_path}")
                    self.stats["skipped"] += 1
                    continue
                
                # Check file size first
                try:
                    if os.path.getsize(file_path) > self.file_processor.max_file_size:
                        logger.warning(f"File too large to process: {file_path}")
                        self.stats["too_large"] += 1
                        self.stats["skipped"] += 1
                        continue
                except Exception as e:
                    logger.error(f"Error checking file size {file_path}: {e}")
                    self.stats["failed"] += 1
                    continue
                
                # Check if the file format is supported by textract
                textract_supported = [
                    '.csv', '.doc', '.docx', '.eml', '.epub', '.gif', '.htm', '.html',
                    '.jpeg', '.jpg', '.json', '.log', '.mp3', '.msg', '.odt', '.ogg',
                    '.pdf', '.png', '.pptx', '.ps', '.psv', '.rtf', '.tab', '.tff',
                    '.tif', '.tiff', '.tsv', '.txt', '.wav', '.xls', '.xlsx'
                ]
                
                if ext not in textract_supported:
                    logger.warning(f"File extension {ext} not supported by textract: {file_path}")
                    self.stats["unsupported_format"] += 1
                    self.stats["skipped"] += 1
                    continue
                
                # Skip files that can't be processed for other reasons
                if not self.file_processor.can_process(file_path):
                    self.stats["skipped"] += 1
                    continue
                    
                file_paths.append(file_path)
        
        self.stats["total"] = len(file_paths)
        logger.info(f"Found {total_files_found} total files")
        logger.info(f"Skipped {self.stats['skipped']} files (too large: {self.stats['too_large']}, unsupported format: {self.stats['unsupported_format']})")
        logger.info(f"Will process {len(file_paths)} files")
        
        # Process files with progress bar
        documents_batch = []
        with tqdm(total=len(file_paths), desc="Processing files") as pbar:
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
            
        logger.info(f"Processing complete. Stats: {self.stats}")
        return self.stats
