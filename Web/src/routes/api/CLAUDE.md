# routes/api/ Directory

## Purpose
Server-side API endpoints for the SvelteKit application.

## Files
- `search/+server.ts` - Document search proxy to backend API

## API Routes
- `GET /api/search?q=query` - Search documents via backend

## Implementation
- Proxies requests to Python FastAPI backend
- Handles CORS and authentication headers
- Formats responses for frontend consumption
- Error handling and logging

## Backend Integration
- Connects to FastAPI server on port 8000
- Forwards search queries to Elasticsearch
- Returns formatted document results
- Maintains stateless API design