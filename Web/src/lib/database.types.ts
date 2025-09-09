// Basic database types for Supabase
export interface Database {
  public: {
    Tables: {
      cases: {
        Row: {
          id: string
          user_id: string
          case_name: string
          case_docs?: string
          created_at: string
          updated_at: string
        }
        Insert: {
          id?: string
          user_id: string
          case_name: string
          case_docs?: string
          created_at?: string
          updated_at?: string
        }
        Update: {
          id?: string
          user_id?: string
          case_name?: string
          case_docs?: string
          created_at?: string
          updated_at?: string
        }
      }
      case_documents: {
        Row: {
          id: string
          case_id: string
          document_ids: string
          notes?: string
          added_at: string
        }
        Insert: {
          id?: string
          case_id: string
          document_ids: string
          notes?: string
          added_at?: string
        }
        Update: {
          id?: string
          case_id?: string
          document_ids?: string
          notes?: string
          added_at?: string
        }
      }
    }
    Views: {
      [_ in never]: never
    }
    Functions: {
      [_ in never]: never
    }
    Enums: {
      [_ in never]: never
    }
    CompositeTypes: {
      [_ in never]: never
    }
  }
}