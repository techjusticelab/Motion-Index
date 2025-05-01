import type { PageServerLoad } from './$types';

export const load: PageServerLoad = async ({ locals: { supabase, getSession } }) => {
    // Get the current authenticated user's session
    const session = await getSession();

    if (!session) {
        // Return empty arrays if no user is authenticated
        return {
            cases: [],
            caseDocuments: []
        };
    }

    // Fetch cases that belong to the current user
    const { data: cases, error: casesError } = await supabase
        .from('cases')
        .select('id, case_docs, created_at, updated_at')
        .eq('user_id', session.user.id);

    if (casesError) {
        console.error('Error fetching cases:', casesError);
        return { cases: [], caseDocuments: [] };
    }

    // Get all case IDs to use for the next query
    const caseIds = cases?.map(c => c.id) || [];

    // Only fetch case documents if we have cases
    let caseDocuments = [];
    if (caseIds.length > 0) {
        // Fetch documents associated with the user's cases
        const { data: documents, error: documentsError } = await supabase
            .from('case_documents')
            .select('id, case_id, document_ids, added_at, notes')
            .in('case_id', caseIds);

        if (documentsError) {
            console.error('Error fetching case documents:', documentsError);
        } else {
            caseDocuments = documents || [];
        }
    }

    return {
        cases: cases || [],
        caseDocuments: caseDocuments
    };
};