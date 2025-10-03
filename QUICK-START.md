# Quick Start Guide - Motion-Index Local Document Server

## Overview

This setup provides local document serving to replace S3 functionality while connecting to your existing remote Elasticsearch cluster. The system will serve documents that were previously stored in S3 from local storage.

## Prerequisites

- Docker and Docker Compose installed
- Remote Elasticsearch cluster access
- 10GB of free disk space for documents

## Configuration

### 1. Update Environment Variables

Copy and update your environment file:
```bash
cp .env.template .env
```

Edit `.env` with your credentials:
```env
# Supabase Configuration
SUPABASE_URL=https://your-project.supabase.co
SUPABASE_ANON_KEY=your-anon-key
SUPABASE_SERVICE_KEY=your-service-key

# Storage Configuration
USE_LOCAL_STORAGE=true
STORAGE_PATH=/app/data

# Remote Elasticsearch Configuration
ES_HOST=my-elasticsearch-project-fe3a6c.es.us-east-1.aws.elastic.cloud
ES_PORT=443
ES_USE_SSL=true
ES_INDEX=cpda
# Add authentication if needed:
# ES_USERNAME=your-username
# ES_PASSWORD=your-password
# ES_API_KEY=your-api-key
```

### 2. Document Setup

Place your documents in the local storage structure that matches your S3 bucket structure:

```bash
# Example S3 structure: s3://cpda-documents/memorandum/2025/05/01/document.pdf
# Local structure: API/data/documents/memorandum/2025/05/01/document.pdf

# Create directories as needed
mkdir -p API/data/documents/memorandum/2025/05/01
mkdir -p API/data/documents/brief/2025/05/01
mkdir -p API/data/documents/order/2025/05/01

# Copy your PDF documents to appropriate directories
# cp /path/to/your/documents/*.pdf API/data/documents/memorandum/2025/05/01/
```

### 3. Start the Services

```bash
# Build and start the document server
docker-compose -f docker-compose.simple.yml up --build

# Or run in background
docker-compose -f docker-compose.simple.yml up --build -d
```

### 4. Access Points

- **Web Interface**: http://localhost:3000
- **API Server**: http://localhost:8000
- **API Documentation**: http://localhost:8000/docs
- **Remote Elasticsearch**: (via your configured endpoint)

### 5. Test Document Serving

The system will automatically serve documents using S3 URIs from your Elasticsearch data:

```bash
# Test S3 URI serving (example)
curl -I "http://localhost:8000/api/documents/s3%3A//cpda-documents/memorandum/2025/05/01/test-document.pdf"

# Test direct path serving
curl -I "http://localhost:8000/api/documents/memorandum/2025/05/01/test-document.pdf"
```

### 6. Document Structure

Your local storage should mirror your S3 bucket structure:

```
API/data/documents/
├── memorandum/
│   └── 2025/
│       └── 05/
│           └── 01/
│               ├── document1.pdf
│               └── document2.pdf
├── brief/
│   └── 2025/
│       └── 05/
└── order/
    └── 2025/
        └── 05/
```

## Option 2: Local Development (No Docker)

### 1. Run API Server Locally

```bash
cd API

# Install dependencies
pip install -r requirements.simple.txt

# Set environment variables for remote Elasticsearch
export USE_LOCAL_STORAGE=true
export STORAGE_PATH=./data
export ES_HOST=my-elasticsearch-project-fe3a6c.es.us-east-1.aws.elastic.cloud
export ES_PORT=443
export ES_USE_SSL=true
export ES_INDEX=cpda

# Run the server
python server.py
```

### 2. Run Web Frontend Locally

```bash
cd Web

# Install dependencies
npm install

# Set environment variables
export PUBLIC_API_URL=http://localhost:8000
export PUBLIC_SUPABASE_URL=your-supabase-url
export PUBLIC_SUPABASE_ANON_KEY=your-anon-key

# Start development server
npm run dev
```

## Troubleshooting

### Common Issues

1. **Port Conflicts**
   ```bash
   # Check if ports are in use
   netstat -tulpn | grep -E ':3000|:8000'

   # Stop conflicting services
   docker-compose -f docker-compose.simple.yml down
   ```

2. **Permission Issues with Document Storage**
   ```bash
   # Fix permissions for data directory
   sudo chown -R $USER:$USER API/data
   chmod -R 755 API/data
   ```

3. **Remote Elasticsearch Connection Issues**
   ```bash
   # Test connection to remote Elasticsearch
   curl -k "https://my-elasticsearch-project-fe3a6c.es.us-east-1.aws.elastic.cloud:443"

   # Check API logs for ES connection errors
   docker-compose -f docker-compose.simple.yml logs api
   ```

4. **Docker Build Issues**
   ```bash
   # Clean rebuild
   docker-compose -f docker-compose.simple.yml down --volumes
   docker-compose -f docker-compose.simple.yml build --no-cache
   docker-compose -f docker-compose.simple.yml up
   ```

5. **Document Not Found Errors**
   ```bash
   # Check if document exists in local storage
   ls -la API/data/documents/memorandum/2025/05/01/

   # Check API logs for file serving attempts
   docker-compose -f docker-compose.simple.yml logs api | grep "Serving document"
   ```

### Check Service Health

```bash
# Check container status
docker-compose -f docker-compose.simple.yml ps

# Check logs
docker-compose -f docker-compose.simple.yml logs api
docker-compose -f docker-compose.simple.yml logs web

# Check API health
curl http://localhost:8000/docs

# Test remote Elasticsearch connection
curl -k "https://my-elasticsearch-project-fe3a6c.es.us-east-1.aws.elastic.cloud:443/_cluster/health"
```

### Test File Serving

```bash
# Test serving a document that exists in your Elasticsearch data
# Example S3 URI from your data
curl -I "http://localhost:8000/api/documents/s3%3A//cpda-documents/memorandum/2025/05/01/Cedillo%201538.5%20Supp_f12b5be9.pdf"

# Test direct path serving
curl -I "http://localhost:8000/api/documents/memorandum/2025/05/01/test-document.pdf"
```

## Configuration Details

### Document Storage Structure
Your local documents should mirror the S3 bucket structure:
```
API/data/documents/
├── memorandum/
│   └── 2025/
│       └── 05/
│           └── 01/
│               └── Cedillo 1538.5 Supp_f12b5be9.pdf
├── brief/
└── order/
```

### Environment Variables
- `USE_LOCAL_STORAGE=true` - Enables local document serving
- `STORAGE_PATH=/app/data` - Document storage location in container
- `ES_HOST=my-elasticsearch-project-fe3a6c.es.us-east-1.aws.elastic.cloud` - Remote ES host
- `ES_PORT=443` - ES port with SSL
- `ES_USE_SSL=true` - Enable SSL for ES connection
- `ES_INDEX=cpda` - Your Elasticsearch index name

### API Endpoints
- `GET /api/documents/{path}` - Serve documents (handles S3 URIs)
- `POST /search` - Search documents in remote Elasticsearch
- `POST /categorise` - Upload and process new documents
- `GET /docs` - API documentation

## Document Migration

To migrate documents from S3 to local storage:

1. **Download from S3**: Download your S3 documents preserving directory structure
2. **Match Structure**: Place files in `API/data/documents/` matching S3 paths
3. **Test Access**: Verify documents are accessible via the API endpoints

## Demo Workflow

1. **Prepare documents**: Place PDFs in `API/data/documents/` directory structure
2. **Start services**: `docker-compose -f docker-compose.simple.yml up -d`
3. **Access web**: http://localhost:3000
4. **Search existing documents**: Your Elasticsearch data will show documents
5. **View documents**: Click search results - they'll be served from local storage
6. **Upload new documents**: New uploads will be stored locally

Your system now serves documents locally while using your existing Elasticsearch cluster!