package classifier

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

// ollamaClassifier implements classification using Ollama local models
type ollamaClassifier struct {
	baseURL    string
	model      string
	httpClient *http.Client
}

// NewOllamaClassifier creates a new Ollama-based classifier
func NewOllamaClassifier(config *OllamaConfig) (Classifier, error) {
	baseURL := config.BaseURL
	if baseURL == "" {
		baseURL = "http://localhost:11434"
	}

	model := config.Model
	if model == "" {
		model = "gpt-oss:20b"
	}

	timeout := config.Timeout
	if timeout == 0 {
		timeout = 120 * time.Second
	}

	return &ollamaClassifier{
		baseURL: baseURL,
		model:   model,
		httpClient: &http.Client{
			Timeout: timeout,
		},
	}, nil
}

// Classify analyzes document text using Ollama and returns classification results
func (o *ollamaClassifier) Classify(ctx context.Context, text string, metadata *DocumentMetadata) (*ClassificationResult, error) {
	// Create the classification prompt
	prompt := o.buildClassificationPrompt(text, metadata)

	// Make request to Ollama
	response, err := o.makeOllamaRequest(ctx, prompt)
	if err != nil {
		return nil, NewClassificationError("ollama_request", "failed to classify document", err)
	}

	// Parse the response
	result, err := o.parseClassificationResponse(response)
	if err != nil {
		return nil, NewClassificationError("response_parsing", "failed to parse classification response", err)
	}

	return result, nil
}

// GetSupportedCategories returns the categories this classifier can identify
func (o *ollamaClassifier) GetSupportedCategories() []string {
	return GetDefaultCategories()
}

// IsConfigured returns true if the classifier is properly configured
func (o *ollamaClassifier) IsConfigured() bool {
	return o.baseURL != "" && o.model != ""
}

// buildClassificationPrompt creates a prompt for Ollama classification
func (o *ollamaClassifier) buildClassificationPrompt(text string, metadata *DocumentMetadata) string {
	// Limit text length for Ollama (smaller context window)
	maxTextLength := 4000
	if len(text) > maxTextLength {
		text = text[:maxTextLength] + "..."
	}

	fileName := "unknown"
	if metadata != nil && metadata.FileName != "" {
		fileName = metadata.FileName
	}

	prompt := fmt.Sprintf(`You are an expert legal document analyzer. Analyze this legal document and classify it.

Document: %s

Text:
%s

Available document types: %s

Respond with ONLY a JSON object:
{
  "document_type": "<type from available types>",
  "legal_category": "<criminal|civil|traffic|family>",
  "subject": "<brief subject line>",
  "summary": "<legal summary in 2-3 sentences>",
  "confidence": <0.0 to 1.0>,
  "keywords": ["<key legal terms>"],
  "legal_tags": ["<legal categories>"],
  "case_info": {
    "case_number": "<case number or null>",
    "case_name": "<case title or null>",
    "case_type": "<criminal|civil|traffic|family>"
  },
  "filing_date": "<YYYY-MM-DD or null>",
  "entities": [
    {
      "text": "<entity>",
      "type": "<PERSON|ORGANIZATION|LOCATION|DATE|STATUTE>",
      "confidence": <0.0 to 1.0>
    }
  ]
}`,
		fileName,
		text,
		strings.Join(GetDefaultDocumentTypes(), ", "),
	)

	return prompt
}

// Ollama API request/response structures
type ollamaRequest struct {
	Model  string `json:"model"`
	Prompt string `json:"prompt"`
	Stream bool   `json:"stream"`
}

type ollamaResponse struct {
	Response string `json:"response"`
	Done     bool   `json:"done"`
	Error    string `json:"error,omitempty"`
}

// makeOllamaRequest sends a request to Ollama's API
func (o *ollamaClassifier) makeOllamaRequest(ctx context.Context, prompt string) (string, error) {
	reqBody := ollamaRequest{
		Model:  o.model,
		Prompt: prompt,
		Stream: false, // We want the complete response
	}

	jsonBody, err := json.Marshal(reqBody)
	if err != nil {
		return "", fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", o.baseURL+"/api/generate", bytes.NewBuffer(jsonBody))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := o.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("Ollama API returned status %d: %s", resp.StatusCode, string(body))
	}

	var ollamaResp ollamaResponse
	if err := json.Unmarshal(body, &ollamaResp); err != nil {
		return "", fmt.Errorf("failed to unmarshal response: %w", err)
	}

	if ollamaResp.Error != "" {
		return "", fmt.Errorf("Ollama API error: %s", ollamaResp.Error)
	}

	if !ollamaResp.Done {
		return "", fmt.Errorf("Ollama response not complete")
	}

	return ollamaResp.Response, nil
}

// parseClassificationResponse parses the Ollama response into a ClassificationResult
func (o *ollamaClassifier) parseClassificationResponse(response string) (*ClassificationResult, error) {
	// Try to extract JSON from the response
	jsonStart := strings.Index(response, "{")
	jsonEnd := strings.LastIndex(response, "}") + 1

	if jsonStart == -1 || jsonEnd <= jsonStart {
		return nil, fmt.Errorf("no valid JSON found in response")
	}

	jsonStr := response[jsonStart:jsonEnd]

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
		result.Confidence = 0.3 // Lower default confidence for local model
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
		result.Confidence = result.Confidence * 0.7 // Further reduce confidence for fallback
	}

	// Ensure subject is provided if summary exists but subject doesn't
	if result.Subject == "" && result.Summary != "" {
		// Extract first sentence as subject if available
		sentences := strings.Split(result.Summary, ".")
		if len(sentences) > 0 && len(sentences[0]) > 0 {
			words := strings.Fields(sentences[0])
			if len(words) > 10 {
				result.Subject = strings.Join(words[:10], " ") + "..."
			} else {
				result.Subject = sentences[0]
			}
		}
	}

	// Set success flag
	result.Success = true

	return &result, nil
}