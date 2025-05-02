from fastapi import FastAPI, Request
from fastapi.middleware.cors import CORSMiddleware
from fastapi.routing import APIRoute
import sys
import os

# Add the parent directory to sys.path to import from server.py
sys.path.append(os.path.dirname(os.path.dirname(os.path.abspath(__file__))))

# Import the FastAPI app from server.py
from server import app as main_app

# Create a new FastAPI app for Vercel
app = FastAPI()

# Function to modify route paths to work with Vercel's routing
def create_route_with_prefix(route):
    # Create a new route with the same endpoint but prefixed path
    new_route = APIRoute(
        path=route.path,  # The original path is kept for now
        endpoint=route.endpoint,
        methods=route.methods,
        name=route.name,
        response_model=route.response_model,
        tags=route.tags,
        dependencies=route.dependencies,
        description=route.description,
        summary=route.summary,
        response_description=route.response_description,
    )
    return new_route

# Copy all routes from the main app
for route in main_app.routes:
    app.routes.append(create_route_with_prefix(route))

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
    # For debugging
    print(f"Received request: {request.url.path}")
    return app(request.scope, request.receive, request.send)
