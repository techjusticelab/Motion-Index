# Motion-Index Web Frontend

Modern SvelteKit application for legal document search, upload, and case management.

## Quick Start

```bash
# Setup
npm install
cp .env.example .env.local  # Configure environment variables

# Development
npm run dev        # Start development server (port 5173)
npm run test       # Run tests
npm run typecheck  # Type checking
npm run lint       # Linting and formatting check
```

## Features

- **Document Search**: Full-text search with advanced filtering (court, judge, legal area)
- **User Authentication**: Secure login/registration via Supabase
- **Document Upload**: Multi-format file upload with AI processing
- **Case Management**: Organize documents by legal cases
- **Document Viewing**: In-browser PDF and document preview with comprehensive metadata
- **Responsive Design**: Mobile-first interface using Tailwind CSS

## Technology Stack

- **Framework**: SvelteKit 2.x with TypeScript
- **Styling**: Tailwind CSS v4 with custom design system
- **Authentication**: Supabase Auth with SSR
- **API**: Go Fiber backend integration
- **Testing**: Vitest + Testing Library
- **Type Safety**: Strict TypeScript configuration

## Project Structure

```
src/
├── routes/              # File-based routing
│   ├── +layout.svelte   # Root layout
│   ├── +page.svelte     # Home/search page
│   ├── auth/           # Authentication pages
│   ├── account/        # User profile and case management
│   ├── upload/         # Document upload
│   └── api/            # Server-side API routes
├── lib/
│   ├── components/     # Reusable UI components
│   │   ├── ui/         # Base UI components (Button, Modal, etc.)
│   │   ├── document/   # Document-specific components
│   │   └── forms/      # Form components
│   ├── api/           # API client modules
│   ├── hooks/         # Svelte composition hooks
│   ├── types/         # TypeScript type definitions
│   └── utils/         # Utility functions
├── app.html           # HTML template
└── app.css           # Global styles
tests/                 # Test files mirroring src structure
```

## Scripts

```bash
npm run dev          # Development server
npm run build        # Production build
npm run preview      # Preview production build
npm run test         # Run tests once
npm run test:watch   # Run tests in watch mode
npm run typecheck    # TypeScript type checking
npm run check        # Svelte component checking
npm run lint         # ESLint + Prettier checking
npm run format       # Format code with Prettier
```

## Environment Configuration

Copy `.env.example` to `.env.local` and configure:

```env
# Required
PUBLIC_SUPABASE_URL=https://your-project.supabase.co
PUBLIC_SUPABASE_ANON_KEY=your-anon-key
PUBLIC_API_URL=http://localhost:8003

# Optional
SUPABASE_SERVICE_ROLE_KEY=your-service-key
```

## Code Quality

- **ESLint**: Strict TypeScript rules with no-unused-vars, no-explicit-any
- **Prettier**: Consistent code formatting with tabs, single quotes
- **TypeScript**: Strict mode with explicit return types for exports
- **Testing**: Component and utility testing with Vitest
- **Environment Validation**: Startup validation of required variables

## Development Conventions

- **Imports**: Direct component imports preferred over barrel exports
- **Types**: Centralized in `src/lib/types/` with explicit exports
- **API**: Fetch-based client with consistent error handling
- **Components**: Colocated styles, TypeScript props, responsive design
- **Routes**: File-based with proper +page.server.ts for server logic

## Architecture Principles

- **UNIX Philosophy**: Small, focused, composable modules
- **KISS**: Simple solutions over abstractions
- **Type Safety**: Explicit types and runtime validation
- **Performance**: Minimal runtime, optimized builds
- **Accessibility**: Semantic HTML, proper ARIA attributes

Access the application at `http://localhost:5173` during development.
