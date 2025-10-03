# Motion-Index Web Frontend

SvelteKit-based user interface for legal document search, upload, and case management.

## Overview

Modern, responsive web application providing intuitive access to the legal document search system. Built with SvelteKit for optimal performance and developer experience.

## Features

- **Document Search**: Full-text search with advanced filtering (court, judge, legal area)
- **User Authentication**: Secure login/registration via Supabase
- **Document Upload**: Multi-format file upload with processing status
- **Case Management**: Organize documents by legal cases
- **Document Viewing**: In-browser PDF and document preview
- **Responsive Design**: Mobile-first interface using Tailwind CSS

## Technology Stack

- **Framework**: SvelteKit (SSR + client-side hydration)
- **Styling**: Tailwind CSS v4
- **Authentication**: Supabase Auth
- **API Communication**: REST API integration with FastAPI backend
- **Type Safety**: Full TypeScript support

## Development

```bash
npm install
npm run dev       # Development server (port 5173)
npm run build     # Production build
npm run preview   # Preview production build
npm run check     # Type checking
npm run lint      # Linting
npm run format    # Code formatting
```

## Configuration

Environment variables:
- `PUBLIC_API_URL` - Backend API endpoint
- `PUBLIC_SUPABASE_URL` - Supabase project URL
- `PUBLIC_SUPABASE_ANON_KEY` - Supabase anonymous key

## Architecture

- **File-based routing** in `src/routes/`
- **Component library** in `src/routes/lib/components/`
- **Shared utilities** in `src/lib/`
- **Type definitions** for database schemas and API responses

Access the application at `http://localhost:5173` during development.
