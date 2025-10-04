# Motion Index Web - Deployment Guide

## Docker Setup for DigitalOcean

This guide explains how to deploy the SvelteKit frontend to DigitalOcean App Platform using Docker.

## Files Created

- `Dockerfile.prod` - Production Docker image (already existed)
- `.dockerignore` - Docker ignore patterns
- `docker-compose.yml` - Local testing with Docker
- `app.yaml` - DigitalOcean App Platform specification
- Updated `package.json` with Docker scripts

## Local Docker Testing

### Setup environment variables first:
```bash
# Copy template and edit with your values
cp template.env .env
# Edit .env with your actual Supabase and API URLs
```

### Build and run production image:
```bash
# Make sure you have set environment variables in .env first
source .env
yarn docker:build
yarn docker:run
```

### Using docker-compose:
```bash
# Production mode
yarn docker:prod

# Development mode  
yarn docker:dev
```

### Manual Docker commands:
```bash
# Build production image
docker build -f Dockerfile.prod -t motion-index-web .

# Run container
docker run -p 3000:3000 \
  -e PUBLIC_SUPABASE_URL=your-supabase-url \
  -e PUBLIC_SUPABASE_ANON_KEY=your-anon-key \
  -e PUBLIC_API_URL=your-api-url \
  motion-index-web
```

## DigitalOcean App Platform Deployment

### Option 1: Using doctl CLI

1. Install doctl and authenticate:
```bash
# Install doctl
snap install doctl

# Authenticate
doctl auth init
```

2. Deploy using app spec:
```bash
# Create new app
doctl apps create --spec app.yaml

# Update existing app
doctl apps update <app-id> --spec app.yaml
```

### Option 2: Using DigitalOcean Dashboard

1. Go to [DigitalOcean Apps](https://cloud.digitalocean.com/apps)
2. Click "Create App"
3. Connect your GitHub repository
4. Select this repository and the `main` branch
5. DigitalOcean will auto-detect it as a Node.js app
6. Configure environment variables (see below)
7. Deploy

### Option 3: Using Docker Registry

1. Build and push to registry:
```bash
# Build for production
docker build -f Dockerfile.prod -t motion-index-web .

# Tag for registry
docker tag motion-index-web registry.digitalocean.com/your-registry/motion-index-web

# Push to DigitalOcean Container Registry
docker push registry.digitalocean.com/your-registry/motion-index-web
```

2. Deploy from Container Registry in DigitalOcean Apps

## Environment Variables

Set these in DigitalOcean App Platform:

### Required Variables:
- `PUBLIC_SUPABASE_URL` - Your Supabase project URL
- `PUBLIC_SUPABASE_ANON_KEY` - Supabase anonymous/public key
- `SUPABASE_SERVICE_ROLE_KEY` - Supabase service role key (server-side)
- `PUBLIC_API_URL` - URL to your Go Fiber API backend

### Optional Variables:
- `NODE_ENV=production` (auto-set)
- `PORT=3000` (auto-set)
- `HOST=0.0.0.0` (auto-set)

## App Platform Configuration

The `app.yaml` file includes:

- **Build**: `yarn install --frozen-lockfile && yarn build`
- **Run**: `node build`
- **Health Check**: HTTP check on `/`
- **Instance**: `basic-xxs` (1 vCPU, 512MB RAM)
- **Auto-scaling**: CPU threshold at 80%
- **Port**: 3000

## Customization

### Update GitHub Repository
Edit `app.yaml` line 8-9:
```yaml
github:
  repo: your-username/motion-index-web
  branch: main
```

### Custom Domain
Uncomment and configure domains section in `app.yaml`:
```yaml
domains:
  - domain: yourdomain.com
    type: PRIMARY
  - domain: www.yourdomain.com
    type: ALIAS
```

### Resource Scaling
Adjust instance size in `app.yaml`:
```yaml
instance_size_slug: basic-s  # 1 vCPU, 1GB RAM
# or
instance_size_slug: basic-m  # 1 vCPU, 2GB RAM
```

## Troubleshooting

### Build Issues
- Ensure all dependencies are in `package.json`
- Check for missing environment variables
- Verify Dockerfile.prod builds locally

### Runtime Issues
- Check app logs in DigitalOcean dashboard
- Verify environment variables are set
- Ensure API backend is accessible
- Check health check endpoint

### CORS Issues
- Ensure Go Fiber API has proper CORS headers
- Verify `PUBLIC_API_URL` points to correct API
- Check that API allows origin from your app domain

## Monitoring

DigitalOcean Apps provides:
- Real-time logs
- Performance metrics
- Health check status
- Auto-scaling events
- Deployment history

Access via: Apps → Your App → Runtime Logs / Insights