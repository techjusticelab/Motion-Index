from elasticsearch import Elasticsearch, helpers
import logging
from typing import List, Dict, Any, Optional
import json

# Import constants
from constants import (
    ES_DEFAULT_HOST,
    ES_DEFAULT_PORT, 
    ES_DEFAULT_INDEX,
    ES_DOCUMENT_MAPPING,
    ES_BULK_CHUNK_SIZE
)

# Configure logging
logging.basicConfig(
    level=logging.INFO,
    format='%(asctime)s - %(name)s - %(levelname)s - %(message)s',
    handlers=[
        logging.FileHandler("elasticsearch.log"),
        logging.StreamHandler()
    ]
)
logger = logging.getLogger("elasticsearch_handler")

class ElasticsearchHandler:
    """Handles interactions with Elasticsearch for document indexing and retrieval"""
    
    def __init__(
        self, 
        host: str = ES_DEFAULT_HOST, 
        port: int = ES_DEFAULT_PORT, 
        index_name: str = ES_DEFAULT_INDEX,
        username: Optional[str] = None,
        password: Optional[str] = None,
        use_ssl: bool = False
    ):
        """Initialize Elasticsearch connection"""
        self.index_name = index_name
        
        # Configure ES client
        es_config = {
            'hosts': [f"http{'s' if use_ssl else ''}://{host}:{port}"]
        }
        
        # Add authentication if provided
        if username and password:
            es_config['basic_auth'] = (username, password)
            
        try:
            self.es = Elasticsearch(**es_config)
            if not self.es.ping():
                logger.error("Could not connect to Elasticsearch")
                raise ConnectionError("Failed to connect to Elasticsearch")
            logger.info(f"Connected to Elasticsearch at {host}:{port}")
        except Exception as e:
            logger.error(f"Error connecting to Elasticsearch: {e}")
            raise
    
    def create_index(self, mapping: Dict[str, Any] = None) -> bool:
        """Create the document index with proper mappings"""
        if mapping is None:
            mapping = ES_DOCUMENT_MAPPING
            
        try:
            if not self.es.indices.exists(index=self.index_name):
                self.es.indices.create(index=self.index_name, body=mapping)
                logger.info(f"Created index '{self.index_name}'")
                return True
            else:
                logger.info(f"Index '{self.index_name}' already exists")
                return False
        except Exception as e:
            logger.error(f"Error creating index: {e}")
            return False
    
    def index_document(self, document: Document) -> bool:
        """Index a single document into Elasticsearch"""
        try:
            doc_dict = document.to_dict()
            self.es.index(index=self.index_name, document=doc_dict, id=document.hash)
            logger.info(f"Indexed document: {document.file_name}")
            return True
        except Exception as e:
            logger.error(f"Error indexing document {document.file_name}: {e}")
            return False
    
    def bulk_index_documents(self, documents: List[Document], chunk_size: int = ES_BULK_CHUNK_SIZE) -> tuple:
        """Bulk index multiple documents into Elasticsearch"""
        success_count = 0
        error_count = 0
        
        # Convert documents to ES actions
        actions = [
            {
                "_index": self.index_name,
                "_id": doc.hash,
                "_source": doc.to_dict()
            }
            for doc in documents
        ]
        
        # Bulk index in chunks
        for i in range(0, len(actions), chunk_size):
            chunk = actions[i:i + chunk_size]
            try:
                success, errors = helpers.bulk(
                    self.es, 
                    chunk, 
                    stats_only=False,
                    raise_on_error=False
                )
                success_count += success
                error_count += len(errors)
                if errors:
                    for error in errors:
                        logger.error(f"Bulk indexing error: {json.dumps(error)}")
            except Exception as e:
                logger.error(f"Error during bulk indexing: {e}")
                error_count += len(chunk)
        
        logger.info(f"Bulk indexing complete. Success: {success_count}, Errors: {error_count}")
        return success_count, error_count
    
    def document_exists(self, doc_hash: str) -> bool:
        """Check if a document with the given hash already exists"""
        try:
            return self.es.exists(index=self.index_name, id=doc_hash)
        except Exception as e:
            logger.error(f"Error checking document existence: {e}")
            return False
    
    def search_documents(self, query: str, size: int = 10) -> List[Dict[str, Any]]:
        """Search for documents matching the query"""
        try:
            response = self.es.search(
                index=self.index_name,
                body={
                    "query": {
                        "multi_match": {
                            "query": query,
                            "fields": ["text", "metadata.subject", "metadata.case_name", "file_name"],
                            "type": "best_fields"
                        }
                    },
                    "size": size
                }
            )
            hits = response["hits"]["hits"]
            return [hit["_source"] for hit in hits]
        except Exception as e:
            logger.error(f"Error searching documents: {e}")
            return []