# Motion-Index

A comprehensive legal document processing system with serverless cloud deployment on DigitalOcean.

## 🚀 Quick Deploy - DigitalOcean (Recommended)

**Deploy to DigitalOcean App Platform with automatic scaling and CDN:**

1. **Fork this repository** to your GitHub account
2. **Set up DigitalOcean Spaces** for document storage 
3. **Configure GitHub Secrets** with your credentials
4. **Push to main branch** - automatic deployment via GitHub Actions

📖 **[Complete Setup Guide](SETUP-DIGITALOCEAN.md)**

**Features:**
- ✅ **Serverless auto-scaling** - handles traffic spikes automatically
- ✅ **Global CDN** - fast document delivery worldwide  
- ✅ **23,500+ PDF documents** with cloud storage migration
- ✅ **Managed database** - Elasticsearch with automatic backups
- ✅ **$35/month** total cost - 70% cheaper than AWS equivalent

## 🏠 Local Development

**For development and testing:**

```bash
# 1. Start local Elasticsearch
cd Database/es && docker-compose -f docker-compose.standalone.yml up -d

# 2. Configure environment
cd ../../API && cp .env.local.example .env  # Edit with your credentials

# 3. Start API
python server.py

# 4. Start Web frontend
cd ../Web && npm install && npm run dev
```

**Access**: http://localhost:5173

## 📁 Repository Structure

- **`API/`** - FastAPI backend with unified storage (local/cloud)
- **`Web/`** - SvelteKit frontend application  
- **`Database/`** - Elasticsearch setup for local development
- **`scripts/`** - Migration and deployment utilities
- **`.github/workflows/`** - Automated deployment to DigitalOcean

## 🏗️ Architecture

### Cloud Deployment (Production)
- **API Backend**: FastAPI on DigitalOcean App Platform (auto-scaling)
- **Web Frontend**: SvelteKit with SSR (auto-scaling)
- **Database**: Managed Elasticsearch with automatic backups
- **Storage**: DigitalOcean Spaces with global CDN
- **Authentication**: Supabase for user management

### Local Development  
- **API Backend**: FastAPI with local file storage (port 8000)
- **Web Frontend**: SvelteKit development server (port 5173)
- **Database**: Local Elasticsearch + Kibana (ports 9200/5601)
- **Storage**: Local filesystem (API/data/documents/)

### Key Features
- ✅ **23,500+ legal documents** with full-text search
- ✅ **Unified storage handler** - seamlessly switch between local/cloud
- ✅ **Automated deployments** via GitHub Actions
- ✅ **Global CDN** for fast document delivery
- ✅ **Auto-scaling** handles traffic spikes automatically
- ✅ **Full-text search** in legal documents
- ✅ **Document classification** and metadata extraction
- ✅ **Case management** and user authentication

## 📋 Prerequisites

- Docker and Docker Compose
- Node.js 18+ and npm
- Python 3.9+
- 10GB+ free disk space
- Supabase account (for authentication only)

## 🔧 Configuration

Create `.env` files in both API and Web directories:

**API/.env:**
```env
SUPABASE_URL=https://your-project.supabase.co
SUPABASE_ANON_KEY=your-anon-key
USE_LOCAL_STORAGE=true
ES_HOST=localhost
ES_PORT=9200
```

**Web environment:**
```bash
export PUBLIC_API_URL=http://localhost:8000
export PUBLIC_SUPABASE_URL=https://your-project.supabase.co
```

## 📚 Documentation

- **[LOCAL-HOSTING.md](LOCAL-HOSTING.md)** - Complete local setup guide
- **[QUICK-START.md](QUICK-START.md)** - Legacy S3 + remote ES setup
- **[API/README.simple.md](API/README.simple.md)** - API documentation
- **[Database/es/README.md](Database/es/README.md)** - Elasticsearch migration

## 🔄 Migration Status

**Elasticsearch:** ✅ Migrated (20,362/23,862 documents)
**Document Storage:** ✅ Local filesystem (API/data/documents/)
**Database:** ✅ Local Elasticsearch + Kibana
**Cloud Dependencies:** Only Supabase for authentication

## 💾 Data Structure

```
API/data/documents/
├── memorandum/2025/05/01/document.pdf
├── brief/2025/05/01/document.pdf
└── order/2025/05/01/document.pdf
```

## 🛠️ Development

### Local Development (No Docker)
```bash
# 1. Start Elasticsearch
cd Database/es && docker-compose -f docker-compose.standalone.yml up -d

# 2. API Development
cd ../../API
pip install -r requirements.simple.txt
python server.py

# 3. Web Development
cd ../Web
npm run dev
```

### Docker Development
```bash
# Start all services with Docker
docker network create motion-index-network
cd Database/es && docker-compose -f docker-compose.standalone.yml up -d
cd ../../API && docker-compose -f docker-compose.api.yml up -d
cd ../Web && npm run dev  # Web still runs locally for development
```

### Database Management
```bash
cd Database/es
./check.sh  # Check Elasticsearch status
docker-compose -f docker-compose.standalone.yml logs  # View logs
```

## 🌐 Access Points

- **Web Interface**: http://localhost:5173
- **API Documentation**: http://localhost:8000/docs
- **Elasticsearch**: http://localhost:9200
- **Kibana**: http://localhost:5601

## 📊 Performance

- **Search Response**: <100ms for most queries
- **Document Serving**: Direct filesystem access
- **Migration Speed**: 340 documents/second
- **Storage Efficiency**: 653MB for 20K documents

## 🔒 Security Notes

- Documents served through authenticated API endpoints
- Path traversal protection implemented
- Supabase handles user authentication
- Local storage isolated from external access

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch
3. Make changes with local testing
4. Submit a pull request

## 📄 License

See [LICENSE](LICENSE) file for details.