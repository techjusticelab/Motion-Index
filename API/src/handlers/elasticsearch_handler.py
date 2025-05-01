"""
Elasticsearch handler for document indexing and search.
"""
import json
import logging
from typing import List, Dict, Any, Optional, Tuple
from elasticsearch import Elasticsearch, helpers

from src.models.document import Document, Metadata
from src.utils.constants import (
    ES_DEFAULT_HOST,
    ES_DEFAULT_PORT,
    ES_DEFAULT_INDEX,
    ES_DOCUMENT_MAPPING,
    ES_BULK_CHUNK_SIZE
)
from src.utils.text_normalizer import normalize_court_name, group_similar_court_names

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
    
    def index_document(self, document: Document) -> str:
        """
        Index a single document into Elasticsearch and return its ID.
    
        Args:
            document: Document object to index

        Returns:
            Document ID if successful, raises exception otherwise
        """
        try:
            # Convert document to dictionary - use to_dict() instead of dict()
            doc_dict = document.to_dict()

            # Normalize court name if present in metadata
            if doc_dict.get('metadata', {}).get('court'):
                doc_dict['metadata']['court'] = normalize_court_name(doc_dict['metadata']['court'])

            # Get the document hash for the ID
            doc_id = getattr(document, 'hash', None) or doc_dict.get('hash_value')
            print(f"Document ID: {doc_id}")
            # Index the document
            response = self.es.index(
                index=self.index_name, 
                body=doc_dict, 
                id=doc_id,
                refresh=True  # Ensure the document is immediately available for search
            )

            logger.info(f"Indexed document: {getattr(document, 'file_name', 'unknown')}")
            return response['_id']
        except Exception as e:
            logger.error(f"Error indexing document: {e}")
            raise
       
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
        actions = []
        for doc in documents:
            # Convert document to dictionary using to_dict()
            doc_dict = doc.to_dict()
            
            # Normalize court name if present in metadata
            if doc_dict.get('metadata', {}).get('court'):
                doc_dict['metadata']['court'] = normalize_court_name(doc_dict['metadata']['court'])
            
            # Get the document hash for the ID
            doc_id = getattr(doc, 'hash', None) or doc_dict.get('hash_value')
            
            actions.append({
                "_index": self.index_name,
                "_id": doc_id,
                "_source": doc_dict
            })
        
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
    
    def document_exists_by_id(self, document_id: str) -> bool:
        """
        Check if a document exists by ID.
        
        Args:
            document_id: ID of the document
            
        Returns:
            True if document exists, False otherwise
        """
        try:
            return self.es.exists(index=self.index_name, id=document_id)
        except Exception as e:
            logger.error(f"Error checking document existence by ID: {e}")
            return False
    
    def get_document(self, document_id: str) -> Dict[str, Any]:
        """
        Retrieve a document by ID.
        
        Args:
            document_id: ID of the document
            
        Returns:
            Document data with added ID field
        """
        try:
            response = self.es.get(
                index=self.index_name,
                id=document_id
            )
            
            # Combine document data with its ID
            doc_data = response['_source']
            doc_data['id'] = response['_id']
            
            return doc_data
        except Exception as e:
            logger.error(f"Error retrieving document: {e}")
            raise
    
    def update_document_metadata(self, document_id: str, metadata: Metadata) -> bool:
        """
        Update metadata fields for a document.
        
        Args:
            document_id: ID of the document
            metadata: Dictionary of metadata fields to update
            
        Returns:
            True if successful, False otherwise
        """
        try:
            # Prepare the update body
            update_body = {
                "doc": {
                    "metadata": metadata
                }
            }
            print(metadata)
            # If the metadata includes a new file name, update it
                
            # If the metadata includes document type or category, update those too
            if "doc_type" in metadata:
                update_body["doc"]["doc_type"] = metadata["doc_type"]
            
            if "category" in metadata:
                update_body["doc"]["category"] = metadata["category"]
            
            # Normalize court name if present
            if metadata.get('court'):
                update_body["doc"]["metadata"]["court"] = normalize_court_name(metadata['court'])
            
            if metadata.get('legal_tags'):
                # Normalize legal tags if needed
                update_body["doc"]["metadata"]["legal_tags"] = [tag.strip() for tag in metadata['legal_tags']]
            if metadata.get('judge'):
                # Normalize judge name if needed
                update_body["doc"]["metadata"]["judge"] = metadata['judge'].strip()

            print(f"Updating document {document_id} with metadata: {update_body}")
            self.es.update(
                index=self.index_name,
                id=document_id,
                body=update_body,
                refresh=True
            )
            
            logger.info(f"Updated metadata for document {document_id}")
            return True
        except Exception as e:
            logger.error(f"Error updating document metadata: {e}")
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
                        use_fuzzy: bool = False,
                        legal_tags_match_all: bool = False) -> Dict[str, Any]:
        """
        Search for documents with advanced filtering options.
        
        Args:
            query: Optional text query string
            doc_type: Optional document type to filter by (e.g., "Motion", "Affidavit")
            metadata_filters: Optional dictionary of metadata field filters
                Example: {"judge": "Smith", "court": "District Court", "legal_tags": ["Criminal", "Family"]}
            date_range: Optional date range filter for document timestamps
                Example: {"start": "2023-01-01", "end": "2023-12-31"}
            size: Maximum number of results to return
            from_value: Starting index for pagination
            sort_by: Field to sort results by (e.g., "created_at", "metadata.timestamp")
            sort_order: Sort order ("asc" or "desc")
            use_fuzzy: Whether to use fuzzy matching for the query (default: False)
            legal_tags_match_all: Whether to match all tags (AND logic) or any tag (OR logic) when filtering by legal_tags
        
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
                        must_clauses.append({
                            "multi_match": {
                                "query": query,
                                "fields": [
                                    "text^1", 
                                    "metadata.subject^2",
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
                                    "text^2",
                                    "metadata.subject^3",
                                    "metadata.case_name^3",
                                    "file_name^2.5"
                                ],
                                "type": "phrase",
                                "boost": 2.0
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
                    if value is not None and value != "" and value != []:
                        # IMPORTANT: Special handling for legal_tags field
                        if field == 'legal_tags':
                            # Print debug information
                            print(f"Applying legal_tags filter: {value}, match_all: {legal_tags_match_all}")
                            
                            if isinstance(value, list) and len(value) > 0:
                                if legal_tags_match_all:
                                    # AND logic - all tags must match
                                    for tag in value:
                                        filter_clauses.append({
                                            "term": {"metadata.legal_tags.keyword": tag}
                                        })
                                else:
                                    # OR logic - any tag can match
                                    filter_clauses.append({
                                        "terms": {"metadata.legal_tags.keyword": value}
                                    })
                                    
                                    # When using OR logic, add a function score to sort by number of matching tags
                                    # This will prioritize documents with more matching tags
                                    if not sort_by:  # Only apply if no explicit sort is requested
                                        function_score = {
                                            "function_score": {
                                                "query": {"match_all": {}},
                                                "functions": [
                                                    {
                                                        "script_score": {
                                                            "script": {
                                                                "source": """
                                                                    def tags = params.tags;
                                                                    def doc_tags = doc['metadata.legal_tags.keyword'];
                                                                    int matches = 0;
                                                                    if (!doc_tags.empty) {
                                                                        for (tag in tags) {
                                                                            if (doc_tags.contains(tag)) {
                                                                                matches++;
                                                                            }
                                                                        }
                                                                    }
                                                                    return matches;
                                                                """,
                                                                "params": {
                                                                    "tags": value
                                                                }
                                                            }
                                                        }
                                                    }
                                                ],
                                                "boost_mode": "replace"
                                            }
                                        }
                                        must_clauses.append(function_score)
                            elif isinstance(value, str) and value:
                                # Create a term query for a single tag
                                filter_clauses.append({
                                    "term": {"metadata.legal_tags.keyword": value}
                                })
                        # Special handling for court field with normalization
                        elif field == 'court':
                            if isinstance(value, list):
                                # Handle multiple court selections
                                court_should_clauses = []
                                for court in value:
                                    normalized_court = normalize_court_name(court)
                                    court_should_clauses.append({
                                        "bool": {
                                            "should": [
                                                {"term": {"metadata.court.keyword": normalized_court}},
                                                {"term": {"metadata.court.keyword": court}},
                                                {"match": {"metadata.court": {"query": normalized_court, "operator": "and"}}}
                                            ],
                                            "minimum_should_match": 1
                                        }
                                    })
                                
                                if court_should_clauses:
                                    filter_clauses.append({
                                        "bool": {
                                            "should": court_should_clauses,
                                            "minimum_should_match": 1
                                        }
                                    })
                            else:
                                # Single court selection
                                normalized_court = normalize_court_name(value)
                                filter_clauses.append({
                                    "bool": {
                                        "should": [
                                            {"term": {"metadata.court.keyword": normalized_court}},
                                            {"term": {"metadata.court.keyword": value}},
                                            {"match": {"metadata.court": {"query": normalized_court, "operator": "and"}}}
                                        ],
                                        "minimum_should_match": 1
                                    }
                                })
                        # Handle other fields
                        elif isinstance(value, list):
                            if len(value) > 0:
                                # For list values, use terms query (OR condition)
                                filter_clauses.append({
                                    "terms": {f"metadata.{field}.keyword": value}
                                })
                        else:
                            # For single values, use term query
                            filter_clauses.append({
                                "term": {f"metadata.{field}.keyword": value}
                            })
            
            # Date range filter
            if date_range:
                date_filter = {"range": {"metadata.timestamp": {}}}
                if "start" in date_range and date_range["start"]:
                    date_filter["range"]["metadata.timestamp"]["gte"] = date_range["start"]
                if "end" in date_range and date_range["end"]:
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
                if sort_by == "created_at":
                    sort_field = "metadata.timestamp"
                else:
                    sort_field = sort_by
                    if not sort_by.endswith(".keyword") and sort_by != "metadata.timestamp":
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
            
            # Print the full search body for debugging
            print("Search query:")
            print(json.dumps(search_body, indent=2))
            
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
            logger.exception(e)  # Log the full exception traceback
            return {
                "total": 0,
                "hits": [],
                "page_size": size,
                "from": from_value
            }
    def get_legal_tags(self) -> List[str]:
        """
        Get a list of all legal tags (categories) used in the documents.
        
        Returns:
            List of unique legal tags
        """
        try:
            response = self.es.search(
                index=self.index_name,
                body={
                    "size": 0,
                    "aggs": {
                        "legal_tags": {
                            "terms": {
                                "field": "metadata.legal_tags.keyword",
                                "size": 50  # Get up to 50 different legal tags
                            }
                        }
                    }
                }
            )
            
            tags = []
            for bucket in response.get("aggregations", {}).get("legal_tags", {}).get("buckets", []):
                tags.append(bucket["key"])
                
            return tags
        except Exception as e:
            logger.error(f"Error getting legal tags: {e}")
            return []
        
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
            field: Metadata field name (e.g., "judge", "court", "legal_tags")
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
                
            # For court field, we need more results to properly normalize and deduplicate
            agg_size = size * 3 if field == "court" else size
            
            search_body = {
                "size": 0,
                "aggs": {
                    "field_values": {
                        "terms": {
                            "field": f"{agg_field}.keyword",
                            "size": agg_size
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
            
            # For court field, normalize and deduplicate the values
            if field == "court":
                values = group_similar_court_names(values)
                # Sort alphabetically for better UX
                values.sort()
                # Limit to the requested size after deduplication
                values = values[:size]
                
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
            
            # Get date range of documents from metadata.timestamp
            date_range_response = self.es.search(
                index=self.index_name,
                body={
                    "size": 0,
                    "aggs": {
                        "min_date": {"min": {"field": "metadata.timestamp"}},
                        "max_date": {"max": {"field": "metadata.timestamp"}}
                    }
                }
            )
            
            aggs = date_range_response.get("aggregations", {})
            min_date = aggs.get("min_date", {}).get("value_as_string")
            max_date = aggs.get("max_date", {}).get("value_as_string")
            
            # Get legal tags statistics
            legal_tags = self.get_legal_tags()
            
            return {
                "total_documents": total_docs,
                "document_types": doc_types,
                "legal_tags": legal_tags,
                "date_range": {
                    "oldest": min_date,
                    "newest": max_date
                }
            }
        
        except Exception as e:
            logger.error(f"Error getting document stats: {e}")
            return {"total_documents": 0}