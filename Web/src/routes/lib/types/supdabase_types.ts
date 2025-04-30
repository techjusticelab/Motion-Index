export type Json =
    | string
    | number
    | boolean
    | null
    | { [key: string]: Json | undefined }
    | Json[]

export interface Database {
    public: {
        Tables: {
            profiles: {
                Row: {
                    id: string
                    email: string
                    cpda_id: string | null
                    updated_at: string | null
                    created_at: string
                }
                Insert: {
                    id: string
                    email: string
                    cpda_id?: string | null
                    updated_at?: string | null
                    created_at?: string
                }
                Update: {
                    id?: string
                    email?: string
                    cpda_id?: string | null
                    updated_at?: string | null
                    created_at?: string
                }
            }
            cases: {
                Row: {
                    id: string
                    user_id: string
                    case_name: string
                    created_at: string
                }
                Insert: {
                    id?: string
                    user_id: string
                    case_name: string
                    created_at?: string
                }
                Update: {
                    id?: string
                    user_id?: string
                    case_name?: string
                    created_at?: string
                }
            }
            documents: {
                Row: {
                    id: string
                    case_id: string
                    name: string
                    file_path: string
                    document_type: string | null
                    created_at: string
                }
                Insert: {
                    id?: string
                    case_id: string
                    name: string
                    file_path: string
                    document_type?: string | null
                    created_at?: string
                }
                Update: {
                    id?: string
                    case_id?: string
                    name?: string
                    file_path?: string
                    document_type?: string | null
                    created_at?: string
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