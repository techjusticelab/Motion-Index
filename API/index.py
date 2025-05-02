import os
import logging
import json
from typing import Dict, List, Optional, Any, Union
from fastapi import FastAPI, Query, HTTPException, File, UploadFile, Request
from fastapi.responses import JSONResponse
from fastapi.middleware.cors import CORSMiddleware
from pydantic import BaseModel, Field

# Setup logging
logger = logging.getLogger(__name__)

# Initialize FastAPI app
app = FastAPI(
    title="Motion-Index API",
    description="API for searching legal documents in Elasticsearch",
    version="1.0.0"
)

# Add CORS middleware - allow all origins for Vercel deployment
app.add_middleware(
    CORSMiddleware,
    allow_origins=["*"],
    allow_credentials=True,
    allow_methods=["*"],
    allow_headers=["*"],
)

# Pydantic models for request/response
class SearchRequest(BaseModel):
    query: Optional[str] = None
    doc_type: Optional[str] = None
    case_number: Optional[str] = None
    case_name: Optional[str] = None
    judge: Optional[Union[str, List[str]]] = None
    court: Optional[Union[str, List[str]]] = None
    author: Optional[str] = None
    status: Optional[str] = None
    legal_tags: Optional[Union[str, List[str]]] = None
    legal_tags_match_all: bool = Field(default=False, description="Whether to match all tags (AND) or any tag (OR)")
    date_range: Optional[Dict[str, str]] = None
    size: int = Field(default=10, ge=1, le=100)
    sort_by: Optional[str] = None
    sort_order: str = Field(default="desc", pattern="^(asc|desc)$")
    page: int = Field(default=1)
    use_fuzzy: bool = Field(default=False, description="Whether to use fuzzy matching for search queries")

class MetadataFieldRequest(BaseModel):
    field: str
    prefix: Optional[str] = None
    size: int = Field(default=20, ge=1, le=100)

# Import real API components
from src.handlers.elasticsearch_handler import ElasticsearchHandler
from src.utils.constants import (
    ES_DEFAULT_HOST,
    ES_DEFAULT_PORT,
    ES_DEFAULT_INDEX
)
from src.core.document_processor import DocumentProcessor

# Initialize Elasticsearch handler
es_handler = ElasticsearchHandler(
    host=os.environ.get("ES_HOST", ES_DEFAULT_HOST),
    port=int(os.environ.get("ES_PORT", ES_DEFAULT_PORT)),
    index_name=os.environ.get("ES_INDEX", ES_DEFAULT_INDEX),
    username=os.environ.get("ES_USERNAME"),
    password=os.environ.get("ES_PASSWORD"),
    api_key=os.environ.get("ES_API_KEY"),
    cloud_id=os.environ.get("ES_CLOUD_ID"),
    use_ssl=os.environ.get("ES_USE_SSL", "True").lower() == "true"
)

# Initialize document processor
processor = DocumentProcessor(
    # Elasticsearch settings
    es_username=os.environ.get("ES_USERNAME"),
    es_password=os.environ.get("ES_PASSWORD"),
    es_api_key=os.environ.get("ES_API_KEY"),
    es_host=os.environ.get("ES_HOST", ES_DEFAULT_HOST),
    es_port=int(os.environ.get("ES_PORT", ES_DEFAULT_PORT)),
    es_index=os.environ.get("ES_INDEX", ES_DEFAULT_INDEX),
    es_cloud_id=os.environ.get("ES_CLOUD_ID"),
    es_use_ssl=os.environ.get("ES_USE_SSL", "True").lower() == "true",
    # S3 settings
    s3_bucket=os.environ.get("S3_BUCKET_NAME"),
    s3_region=os.environ.get("AWS_REGION"),
    s3_access_key=os.environ.get("AWS_ACCESS_KEY_ID"),
    s3_secret_key=os.environ.get("AWS_SECRET_ACCESS_KEY"),
    # Processing settings
    max_workers=os.environ.get("MAX_WORKERS", 4),
    batch_size=os.environ.get("BATCH_SIZE", 100),
)

logger.info("Using real API components")

@app.get("/document-types")
async def get_document_types():
    """Get a list of all document types and their counts."""
    try:
        return es_handler.get_document_types()
    except Exception as e:
        logger.error(f"Error in get_document_types: {str(e)}")
        return JSONResponse(status_code=500, content={"error": str(e)})

@app.get("/legal-tags")
async def get_legal_tags():
    """Get a list of all legal tags."""
    try:
        return es_handler.get_legal_tags()
    except Exception as e:
        logger.error(f"Error in get_legal_tags: {str(e)}")
        return JSONResponse(status_code=500, content={"error": str(e)})

@app.post("/search")
async def search_documents(search_request: SearchRequest):
    """Search for documents with advanced filtering options."""
    try:
        return es_handler.search_documents(search_request)
    except Exception as e:
        logger.error(f"Error in search_documents: {str(e)}")
        return JSONResponse(status_code=500, content={"error": str(e)})

@app.get("/metadata-fields")
async def get_metadata_fields():
    """Get a list of all available metadata fields for filtering."""
    try:
        # Return the standard metadata fields structure
        return {
            "fields": [
                {"id": "doc_type", "name": "Document Type", "type": "string"},
                {"id": "category", "name": "Category", "type": "string"},
                {"id": "metadata.case_number", "name": "Case Number", "type": "string"},
                {"id": "metadata.case_name", "name": "Case Name", "type": "string"},
                {"id": "metadata.judge", "name": "Judge", "type": "string"},
                {"id": "metadata.court", "name": "Court", "type": "string"},
                {"id": "metadata.legal_tags", "name": "Legal Tags", "type": "string"},
                {"id": "metadata.author", "name": "Author", "type": "string"},
                {"id": "metadata.status", "name": "Status", "type": "string"},
                {"id": "created_at", "name": "Date", "type": "date"}
            ]
        }
    except Exception as e:
        logger.error(f"Error in get_metadata_fields: {str(e)}")
        return JSONResponse(status_code=500, content={"error": str(e)})

@app.get("/all-field-options")
async def get_all_field_options():
    """Get all available options for multiple fields at once."""
    try:
        # Fields to get options for
        fields = [
            "doc_type",
            "category",
            "case_number",
            "judge",
            "court",
            "legal_tags",
            "status"
        ]
        
        # Get options for each field
        result = {}
        for field in fields:
            size = 10000  # Use a larger size to get more options
            options = es_handler.get_metadata_field_values(field, size=size)
            result[field] = options
            
        return result
    except Exception as e:
        logger.error(f"Error in get_all_field_options: {str(e)}")
        return JSONResponse(status_code=500, content={"error": str(e)})

@app.get("/health")
async def health_check():
    """Health check endpoint."""
    try:
        # Check Elasticsearch connection
        es_status = "unknown"
        try:
            # Attempt to ping Elasticsearch
            es_status = "connected" if es_handler.ping() else "disconnected"
        except Exception as e:
            es_status = f"error: {str(e)}"
        
        return {
            "status": "healthy", 
            "environment": "vercel" if os.environ.get('VERCEL') == '1' else "local",
            "elasticsearch": es_status
        }
    except Exception as e:
        logger.error(f"Error in health_check: {str(e)}")
        return JSONResponse(status_code=500, content={"error": str(e)})

@app.get("/")
async def root():
    """Root endpoint to check if the API is running."""
    return {"message": "Motion-Index API is running"}

