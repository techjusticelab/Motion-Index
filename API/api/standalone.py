from fastapi import FastAPI, Request
from fastapi.middleware.cors import CORSMiddleware
from fastapi.responses import JSONResponse
import os
import json

# Create a new FastAPI app for Vercel
app = FastAPI()

# Add CORS middleware with permissive settings
app.add_middleware(
    CORSMiddleware,
    allow_origins=["*"],  # Allow all origins for now
    allow_credentials=True,
    allow_methods=["*"],
    allow_headers=["*"],
)

# Simplified mock endpoints for Vercel deployment
@app.get("/document-types")
async def get_document_types():
    """Mock endpoint for document types"""
    return {
        "Motion": 25,
        "Brief": 18,
        "Opinion": 12,
        "Order": 10,
        "Complaint": 8
    }

@app.get("/legal-tags")
async def get_legal_tags():
    """Mock endpoint for legal tags"""
    return [
        "Civil Procedure",
        "Constitutional Law",
        "Criminal Law",
        "Evidence",
        "Intellectual Property",
        "Contract Law",
        "Tort Law"
    ]

@app.post("/search")
async def search_documents(request: Request):
    """Mock search endpoint"""
    try:
        data = await request.json()
        # Return mock search results
        return {
            "total": 5,
            "hits": [
                {
                    "id": "doc1",
                    "file_name": "sample_motion.pdf",
                    "file_path": "/documents/sample_motion.pdf",
                    "text": "This is a sample motion text...",
                    "doc_type": "Motion",
                    "metadata": {
                        "document_name": "Motion to Dismiss",
                        "subject": "Civil Procedure",
                        "case_name": "Smith v. Jones",
                        "case_number": "CV-2023-123",
                        "author": "John Smith",
                        "judge": "Judge Wilson",
                        "legal_tags": ["Civil Procedure"],
                        "court": "District Court"
                    },
                    "created_at": "2023-05-15T10:30:00Z"
                },
                {
                    "id": "doc2",
                    "file_name": "sample_brief.pdf",
                    "file_path": "/documents/sample_brief.pdf",
                    "text": "This is a sample brief text...",
                    "doc_type": "Brief",
                    "metadata": {
                        "document_name": "Appellant Brief",
                        "subject": "Constitutional Law",
                        "case_name": "Brown v. State",
                        "case_number": "CV-2023-456",
                        "author": "Jane Doe",
                        "judge": "Judge Martinez",
                        "legal_tags": ["Constitutional Law"],
                        "court": "Appeals Court"
                    },
                    "created_at": "2023-06-20T14:45:00Z"
                }
            ]
        }
    except Exception as e:
        return JSONResponse(status_code=500, content={"error": str(e)})

@app.get("/metadata-fields")
async def get_metadata_fields():
    """Get a list of all available metadata fields for filtering."""
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
    """Get all available options for multiple fields at once."""
    return {
        "doc_type": ["Motion", "Brief", "Opinion", "Order", "Complaint"],
        "category": ["Civil", "Criminal", "Administrative", "Constitutional"],
        "case_number": ["CV-2023-123", "CV-2023-456", "CR-2023-789"],
        "judge": ["Judge Wilson", "Judge Martinez", "Judge Johnson"],
        "court": ["District Court", "Appeals Court", "Supreme Court"],
        "legal_tags": ["Civil Procedure", "Constitutional Law", "Criminal Law", "Evidence"],
        "status": ["Active", "Closed", "Pending"]
    }

@app.get("/health")
async def health_check():
    """Health check endpoint for Vercel."""
    return {"status": "healthy", "environment": "vercel"}

@app.get("/")
async def root():
    """Root endpoint to check if the API is running."""
    return {"message": "Motion-Index API is running (standalone version for Vercel)"}

# Vercel serverless function handler
def handler(request: Request):
    print(f"Received request: {request.url.path}")
    return app(request.scope, request.receive, request.send)
