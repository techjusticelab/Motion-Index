from fastapi import FastAPI, Request
from fastapi.middleware.cors import CORSMiddleware
from fastapi.responses import JSONResponse
import sys
import os
import json

# Add the parent directory to sys.path to import from server.py
sys.path.append(os.path.dirname(os.path.dirname(os.path.abspath(__file__))))

# Import the FastAPI app and handlers from server.py
from server import app as main_app, es_handler

# Create a new FastAPI app for Vercel
app = FastAPI()

# Add CORS middleware with permissive settings for development
app.add_middleware(
    CORSMiddleware,
    allow_origins=["*"],  # Allow all origins for now
    allow_credentials=True,
    allow_methods=["*"],
    allow_headers=["*"],
)

# Direct API endpoint implementations for Vercel
@app.get("/document-types")
async def get_document_types():
    try:
        return es_handler.get_document_types()
    except Exception as e:
        return JSONResponse(status_code=500, content={"error": str(e)})

@app.get("/legal-tags")
async def get_legal_tags():
    try:
        return es_handler.get_legal_tags()
    except Exception as e:
        return JSONResponse(status_code=500, content={"error": str(e)})

@app.post("/search")
async def search_documents(request: Request):
    try:
        data = await request.json()
        # Calculate offset from page number
        from_value = (data.get('page', 1) - 1) * data.get('size', 10)
        
        # Build metadata filters
        metadata_filters = {}
        for field in ['case_number', 'case_name', 'judge', 'court', 'author', 'status', 'legal_tags']:
            if field in data and data[field]:
                metadata_filters[field] = data[field]
        
        # Execute search
        results = es_handler.search_documents(
            query=data.get('query'),
            doc_type=data.get('doc_type'),
            metadata_filters=metadata_filters if metadata_filters else None,
            date_range=data.get('date_range'),
            size=data.get('size', 10),
            sort_by=data.get('sort_by'),
            sort_order=data.get('sort_order', 'desc'),
            use_fuzzy=data.get('use_fuzzy', False),
            from_value=from_value,
            legal_tags_match_all=data.get('legal_tags_match_all', False)
        )
        return results
    except Exception as e:
        return JSONResponse(status_code=500, content={"error": str(e)})

@app.get("/health")
async def health_check():
    """Health check endpoint for Vercel."""
    return {"status": "healthy", "environment": "vercel"}

# Vercel serverless function handler
def handler(request: Request):
    # Extract the actual path from the request
    path = request.url.path
    print(f"Received request: {path}")
    
    # Strip /api prefix if present
    if path.startswith('/api/'):
        path_parts = path.split('/api/', 1)
        if len(path_parts) > 1:
            # Modify the request scope to change the path
            request.scope["path"] = f"/{path_parts[1]}"
            print(f"Modified path: {request.scope['path']}")
    
    return app(request.scope, request.receive, request.send)
