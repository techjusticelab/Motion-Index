from typing import Dict, Set, Any

# File Processing Constants
MAX_FILE_SIZE_DEFAULT = 50 * 1024 * 1024  # 50MB
SUPPORTED_FORMATS: Set[str] = {
    '.pdf', '.docx', '.doc', '.txt', '.rtf', '.odt', '.html', '.htm',
    '.wpd', '.wp', '.wp5', '.pptx', '.ppt', '.xlsx', '.xls'
}

# Metadata Extraction Patterns
METADATA_PATTERNS: Dict[str, str] = {
    'case_name': r'(?:Case[:\s]+|Matter of:?|In re:?)([A-Za-z0-9\s\.,]+v\.?[A-Za-z0-9\s\.,]+)',
    'court': r'(?:COURT:?|Court of:?|IN THE)([A-Za-z\s]+COURT[A-Za-z\s]*)',
    'judge': r'(?:Judge|JUDGE|Hon\.|Honorable)[:\s]+([A-Za-z\s\.]+)',
    'date': r'(?:Date[d]?:?|Filed on:?)[\s]+(\d{1,2}[/-]\d{1,2}[/-]\d{2,4}|\w+ \d{1,2},? \d{4})',
    'subject': r'(?:Subject:?|RE:|REGARDING:?)[\s]+([A-Za-z0-9\s\.,]+)'
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

# Path-based Document Categories
PATH_CATEGORIES: Dict[str, str] = {
    'motion': 'Motion',
    'petition': 'Petition',
    'order': 'Order',
    'brief': 'Brief',
    'report': 'Report',
    'exhibit': 'Exhibit'
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
            "metadata": {
                "properties": {
                    "document_name": {"type": "text", "fields": {"keyword": {"type": "keyword"}}},
                    "subject": {"type": "text", "fields": {"keyword": {"type": "keyword"}}},
                    "status": {"type": "keyword"},
                    "case_name": {"type": "text", "fields": {"keyword": {"type": "keyword"}}},
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