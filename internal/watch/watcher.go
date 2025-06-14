package watch

import (
	"fmt"
	"path/filepath"
	"sync"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/singoesdeep/templater/internal/engine"
)

// Watcher represents a file system watcher for templates and data files
type Watcher struct {
	watcher    *fsnotify.Watcher
	template   string
	data       string
	output     string
	interval   time.Duration
	lastUpdate time.Time
	mu         sync.Mutex
	stopChan   chan struct{}
	statusChan chan Status
	errorChan  chan error
}

// Status represents the current status of the watcher
type Status struct {
	TemplatePath string
	DataPath     string
	OutputPath   string
	LastUpdate   time.Time
	IsWatching   bool
}

// NewWatcher creates a new file system watcher
func NewWatcher(template, data, output string, interval time.Duration) (*Watcher, error) {
	w, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, fmt.Errorf("error creating watcher: %w", err)
	}

	return &Watcher{
		watcher:    w,
		template:   template,
		data:       data,
		output:     output,
		interval:   interval,
		stopChan:   make(chan struct{}),
		statusChan: make(chan Status, 1),
		errorChan:  make(chan error, 1),
	}, nil
}

// Start begins watching for file changes
func (w *Watcher) Start() error {
	// Add template file to watcher
	if err := w.watcher.Add(filepath.Dir(w.template)); err != nil {
		return fmt.Errorf("error watching template directory: %w", err)
	}

	// Add data file to watcher if it exists
	if w.data != "" {
		if err := w.watcher.Add(filepath.Dir(w.data)); err != nil {
			return fmt.Errorf("error watching data directory: %w", err)
		}
	}

	// Start watching in a goroutine
	go w.watch()

	return nil
}

// Stop stops the watcher
func (w *Watcher) Stop() {
	close(w.stopChan)
	w.watcher.Close()
}

// Status returns the current status of the watcher
func (w *Watcher) Status() Status {
	w.mu.Lock()
	defer w.mu.Unlock()

	return Status{
		TemplatePath: w.template,
		DataPath:     w.data,
		OutputPath:   w.output,
		LastUpdate:   w.lastUpdate,
		IsWatching:   true,
	}
}

// StatusChannel returns the channel for status updates
func (w *Watcher) StatusChannel() <-chan Status {
	return w.statusChan
}

// ErrorChannel returns the channel for error updates
func (w *Watcher) ErrorChannel() <-chan error {
	return w.errorChan
}

// watch monitors files for changes and triggers regeneration
func (w *Watcher) watch() {
	for {
		select {
		case event, ok := <-w.watcher.Events:
			if !ok {
				return
			}

			// Check if the event is for our watched files
			if w.isRelevantEvent(event) {
				// Debounce changes
				time.Sleep(w.interval)

				// Process the change
				if err := w.processChange(); err != nil {
					w.errorChan <- err
				} else {
					w.updateStatus()
				}
			}

		case err, ok := <-w.watcher.Errors:
			if !ok {
				return
			}
			w.errorChan <- fmt.Errorf("watcher error: %w", err)

		case <-w.stopChan:
			return
		}
	}
}

// isRelevantEvent checks if the event is for our watched files
func (w *Watcher) isRelevantEvent(event fsnotify.Event) bool {
	// Check if the event is for our template or data file
	return (event.Name == w.template || event.Name == w.data) &&
		(event.Op&fsnotify.Write == fsnotify.Write || event.Op&fsnotify.Create == fsnotify.Create)
}

// processChange handles file changes by regenerating the output
func (w *Watcher) processChange() error {
	// Load data if available
	var data map[string]string
	var err error
	if w.data != "" {
		data, err = engine.LoadData(w.data)
		if err != nil {
			return fmt.Errorf("error loading data: %w", err)
		}
	}

	// Render template
	result, err := engine.RenderTemplate(w.template, data)
	if err != nil {
		return fmt.Errorf("error rendering template: %w", err)
	}

	// Write output
	if err := engine.WriteToFile(w.output, result); err != nil {
		return fmt.Errorf("error writing output: %w", err)
	}

	// Update last update time
	w.mu.Lock()
	w.lastUpdate = time.Now()
	w.mu.Unlock()

	return nil
}

// updateStatus sends the current status to the status channel
func (w *Watcher) updateStatus() {
	w.mu.Lock()
	status := Status{
		TemplatePath: w.template,
		DataPath:     w.data,
		OutputPath:   w.output,
		LastUpdate:   w.lastUpdate,
		IsWatching:   true,
	}
	w.mu.Unlock()

	select {
	case w.statusChan <- status:
	default:
		// Channel is full, skip update
	}
}
