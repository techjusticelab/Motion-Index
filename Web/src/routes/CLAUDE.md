# routes/ Directory

## Purpose
SvelteKit file-based routing system. Each directory/file corresponds to a route in the application.

## Structure
- `+layout.svelte` - Root layout component (header, footer, nav)
- `+page.svelte` - Home page with document search interface
- `auth/` - Authentication pages (login, register, logout)
- `account/` - User account management and case management
- `upload/` - Document upload interface
- `api/` - Server-side API routes
- Various `+page.server.ts` files for server-side logic

## Key Files
- `+layout.svelte` - Main navigation, authentication state management
- `+page.svelte` - Search interface for legal documents
- `+layout.server.ts` - Server-side layout logic, session handling
- `hooks.server.ts` - SvelteKit hooks for authentication

## Route Patterns
- `/` - Document search (main page)
- `/auth/*` - Authentication flows
- `/account` - User profile and case management  
- `/upload` - Document upload
- `/api/*` - REST API endpoints