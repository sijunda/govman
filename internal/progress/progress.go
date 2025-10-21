package progress

import (
	"fmt"
	"strings"
	"sync"
	"time"

	util "github.com/sijunda/govman/internal/util"
)

// Pre-allocated buffers to reduce allocations
const (
	defaultBarWidth = 50
	fillChar        = "█"
	emptyChar       = "░"
)

type ProgressBar struct {
	total       int64
	current     int64
	width       int
	description string
	startTime   time.Time
	lastUpdate  time.Time
	mutex       sync.Mutex
	finished    bool
}

// New constructs a new ProgressBar with a total byte count and a description.
// Parameters: total is the total size to track; description is a label shown with the bar.
// Returns a *ProgressBar initialized with default width and timestamps.
func New(total int64, description string) *ProgressBar {
	return &ProgressBar{
		total:       total,
		current:     0,
		width:       defaultBarWidth,
		description: description,
		startTime:   time.Now(),
		lastUpdate:  time.Now(),
	}
}

// Write implements io.Writer for ProgressBar by adding the number of bytes written to progress.
// Parameter p is the byte slice written. Returns the length of p and a nil error.
func (pb *ProgressBar) Write(p []byte) (n int, err error) {
	n = len(p)
	pb.Add(int64(n))
	return
}

// Add increases the current progress by n bytes and throttles rendering for performance.
// Parameter n is the increment amount. No return value.
func (pb *ProgressBar) Add(n int64) {
	pb.mutex.Lock()
	defer pb.mutex.Unlock()

	pb.current += n
	if pb.current > pb.total {
		pb.current = pb.total
	}

	now := time.Now()
	if now.Sub(pb.lastUpdate) > 100*time.Millisecond || pb.current == pb.total {
		pb.render()
		pb.lastUpdate = now
	}
}

// Set updates the current progress to a specific value and triggers a render.
// Parameter current is the new progress position. No return value.
func (pb *ProgressBar) Set(current int64) {
	pb.mutex.Lock()
	defer pb.mutex.Unlock()

	pb.current = current
	if pb.current > pb.total {
		pb.current = pb.total
	}
	pb.render()
}

// Finish marks the progress as complete, renders the final state, and prints a newline.
// No parameters. No return value.
func (pb *ProgressBar) Finish() {
	pb.mutex.Lock()
	defer pb.mutex.Unlock()

	if pb.finished {
		return
	}

	pb.current = pb.total
	pb.finished = true
	pb.render()
	fmt.Println()
}

// render draws the progress bar with percentage, speed, and ETA.
// Internal helper; respects total <= 0 and throttling logic from Add/Set. No return value.
func (pb *ProgressBar) render() {
	if pb.total <= 0 {
		return
	}

	percentage := float64(pb.current) / float64(pb.total) * 100
	filledWidth := int(float64(pb.width) * float64(pb.current) / float64(pb.total))

	// Optimized string building using Builder with pre-allocated capacity
	var bar strings.Builder
	bar.Grow(pb.width * 3) // Pre-allocate for UTF-8 characters

	// Use more efficient string building
	for i := 0; i < filledWidth; i++ {
		bar.WriteString(fillChar)
	}

	for i := filledWidth; i < pb.width; i++ {
		bar.WriteString(emptyChar)
	}

	elapsed := time.Since(pb.startTime)
	var speedStr, etaStr string

	if elapsed.Seconds() > 1 {
		speed := float64(pb.current) / elapsed.Seconds()
		speedStr = util.FormatBytes(int64(speed)) + "/s"

		if speed > 0 && pb.current < pb.total {
			remaining := pb.total - pb.current
			eta := time.Duration(float64(remaining)/speed) * time.Second
			etaStr = util.FormatDuration(eta)
		}
	}

	currentStr := util.FormatBytes(pb.current)
	totalStr := util.FormatBytes(pb.total)

	// Build status string more efficiently
	var status strings.Builder
	status.Grow(120) // Pre-allocate typical status line length

	status.WriteString("\r")
	status.WriteString(pb.description)
	status.WriteString(" [")
	status.WriteString(bar.String())
	status.WriteString(fmt.Sprintf("] %.1f%% (%s/%s)", percentage, currentStr, totalStr))

	if speedStr != "" {
		status.WriteString(" ")
		status.WriteString(speedStr)
	}

	if etaStr != "" {
		status.WriteString(" ETA: ")
		status.WriteString(etaStr)
	}

	// Pad to 80 characters for consistent terminal display
	statusStr := status.String()
	if len(statusStr) < 80 {
		statusStr += strings.Repeat(" ", 80-len(statusStr))
	}

	fmt.Print(statusStr)
}

type MultiProgress struct {
	bars   []*ProgressBar
	mutex  sync.Mutex
	active bool
}

// NewMultiProgress creates a MultiProgress manager to track multiple progress bars.
// No parameters. Returns a *MultiProgress with active set to true.
func NewMultiProgress() *MultiProgress {
	return &MultiProgress{
		active: true,
	}
}

// AddBar adds a new ProgressBar to the MultiProgress manager.
// Parameters: total is the bar's total size; description labels the bar. Returns the created *ProgressBar.
func (mp *MultiProgress) AddBar(total int64, description string) *ProgressBar {
	mp.mutex.Lock()
	defer mp.mutex.Unlock()

	bar := &ProgressBar{
		total:       total,
		current:     0,
		width:       40,
		description: description,
		startTime:   time.Now(),
		lastUpdate:  time.Now(),
	}

	mp.bars = append(mp.bars, bar)

	return bar
}

// Stop stops the MultiProgress manager and finishes all tracked progress bars.
// No parameters. No return value.
func (mp *MultiProgress) Stop() {
	mp.mutex.Lock()
	defer mp.mutex.Unlock()

	mp.active = false
	for _, bar := range mp.bars {
		bar.Finish()
	}
}
