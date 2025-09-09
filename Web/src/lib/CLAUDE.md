# lib/ Directory

## Purpose
Shared libraries, components, and utilities used across the application.

## Files
- `auth.ts` - Custom authentication wrapper and state management
- `supabase.ts` - Case management database operations
- `styles.css` - Custom CSS components and variables (if needed)

## Key Libraries

### auth.ts
- Custom auth state management with Svelte stores
- Cookie-based session persistence
- Wrapper around Supabase Auth for consistency
- Authentication helpers and utilities

### supabase.ts
- CaseManager class for database operations
- CRUD operations for cases and case documents
- Supabase client wrapper methods
- TypeScript interfaces for data models

## Usage
```typescript
import { isAuthenticated, currentUser } from '$lib/auth';
import { CaseManager } from '$lib/supabase';
```

## Design Patterns
- Single responsibility principle
- TypeScript interfaces for type safety
- Async/await for database operations
- Error handling and logging