# lib/ Directory

## Purpose
Shared libraries, components, and utilities used across the application.

## Files
- `supabase-client.ts` - Supabase browser client configuration
- `supabase.ts` - Case management database operations
- `api/` - Go Fiber API client modules
- `components/` - Reusable UI components
- `styles.css` - Custom CSS components and variables (if needed)

## Key Libraries

### supabase-client.ts
- Browser-side Supabase client for authentication
- Used in client-side layout load functions
- Configured with environment variables

### supabase.ts
- CaseManager class for database operations
- CRUD operations for cases and case documents
- Supabase client wrapper methods
- TypeScript interfaces for data models

### api/
- Modular API client for Go Fiber backend
- `config.ts` - API URLs and configuration (port 6000)
- `search.ts` - Document search operations
- `documents.ts` - Document upload and management
- `redaction.ts` - PDF redaction analysis
- `types.ts` - TypeScript interfaces
- `index.ts` - Main exports

## Usage
```typescript
import { createSupabaseLoadClient } from '$lib/supabase-client';
import { CaseManager } from '$lib/supabase';
import { searchDocuments, categoriseDocument } from '$lib/api';
```

## Authentication
Uses Supabase SSR (Server-Side Rendering) authentication:
- Server hooks handle auth state management
- Session data passed through layout
- No custom authentication wrapper needed

## CORS Configuration Requirements

For the Go Fiber API backend to work with browser requests:

### Required CORS Headers
The Go Fiber API must include these headers:
```
Access-Control-Allow-Origin: http://localhost:5173
Access-Control-Allow-Methods: GET, POST, PUT, DELETE, OPTIONS
Access-Control-Allow-Headers: Content-Type, Authorization, ngrok-skip-browser-warning
Access-Control-Allow-Credentials: true
```

### Common Issues
1. **Browser Shows "Blocked"**: Usually indicates CORS misconfiguration
2. **Network Error**: API server not running or wrong port
3. **Status 0**: Browser blocked the request before it reached the server

### Go Fiber CORS Middleware
The API should use Go Fiber's CORS middleware:
```go
app.Use(cors.New(cors.Config{
    AllowOrigins: "http://localhost:5173",
    AllowMethods: "GET,POST,PUT,DELETE,OPTIONS",
    AllowHeaders: "Content-Type,Authorization,ngrok-skip-browser-warning",
    AllowCredentials: true,
}))
```