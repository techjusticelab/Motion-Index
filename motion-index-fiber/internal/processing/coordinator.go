package processing

import (
	"context"
	"fmt"
	"log"
	"sync/atomic"
	"time"

	"motion-index-fiber/internal/config"
	"motion-index-fiber/internal/hardware"
	"motion-index-fiber/pkg/cloud/digitalocean"
	"motion-index-fiber/pkg/processing/classifier"
	"motion-index-fiber/pkg/processing/extractor"
	"motion-index-fiber/pkg/processing/gpu"
	"motion-index-fiber/pkg/processing/queue"
	"motion-index-fiber/pkg/storage"
)

// DocumentCoordinator orchestrates document processing workflows
// Following UNIX philosophy: coordinate small, focused tools
type DocumentCoordinator struct {
	config       *config.Config
	workerConfig *hardware.WorkerConfig
	hardwareInfo *hardware.Analysis
	
	// Services
	storage    storage.Service
	extractor  extractor.Service
	classifier classifier.Service
	search     interface{} // Search service interface
	gpu        gpu.GPUAccelerator
	
	// Queue management
	queueManager queue.QueueManager
	
	// Metrics
	processed *int64
	errors    *int64
	skipped   *int64
}

// CoordinatorConfig holds configuration for document coordination
type CoordinatorConfig struct {
	Config       *config.Config
	WorkerConfig *hardware.WorkerConfig
	HardwareInfo *hardware.Analysis
}

// NewDocumentCoordinator creates a new document processing coordinator
func NewDocumentCoordinator(cfg *CoordinatorConfig) (*DocumentCoordinator, error) {
	coordinator := &DocumentCoordinator{
		config:       cfg.Config,
		workerConfig: cfg.WorkerConfig,
		hardwareInfo: cfg.HardwareInfo,
		processed:    new(int64),
		errors:       new(int64),
		skipped:      new(int64),
	}
	
	if err := coordinator.initializeServices(); err != nil {
		return nil, fmt.Errorf("failed to initialize services: %w", err)
	}
	
	if err := coordinator.initializeQueues(); err != nil {
		return nil, fmt.Errorf("failed to initialize queues: %w", err)
	}
	
	return coordinator, nil
}

// initializeServices initializes all required services
func (c *DocumentCoordinator) initializeServices() error {
	// Initialize DigitalOcean provider
	provider, err := digitalocean.NewProviderFromEnvironment()
	if err != nil {
		return fmt.Errorf("failed to create DigitalOcean provider: %w", err)
	}
	
	if err := provider.Initialize(); err != nil {
		return fmt.Errorf("failed to initialize provider: %w", err)
	}
	
	services := provider.GetServices()
	c.storage = services.Storage
	c.search = services.Search
	
	// Initialize extraction service with enhanced capabilities
	c.extractor = extractor.NewExtractionService(&extractor.ServiceConfig{
		Type:          extractor.ServiceTypeAuto,
		EnableOCR:     true,
		EnableDslipak: true,
	})
	
	// Initialize classification service
	c.classifier = classifier.NewOpenAIClassifier(&classifier.OpenAIConfig{
		APIKey: c.config.OpenAI.APIKey,
		Model:  c.config.OpenAI.Model,
	})
	
	// Initialize GPU accelerator
	c.gpu = gpu.NewNVIDIAAccelerator()
	
	return nil
}

// initializeQueues sets up processing queues with optimal configuration
func (c *DocumentCoordinator) initializeQueues() error {
	c.queueManager = queue.NewQueueManager()
	
	ctx := context.Background()
	
	// Extraction queue - CPU intensive
	_, err := c.queueManager.CreateQueue(&queue.QueueConfig{
		Name:           "extraction",
		Type:           queue.QueueTypeExtraction,
		MaxSize:        1000,
		WorkerCount:    c.workerConfig.Extraction,
		ProcessTimeout: 5 * time.Minute,
		RetryAttempts:  3,
		RetryDelay:     10 * time.Second,
		EnableMetrics:  true,
	}, c.createExtractionProcessor())
	
	if err != nil {
		return fmt.Errorf("failed to create extraction queue: %w", err)
	}
	
	// Classification queue - API rate limited
	_, err = c.queueManager.CreateQueue(&queue.QueueConfig{
		Name:            "classification",
		Type:            queue.QueueTypeClassification,
		MaxSize:         500,
		WorkerCount:     c.workerConfig.Classification,
		ProcessTimeout:  2 * time.Minute,
		RetryAttempts:   5,
		RetryDelay:      30 * time.Second,
		EnableRateLimit: true,
		RateLimit:       50, // 50 requests per minute
		BurstSize:       10,
		EnableMetrics:   true,
	}, c.createClassificationProcessor())
	
	if err != nil {
		return fmt.Errorf("failed to create classification queue: %w", err)
	}
	
	// Indexing queue - Network I/O
	_, err = c.queueManager.CreateQueue(&queue.QueueConfig{
		Name:           "indexing",
		Type:           queue.QueueTypeIndexing,
		MaxSize:        1000,
		WorkerCount:    c.workerConfig.Indexing,
		ProcessTimeout: 1 * time.Minute,
		RetryAttempts:  3,
		RetryDelay:     5 * time.Second,
		EnableMetrics:  true,
	}, c.createIndexingProcessor())
	
	if err != nil {
		return fmt.Errorf("failed to create indexing queue: %w", err)
	}
	
	// Start queue manager
	return c.queueManager.Start(ctx)
}

// ClassifyAll processes all documents in storage
func (c *DocumentCoordinator) ClassifyAll(ctx context.Context) error {
	log.Println("📋 Listing all documents from storage...")
	
	objects, err := c.storage.ListObjects(ctx, "")
	if err != nil {
		return fmt.Errorf("failed to list objects: %w", err)
	}
	
	log.Printf("🚀 Processing %d documents for classification...", len(objects))
	
	return c.processDocuments(ctx, objects)
}

// ClassifyBatch processes documents in specified batch size
func (c *DocumentCoordinator) ClassifyBatch(ctx context.Context, batchSize int) error {
	log.Printf("📋 Listing documents for batch processing (limit: %d)...", batchSize)
	
	objects, err := c.storage.ListObjects(ctx, "")
	if err != nil {
		return fmt.Errorf("failed to list objects: %w", err)
	}
	
	// Limit to batch size
	if len(objects) > batchSize {
		objects = objects[:batchSize]
	}
	
	log.Printf("🚀 Processing %d documents in batch...", len(objects))
	
	return c.processDocuments(ctx, objects)
}

// ClassifyFiles processes specific files
func (c *DocumentCoordinator) ClassifyFiles(ctx context.Context, files []string) error {
	log.Printf("🚀 Processing %d specific files...", len(files))
	
	// Convert file paths to storage objects
	objects := make([]*storage.StorageObject, len(files))
	for i, file := range files {
		objects[i] = &storage.StorageObject{
			Key:  file,
			Path: file,
		}
	}
	
	return c.processDocuments(ctx, objects)
}

// processDocuments handles the core document processing workflow
func (c *DocumentCoordinator) processDocuments(ctx context.Context, objects []*storage.StorageObject) error {
	totalDocuments := len(objects)
	if totalDocuments == 0 {
		log.Println("✅ No documents to process")
		return nil
	}
	
	// Start progress monitoring
	stopMonitor := make(chan bool)
	go c.monitorProgress(stopMonitor, int64(totalDocuments))
	defer func() { stopMonitor <- true }()
	
	// Process documents through the pipeline
	for _, obj := range objects {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			// Queue document for processing
			if err := c.queueDocument(ctx, obj); err != nil {
				log.Printf("⚠️ Failed to queue document %s: %v", obj.Key, err)
				atomic.AddInt64(c.errors, 1)
				continue
			}
		}
	}
	
	// Wait for all processing to complete
	if err := c.waitForCompletion(ctx); err != nil {
		return fmt.Errorf("processing incomplete: %w", err)
	}
	
	c.printSummary()
	return nil
}

// queueDocument adds a document to the processing pipeline
func (c *DocumentCoordinator) queueDocument(ctx context.Context, obj *storage.StorageObject) error {
	// Create processing job
	job := &ProcessingJob{
		DocumentKey: obj.Key,
		DocumentPath: obj.Path,
		Timestamp:   time.Now(),
	}
	
	// Start with extraction
	queue, err := c.queueManager.GetQueue("extraction")
	if err != nil {
		return err
	}
	
	return queue.Enqueue(ctx, job)
}

// waitForCompletion waits for all queues to finish processing
func (c *DocumentCoordinator) waitForCompletion(ctx context.Context) error {
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()
	
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
			if c.queueManager.IsEmpty() {
				return nil
			}
		}
	}
}

// monitorProgress provides real-time progress updates
func (c *DocumentCoordinator) monitorProgress(stop <-chan bool, total int64) {
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()
	
	for {
		select {
		case <-stop:
			return
		case <-ticker.C:
			processed := atomic.LoadInt64(c.processed)
			errors := atomic.LoadInt64(c.errors)
			skipped := atomic.LoadInt64(c.skipped)
			
			progress := float64(processed+errors+skipped) / float64(total) * 100
			
			log.Printf("📊 Progress: %.1f%% (%d/%d) | ✅ %d processed | ❌ %d errors | ⏭️ %d skipped",
				progress, processed+errors+skipped, total, processed, errors, skipped)
		}
	}
}

// printSummary displays final processing statistics
func (c *DocumentCoordinator) printSummary() {
	processed := atomic.LoadInt64(c.processed)
	errors := atomic.LoadInt64(c.errors)
	skipped := atomic.LoadInt64(c.skipped)
	
	fmt.Println("\n📊 Processing Summary")
	fmt.Println("====================")
	fmt.Printf("✅ Processed: %d\n", processed)
	fmt.Printf("❌ Errors: %d\n", errors)
	fmt.Printf("⏭️ Skipped: %d\n", skipped)
	fmt.Printf("📈 Total: %d\n", processed+errors+skipped)
}

// Close shuts down the coordinator and cleans up resources
func (c *DocumentCoordinator) Close() error {
	if c.queueManager != nil {
		return c.queueManager.Stop(context.Background())
	}
	return nil
}

// ProcessingJob represents a document processing job
type ProcessingJob struct {
	DocumentKey  string    `json:"document_key"`
	DocumentPath string    `json:"document_path"`
	Timestamp    time.Time `json:"timestamp"`
	
	// Processing results
	ExtractedText string                 `json:"extracted_text,omitempty"`
	Classification string               `json:"classification,omitempty"`
	Metadata      map[string]interface{} `json:"metadata,omitempty"`
}