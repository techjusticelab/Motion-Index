"""
FastAPI server for Motion-Index document search API.
"""
import os
import logging
from typing import Dict, List, Optional, Any, Union
from fastapi import FastAPI, Query, HTTPException, File, UploadFile
from fastapi.middleware.cors import CORSMiddleware
import dotenv
import uvicorn
from pydantic import BaseModel, Field

from src.handlers.elasticsearch_handler import ElasticsearchHandler
from src.utils.constants import (
    ES_DEFAULT_HOST,
    ES_DEFAULT_PORT,
    ES_DEFAULT_INDEX
)
from src.core.document_processor import DocumentProcessor
from src.models.document import Document

# Setup logging
logger = logging.getLogger(__name__)

class MetadataUpdateRequest(BaseModel):
    document_id: str
    metadata: Dict[str, Any]

# Load environment variables
dotenv.load_dotenv()

# Initialize FastAPI app
app = FastAPI(
    title="Motion-Index API",
    description="API for searching legal documents in Elasticsearch",
    version="1.0.0"
)

# Add CORS middleware
app.add_middleware(
    CORSMiddleware,
    allow_origins=["*"],  # In production, replace with specific origins
    allow_credentials=True,
    allow_methods=["*"],
    allow_headers=["*"],
)

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

# Initialize DocumentProcessor
document_processor = DocumentProcessor(
    es_host=os.environ.get("ES_HOST", ES_DEFAULT_HOST),
    es_port=int(os.environ.get("ES_PORT", ES_DEFAULT_PORT)),
    es_index=os.environ.get("ES_INDEX", ES_DEFAULT_INDEX),
    es_username=os.environ.get("ES_USERNAME"),
    es_password=os.environ.get("ES_PASSWORD"),
    es_api_key=os.environ.get("ES_API_KEY"),
    es_cloud_id=os.environ.get("ES_CLOUD_ID"),
    use_llm_classification=True
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
    date_range: Optional[Dict[str, str]] = None
    size: int = Field(default=10, ge=1, le=100)
    sort_by: Optional[str] = None
    sort_order: str = Field(default="desc", pattern="^(asc|desc)$")
    page: int = Field(default=1, ge=1)
    use_fuzzy: bool = Field(default=False, description="Whether to use fuzzy matching for search queries")

class MetadataFieldRequest(BaseModel):
    field: str
    prefix: Optional[str] = None
    size: int = Field(default=20, ge=1, le=100)

@app.get("/")
async def root():
    """Root endpoint to check if the API is running."""
    return {"message": "Motion-Index API is running"}

@app.get("/health")
async def health_check():
    """Health check endpoint."""
    try:
        # Check Elasticsearch connection
        if es_handler.es.ping():
            return {"status": "healthy", "elasticsearch": "connected"}
        else:
            return {"status": "unhealthy", "elasticsearch": "disconnected"}
    except Exception as e:
        return {"status": "unhealthy", "error": str(e)}

@app.post("/search")
async def search_documents(search_request: SearchRequest):
    """
    Search for documents with advanced filtering options.
    """
    try:
        # Calculate offset from page number
        from_value = (search_request.page - 1) * search_request.size
        
        # Build metadata filters from individual parameters
        metadata_filters = {}
        if search_request.case_number:
            metadata_filters["case_number"] = search_request.case_number
        if search_request.case_name:
            metadata_filters["case_name"] = search_request.case_name
        if search_request.judge:
            metadata_filters["judge"] = search_request.judge
        if search_request.court:
            metadata_filters["court"] = search_request.court
        if search_request.author:
            metadata_filters["author"] = search_request.author
        if search_request.status:
            metadata_filters["status"] = search_request.status
            
        # Execute search
        results = es_handler.search_documents(
            query=search_request.query,
            doc_type=search_request.doc_type,
            metadata_filters=metadata_filters if metadata_filters else None,
            date_range=search_request.date_range,
            size=search_request.size,
            sort_by=search_request.sort_by,
            sort_order=search_request.sort_order,
            use_fuzzy=search_request.use_fuzzy
        )
        return results
    except Exception as e:
        raise HTTPException(status_code=500, detail=str(e))


@app.get("/legal-tags")
async def get_legal_tags():
    """
    Get a list of all legal types and their counts.
    """
    try:
        return es_handler.get_legal_tags()
    except Exception as e:
        raise HTTPException(status_code=500, detail=str(e))
    
@app.get("/document-types")
async def get_document_types():
    """
    Get a list of all document types and their counts.
    """
    try:
        return es_handler.get_document_types()
    except Exception as e:
        raise HTTPException(status_code=500, detail=str(e))

@app.post("/metadata-field-values")
async def get_metadata_field_values(request: MetadataFieldRequest):
    """
    Get unique values for a specific metadata field, optionally filtered by prefix.
    """
    try:
        return es_handler.get_metadata_field_values(
            field=request.field,
            prefix=request.prefix,
            size=request.size
        )
    except Exception as e:
        raise HTTPException(status_code=500, detail=str(e))

@app.get("/document-stats")
async def get_document_stats():
    """
    Get statistics about the indexed documents.
    """
    try:
        return es_handler.get_document_stats()
    except Exception as e:
        raise HTTPException(status_code=500, detail=str(e))

@app.get("/metadata-fields")
async def get_metadata_fields():
    """
    Get a list of all available metadata fields for filtering.
    """
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

@app.get("/all-field-options")
async def get_all_field_options():
    """
    Get all available options for multiple fields at once.
    This is used to populate dropdowns and filter options in the UI.
    """
    try:
        # Fields to get options for (excluding case_name and author as requested)
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
            # Use a larger size to get more options
            options = es_handler.get_metadata_field_values(field, size=100)
            result[field] = options
            
        return result
    except Exception as e:
        raise HTTPException(status_code=500, detail=str(e))

@app.post("/categorise")
async def categorise_document(file: UploadFile = File(...)):
    """
    Upload and categorise a document using the document processor.
    Returns the full document data from Elasticsearch including all metadata fields.
    """
    try:
        # Save the uploaded file temporarily
        temp_file_path = f"/tmp/{file.filename}"
        with open(temp_file_path, "wb") as temp_file:
            temp_file.write(await file.read())

        # Process and categorise the file using DocumentProcessor
        document = document_processor.process_file(temp_file_path)
        
        if not document:
            raise HTTPException(status_code=400, detail="Failed to process the document")
        
        # Index the document in Elasticsearch to get an ID
        doc_id = es_handler.index_document(document)
        
        # Retrieve the full document data from Elasticsearch
        full_document = es_handler.get_document(doc_id)
        
        # Clean up the temporary file
        os.remove(temp_file_path)

        return {
            "message": "Document categorised successfully",
            "document": full_document
        }
    except Exception as e:
        raise HTTPException(status_code=500, detail=f"Error categorising document: {str(e)}")

@app.post("/update-metadata")
async def update_document_metadata(request: MetadataUpdateRequest):
    """
    Update metadata fields for a document.
    """
    try:
        # Check if document exists
        if not es_handler.document_exists_by_id(request.document_id):
            raise HTTPException(status_code=404, detail="Document not found")
        
        # Update the document metadata
        success = es_handler.update_document_metadata(request.document_id, request.metadata)
        
        if not success:
            raise HTTPException(status_code=500, detail="Failed to update document metadata")
        
        # Retrieve the updated document
        updated_document = es_handler.get_document(request.document_id)
        
        return {
            "message": "Document metadata updated successfully",
            "document": updated_document
        }
    except HTTPException:
        raise
    except Exception as e:
        raise HTTPException(status_code=500, detail=f"Error updating metadata: {str(e)}")



if __name__ == "__main__":
    uvicorn.run("server:app", host="0.0.0.0", port=8000, reload=True)