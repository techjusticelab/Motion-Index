# routes/auth/ Directory

## Purpose
Authentication pages and flows using Supabase Auth.

## Files
- `login/+page.svelte` - Login form with email/password
- `login/+page.server.ts` - Server-side login logic
- `register/+page.svelte` - User registration form
- `register/+page.server.ts` - Server-side registration logic
- `logout/+page.svelte` - Logout confirmation and cleanup
- `confirm/+page.server.ts` - Email confirmation handling
- `forgot/+page.svelte` - Password reset form
- `reset/+page.svelte` - Password reset completion

## Authentication Flow
1. User submits login/register form
2. Server-side validation and Supabase API calls
3. Session creation and cookie management
4. Redirect to requested page or home
5. Client-side auth state management

## Features
- Email/password authentication
- Password reset via email
- Email confirmation for new accounts
- Session persistence with secure cookies
- Custom auth wrapper for API consistency