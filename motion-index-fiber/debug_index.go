package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/joho/godotenv"

	"motion-index-fiber/internal/config"
	"motion-index-fiber/pkg/cloud/digitalocean"
	searchModels "motion-index-fiber/pkg/search/models"
)

func main() {
	log.Printf("🔧 Debug OpenSearch Indexing Issue")

	// Load environment
	if err := godotenv.Load(); err != nil {
		log.Printf("Warning: .env file not found: %v", err)
	}

	// Load configuration
	_, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Initialize DigitalOcean provider
	provider, err := digitalocean.NewProviderFromEnvironment()
	if err != nil {
		log.Fatalf("Failed to create DigitalOcean provider: %v", err)
	}

	if err := provider.Initialize(); err != nil {
		log.Fatalf("Failed to initialize provider: %v", err)
	}

	services := provider.GetServices()

	// Test search service health in detail
	log.Printf("🔍 Checking search service health...")
	ctx := context.Background()
	
	healthStatus, err := services.Search.Health(ctx)
	if err != nil {
		log.Printf("❌ Health check failed: %v", err)
	} else {
		log.Printf("✅ Health check passed:")
		log.Printf("   Status: %s", healthStatus.Status)
		log.Printf("   Cluster: %s", healthStatus.ClusterName)
		log.Printf("   Nodes: %d", healthStatus.NumberOfNodes)
		log.Printf("   Active Shards: %d", healthStatus.ActiveShards)
		log.Printf("   Index Exists: %v", healthStatus.IndexExists)
		log.Printf("   Index Health: %s", healthStatus.IndexHealth)
	}

	// Try a simple document first
	log.Printf("📄 Creating minimal test document...")
	now := time.Now()
	
	simpleDoc := &searchModels.Document{
		ID:          fmt.Sprintf("simple-test-%d", now.Unix()),
		FileName:    "simple-test.txt",
		FilePath:    "test/simple.txt",
		Text:        "This is a simple test document.",
		DocType:     "other",
		Category:    "Test",
		Hash:        "simple_hash",
		CreatedAt:   now,
		UpdatedAt:   now,
		ContentType: "text/plain",
		Size:        31,
		Metadata: &searchModels.DocumentMetadata{
			DocumentName: "Simple Test",
			Subject:      "Test Document",
			DocumentType: searchModels.DocTypeOther,
			ProcessedAt:  now,
			AIClassified: false,
		},
	}

	// Print the document as JSON to see the structure
	docJSON, _ := json.MarshalIndent(simpleDoc, "", "  ")
	log.Printf("📋 Document JSON structure:")
	log.Printf("%s", string(docJSON))

	// Try to index the simple document
	log.Printf("🔄 Indexing simple document...")
	indexID, err := services.Search.IndexDocument(ctx, simpleDoc)
	if err != nil {
		log.Printf("❌ Simple document indexing failed: %v", err)
		
		// Let's check if there are existing documents to understand the expected format
		log.Printf("🔍 Checking existing documents...")
		searchReq := &searchModels.SearchRequest{
			Query: "*",
			Size:  1,
		}
		
		searchResult, err := services.Search.SearchDocuments(ctx, searchReq)
		if err != nil {
			log.Printf("❌ Search failed: %v", err)
		} else {
			log.Printf("📊 Found %d existing documents", len(searchResult.Documents))
			if len(searchResult.Documents) > 0 {
				existingDoc := searchResult.Documents[0]
				log.Printf("📄 Example existing document structure:")
				existingJSON, _ := json.MarshalIndent(existingDoc.Document, "", "  ")
				log.Printf("%s", string(existingJSON))
			}
		}
	} else {
		log.Printf("✅ Simple document indexed successfully!")
		log.Printf("   Index ID: %s", indexID)
		
		// Wait and try to retrieve
		time.Sleep(2 * time.Second)
		retrievedDoc, err := services.Search.GetDocument(ctx, simpleDoc.ID)
		if err != nil {
			log.Printf("❌ Failed to retrieve: %v", err)
		} else {
			log.Printf("✅ Retrieved successfully: %s", retrievedDoc.ID)
		}
	}

	// Get document statistics
	log.Printf("📊 Current index statistics...")
	stats, err := services.Search.GetDocumentStats(ctx)
	if err != nil {
		log.Printf("❌ Failed to get stats: %v", err)
	} else {
		log.Printf("📊 Index has %d total documents", stats.TotalDocuments)
	}
}