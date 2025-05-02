// api.ts
import type { SearchParams, SearchResponse, MetadataField, DocumentStats } from "../../utils/search_types"

export async function getDocumentTypes(): Promise<Record<string, number>> {
    const response = await fetch('/api/document-types');
    return response.json();
}

export async function getMetadataFields(): Promise<{ fields: MetadataField[] }> {
    const response = await fetch('/api/metadata-fields');
    return response.json();
}

export async function getDocumentStats(): Promise<DocumentStats> {
    const response = await fetch('/api/document-stats');
    return response.json();
}

export async function getAllFieldOptions(): Promise<Record<string, string[]>> {
    const response = await fetch('/api/field-options');
    return response.json();
}

export async function searchDocuments(params: SearchParams): Promise<SearchResponse> {
    const response = await fetch('/api/search', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json'
        },
        body: JSON.stringify(params)
    });
    return response.json();
}