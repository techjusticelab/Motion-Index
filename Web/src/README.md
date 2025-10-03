# Web Frontend Source

SvelteKit application providing the user interface for legal document search and management.

## Structure

### Core Application
- **app.html** - Main HTML template
- **app.css** - Global styles with Tailwind CSS v4
- **hooks.server.ts** - Server-side request hooks

### Libraries
- **lib/** - Shared utilities and configurations
  - `supabase.ts` - Authentication client configuration
  - `database.types.ts` - TypeScript type definitions

### Routes (File-based Routing)
- **routes/** - SvelteKit file-based routing structure
  - `/` - Main search interface
  - `/auth/*` - Authentication flows (login, register, forgot password)
  - `/upload` - Document upload interface
  - `/account` - User account and case management
  - `/cases/[caseId]` - Individual case details
  - `/help` - User documentation

### Components
- **routes/lib/components/** - Reusable UI components
  - `Document/` - Document viewing and display components
  - `Search/` - Search form, filters, and results components

### API Integration
- **api.ts** - API client for backend communication
- **routes/api/** - SvelteKit API routes for proxying requests

## Design Principles

- **Component-based**: Modular, reusable Svelte components
- **Type Safety**: Full TypeScript integration
- **Responsive Design**: Tailwind CSS for mobile-first design
- **Progressive Enhancement**: Works without JavaScript for core features

## Development

The frontend communicates with the FastAPI backend and Supabase for authentication, providing a complete document management interface.