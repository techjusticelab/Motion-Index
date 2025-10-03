// Legacy client.ts - Re-exports from the modular API structure
// This file exists for backward compatibility

// Re-export all functionality from the new modular API structure
export * from './search';
export * from './documents';
export * from './redaction';
export * from './types';
export { API_URL } from './config';