# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

Motion-Index is a legal document processing system designed for serverless cloud deployment on DigitalOcean. The system consists of three main components:

1. **API** - Python/FastAPI backend with unified storage (cloud/local) and Elasticsearch integration
2. **Web** - SvelteKit frontend for document search and management
3. **Database** - Managed Elasticsearch with 23,500+ migrated legal documents

## Architecture

### Production Implementation (DigitalOcean)
- **Document Storage**: DigitalOcean Spaces with global CDN
- **Search Engine**: DigitalOcean Managed OpenSearch cluster
- **Authentication**: Supabase for user management
- **File Processing**: Cloud-based text extraction and AI classification
- **Deployment**: DigitalOcean App Platform with auto-scaling

### Local Development
- **Document Storage**: Local filesystem (`API/data/documents/`)
- **Search Engine**: DigitalOcean Managed OpenSearch (same as production)
- **Authentication**: Supabase (same as production)
- **File Processing**: Local text extraction (textract) and optional AI classification (OpenAI)

### API Backend Structure

The API uses a modular handler-based architecture:

**Active Handlers** (`src/handlers/`):
- `elasticsearch_handler.py`: Search operations and document indexing (OpenSearch client)
- `storage_handler.py`: Unified storage (DigitalOcean Spaces + local filesystem)
- `file_processor.py`: Multi-format text extraction
- `document_classifier.py`: AI-powered document classification
- `redaction_handler.py`: PDF redaction detection and processing

**Core Components**:
- `src/core/document_processor.py`: Processing pipeline orchestration
- `src/middleware/auth.py`: JWT authentication with Supabase
- `src/models/`: Data models and schemas
- `server.py`: Main FastAPI application

### Web Frontend Structure

SvelteKit application with:
- **Routes** (`src/routes/`): File-based routing
  - `/`: Document search interface
  - `/auth/*`: Authentication flows
  - `/upload`: Document upload
  - `/account`: User and case management
- **Components** (`src/routes/lib/components/`): Reusable UI components
- **Styling**: Tailwind CSS v4

## Development Commands

### DigitalOcean Deployment (Production)
```bash
# Deploy via GitHub Actions (automatic on push to main)
git push origin main

# Manual deployment with doctl
doctl apps create --spec app.yaml
doctl apps update <app-id> --spec app.yaml

# Migrate PDFs to Spaces
python scripts/migrate_to_spaces.py --base-path ./API/data
```

### Local Development
```bash
# 1. Configure environment (OpenSearch connection)
cd API && cp ../.env.local.example .env  # Edit with OpenSearch credentials

# 2. No local Elasticsearch needed - uses DigitalOcean OpenSearch

# 2. Start API
python server.py  # Port 8000

# 3. Start Web frontend  
cd ../Web && npm install && npm run dev  # Port 5173
```

### Web Frontend Commands
```bash
cd Web
npm install
npm run dev       # Development server (port 5173)
npm run build     # Production build
npm run preview   # Preview production build
npm run check     # Type checking
npm run lint      # Linting
npm run format    # Code formatting
```

### API Backend Commands
```bash
cd API
pip install -r requirements.txt
python server.py       # Run server (port 8000)

# Alternative with uvicorn:
uvicorn server:app --reload --host 0.0.0.0 --port 8000
```

Access local development at http://localhost:5173

## API Endpoints

Main endpoints in `server.py`:

- `GET /health` - System health status
- `POST /search` - Document search
- `POST /categorise` - Document upload and processing
- `GET /api/documents/{path}` - Serve documents from local storage
- `POST /update-metadata` - Update document metadata (auth required)
- `POST /analyze-redactions` - Analyze PDF redactions
- `GET /document-stats` - Index statistics
- `GET /metadata-fields` - Available search fields

## Environment Configuration

### Production (DigitalOcean)
Set via App Platform dashboard or GitHub Secrets:
```env
# DigitalOcean Spaces
DO_SPACES_KEY=your-spaces-key
DO_SPACES_SECRET=your-spaces-secret
DO_SPACES_BUCKET=motion-index-docs
DO_SPACES_REGION=nyc3
USE_CLOUD_STORAGE=true
STORAGE_BACKEND=spaces

# Supabase
SUPABASE_URL=https://your-project.supabase.co
SUPABASE_ANON_KEY=your-anon-key
SUPABASE_SERVICE_KEY=your-service-key

# DigitalOcean Managed OpenSearch
ES_HOST=your-opensearch-host.k.db.ondigitalocean.com
ES_PORT=25060
ES_USERNAME=doadmin
ES_PASSWORD=your-opensearch-password
ES_USE_SSL=true
ES_INDEX=documents
```

### Local Development (.env)
```env
# Local storage
USE_CLOUD_STORAGE=false
STORAGE_BACKEND=local
STORAGE_PATH=./data

# DigitalOcean Managed OpenSearch (same as production)
ES_HOST=your-opensearch-host.k.db.ondigitalocean.com
ES_PORT=25060
ES_USERNAME=doadmin
ES_PASSWORD=your-opensearch-password
ES_USE_SSL=true
ES_INDEX=documents

# Supabase (same as production)
SUPABASE_URL=https://your-project.supabase.co
SUPABASE_ANON_KEY=your-anon-key
SUPABASE_SERVICE_KEY=your-service-key
```

### Web Frontend (.env.local)
```env
PUBLIC_API_URL=http://localhost:8000  # Local dev
# PUBLIC_API_URL=https://your-app.ondigitalocean.app  # Production

PUBLIC_SUPABASE_URL=https://your-project.supabase.co
PUBLIC_SUPABASE_ANON_KEY=your-anon-key
```

## Key Implementation Notes

- **Unified Storage**: `storage_handler.py` automatically switches between DigitalOcean Spaces (production) and local filesystem (development)
- **Document Processing**: Supports PDF, DOCX, TXT, RTF via textract with cloud/local processing
- **Authentication**: Flexible JWT verification supporting multiple algorithms via Supabase
- **Search**: Full-text search with metadata filtering via DigitalOcean Managed OpenSearch
- **File Serving**: CDN redirects for cloud storage, direct serving for local development
- **Auto-scaling**: DigitalOcean App Platform automatically scales based on traffic
- **Migration**: `migrate_to_spaces.py` handles bulk PDF upload and OpenSearch URL updates

## Testing & Quality

Currently no automated tests configured. For development:
- API: Check endpoints at http://localhost:8000/docs (Swagger)
- Frontend: Manual testing with `npm run check` for type safety
- OpenSearch: Use `python test_opensearch.py` script or OpenSearch Dashboard

## Performance Specs

- **Database**: 23,500+ documents indexed  
- **Search Response**: <100ms typical (DigitalOcean Managed OpenSearch)
- **Migration Speed**: 50-100 docs/sec (Python batch processing)
- **Storage**: 7.6GB for complete document set
- **CDN**: Global document delivery with edge caching
- **Auto-scaling**: Handles traffic spikes automatically