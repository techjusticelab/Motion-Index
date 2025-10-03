package classifier

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"
)

// openaiClassifier implements classification using OpenAI's API
type openaiClassifier struct {
	apiKey     string
	model      string
	httpClient *http.Client
}

// NewOpenAIClassifier creates a new OpenAI-based classifier
func NewOpenAIClassifier(config *Config) (Classifier, error) {
	if config.APIKey == "" {
		return nil, fmt.Errorf("OpenAI API key is required")
	}

	model := config.Model
	if model == "" {
		model = "gpt-4" // Default model
	}

	timeout := config.Timeout
	if timeout == 0 {
		timeout = 30 * time.Second
	}

	return &openaiClassifier{
		apiKey: config.APIKey,
		model:  model,
		httpClient: &http.Client{
			Timeout: timeout,
		},
	}, nil
}

// Classify analyzes document text using OpenAI and returns classification results
func (c *openaiClassifier) Classify(ctx context.Context, text string, metadata *DocumentMetadata) (*ClassificationResult, error) {
	// Create the classification prompt
	prompt := c.buildClassificationPrompt(text, metadata)

	// Make request to OpenAI
	response, err := c.makeOpenAIRequest(ctx, prompt)
	if err != nil {
		return nil, NewClassificationError("openai_request", "failed to classify document", err)
	}

	// Parse the response
	result, err := c.parseClassificationResponse(response)
	if err != nil {
		return nil, NewClassificationError("response_parsing", "failed to parse classification response", err)
	}

	return result, nil
}

// GetSupportedCategories returns the categories this classifier can identify
func (c *openaiClassifier) GetSupportedCategories() []string {
	return GetDefaultCategories()
}

// IsConfigured returns true if the classifier is properly configured
func (c *openaiClassifier) IsConfigured() bool {
	return c.apiKey != "" && c.model != ""
}

// buildClassificationPrompt creates an enhanced prompt for OpenAI classification
func (c *openaiClassifier) buildClassificationPrompt(text string, metadata *DocumentMetadata) string {
	// Dynamically adjust text truncation based on document characteristics
	maxTextLength := c.calculateOptimalTextLength(metadata)
	if len(text) > maxTextLength {
		text = text[:maxTextLength] + "..."
	}

	// Generate document-specific analysis context
	contextualPrompt := c.generateContextualPrompt(metadata)

	prompt := fmt.Sprintf(`You are an expert legal document analyzer specializing in California criminal law and civil litigation. 

Analyze the following legal document and provide comprehensive classification and extraction.

Document metadata:
- File name: %s
- File type: %s
- Word count: %d words
- Page count: %d pages
- Source system: %s

%s

CRITICAL INSTRUCTIONS:
1. Classify document type from: %s
2. Provide SUBSTANTIVE legal summary based on document type
3. Extract ALL legal entities with high precision
4. Identify case information, parties, and procedural context

DOCUMENT-SPECIFIC SUMMARIZATION REQUIREMENTS:

FOR MOTIONS (motion_to_suppress, motion_to_dismiss, etc.):
- Motion type and specific relief sought (3-4 sentences)
- Key legal arguments and constitutional/statutory authorities cited (2-3 sentences)
- Factual basis and procedural posture (2 sentences)
- Potential impact on case progression (1 sentence)

FOR ORDERS/RULINGS (order, ruling, judgment):
- Court's holding and primary reasoning (3-4 sentences)
- Key legal precedents and statutes applied (2 sentences)
- Impact on pending motions and case status (2 sentences)
- Practical implications for parties (1-2 sentences)

FOR BRIEFS (brief, reply):
- Main legal arguments and theory of the case (4-5 sentences)
- Factual background and procedural history (2-3 sentences)
- Authorities relied upon and distinguishing cases (2-3 sentences)
- Relief requested and strategic positioning (1-2 sentences)

FOR PLEADINGS (complaint, answer, plea):
- Claims/charges and factual allegations (3-4 sentences)
- Legal theories and causes of action (2-3 sentences)
- Defenses raised and procedural responses (2 sentences)
- Stakes and potential outcomes (1 sentence)

Document text:
%s

Respond with ONLY a JSON object in this exact format:
{
  "document_type": "<one of the available document types>",
  "legal_category": "<primary legal area>",
  "subject": "<concise 8-12 word subject line>",
  "summary": "<comprehensive legal summary following document-specific requirements above>",
  "confidence": <float between 0 and 1>,
  "keywords": ["<key legal terms and procedural elements>"],
  "legal_tags": ["<relevant legal doctrine tags>"],
  "case_info": {
    "case_number": "<case number if found>",
    "case_name": "<case title if found>",
    "case_type": "<criminal|civil|traffic|family>",
    "docket": "<full docket number>"
  },
  "court_info": {
    "court_name": "<court name>",
    "jurisdiction": "<federal|state|local>",
    "level": "<trial|appellate|supreme>",
    "county": "<county if applicable>"
  },
  "parties": [
    {
      "name": "<party name>",
      "role": "<defendant|plaintiff|appellant|respondent>",
      "party_type": "<individual|corporation|government>"
    }
  ],
  "attorneys": [
    {
      "name": "<attorney name>",
      "role": "<defense|prosecution|counsel>",
      "organization": "<law firm or agency>"
    }
  ],
  "judge": {
    "name": "<judge name>",
    "title": "<title if specified>"
  },
  "charges": [
    {
      "statute": "<statute number>",
      "description": "<charge description>",
      "grade": "<felony|misdemeanor>",
      "class": "<A|B|C>"
    }
  ],
  "authorities": [
    {
      "citation": "<legal citation>",
      "case_title": "<case name>",
      "type": "<case_law|statute|regulation>",
      "precedent": <true|false>
    }
  ],
  "filing_date": "<YYYY-MM-DD format or null>",
  "event_date": "<YYYY-MM-DD format or null>",
  "status": "<filed|granted|denied|pending|served>",
  "entities": [
    {
      "text": "<entity text>",
      "type": "<PERSON|ORGANIZATION|LOCATION|DATE|MONEY|LEGAL_CITATION|CASE_NUMBER|STATUTE>",
      "confidence": <float between 0 and 1>
    }
  ]
}

Use null for any field that cannot be determined from the document text.`,
		getStringValue(metadata, "file_name"),
		getStringValue(metadata, "file_type"),
		getIntValue(metadata, "word_count"),
		getIntValue(metadata, "page_count"),
		getStringValue(metadata, "source_system"),
		contextualPrompt,
		strings.Join(GetDefaultDocumentTypes(), ", "),
		text,
	)

	return prompt
}

// OpenAI API request/response structures
type openaiRequest struct {
	Model       string          `json:"model"`
	Messages    []openaiMessage `json:"messages"`
	Temperature float64         `json:"temperature"`
	MaxTokens   int             `json:"max_tokens"`
}

type openaiMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type openaiResponse struct {
	Choices []struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
	Error *struct {
		Message string `json:"message"`
		Type    string `json:"type"`
	} `json:"error,omitempty"`
}

// makeOpenAIRequest sends a request to OpenAI's API with retry logic and rate limiting
func (c *openaiClassifier) makeOpenAIRequest(ctx context.Context, prompt string) (string, error) {
	const (
		maxRetries = 5
		baseDelay  = 2 * time.Second
		maxDelay   = 60 * time.Second
	)

	var lastErr error
	for attempt := 0; attempt < maxRetries; attempt++ {
		if attempt > 0 {
			// Exponential backoff with jitter
			delay := time.Duration(float64(baseDelay) * (1.5 * float64(attempt)))
			if delay > maxDelay {
				delay = maxDelay
			}
			
			select {
			case <-ctx.Done():
				return "", ctx.Err()
			case <-time.After(delay):
				// Continue with retry
			}
		}

		response, err := c.doOpenAIRequest(ctx, prompt)
		if err == nil {
			return response, nil
		}

		lastErr = err
		
		// Check if this is a retryable error
		if !c.isRetryableError(err) {
			return "", err
		}
	}

	return "", fmt.Errorf("failed after %d attempts: %w", maxRetries, lastErr)
}

// doOpenAIRequest performs a single request to OpenAI's API
func (c *openaiClassifier) doOpenAIRequest(ctx context.Context, prompt string) (string, error) {
	reqBody := openaiRequest{
		Model: c.model,
		Messages: []openaiMessage{
			{
				Role:    "user",
				Content: prompt,
			},
		},
		Temperature: 0.1, // Low temperature for consistent results
		MaxTokens:   1500, // Increased for comprehensive responses
	}

	jsonBody, err := json.Marshal(reqBody)
	if err != nil {
		return "", fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", "https://api.openai.com/v1/chat/completions", bytes.NewBuffer(jsonBody))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.apiKey)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("OpenAI API returned status %d: %s", resp.StatusCode, string(body))
	}

	var openaiResp openaiResponse
	if err := json.Unmarshal(body, &openaiResp); err != nil {
		return "", fmt.Errorf("failed to unmarshal response: %w", err)
	}

	if openaiResp.Error != nil {
		return "", fmt.Errorf("OpenAI API error: %s", openaiResp.Error.Message)
	}

	if len(openaiResp.Choices) == 0 {
		return "", fmt.Errorf("no choices returned from OpenAI")
	}

	return openaiResp.Choices[0].Message.Content, nil
}

// parseClassificationResponse parses the enhanced OpenAI response into a ClassificationResult
func (c *openaiClassifier) parseClassificationResponse(response string) (*ClassificationResult, error) {
	// Add debugging logs to see actual response
	log.Printf("[OPENAI] Raw response length: %d chars", len(response))
	log.Printf("[OPENAI] Raw response preview: %.200s...", response)
	
	// Try to extract JSON from the response (in case there's extra text)
	jsonStart := strings.Index(response, "{")
	jsonEnd := strings.LastIndex(response, "}") + 1

	log.Printf("[OPENAI] JSON bounds: start=%d, end=%d", jsonStart, jsonEnd)

	if jsonStart == -1 || jsonEnd <= jsonStart {
		log.Printf("[OPENAI] ❌ No valid JSON brackets found in response")
		return nil, fmt.Errorf("no valid JSON found in response")
	}

	jsonStr := response[jsonStart:jsonEnd]
	log.Printf("[OPENAI] Extracted JSON length: %d chars", len(jsonStr))
	log.Printf("[OPENAI] Extracted JSON preview: %.200s...", jsonStr)

	var result ClassificationResult
	if err := json.Unmarshal([]byte(jsonStr), &result); err != nil {
		return nil, fmt.Errorf("failed to parse JSON response: %w", err)
	}

	// Validate and set defaults
	if result.DocumentType == "" {
		result.DocumentType = DocumentTypeOther
	}
	if result.LegalCategory == "" {
		result.LegalCategory = LegalCategoryCivil
	}
	if result.Confidence == 0 {
		result.Confidence = 0.5 // Default confidence
	}

	// Validate document type against known types
	validTypes := GetDefaultDocumentTypes()
	isValidType := false
	for _, validType := range validTypes {
		if result.DocumentType == validType {
			isValidType = true
			break
		}
	}
	if !isValidType {
		result.DocumentType = DocumentTypeOther
		result.Confidence = result.Confidence * 0.8 // Reduce confidence for fallback
	}

	// Ensure subject is provided if summary exists but subject doesn't
	if result.Subject == "" && result.Summary != "" {
		// Extract first sentence as subject if available
		sentences := strings.Split(result.Summary, ".")
		if len(sentences) > 0 && len(sentences[0]) > 0 {
			words := strings.Fields(sentences[0])
			if len(words) > 12 {
				result.Subject = strings.Join(words[:12], " ") + "..."
			} else {
				result.Subject = sentences[0]
			}
		}
	}

	// Set success flag
	result.Success = true

	return &result, nil
}

// Helper functions for metadata access
func getStringValue(metadata *DocumentMetadata, key string) string {
	if metadata == nil {
		return ""
	}

	switch key {
	case "file_name":
		return metadata.FileName
	case "file_type":
		return metadata.FileType
	case "source_system":
		return metadata.SourceSystem
	default:
		if metadata.Properties != nil {
			return metadata.Properties[key]
		}
		return ""
	}
}

func getIntValue(metadata *DocumentMetadata, key string) int {
	if metadata == nil {
		return 0
	}

	switch key {
	case "word_count":
		return metadata.WordCount
	case "page_count":
		return metadata.PageCount
	case "size":
		return int(metadata.Size)
	default:
		return 0
	}
}

// calculateOptimalTextLength determines the best text length based on document characteristics
func (c *openaiClassifier) calculateOptimalTextLength(metadata *DocumentMetadata) int {
	baseLength := 8000 // Conservative baseline
	
	if metadata == nil {
		return baseLength
	}
	
	// Adjust based on document size and page count
	wordCount := metadata.WordCount
	pageCount := metadata.PageCount
	
	switch {
	case wordCount < 500: // Short documents
		return baseLength // Use full text
	case wordCount < 2000: // Medium documents  
		return baseLength + 2000 // Allow more text
	case wordCount > 10000: // Large documents
		return baseLength + 4000 // Increase significantly for complex docs
	case pageCount > 20: // Multi-page documents
		return baseLength + 3000 // More context for long documents
	default:
		return baseLength
	}
}

// generateContextualPrompt creates document-specific analysis instructions
func (c *openaiClassifier) generateContextualPrompt(metadata *DocumentMetadata) string {
	if metadata == nil {
		return "CRITICAL INSTRUCTIONS:\n1. Classify document type from available types\n2. Provide substantive legal summary\n3. Extract legal entities with precision"
	}
	
	wordCount := metadata.WordCount
	pageCount := metadata.PageCount
	
	contextPrompt := "DOCUMENT ANALYSIS CONTEXT:\n"
	
	// Add analysis guidance based on document characteristics
	switch {
	case wordCount < 300:
		contextPrompt += "- SHORT DOCUMENT: Focus on key identifying elements and brief classification\n"
		contextPrompt += "- Prioritize document type identification over detailed extraction\n"
	case wordCount > 5000:
		contextPrompt += "- COMPREHENSIVE DOCUMENT: Perform detailed analysis and full entity extraction\n"
		contextPrompt += "- Extract maximum legal detail including all parties, dates, and authorities\n"
	case pageCount > 10:
		contextPrompt += "- MULTI-PAGE DOCUMENT: Analyze structure and extract section-specific information\n"
		contextPrompt += "- Look for procedural progression and case development over multiple sections\n"
	default:
		contextPrompt += "- STANDARD DOCUMENT: Perform balanced analysis with focus on legal substance\n"
	}
	
	// Add specific guidance based on file type
	fileType := strings.ToLower(metadata.FileType)
	switch {
	case strings.Contains(fileType, "pdf"):
		contextPrompt += "- PDF DOCUMENT: May contain formatted legal text, pay attention to structure\n"
	case strings.Contains(fileType, "docx"):
		contextPrompt += "- WORD DOCUMENT: Likely draft or working document, analyze for intent and completeness\n"
	case strings.Contains(fileType, "txt"):
		contextPrompt += "- TEXT DOCUMENT: May lack formatting, focus on content analysis\n"
	}
	
	return contextPrompt
}

// isRetryableError determines if an error should trigger a retry
func (c *openaiClassifier) isRetryableError(err error) bool {
	if err == nil {
		return false
	}

	errStr := err.Error()
	
	// Retry on rate limit errors (429)
	if strings.Contains(errStr, "status 429") {
		return true
	}
	
	// Retry on server errors (5xx)
	if strings.Contains(errStr, "status 5") {
		return true
	}
	
	// Retry on timeout errors
	if strings.Contains(errStr, "timeout") || strings.Contains(errStr, "context deadline exceeded") {
		return true
	}
	
	// Retry on connection errors
	if strings.Contains(errStr, "connection") || strings.Contains(errStr, "network") {
		return true
	}
	
	// Don't retry on client errors (4xx except 429), authentication errors, etc.
	return false
}
