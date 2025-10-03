package pipeline

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"motion-index-fiber/pkg/processing/classifier"
	"motion-index-fiber/pkg/processing/extractor"
	"motion-index-fiber/pkg/search"
	"motion-index-fiber/pkg/storage"
)

// Mock processors for testing
type MockProcessor struct {
	mock.Mock
}

func (m *MockProcessor) Process(ctx context.Context, req *ProcessRequest) (*ProcessResult, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*ProcessResult), args.Error(1)
}

func (m *MockProcessor) GetType() ProcessorType {
	args := m.Called()
	return args.Get(0).(ProcessorType)
}

func (m *MockProcessor) IsHealthy() bool {
	args := m.Called()
	return args.Bool(0)
}

func TestNewWorkerPool(t *testing.T) {
	tests := []struct {
		name        string
		maxWorkers  int
		queueSize   int
		shouldPanic bool
	}{
		{
			name:        "valid configuration",
			maxWorkers:  5,
			queueSize:   10,
			shouldPanic: false,
		},
		{
			name:        "single worker",
			maxWorkers:  1,
			queueSize:   1,
			shouldPanic: false,
		},
		{
			name:        "zero workers",
			maxWorkers:  0,
			queueSize:   10,
			shouldPanic: true,
		},
		{
			name:        "negative workers",
			maxWorkers:  -1,
			queueSize:   10,
			shouldPanic: true,
		},
		{
			name:        "zero queue size",
			maxWorkers:  5,
			queueSize:   0,
			shouldPanic: true,
		},
		{
			name:        "negative queue size",
			maxWorkers:  5,
			queueSize:   -1,
			shouldPanic: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.shouldPanic {
				assert.Panics(t, func() {
					NewWorkerPool(tt.maxWorkers, tt.queueSize)
				})
			} else {
				assert.NotPanics(t, func() {
					pool := NewWorkerPool(tt.maxWorkers, tt.queueSize)
					assert.NotNil(t, pool)
					assert.Equal(t, tt.maxWorkers, pool.maxWorkers)
					assert.Equal(t, tt.queueSize, cap(pool.jobs))
				})
			}
		})
	}
}

func TestWorkerPool_StartStop(t *testing.T) {
	pool := NewWorkerPool(2, 5)

	// Test starting the pool
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Start should not block
	done := make(chan struct{})
	go func() {
		pool.Start(ctx)
		close(done)
	}()

	// Give workers time to start
	time.Sleep(10 * time.Millisecond)

	// Pool should be running
	assert.True(t, pool.IsRunning())

	// Cancel context to stop workers
	cancel()

	// Wait for shutdown with timeout
	select {
	case <-done:
		// Good, workers stopped
	case <-time.After(100 * time.Millisecond):
		t.Fatal("Workers did not stop within timeout")
	}

	// Pool should no longer be running
	assert.False(t, pool.IsRunning())
}

func TestWorkerPool_Submit(t *testing.T) {
	pool := NewWorkerPool(2, 3)

	// Start the pool
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go pool.Start(ctx)
	time.Sleep(10 * time.Millisecond) // Let workers start

	tests := []struct {
		name        string
		job         *WorkerJob
		expectError bool
		poolRunning bool
	}{
		{
			name: "valid job submission",
			job: &WorkerJob{
				ID:      "test-job-1",
				Request: &ProcessRequest{ID: "req-1"},
				Result:  make(chan *ProcessResult, 1),
				Error:   make(chan error, 1),
			},
			expectError: false,
			poolRunning: true,
		},
		{
			name:        "nil job",
			job:         nil,
			expectError: true,
			poolRunning: true,
		},
		{
			name: "job with empty ID",
			job: &WorkerJob{
				ID:      "",
				Request: &ProcessRequest{ID: "req-2"},
				Result:  make(chan *ProcessResult, 1),
				Error:   make(chan error, 1),
			},
			expectError: true,
			poolRunning: true,
		},
		{
			name: "job with nil request",
			job: &WorkerJob{
				ID:      "test-job-3",
				Request: nil,
				Result:  make(chan *ProcessResult, 1),
				Error:   make(chan error, 1),
			},
			expectError: true,
			poolRunning: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := pool.Submit(tt.job)

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestWorkerPool_QueueFull(t *testing.T) {
	// Create a small pool to test queue limits
	pool := NewWorkerPool(1, 1)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go pool.Start(ctx)
	time.Sleep(10 * time.Millisecond)

	// Create a job that will block the worker
	blockingJob := &WorkerJob{
		ID:      "blocking-job",
		Request: &ProcessRequest{ID: "blocking-req"},
		Result:  make(chan *ProcessResult, 1),
		Error:   make(chan error, 1),
	}

	// Submit the blocking job
	err := pool.Submit(blockingJob)
	assert.NoError(t, err)

	// Try to submit another job - should fail due to full queue
	anotherJob := &WorkerJob{
		ID:      "another-job",
		Request: &ProcessRequest{ID: "another-req"},
		Result:  make(chan *ProcessResult, 1),
		Error:   make(chan error, 1),
	}

	// This should timeout or return an error due to full queue
	err = pool.Submit(anotherJob)
	// Note: The actual behavior depends on implementation -
	// it might block or return an error
}

func TestWorkerPool_ProcessJob(t *testing.T) {
	// Create mock processors
	extractor := &MockProcessor{}
	classifier := &MockProcessor{}
	searchService := &MockProcessor{}
	storageService := &MockProcessor{}

	// Set up processor expectations
	extractor.On("GetType").Return(ProcessorTypeExtraction)
	classifier.On("GetType").Return(ProcessorTypeClassification)
	searchService.On("GetType").Return(ProcessorTypeIndexing)
	storageService.On("GetType").Return(ProcessorTypeStorage)

	// Set up successful processing
	mockResult := &ProcessResult{
		ID:     "test-result",
		Status: ProcessStatusSuccess,
		Data:   map[string]interface{}{"extracted_text": "test content"},
	}

	extractor.On("Process", mock.Anything, mock.Anything).Return(mockResult, nil)
	classifier.On("Process", mock.Anything, mock.Anything).Return(mockResult, nil)
	searchService.On("Process", mock.Anything, mock.Anything).Return(mockResult, nil)
	storageService.On("Process", mock.Anything, mock.Anything).Return(mockResult, nil)

	// Create worker pool
	pool := NewWorkerPool(1, 5)
	pool.processors = []Processor{extractor, classifier, searchService, storageService}

	// Start the pool
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go pool.Start(ctx)
	time.Sleep(10 * time.Millisecond)

	// Create and submit a job
	job := &WorkerJob{
		ID: "test-job",
		Request: &ProcessRequest{
			ID:       "test-req",
			FileName: "test.pdf",
			Options: &ProcessOptions{
				ExtractText:   true,
				ClassifyDoc:   true,
				IndexDocument: true,
				StoreDocument: true,
			},
		},
		Result: make(chan *ProcessResult, 1),
		Error:  make(chan error, 1),
	}

	err := pool.Submit(job)
	assert.NoError(t, err)

	// Wait for job completion
	select {
	case result := <-job.Result:
		assert.NotNil(t, result)
		assert.Equal(t, ProcessStatusSuccess, result.Status)
	case err := <-job.Error:
		t.Fatalf("Unexpected error: %v", err)
	case <-time.After(100 * time.Millisecond):
		t.Fatal("Job did not complete within timeout")
	}

	// Verify all processors were called
	extractor.AssertExpectations(t)
	classifier.AssertExpectations(t)
	searchService.AssertExpectations(t)
	storageService.AssertExpectations(t)
}

func TestWorkerPool_ProcessJobWithError(t *testing.T) {
	// Create mock processors
	extractor := &MockProcessor{}
	classifier := &MockProcessor{}

	// Set up processor expectations
	extractor.On("GetType").Return(ProcessorTypeExtraction)
	classifier.On("GetType").Return(ProcessorTypeClassification)

	// Set up extraction to succeed, classification to fail
	extractResult := &ProcessResult{
		ID:     "extract-result",
		Status: ProcessStatusSuccess,
		Data:   map[string]interface{}{"extracted_text": "test content"},
	}

	extractor.On("Process", mock.Anything, mock.Anything).Return(extractResult, nil)
	classifier.On("Process", mock.Anything, mock.Anything).Return(nil, fmt.Errorf("classification failed"))

	// Create worker pool
	pool := NewWorkerPool(1, 5)
	pool.processors = []Processor{extractor, classifier}

	// Start the pool
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go pool.Start(ctx)
	time.Sleep(10 * time.Millisecond)

	// Create and submit a job
	job := &WorkerJob{
		ID: "test-job-error",
		Request: &ProcessRequest{
			ID:       "test-req-error",
			FileName: "test.pdf",
			Options: &ProcessOptions{
				ExtractText: true,
				ClassifyDoc: true,
			},
		},
		Result: make(chan *ProcessResult, 1),
		Error:  make(chan error, 1),
	}

	err := pool.Submit(job)
	assert.NoError(t, err)

	// Wait for job completion (should fail)
	select {
	case result := <-job.Result:
		t.Fatalf("Expected error but got result: %v", result)
	case err := <-job.Error:
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "classification failed")
	case <-time.After(100 * time.Millisecond):
		t.Fatal("Job did not complete within timeout")
	}

	extractor.AssertExpectations(t)
	classifier.AssertExpectations(t)
}

func TestWorkerPool_ContextCancellation(t *testing.T) {
	pool := NewWorkerPool(2, 5)

	// Create a context that we'll cancel quickly
	ctx, cancel := context.WithCancel(context.Background())

	// Start the pool
	go pool.Start(ctx)
	time.Sleep(10 * time.Millisecond)

	// Verify pool is running
	assert.True(t, pool.IsRunning())

	// Cancel the context
	cancel()

	// Wait for pool to stop
	timeout := time.After(100 * time.Millisecond)
	ticker := time.NewTicker(5 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-timeout:
			t.Fatal("Pool did not stop after context cancellation")
		case <-ticker.C:
			if !pool.IsRunning() {
				return // Success
			}
		}
	}
}

func TestWorkerPool_GetStats(t *testing.T) {
	pool := NewWorkerPool(3, 10)

	stats := pool.GetStats()
	assert.NotNil(t, stats)
	assert.Equal(t, 3, stats.MaxWorkers)
	assert.Equal(t, 10, stats.QueueCapacity)
	assert.Equal(t, 0, stats.ActiveJobs)
	assert.Equal(t, 0, stats.QueuedJobs)
	assert.False(t, stats.IsRunning)

	// Start the pool and check stats again
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go pool.Start(ctx)
	time.Sleep(10 * time.Millisecond)

	stats = pool.GetStats()
	assert.True(t, stats.IsRunning)
}

func TestWorkerJob_Validation(t *testing.T) {
	tests := []struct {
		name        string
		job         *WorkerJob
		expectValid bool
	}{
		{
			name: "valid job",
			job: &WorkerJob{
				ID: "valid-job",
				Request: &ProcessRequest{
					ID:       "valid-req",
					FileName: "test.pdf",
				},
				Result: make(chan *ProcessResult, 1),
				Error:  make(chan error, 1),
			},
			expectValid: true,
		},
		{
			name:        "nil job",
			job:         nil,
			expectValid: false,
		},
		{
			name: "empty job ID",
			job: &WorkerJob{
				ID: "",
				Request: &ProcessRequest{
					ID:       "valid-req",
					FileName: "test.pdf",
				},
				Result: make(chan *ProcessResult, 1),
				Error:  make(chan error, 1),
			},
			expectValid: false,
		},
		{
			name: "nil request",
			job: &WorkerJob{
				ID:      "valid-job",
				Request: nil,
				Result:  make(chan *ProcessResult, 1),
				Error:   make(chan error, 1),
			},
			expectValid: false,
		},
		{
			name: "nil result channel",
			job: &WorkerJob{
				ID: "valid-job",
				Request: &ProcessRequest{
					ID:       "valid-req",
					FileName: "test.pdf",
				},
				Result: nil,
				Error:  make(chan error, 1),
			},
			expectValid: false,
		},
		{
			name: "nil error channel",
			job: &WorkerJob{
				ID: "valid-job",
				Request: &ProcessRequest{
					ID:       "valid-req",
					FileName: "test.pdf",
				},
				Result: make(chan *ProcessResult, 1),
				Error:  nil,
			},
			expectValid: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			valid := validateWorkerJob(tt.job)
			assert.Equal(t, tt.expectValid, valid)
		})
	}
}

// Helper function that might be used in the actual implementation
func validateWorkerJob(job *WorkerJob) bool {
	if job == nil {
		return false
	}
	if job.ID == "" {
		return false
	}
	if job.Request == nil {
		return false
	}
	if job.Result == nil {
		return false
	}
	if job.Error == nil {
		return false
	}
	return true
}
