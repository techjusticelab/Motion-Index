// types.ts
export interface SearchParams {
    query: string;
    doc_type: string;
    case_number: string;
    case_name: string;
    judge: string[];
    court: string[];
    author: string;
    status: string;
    date_range: {
        start: string;
        end: string;
    };
    legal_tags: string[];
    legal_tags_match_all: boolean; // Whether to match all tags (AND) or any tag (OR)
    size: number;
    sort_by: string;
    sort_order: 'asc' | 'desc';
    page: number;
    use_fuzzy: boolean;
}

export interface Document {
    metadata: {
        document_name: string;
        case_number?: string;
        case_name?: string;
        judge?: string;
        court?: string;
        timestamp?: string;
        subject?: string;
        status?: string;
        author?: string;
        legal_tags?: string[];
    };
    file_name: string;
    s3_uri: string;
    doc_type: string;
    text: string;
    created_at: string;
    highlight?: {
        text: string[];
    };
}

export interface SearchResponse {
    total: number;
    hits: Document[];
}

export interface MetadataField {
    name: string;
    type: string;
    description?: string;
}

export interface DocumentStats {
    total_documents: number;
    date_range?: {
        oldest: string;
        newest: string;
    };
}