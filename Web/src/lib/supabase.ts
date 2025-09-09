// Supabase database operations for case management
import type { SupabaseClient } from '@supabase/supabase-js';

export interface Case {
  id: string;
  user_id: string;
  case_name: string;
  case_docs?: string[]; // Array of case_document IDs
  created_at: string;
  updated_at: string;
}

export interface CaseDocument {
  id: string;
  case_id: string;
  document_ids: string;
  notes?: string;
  added_at: string;
}

export class CaseManager {
  constructor(private supabase: SupabaseClient) {}

  // Create a new case
  async createCase(userId: string, caseName: string): Promise<Case | null> {
    // The database has a circular foreign key constraint issue with case_docs field
    // Try to create the case without the case_docs field first
    const { data, error } = await this.supabase
      .from('cases')
      .insert({
        user_id: userId,
        case_name: caseName
      })
      .select()
      .single();

    if (error) {
      console.error('Error creating case:', error);
      return null;
    }
    return data;
  }

  // Get all cases for a user
  async getUserCases(userId: string): Promise<Case[]> {
    const { data, error } = await this.supabase
      .from('cases')
      .select('*')
      .eq('user_id', userId)
      .order('updated_at', { ascending: false });

    if (error) {
      console.error('Error fetching cases:', error);
      return [];
    }
    return data || [];
  }

  // Update case name
  async updateCaseName(caseId: string, caseName: string): Promise<boolean> {
    const { error } = await this.supabase
      .from('cases')
      .update({ 
        case_name: caseName,
        updated_at: new Date().toISOString()
      })
      .eq('id', caseId);

    if (error) {
      console.error('Error updating case:', error);
      return false;
    }
    return true;
  }

  // Delete a case
  async deleteCase(caseId: string): Promise<boolean> {
    // First delete all case documents
    await this.supabase
      .from('case_documents')
      .delete()
      .eq('case_id', caseId);

    // Then delete the case
    const { error } = await this.supabase
      .from('cases')
      .delete()
      .eq('id', caseId);

    if (error) {
      console.error('Error deleting case:', error);
      return false;
    }
    return true;
  }

  // Add document to case
  async addDocumentToCase(caseId: string, documentId: string, notes?: string): Promise<CaseDocument | null> {
    const { data, error } = await this.supabase
      .from('case_documents')
      .insert({
        case_id: caseId,
        document_ids: documentId,
        notes
      })
      .select()
      .single();

    if (error) {
      console.error('Error adding document to case:', error);
      return null;
    }

    // Also add the document ID to the case_docs array
    await this.addDocumentIdToCase(caseId, data.id);

    return data;
  }

  // Get documents for a case
  async getCaseDocuments(caseId: string): Promise<CaseDocument[]> {
    const { data, error } = await this.supabase
      .from('case_documents')
      .select('*')
      .eq('case_id', caseId)
      .order('added_at', { ascending: false });

    if (error) {
      console.error('Error fetching case documents:', error);
      return [];
    }
    return data || [];
  }

  // Remove document from case
  async removeDocumentFromCase(caseDocumentId: string): Promise<boolean> {
    // First get the document to find the case_id
    const { data: docData, error: fetchError } = await this.supabase
      .from('case_documents')
      .select('case_id')
      .eq('id', caseDocumentId)
      .single();

    if (fetchError) {
      console.error('Error fetching document:', fetchError);
      return false;
    }

    // Delete the document
    const { error } = await this.supabase
      .from('case_documents')
      .delete()
      .eq('id', caseDocumentId);

    if (error) {
      console.error('Error removing document from case:', error);
      return false;
    }

    // Also remove the document ID from the case_docs array
    await this.removeDocumentIdFromCase(docData.case_id, caseDocumentId);

    return true;
  }

  // Update document notes
  async updateDocumentNotes(caseDocumentId: string, notes: string): Promise<boolean> {
    const { error } = await this.supabase
      .from('case_documents')
      .update({ notes })
      .eq('id', caseDocumentId);

    if (error) {
      console.error('Error updating document notes:', error);
      return false;
    }
    return true;
  }

  // Add document ID to case_docs array
  async addDocumentIdToCase(caseId: string, documentId: string): Promise<boolean> {
    // First get the current case_docs array
    const { data: caseData, error: fetchError } = await this.supabase
      .from('cases')
      .select('case_docs')
      .eq('id', caseId)
      .single();

    if (fetchError) {
      console.error('Error fetching case:', fetchError);
      return false;
    }

    // Add the document ID to the array if it's not already there
    const currentDocs = caseData.case_docs || [];
    if (!currentDocs.includes(documentId)) {
      const updatedDocs = [...currentDocs, documentId];
      
      const { error } = await this.supabase
        .from('cases')
        .update({
          case_docs: updatedDocs,
          updated_at: new Date().toISOString()
        })
        .eq('id', caseId);

      if (error) {
        console.error('Error updating case_docs array:', error);
        return false;
      }
    }

    return true;
  }

  // Remove document ID from case_docs array
  async removeDocumentIdFromCase(caseId: string, documentId: string): Promise<boolean> {
    // First get the current case_docs array
    const { data: caseData, error: fetchError } = await this.supabase
      .from('cases')
      .select('case_docs')
      .eq('id', caseId)
      .single();

    if (fetchError) {
      console.error('Error fetching case:', fetchError);
      return false;
    }

    // Remove the document ID from the array
    const currentDocs = caseData.case_docs || [];
    const updatedDocs = currentDocs.filter((id: string) => id !== documentId);
    
    const { error } = await this.supabase
      .from('cases')
      .update({
        case_docs: updatedDocs,
        updated_at: new Date().toISOString()
      })
      .eq('id', caseId);

    if (error) {
      console.error('Error updating case_docs array:', error);
      return false;
    }

    return true;
  }
}