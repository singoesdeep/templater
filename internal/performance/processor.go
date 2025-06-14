package performance

import (
	"fmt"
	"runtime"
	"sync"
	"time"

	"github.com/singoesdeep/templater/internal/engine"
)

// ProcessingStats tracks performance metrics
type ProcessingStats struct {
	StartTime      time.Time
	EndTime        time.Time
	TemplateCount  int
	ErrorCount     int
	MemoryUsage    uint64
	ProcessingTime time.Duration
	WorkerCount    int
	CacheHits      int
	CacheMisses    int
}

// ConcurrentProcessor handles parallel template processing
type ConcurrentProcessor struct {
	MaxWorkers int
	Stats      *ProcessingStats
	mu         sync.RWMutex
	workerPool chan struct{}
	stopChan   chan struct{}
}

// NewConcurrentProcessor creates a new processor with optimal worker count
func NewConcurrentProcessor() *ConcurrentProcessor {
	// Use number of CPU cores for optimal performance
	numCPU := runtime.NumCPU()
	// Limit maximum workers to prevent resource exhaustion
	if numCPU > 8 {
		numCPU = 8
	}
	return &ConcurrentProcessor{
		MaxWorkers: numCPU,
		Stats: &ProcessingStats{
			StartTime:   time.Now(),
			WorkerCount: numCPU,
		},
		workerPool: make(chan struct{}, numCPU),
		stopChan:   make(chan struct{}),
	}
}

// SetMaxWorkers adjusts the maximum number of concurrent workers
func (p *ConcurrentProcessor) SetMaxWorkers(count int) {
	if count < 1 {
		count = 1
	}
	if count > 16 {
		count = 16
	}
	p.mu.Lock()
	p.MaxWorkers = count
	p.workerPool = make(chan struct{}, count)
	p.Stats.WorkerCount = count
	p.mu.Unlock()
}

// ProcessTemplates concurrently processes multiple templates
func (p *ConcurrentProcessor) ProcessTemplates(templates []string, data map[string]string) (map[string]string, error) {
	results := make(map[string]string, len(templates))
	var wg sync.WaitGroup
	errors := make(chan error, len(templates))
	resultChan := make(chan struct {
		path   string
		result string
	}, len(templates))

	// Process templates concurrently
	for _, tmpl := range templates {
		select {
		case <-p.stopChan:
			return results, fmt.Errorf("processing stopped")
		case p.workerPool <- struct{}{}: // Acquire worker slot
			wg.Add(1)
			go func(templatePath string) {
				defer wg.Done()
				defer func() { <-p.workerPool }() // Release worker slot

				// Process template
				result, err := engine.RenderTemplate(templatePath, data)
				if err != nil {
					errors <- fmt.Errorf("error processing %s: %w", templatePath, err)
					p.mu.Lock()
					p.Stats.ErrorCount++
					p.mu.Unlock()
					return
				}

				// Send result through channel
				resultChan <- struct {
					path   string
					result string
				}{templatePath, result}

				// Update stats
				p.mu.Lock()
				p.Stats.TemplateCount++
				p.mu.Unlock()
			}(tmpl)
		}
	}

	// Wait for all templates to be processed
	wg.Wait()
	close(errors)
	close(resultChan)

	// Collect results
	for result := range resultChan {
		results[result.path] = result.result
	}

	// Collect errors
	var errs []error
	for err := range errors {
		errs = append(errs, err)
	}

	// Update stats
	p.mu.Lock()
	p.Stats.EndTime = time.Now()
	p.Stats.ProcessingTime = p.Stats.EndTime.Sub(p.Stats.StartTime)
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	p.Stats.MemoryUsage = m.Alloc
	p.mu.Unlock()

	// Trigger garbage collection
	runtime.GC()

	if len(errs) > 0 {
		return results, fmt.Errorf("processing errors: %v", errs)
	}

	return results, nil
}

// Stop stops the processor and cancels any ongoing operations
func (p *ConcurrentProcessor) Stop() {
	close(p.stopChan)
}

// GetStats returns the current processing statistics
func (p *ConcurrentProcessor) GetStats() *ProcessingStats {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return p.Stats
}

// OptimizeMemory triggers garbage collection and memory optimization
func OptimizeMemory() {
	// Force garbage collection
	runtime.GC()

	// Return memory to the OS
	debug := runtime.MemStats{}
	runtime.ReadMemStats(&debug)
	if debug.Alloc > 100*1024*1024 { // If using more than 100MB
		runtime.GC()
	}
}

// MonitorResources starts resource monitoring at specified intervals
func MonitorResources(interval time.Duration, callback func(*ProcessingStats)) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for range ticker.C {
		var m runtime.MemStats
		runtime.ReadMemStats(&m)

		stats := &ProcessingStats{
			MemoryUsage: m.Alloc,
		}

		callback(stats)
	}
}
