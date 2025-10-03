package models

import (
	"time"
)

// Document represents a legal document in the search index
type Document struct {
	ID          string            `json:"id"`
	FileName    string            `json:"file_name"`
	FilePath    string            `json:"file_path"`
	FileURL     string            `json:"file_url,omitempty"`
	S3URI       string            `json:"s3_uri,omitempty"`
	Text        string            `json:"text"`
	DocType     string            `json:"doc_type"`
	Category    string            `json:"category,omitempty"`
	Hash        string            `json:"hash"`
	CreatedAt   time.Time         `json:"created_at"`
	UpdatedAt   time.Time         `json:"updated_at"`
	Metadata    *DocumentMetadata `json:"metadata"`
	Size        int64             `json:"size,omitempty"`
	ContentType string            `json:"content_type,omitempty"`
	
	// For backward compatibility with tests
	Title   string `json:"title,omitempty"`
	Content string `json:"content,omitempty"`
}

// DocumentType represents specific legal document categories
type DocumentType string

const (
	// Motions
	DocTypeMotionToSuppress         DocumentType = "motion_to_suppress"
	DocTypeMotionToDismiss          DocumentType = "motion_to_dismiss"
	DocTypeMotionToCompel           DocumentType = "motion_to_compel"
	DocTypeMotionInLimine           DocumentType = "motion_in_limine"
	DocTypeMotionForSummaryJudgment DocumentType = "motion_summary_judgment"
	DocTypeMotionToStrike           DocumentType = "motion_to_strike"
	DocTypeMotionForReconsideration DocumentType = "motion_for_reconsideration"
	DocTypeMotionToAmend            DocumentType = "motion_to_amend"
	DocTypeMotionForContinuance     DocumentType = "motion_for_continuance"
	
	// Orders and Rulings
	DocTypeOrder                    DocumentType = "order"
	DocTypeRuling                   DocumentType = "ruling"
	DocTypeJudgment                 DocumentType = "judgment"
	DocTypeSentence                 DocumentType = "sentence"
	DocTypeInjunction               DocumentType = "injunction"
	
	// Briefs and Pleadings
	DocTypeBrief                    DocumentType = "brief"
	DocTypeComplaint                DocumentType = "complaint"
	DocTypeAnswer                   DocumentType = "answer"
	DocTypePlea                     DocumentType = "plea"
	DocTypeReply                    DocumentType = "reply"
	
	// Administrative
	DocTypeDocketEntry              DocumentType = "docket_entry"
	DocTypeNotice                   DocumentType = "notice"
	DocTypeStipulation              DocumentType = "stipulation"
	DocTypeCorrespondence           DocumentType = "correspondence"
	DocTypeTranscript               DocumentType = "transcript"
	DocTypeEvidence                 DocumentType = "evidence"
	DocTypeOther                    DocumentType = "other"
	DocTypeUnknown                  DocumentType = "unknown"
)

// DocumentMetadata contains comprehensive legal-specific metadata
type DocumentMetadata struct {
	// Basic Information
	DocumentName string       `json:"document_name"`
	Subject      string       `json:"subject"`
	Summary      string       `json:"summary,omitempty"`          // Enhanced legal summary
	DocumentType DocumentType `json:"document_type"`
	
	// Case Information
	Case         *CaseInfo    `json:"case,omitempty"`
	
	// Court Information
	Court        *CourtInfo   `json:"court,omitempty"`
	
	// People & Parties
	Parties      []Party      `json:"parties,omitempty"`
	Attorneys    []Attorney   `json:"attorneys,omitempty"`
	Judge        *Judge       `json:"judge,omitempty"`
	
	// Dates & Status
	FilingDate   *time.Time   `json:"filing_date,omitempty"`
	EventDate    *time.Time   `json:"event_date,omitempty"`
	Timestamp    *time.Time   `json:"timestamp,omitempty"`        // For backward compatibility
	Status       string       `json:"status,omitempty"`
	
	// Document Properties
	Language     string       `json:"language,omitempty"`
	Pages        int          `json:"pages,omitempty"`
	WordCount    int          `json:"word_count,omitempty"`
	
	// Legal Classification
	LegalTags    []string     `json:"legal_tags,omitempty"`
	Charges      []Charge     `json:"charges,omitempty"`
	Authorities  []Authority  `json:"authorities,omitempty"`
	
	// Processing Metadata
	ProcessedAt  time.Time    `json:"processed_at"`
	Confidence   float64      `json:"confidence,omitempty"`
	AIClassified bool         `json:"ai_classified"`
	
	// Legacy fields for backward compatibility
	CaseName     string       `json:"case_name,omitempty"`
	CaseNumber   string       `json:"case_number,omitempty"`
	Author       string       `json:"author,omitempty"`
}

// CaseInfo contains detailed case information
type CaseInfo struct {
	CaseNumber   string    `json:"case_number"`
	CaseName     string    `json:"case_name"`
	CaseType     string    `json:"case_type,omitempty"`     // "criminal", "civil", "traffic"
	Chapter      string    `json:"chapter,omitempty"`       // Bankruptcy chapter
	Docket       string    `json:"docket,omitempty"`        // Full docket number
	NatureOfSuit string    `json:"nature_of_suit,omitempty"`
}

// CourtInfo contains detailed court information
type CourtInfo struct {
	CourtID      string    `json:"court_id"`
	CourtName    string    `json:"court_name"`
	Jurisdiction string    `json:"jurisdiction"`            // "federal", "state", "local"
	Level        string    `json:"level"`                   // "trial", "appellate", "supreme"
	District     string    `json:"district,omitempty"`
	Division     string    `json:"division,omitempty"`
	County       string    `json:"county,omitempty"`
}

// Party represents a party to the case
type Party struct {
	Name         string    `json:"name"`
	Role         string    `json:"role"`                    // "defendant", "plaintiff", "appellant"
	PartyType    string    `json:"party_type,omitempty"`    // "individual", "corporation", "government"
	Date         *time.Time `json:"date,omitempty"`         // Date associated with this party
}

// Attorney represents legal counsel
type Attorney struct {
	Name         string    `json:"name"`
	BarNumber    string    `json:"bar_number,omitempty"`
	Role         string    `json:"role"`                    // "defense", "prosecution", "counsel"
	Organization string    `json:"organization,omitempty"`
	ContactInfo  string    `json:"contact_info,omitempty"`
}

// Judge represents the presiding judge
type Judge struct {
	Name         string    `json:"name"`
	Title        string    `json:"title,omitempty"`
	JudgeID      string    `json:"judge_id,omitempty"`
}

// Charge represents criminal charges
type Charge struct {
	Statute      string    `json:"statute"`
	Description  string    `json:"description"`
	Grade        string    `json:"grade,omitempty"`         // "felony", "misdemeanor"
	Class        string    `json:"class,omitempty"`         // "A", "B", "C"
	Count        int       `json:"count,omitempty"`         // Count number
}

// Authority represents legal authorities cited
type Authority struct {
	Citation     string    `json:"citation"`
	CaseTitle    string    `json:"case_title,omitempty"`
	Type         string    `json:"type"`                    // "case_law", "statute", "regulation"
	Precedent    bool      `json:"precedent"`
	Page         string    `json:"page,omitempty"`          // Page or section reference
}

// GetDocumentMapping returns the enhanced OpenSearch mapping for legal documents
func GetDocumentMapping() map[string]interface{} {
	return map[string]interface{}{
		"mappings": map[string]interface{}{
			"properties": map[string]interface{}{
				"id": map[string]interface{}{
					"type": "keyword",
				},
				"file_name": map[string]interface{}{
					"type": "text",
					"fields": map[string]interface{}{
						"keyword": map[string]interface{}{
							"type": "keyword",
						},
					},
				},
				"file_path": map[string]interface{}{
					"type": "keyword",
				},
				"file_url": map[string]interface{}{
					"type":  "keyword",
					"index": false,
				},
				"s3_uri": map[string]interface{}{
					"type":  "keyword",
					"index": false,
				},
				"text": map[string]interface{}{
					"type":     "text",
					"analyzer": "legal_analyzer",
				},
				"doc_type": map[string]interface{}{
					"type": "keyword",
				},
				"category": map[string]interface{}{
					"type": "keyword",
				},
				"hash": map[string]interface{}{
					"type": "keyword",
				},
				"created_at": map[string]interface{}{
					"type": "date",
				},
				"updated_at": map[string]interface{}{
					"type": "date",
				},
				"size": map[string]interface{}{
					"type": "long",
				},
				"content_type": map[string]interface{}{
					"type": "keyword",
				},
				"metadata": map[string]interface{}{
					"properties": map[string]interface{}{
						"document_name": map[string]interface{}{
							"type": "text",
							"fields": map[string]interface{}{
								"keyword": map[string]interface{}{
									"type": "keyword",
								},
							},
						},
						"subject": map[string]interface{}{
							"type": "text",
							"analyzer": "legal_analyzer",
						},
						"summary": map[string]interface{}{
							"type": "text",
							"analyzer": "legal_analyzer",
						},
						"document_type": map[string]interface{}{
							"type": "keyword",
						},
						"status": map[string]interface{}{
							"type": "keyword",
						},
						"filing_date": map[string]interface{}{
							"type": "date",
						},
						"event_date": map[string]interface{}{
							"type": "date",
						},
						"timestamp": map[string]interface{}{
							"type": "date",
						},
						"processed_at": map[string]interface{}{
							"type": "date",
						},
						"confidence": map[string]interface{}{
							"type": "float",
						},
						"ai_classified": map[string]interface{}{
							"type": "boolean",
						},
						"case": map[string]interface{}{
							"properties": map[string]interface{}{
								"case_number": map[string]interface{}{
									"type": "keyword",
								},
								"case_name": map[string]interface{}{
									"type": "text",
									"fields": map[string]interface{}{
										"keyword": map[string]interface{}{
											"type": "keyword",
										},
									},
								},
								"case_type": map[string]interface{}{
									"type": "keyword",
								},
								"chapter": map[string]interface{}{
									"type": "keyword",
								},
								"docket": map[string]interface{}{
									"type": "keyword",
								},
								"nature_of_suit": map[string]interface{}{
									"type": "keyword",
								},
							},
						},
						"court": map[string]interface{}{
							"properties": map[string]interface{}{
								"court_id": map[string]interface{}{
									"type": "keyword",
								},
								"court_name": map[string]interface{}{
									"type": "keyword",
								},
								"jurisdiction": map[string]interface{}{
									"type": "keyword",
								},
								"level": map[string]interface{}{
									"type": "keyword",
								},
								"district": map[string]interface{}{
									"type": "keyword",
								},
								"division": map[string]interface{}{
									"type": "keyword",
								},
								"county": map[string]interface{}{
									"type": "keyword",
								},
							},
						},
						"parties": map[string]interface{}{
							"type": "nested",
							"properties": map[string]interface{}{
								"name": map[string]interface{}{
									"type": "keyword",
								},
								"role": map[string]interface{}{
									"type": "keyword",
								},
								"party_type": map[string]interface{}{
									"type": "keyword",
								},
								"date": map[string]interface{}{
									"type": "date",
								},
							},
						},
						"attorneys": map[string]interface{}{
							"type": "nested",
							"properties": map[string]interface{}{
								"name": map[string]interface{}{
									"type": "keyword",
								},
								"bar_number": map[string]interface{}{
									"type": "keyword",
								},
								"role": map[string]interface{}{
									"type": "keyword",
								},
								"organization": map[string]interface{}{
									"type": "keyword",
								},
								"contact_info": map[string]interface{}{
									"type": "keyword",
									"index": false,
								},
							},
						},
						"judge": map[string]interface{}{
							"properties": map[string]interface{}{
								"name": map[string]interface{}{
									"type": "keyword",
								},
								"title": map[string]interface{}{
									"type": "keyword",
								},
								"judge_id": map[string]interface{}{
									"type": "keyword",
								},
							},
						},
						"charges": map[string]interface{}{
							"type": "nested",
							"properties": map[string]interface{}{
								"statute": map[string]interface{}{
									"type": "keyword",
								},
								"description": map[string]interface{}{
									"type": "text",
								},
								"grade": map[string]interface{}{
									"type": "keyword",
								},
								"class": map[string]interface{}{
									"type": "keyword",
								},
								"count": map[string]interface{}{
									"type": "integer",
								},
							},
						},
						"authorities": map[string]interface{}{
							"type": "nested",
							"properties": map[string]interface{}{
								"citation": map[string]interface{}{
									"type": "keyword",
								},
								"case_title": map[string]interface{}{
									"type": "text",
									"fields": map[string]interface{}{
										"keyword": map[string]interface{}{
											"type": "keyword",
										},
									},
								},
								"type": map[string]interface{}{
									"type": "keyword",
								},
								"precedent": map[string]interface{}{
									"type": "boolean",
								},
								"page": map[string]interface{}{
									"type": "keyword",
								},
							},
						},
						"legal_tags": map[string]interface{}{
							"type": "keyword",
						},
						"language": map[string]interface{}{
							"type": "keyword",
						},
						"pages": map[string]interface{}{
							"type": "integer",
						},
						"word_count": map[string]interface{}{
							"type": "integer",
						},
						// Legacy fields for backward compatibility
						"case_name": map[string]interface{}{
							"type": "text",
							"fields": map[string]interface{}{
								"keyword": map[string]interface{}{
									"type": "keyword",
								},
							},
						},
						"case_number": map[string]interface{}{
							"type": "keyword",
						},
						"author": map[string]interface{}{
							"type": "keyword",
						},
					},
				},
			},
		},
		"settings": map[string]interface{}{
			"number_of_shards":   1,
			"number_of_replicas": 0,
			"analysis": map[string]interface{}{
				"analyzer": map[string]interface{}{
					"legal_analyzer": map[string]interface{}{
						"type":      "custom",
						"tokenizer": "standard",
						"filter": []string{
							"lowercase",
							"stop",
							"stemmer",
							"legal_synonyms",
						},
					},
				},
				"filter": map[string]interface{}{
					"legal_synonyms": map[string]interface{}{
						"type": "synonym",
						"synonyms": []string{
							"motion,petition,application",
							"defendant,accused,respondent",
							"plaintiff,petitioner,complainant",
							"attorney,counsel,lawyer",
							"judge,court,magistrate",
							"order,ruling,decision",
							"suppress,exclude,prohibit",
							"dismiss,quash,deny",
						},
					},
				},
			},
		},
	}
}
