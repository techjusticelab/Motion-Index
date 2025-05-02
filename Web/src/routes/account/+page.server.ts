import type { PageServerLoad } from './$types';
import { error } from '@sveltejs/kit';
import type { Actions } from './$types';

export const load: PageServerLoad = async ({ locals }) => {
    // Get the current authenticated user's session using the safe method
    const { session, user } = locals.getSession ? await locals.safeGetSession() : { session: null, user: null };

    if (!session || !user) {
        // Return empty arrays if no user is authenticated
        return {
            cases: [],
            caseDocuments: []
        };
    }

    // Fetch cases that belong to the current user
    const { data: cases, error: casesError } = await locals.supabase
        .from('cases')
        .select('id, case_name, case_docs, created_at, updated_at')
        .eq('user_id', user.id)
        .order('created_at', { ascending: false });

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
        const { data: documents, error: documentsError } = await locals.supabase
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
        caseDocuments: caseDocuments,
        // Add session to the PageData for client-side use
        session: session
    };
};

export const actions: Actions = {
    createCase: async ({ request, locals }) => {
        // Get authenticated user
        const { session, user } = await locals.safeGetSession();

        if (!session || !user) {
            throw error(401, "Unauthorized");
        }

        const formData = await request.formData();
        const caseName = formData.get('case_name')?.toString();

        if (!caseName || caseName.trim() === '') {
            return { success: false, message: 'Case name is required' };
        }

        try {
            const { data, error: createError } = await locals.supabase
                .from('cases')
                .insert({
                    user_id: user.id,
                    case_name: caseName.trim(),
                    case_docs: [],
                    created_at: new Date().toISOString(),
                    updated_at: new Date().toISOString()
                })
                .select();

            if (createError) throw createError;

            return {
                success: true,
                message: 'Case created successfully',
                case: data[0]
            };
        } catch (err) {
            console.error('Error creating case:', err);
            return { success: false, message: 'Failed to create case' };
        }
    },

    deleteCase: async ({ request, locals }) => {
        // Get authenticated user
        const { session, user } = await locals.safeGetSession();

        if (!session || !user) {
            throw error(401, "Unauthorized");
        }

        const formData = await request.formData();
        const caseId = formData.get('id')?.toString();

        if (!caseId) {
            return { success: false, message: 'Case ID is required' };
        }

        try {
            // Verify the case belongs to the user first
            const { data: caseData, error: caseError } = await locals.supabase
                .from('cases')
                .select('id')
                .eq('id', caseId)
                .eq('user_id', user.id)
                .single();

            if (caseError || !caseData) {
                throw error(403, "You don't have permission to delete this case");
            }

            // Delete associated documents first
            const { error: docDeleteError } = await locals.supabase
                .from('case_documents')
                .delete()
                .eq('case_id', caseId);

            if (docDeleteError) throw docDeleteError;

            // Then delete the case
            const { error: caseDeleteError } = await locals.supabase
                .from('cases')
                .delete()
                .eq('id', caseId);

            if (caseDeleteError) throw caseDeleteError;

            return {
                success: true,
                message: 'Case deleted successfully'
            };
        } catch (err) {
            console.error('Error deleting case:', err);
            if (err instanceof Response) throw err;
            return { success: false, message: 'Failed to delete case' };
        }
    }
};