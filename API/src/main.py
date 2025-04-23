#!/usr/bin/env python3
"""
Main entry point for the Motion-Index document processing system.
"""
import os
import sys
import argparse
import logging
import dotenv
from pathlib import Path

# Add the parent directory to sys.path to enable imports
sys.path.append(os.path.dirname(os.path.dirname(os.path.abspath(__file__))))

from src.core.document_processor import DocumentProcessor

# Configure logging
logging.basicConfig(
    level=logging.INFO,
    format='%(asctime)s - %(name)s - %(levelname)s - %(message)s',
    handlers=[
        logging.FileHandler("motion_index.log"),
        logging.StreamHandler()
    ]
)
logger = logging.getLogger("motion_index")


def main():
    """
    Main function to process documents based on command line arguments.
    """
    # Load environment variables
    dotenv.load_dotenv()
    
    # Parse command line arguments
    parser = argparse.ArgumentParser(
        description="Process legal documents, classify them, upload to S3, and index in Elasticsearch"
    )
    parser.add_argument("directory", help="Directory containing documents to process")
    
    # Elasticsearch options
    es_group = parser.add_argument_group('Elasticsearch Options')
    es_group.add_argument("--es-host", default=os.environ.get('ES_HOST', 'localhost'), 
                         help="Elasticsearch host or cloud URL")
    es_group.add_argument("--es-port", type=int, default=int(os.environ.get('ES_PORT', 9200)), 
                         help="Elasticsearch port (usually 443 for Elastic Cloud)")
    es_group.add_argument("--es-index", default=os.environ.get('ES_INDEX', 'documents'), 
                         help="Elasticsearch index name")
    es_group.add_argument("--es-username", default=os.environ.get('ES_USERNAME'), 
                         help="Elasticsearch username (for basic auth)")
    es_group.add_argument("--es-password", default=os.environ.get('ES_PASSWORD'), 
                         help="Elasticsearch password (for basic auth)")
    es_group.add_argument("--es-api-key", default=os.environ.get('ES_API_KEY'), 
                         help="Elasticsearch API key (for Elastic Cloud)")
    es_group.add_argument("--es-cloud-id", default=os.environ.get('ES_CLOUD_ID'), 
                         help="Elasticsearch Cloud ID (optional for Elastic Cloud)")
    es_group.add_argument("--es-ssl", action="store_true", 
                         default=os.environ.get('ES_USE_SSL', 'False').lower() == 'true', 
                         help="Use SSL for Elasticsearch connection (default for Elastic Cloud)")
    
    # S3 options
    s3_group = parser.add_argument_group('S3 Options')
    s3_group.add_argument("--s3-bucket", default=os.environ.get('S3_BUCKET_NAME'), 
                         help="S3 bucket name for document storage")
    s3_group.add_argument("--s3-region", default=os.environ.get('AWS_REGION'), 
                         help="AWS region for S3")
    s3_group.add_argument("--s3-access-key", default=os.environ.get('AWS_ACCESS_KEY_ID'), 
                         help="AWS access key ID")
    s3_group.add_argument("--s3-secret-key", default=os.environ.get('AWS_SECRET_ACCESS_KEY'), 
                         help="AWS secret access key")
    
    # Processing options
    proc_group = parser.add_argument_group('Processing Options')
    proc_group.add_argument("--workers", type=int, default=int(os.environ.get('MAX_WORKERS', 4)), 
                           help="Number of worker threads")
    proc_group.add_argument("--batch-size", type=int, default=int(os.environ.get('BATCH_SIZE', 100)), 
                           help="Batch size for indexing")
    proc_group.add_argument("--max-file-size", type=int, default=int(os.environ.get('MAX_FILE_SIZE', 50)), 
                           help="Maximum file size in MB")
    proc_group.add_argument("--extensions", nargs="+", 
                           help="List of file extensions to process (e.g., pdf docx txt)")
    proc_group.add_argument("--process-all", action="store_true", 
                           help="Process all files regardless of extension")
    proc_group.add_argument("--force", action="store_true", 
                           help="Process files even if they already exist in the index")
    proc_group.add_argument("--no-llm", action="store_true", 
                           help="Disable LLM-based document classification")
    
    args = parser.parse_args()
    
    # Validate directory
    if not os.path.isdir(args.directory):
        logger.error(f"Directory not found: {args.directory}")
        return 1
    
    # Convert extensions to proper format
    extensions = None
    if not args.process_all and args.extensions:
        extensions = [ext if ext.startswith('.') else f'.{ext}' for ext in args.extensions]
    
    # Initialize document processor
    processor = DocumentProcessor(
        # Elasticsearch settings
        es_host=args.es_host,
        es_port=args.es_port,
        es_index=args.es_index,
        es_username=args.es_username,
        es_password=args.es_password,
        es_api_key=args.es_api_key,  # Added API key support
        es_cloud_id=args.es_cloud_id,  # Added Cloud ID support
        es_use_ssl=args.es_ssl,
        # S3 settings
        s3_bucket=args.s3_bucket,
        s3_region=args.s3_region,
        s3_access_key=args.s3_access_key,
        s3_secret_key=args.s3_secret_key,
        # Processing settings
        max_workers=args.workers,
        batch_size=args.batch_size,
        max_file_size=args.max_file_size * 1024 * 1024,  # Convert to bytes
        skip_existing=not args.force,
        use_llm_classification=not args.no_llm
    )
    
    # Process the directory
    logger.info(f"Starting document processing for directory: {args.directory}")
    stats = processor.process_directory(args.directory, file_extensions=extensions)
    
    # Print summary
    print("\nProcessing Summary:")
    print(f"Total files found: {stats['total']}")
    print(f"Files processed: {stats['processed']}")
    print(f"Documents indexed: {stats['indexed']}")
    print(f"Files uploaded to S3: {stats['uploaded']}")
    print(f"Files skipped: {stats['skipped']}")
    print(f"Failed files: {stats['failed']}")
    
    return 0


if __name__ == "__main__":
    exit(main())
