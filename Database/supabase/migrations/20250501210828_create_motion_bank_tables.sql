-- Users table (likely already exists in your Supabase auth system)
-- Supabase handles this automatically, but here's the reference:
-- auth.users (managed by Supabase)

-- Cases table
CREATE TABLE cases (
  id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  user_id UUID NOT NULL REFERENCES auth.users(id) ON DELETE CASCADE,
  title VARCHAR(255) NOT NULL,
  description TEXT,
  created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
  updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
  status VARCHAR(50) DEFAULT 'active'
);

-- Documents table (metadata - actual content likely in Elasticsearch)
CREATE TABLE documents (
  id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  title VARCHAR(255) NOT NULL,
  description TEXT,
  file_path VARCHAR(255),
  elasticsearch_id VARCHAR(255), -- Reference to document in Elasticsearch
  document_type VARCHAR(100),
  uploaded_by UUID REFERENCES auth.users(id),
  upload_date TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
  is_public BOOLEAN DEFAULT true
);

-- Case-Document association (many-to-many relationship)
CREATE TABLE case_documents (
  id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  case_id UUID NOT NULL REFERENCES cases(id) ON DELETE CASCADE,
  document_id UUID NOT NULL REFERENCES documents(id) ON DELETE CASCADE,
  added_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
  notes TEXT,
  UNIQUE(case_id, document_id)
);

-- Create indexes for performance
CREATE INDEX idx_cases_user_id ON cases(user_id);
CREATE INDEX idx_case_documents_case_id ON case_documents(case_id);
CREATE INDEX idx_case_documents_document_id ON case_documents(document_id);