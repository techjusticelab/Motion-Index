from fastapi import FastAPI, Request
from fastapi.middleware.cors import CORSMiddleware
import sys
import os

# Add the parent directory to sys.path to import from server.py
sys.path.append(os.path.dirname(os.path.dirname(os.path.abspath(__file__))))

# Import the FastAPI app from server.py
from server import app as main_app

# Create a new FastAPI app for Vercel
app = FastAPI()

# Copy all routes from the main app
app.routes = main_app.routes

# Add CORS middleware with the same settings as in server.py
app.add_middleware(
    CORSMiddleware,
    allow_origins=["*"],  # Allow all origins for now
    allow_credentials=True,
    allow_methods=["*"],
    allow_headers=["*"],
)

@app.get("/health")
async def health_check():
    """Health check endpoint for Vercel."""
    return {"status": "healthy", "environment": "vercel"}

# This is required for Vercel serverless function
def handler(request: Request):
    return app(request.scope, request.receive, request.send)
