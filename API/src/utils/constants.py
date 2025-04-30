"""
Constants used throughout the application.
"""
from typing import Dict, Set, Any

# File Processing Constants
MAX_FILE_SIZE_DEFAULT = 50 * 1024 * 1024  # 50MB

# These formats are supported by the textract library
# Reference: https://github.com/deanmalmgren/textract
SUPPORTED_FORMATS: Set[str] = {
    '.csv', '.doc', '.docx', '.eml', '.epub', '.gif', '.htm', '.html',
    '.jpeg', '.jpg', '.json', '.log', '.mp3', '.msg', '.odt', '.ogg',
    '.pdf', '.png', '.pptx', '.ps', '.psv', '.rtf', '.tab', '.tff',
    '.tif', '.tiff', '.tsv', '.txt', '.wav', '.xls', '.xlsx', 'wpd'
}

# Document Types
DOCUMENT_TYPES = {
    'motion': 'Motion',
    'petition': 'Petition',
    'order': 'Order',
    'brief': 'Brief',
    'report': 'Report',
    'exhibit': 'Exhibit',
    'memorandum': 'Memorandum',
    'response': 'Response',
    'opposition': 'Opposition',
    'complaint': 'Complaint',
    'answer': 'Answer',
    'discovery_request': 'Discovery Request',
    'discovery_response': 'Discovery Response',
    'notice': 'Notice',
    'declaration': 'Declaration',
    'affidavit': 'Affidavit',
    'judgment': 'Judgment',
    'transcript': 'Transcript',
    'settlement_agreement': 'Settlement Agreement',
    'unknown': 'Unknown'
}

# File Extension Mappings
MIME_TYPES: Dict[str, str] = {
    '.wpd': 'application/x-wordperfect',
    '.wp': 'application/x-wordperfect',
    '.wp5': 'application/x-wordperfect',
    '.mot': 'application/msword',
    '.mtn': 'application/msword',
    '.pet': 'text/plain',
    '.sup': 'text/plain',
    '.wrt': 'text/plain',
    '.reh': 'text/plain'
}

# Document Category Mappings
EXTENSION_CATEGORIES: Dict[str, str] = {
    '.pdf': 'PDF Document',
    '.docx': 'Word Document',
    '.doc': 'Word Document',
    '.wpd': 'WordPerfect Document',
    '.wp': 'WordPerfect Document',
    '.wp5': 'WordPerfect Document',
    '.txt': 'Text Document',
    '.mot': 'Motion',
    '.mtn': 'Motion',
    '.pet': 'Petition',
    '.sup': 'Supplement',
    '.ord': 'Order',
    '.rep': 'Report',
    '.ppt': 'Presentation',
    '.pptx': 'Presentation'
}

# Elasticsearch Constants
ES_DEFAULT_HOST = "localhost"
ES_DEFAULT_PORT = 9200
ES_DEFAULT_INDEX = "documents"
ES_BULK_CHUNK_SIZE = 500

# Elasticsearch Document Mapping
ES_DOCUMENT_MAPPING = {
    "mappings": {
        "properties": {
            "file_path": {"type": "keyword"},
            "file_name": {"type": "text", "fields": {"keyword": {"type": "keyword"}}},
            "category": {"type": "keyword"},
            "chunk_id": {"type": "integer"},
            "text": {"type": "text", "analyzer": "english"},
            "doc_type": {"type": "keyword"},
            "s3_uri": {"type": "keyword"},
            "metadata": {
                "properties": {
                    "document_name": {"type": "text", "fields": {"keyword": {"type": "keyword"}}},
                    "subject": {"type": "text", "fields": {"keyword": {"type": "keyword"}}},
                    "status": {"type": "keyword"},
                    "case_name": {"type": "text", "fields": {"keyword": {"type": "keyword"}}},
                    "case_number": {"type": "keyword"},
                    "author": {"type": "text", "fields": {"keyword": {"type": "keyword"}}},
                    "judge": {"type": "text", "fields": {"keyword": {"type": "keyword"}}},
                    "court": {"type": "text", "fields": {"keyword": {"type": "keyword"}}},
                    "timestamp": {"type": "date"}
                }
            },
            "embedding": {"type": "dense_vector", "dims": 384},
            "hash": {"type": "keyword"},
            "created_at": {"type": "date"}
        }
    }
}

# Processing Constants
DEFAULT_MAX_WORKERS = 4
DEFAULT_BATCH_SIZE = 100

# LLM Constants
OPENAI_MODEL = "gpt-3.5-turbo"
