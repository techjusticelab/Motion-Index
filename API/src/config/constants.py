"""
Shared constants for the application.
"""

# Elasticsearch configuration
ES_DEFAULT_HOST = "localhost"
ES_DEFAULT_PORT = 9200
ES_DEFAULT_INDEX = "documents"
ES_DOCUMENT_MAPPING = {
    "properties": {
        "id": {"type": "keyword"},
        "content": {"type": "text"},
        "metadata": {"type": "object"}
    }
}
ES_BULK_CHUNK_SIZE = 500

# File processing configuration
MAX_FILE_SIZE_MB = 50
SUPPORTED_FILE_TYPES = {
    "pdf": "application/pdf",
    "docx": "application/vnd.openxmlformats-officedocument.wordprocessingml.document",
    "txt": "text/plain"
}

# S3 configuration
S3_BUCKET_NAME = "motion-index-documents"
S3_DEFAULT_REGION = "us-east-1"

# Document classification
DOC_CLASSIFICATION_THRESHOLD = 0.7
MAX_CLASSIFICATION_RETRIES = 3
