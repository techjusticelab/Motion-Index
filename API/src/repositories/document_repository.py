"""
Repository for document-related database operations.
"""

from typing import List, Dict, Optional
from elasticsearch import Elasticsearch
from elasticsearch.helpers import bulk
from src.models.document import Document
from src.config.constants import ES_DEFAULT_INDEX, ES_BULK_CHUNK_SIZE


class DocumentRepository:
    """
    Repository for document-related database operations.
    """
    
    def __init__(self, es_client: Elasticsearch):
        """
        Initialize the document repository.
        
        Args:
            es_client: Elasticsearch client instance
        """
        self.es_client = es_client
        self.index_name = ES_DEFAULT_INDEX
        
    def create_index(self, mappings: Dict) -> bool:
        """
        Create Elasticsearch index with specified mappings.
        
        Args:
            mappings: Elasticsearch index mappings
            
        Returns:
            bool: True if index was created successfully
        """
        try:
            if not self.es_client.indices.exists(index=self.index_name):
                self.es_client.indices.create(
                    index=self.index_name,
                    body={"mappings": mappings}
                )
            return True
        except Exception as e:
            print(f"Error creating index: {e}")
            return False
    
    def index_documents(self, documents: List[Document]) -> int:
        """
        Index multiple documents in bulk.
        
        Args:
            documents: List of Document objects to index
            
        Returns:
            int: Number of successfully indexed documents
        """
        actions = [
            {
                "_index": self.index_name,
                "_id": doc.id,
                "_source": doc.dict()
            }
            for doc in documents
        ]
        
        success, _ = bulk(
            self.es_client,
            actions,
            chunk_size=ES_BULK_CHUNK_SIZE,
            request_timeout=30
        )
        return success
    
    def search_documents(self, query: Dict) -> List[Dict]:
        """
        Search for documents using Elasticsearch query.
        
        Args:
            query: Elasticsearch query dictionary
            
        Returns:
            List[Dict]: List of matching documents
        """
        try:
            response = self.es_client.search(
                index=self.index_name,
                body=query
            )
            return [hit["_source"] for hit in response["hits"]["hits"]]
        except Exception as e:
            print(f"Error searching documents: {e}")
            return []
