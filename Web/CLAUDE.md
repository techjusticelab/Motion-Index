# Web Frontend - Claude Documentation

## Overview
SvelteKit frontend for Motion-Index legal document processing system. Handles user authentication, document search, case management, and file uploads.

## Architecture
- **Framework**: SvelteKit with TypeScript
- **Styling**: Tailwind CSS v4
- **Authentication**: Supabase Auth
- **Database**: Supabase PostgreSQL
- **State Management**: Svelte stores + custom auth wrapper

## Key Technologies
- SvelteKit 2.x
- Supabase client libraries
- Tailwind CSS 4.x
- TypeScript
- Vite build system

## Development Commands
```bash
npm run dev     # Start development server
npm run build   # Build for production
npm run preview # Preview production build
npm run lint    # Run linter
npm run format  # Format code with Prettier
npm run check   # Type-check with svelte-check
```

## Environment Variables
- `PUBLIC_SUPABASE_URL` - Supabase project URL
- `PUBLIC_SUPABASE_ANON_KEY` - Supabase anonymous key
- `SUPABASE_SERVICE_ROLE_KEY` - Server-side Supabase key

## File Structure
- `/src/routes/` - SvelteKit pages and API routes
- `/src/lib/` - Shared components and utilities
- `/src/app.html` - HTML template
- `/src/app.css` - Global styles with Tailwind