"""
Elasticsearch handler for document indexing and search.
"""
import json
import logging
from typing import List, Dict, Any, Optional, Tuple
from elasticsearch import Elasticsearch, helpers

from src.models.document import Document
from src.utils.constants import (
    ES_DEFAULT_HOST,
    ES_DEFAULT_PORT,
    ES_DEFAULT_INDEX,
    ES_DOCUMENT_MAPPING,
    ES_BULK_CHUNK_SIZE
)

# Configure logging
logger = logging.getLogger("elasticsearch_handler")


class ElasticsearchHandler:
    """
    Handles interactions with Elasticsearch for document indexing and retrieval.
    """
    
    def __init__(
        self, 
        host: str = ES_DEFAULT_HOST, 
        port: int = ES_DEFAULT_PORT, 
        index_name: str = ES_DEFAULT_INDEX,
        username: Optional[str] = None,
        password: Optional[str] = None,
        api_key: Optional[str] = None,
        cloud_id: Optional[str] = None,
        use_ssl: bool = True
    ):
        """
        Initialize Elasticsearch connection.
        
        Args:
            host: Elasticsearch host (or cloud URL for Elastic Cloud)
            port: Elasticsearch port (usually 443 for Elastic Cloud)
            index_name: Name of the index to use
            username: Optional username for basic authentication
            password: Optional password for basic authentication
            api_key: Optional API key for Elastic Cloud authentication
            cloud_id: Optional Cloud ID for Elastic Cloud
            use_ssl: Whether to use SSL for connection (default True)
        """
        self.index_name = index_name
        
        # Configure ES client
        es_config = {}
        
        # Handle different connection methods
        if cloud_id:
            # Connect using Cloud ID
            es_config['cloud_id'] = cloud_id
            logger.info(f"Using Elastic Cloud ID: {cloud_id}")
        else:
            # Connect using host/port
            es_config['hosts'] = [f"http{'s' if use_ssl else ''}://{host}:{port}"]
            logger.info(f"Using Elasticsearch at {host}:{port}")
        
        # Handle different authentication methods
        if api_key:
            # API key authentication (preferred for Elastic Cloud)
            es_config['api_key'] = api_key
            logger.info("Using API key authentication")
        elif username and password:
            # Basic authentication
            es_config['basic_auth'] = (username, password)
            logger.info("Using basic authentication")
            
        try:
            self.es = Elasticsearch(**es_config)
            if not self.es.ping():
                logger.error("Could not connect to Elasticsearch")
                raise ConnectionError("Failed to connect to Elasticsearch")
            logger.info("Successfully connected to Elasticsearch")
        except Exception as e:
            logger.error(f"Error connecting to Elasticsearch: {e}")
            raise
    
    def create_index(self, mapping: Optional[Dict[str, Any]] = None) -> bool:
        """
        Create the document index with proper mappings.
        
        Args:
            mapping: Optional custom mapping for the index
            
        Returns:
            True if index was created or already exists, False on error
        """
        if mapping is None:
            mapping = ES_DOCUMENT_MAPPING
            
        try:
            if not self.es.indices.exists(index=self.index_name):
                self.es.indices.create(index=self.index_name, body=mapping)
                logger.info(f"Created index '{self.index_name}'")
                return True
            else:
                logger.info(f"Index '{self.index_name}' already exists")
                return True
        except Exception as e:
            logger.error(f"Error creating index: {e}")
            return False
    
    def index_document(self, document: Document) -> bool:
        """
        Index a single document into Elasticsearch.
        
        Args:
            document: Document object to index
            
        Returns:
            True if successful, False otherwise
        """
        try:
            doc_dict = document.to_dict()
            self.es.index(index=self.index_name, document=doc_dict, id=document.hash)
            logger.info(f"Indexed document: {document.file_name}")
            return True
        except Exception as e:
            logger.error(f"Error indexing document {document.file_name}: {e}")
            return False
    
    def bulk_index_documents(self, documents: List[Document], chunk_size: int = ES_BULK_CHUNK_SIZE) -> Tuple[int, int]:
        """
        Bulk index multiple documents into Elasticsearch.
        
        Args:
            documents: List of Document objects to index
            chunk_size: Size of chunks for bulk indexing
            
        Returns:
            Tuple of (success_count, error_count)
        """
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
        """
        Check if a document with the given hash already exists.
        
        Args:
            doc_hash: Hash of the document
            
        Returns:
            True if document exists, False otherwise
        """
        try:
            return self.es.exists(index=self.index_name, id=doc_hash)
        except Exception as e:
            logger.error(f"Error checking document existence: {e}")
            return False
    
    def search_documents(self, 
                        query: Optional[str] = None, 
                        doc_type: Optional[str] = None,
                        metadata_filters: Optional[Dict[str, Any]] = None,
                        date_range: Optional[Dict[str, str]] = None,
                        size: int = 10,
                        from_value: int = 0,
                        sort_by: Optional[str] = None,
                        sort_order: str = "desc",
                        use_fuzzy: bool = False) -> Dict[str, Any]:
        """
        Search for documents with advanced filtering options.
        
        Args:
            query: Optional text query string
            doc_type: Optional document type to filter by (e.g., "Motion", "Affidavit")
            metadata_filters: Optional dictionary of metadata field filters
                Example: {"judge": "Smith", "court": "District Court"}
            date_range: Optional date range filter for document timestamps
                Example: {"start": "2023-01-01", "end": "2023-12-31"}
            size: Maximum number of results to return
            from_value: Starting index for pagination
            sort_by: Field to sort results by (e.g., "created_at", "metadata.timestamp")
            sort_order: Sort order ("asc" or "desc")
            use_fuzzy: Whether to use fuzzy matching for the query (default: False)
                When False, performs exact matching which is better for specific terms like "DUI"
                When True, allows for typos and variations using Elasticsearch's fuzzy matching
        
        Returns:
            Dictionary containing search results and metadata
        """
        try:
            # Build the query
            search_body = {
                "size": size,
                "from": from_value
            }
            must_clauses = []
            filter_clauses = []
            
            # Text search query
            if query:
                # Check if the query contains special operators
                has_operators = any(op in query for op in ['OR', 'AND', '"', '*', '~'])
                
                if has_operators:
                    # Use query_string for advanced query syntax
                    must_clauses.append({
                        "query_string": {
                            "query": query,
                            "fields": [
                                "text^1",  # Text content with normal weight
                                "metadata.subject^2",  # Subject with higher weight
                                "metadata.case_name^2", 
                                "file_name^1.5"
                            ],
                            "default_operator": "AND",
                            "analyze_wildcard": True,
                            "allow_leading_wildcard": True
                        }
                    })
                else:
                    # For exact matching, use a combination of match and match_phrase
                    # The match_phrase ensures exact phrases are prioritized
                    if use_fuzzy:
                        # Use fuzzy matching only when explicitly requested
                        must_clauses.append({
                            "multi_match": {
                                "query": query,
                                "fields": [
                                    "text^1",  # Text content with normal weight
                                    "metadata.subject^2",  # Subject with higher weight
                                    "metadata.case_name^2", 
                                    "file_name^1.5"
                                ],
                                "type": "best_fields",
                                "fuzziness": "AUTO"
                            }
                        })
                    else:
                        # For non-fuzzy search, use exact matching
                        # First add a standard multi_match without fuzziness
                        must_clauses.append({
                            "multi_match": {
                                "query": query,
                                "fields": [
                                    "text^1",  # Text content with normal weight
                                    "metadata.subject^2",  # Subject with higher weight
                                    "metadata.case_name^2", 
                                    "file_name^1.5"
                                ],
                                "type": "best_fields",
                                "operator": "AND"
                            }
                        })
                        
                        # Also add a match_phrase query to prioritize exact phrases
                        must_clauses.append({
                            "multi_match": {
                                "query": query,
                                "fields": [
                                    "text^2",  # Higher weight for exact matches
                                    "metadata.subject^3",
                                    "metadata.case_name^3",
                                    "file_name^2.5"
                                ],
                                "type": "phrase",
                                "boost": 2.0  # Give exact matches a higher boost
                            }
                        })
            
            # Document type filter
            if doc_type:
                filter_clauses.append({
                    "term": {"doc_type.keyword": doc_type}
                })
            
            # Metadata filters
            if metadata_filters:
                for field, value in metadata_filters.items():
                    if value is not None:
                        # Handle different field types appropriately
                        if isinstance(value, list):
                            # For list values, use terms query (OR condition)
                            filter_clauses.append({
                                "terms": {f"metadata.{field}": value}
                            })
                        else:
                            # For single values, use term query
                            filter_clauses.append({
                                "term": {f"metadata.{field}": value}
                            })
            
            # Date range filter
            if date_range:
                date_filter = {"range": {"metadata.timestamp": {}}}
                if "start" in date_range:
                    date_filter["range"]["metadata.timestamp"]["gte"] = date_range["start"]
                if "end" in date_range:
                    date_filter["range"]["metadata.timestamp"]["lte"] = date_range["end"]
                filter_clauses.append(date_filter)
            
            # Combine all query parts
            if must_clauses or filter_clauses:
                search_body["query"] = {"bool": {}}
                if must_clauses:
                    search_body["query"]["bool"]["must"] = must_clauses
                if filter_clauses:
                    search_body["query"]["bool"]["filter"] = filter_clauses
            else:
                # If no specific query, match all documents
                search_body["query"] = {"match_all": {}}
            
            # Add sorting if specified
            if sort_by:
                # Use .keyword suffix for text fields when sorting
                sort_field = sort_by
                if sort_by not in ["created_at"] and not sort_by.endswith(".keyword"):
                    sort_field = f"{sort_by}.keyword"
                search_body["sort"] = [{sort_field: {"order": sort_order}}]
                
            # Add highlighting for search results
            if query:
                search_body["highlight"] = {
                    "fields": {
                        "text": {
                            "fragment_size": 150,
                            "number_of_fragments": 3,
                            "pre_tags": ["<strong>"],
                            "post_tags": ["</strong>"]
                        },
                        "metadata.subject": {
                            "fragment_size": 150,
                            "number_of_fragments": 1,
                            "pre_tags": ["<strong>"],
                            "post_tags": ["</strong>"]
                        }
                    }
                }
            
            # Execute the search
            response = self.es.search(
                index=self.index_name,
                body=search_body
            )
            
            # Extract hits and total count
            hits = response["hits"]["hits"]
            total = response["hits"]["total"]["value"]
            
            # Format results with highlighting if available
            formatted_hits = []
            for hit in hits:
                doc = hit["_source"]
                
                # Add highlighting if available
                if "highlight" in hit:
                    doc["highlight"] = hit["highlight"]
                
                formatted_hits.append(doc)
            
            # Return structured response with total and hits
            return {
                "total": total,
                "hits": formatted_hits,
                "page_size": size,
                "from": from_value
            }
        except Exception as e:
            logger.error(f"Error searching documents: {e}")
            return {
                "total": 0,
                "hits": [],
                "page_size": size,
                "from": from_value
            }
    
    def get_document_types(self) -> Dict[str, int]:
        """
        Get a list of all document types and their counts.
        
        Returns:
            Dictionary mapping document types to their counts
        """
        try:
            response = self.es.search(
                index=self.index_name,
                body={
                    "size": 0,  # We only want aggregations, not actual documents
                    "aggs": {
                        "doc_types": {
                            "terms": {
                                "field": "doc_type.keyword",
                                "size": 100  # Get up to 100 different document types
                            }
                        }
                    }
                }
            )
            
            result = {}
            for bucket in response.get("aggregations", {}).get("doc_types", {}).get("buckets", []):
                result[bucket["key"]] = bucket["doc_count"]
                
            return result
        except Exception as e:
            logger.error(f"Error getting document types: {e}")
            return {}
    
    def get_metadata_field_values(self, field: str, prefix: Optional[str] = None, size: int = 20) -> List[str]:
        """
        Get unique values for a specific metadata field, optionally filtered by prefix.
        Useful for autocomplete functionality in search interfaces.
        
        Args:
            field: Metadata field name (e.g., "judge", "court")
            prefix: Optional prefix to filter values (for autocomplete)
            size: Maximum number of values to return
            
        Returns:
            List of unique field values
        """
        try:
            # Build the aggregation query
            agg_field = f"metadata.{field}"
            if field in ["doc_type", "category"]:
                agg_field = field
                
            search_body = {
                "size": 0,
                "aggs": {
                    "field_values": {
                        "terms": {
                            "field": f"{agg_field}.keyword",
                            "size": size
                        }
                    }
                }
            }
            
            # Add prefix filter if provided
            if prefix:
                search_body["query"] = {
                    "prefix": {f"{agg_field}.keyword": prefix}
                }
                
            response = self.es.search(
                index=self.index_name,
                body=search_body
            )
            
            values = []
            for bucket in response.get("aggregations", {}).get("field_values", {}).get("buckets", []):
                if bucket["key"] and bucket["key"] != "null" and bucket["key"] != "None":
                    values.append(bucket["key"])
                    
            return values
        except Exception as e:
            logger.error(f"Error getting metadata field values: {e}")
            return []
    
    def get_document_stats(self) -> Dict[str, Any]:
        """
        Get statistics about the indexed documents.
        
        Returns:
            Dictionary with document statistics
        """
        try:
            # Get total document count
            count_response = self.es.count(index=self.index_name)
            total_docs = count_response.get("count", 0)
            
            # Get document type breakdown
            doc_types = self.get_document_types()
            
            # Get date range of documents
            date_range_response = self.es.search(
                index=self.index_name,
                body={
                    "size": 0,
                    "aggs": {
                        "min_date": {"min": {"field": "created_at"}},
                        "max_date": {"max": {"field": "created_at"}}
                    }
                }
            )
            
            aggs = date_range_response.get("aggregations", {})
            min_date = aggs.get("min_date", {}).get("value_as_string")
            max_date = aggs.get("max_date", {}).get("value_as_string")
            
            return {
                "total_documents": total_docs,
                "document_types": doc_types,
                "date_range": {
                    "oldest": min_date,
                    "newest": max_date
                }
            }
        except Exception as e:
            logger.error(f"Error getting document stats: {e}")
            return {"total_documents": 0}
