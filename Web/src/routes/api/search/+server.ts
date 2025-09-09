// API route to search documents
import { json } from '@sveltejs/kit';
import type { RequestHandler } from './$types';

export const GET: RequestHandler = async ({ url, fetch }) => {
	const query = url.searchParams.get('q');
	
	if (!query) {
		return json({ results: [] });
	}

	try {
		// Call your backend API for document search
		// This should be your existing search endpoint
		const response = await fetch(`http://localhost:8000/search?q=${encodeURIComponent(query)}`);
		
		if (!response.ok) {
			throw new Error(`API responded with status: ${response.status}`);
		}
		
		const data = await response.json();
		
		return json({
			results: data.results || data.documents || []
		});
	} catch (error) {
		console.error('Search API error:', error);
		return json({ 
			error: 'Search failed',
			results: []
		}, { status: 500 });
	}
};