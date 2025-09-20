package progress

import (
	"fmt"
	"strings"
	"sync"
	"time"

	util "github.com/sijunda/govman/internal/util"
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

func New(total int64, description string) *ProgressBar {
	return &ProgressBar{
		total:       total,
		current:     0,
		width:       50,
		description: description,
		startTime:   time.Now(),
		lastUpdate:  time.Now(),
	}
}

func (pb *ProgressBar) Write(p []byte) (n int, err error) {
	n = len(p)
	pb.Add(int64(n))
	return
}

func (pb *ProgressBar) Add(n int64) {
	pb.mutex.Lock()
	defer pb.mutex.Unlock()

	pb.current += n
	if pb.current > pb.total {
		pb.current = pb.total
	}

	// Update display every 100ms to avoid too frequent updates
	now := time.Now()
	if now.Sub(pb.lastUpdate) > 100*time.Millisecond || pb.current == pb.total {
		pb.render()
		pb.lastUpdate = now
	}
}

func (pb *ProgressBar) Set(current int64) {
	pb.mutex.Lock()
	defer pb.mutex.Unlock()

	pb.current = current
	if pb.current > pb.total {
		pb.current = pb.total
	}
	pb.render()
}

func (pb *ProgressBar) Finish() {
	pb.mutex.Lock()
	defer pb.mutex.Unlock()

	if pb.finished {
		return
	}

	pb.current = pb.total
	pb.finished = true
	pb.render()
	fmt.Println() // New line after completion
}

func (pb *ProgressBar) render() {
	if pb.total <= 0 {
		return
	}

	percentage := float64(pb.current) / float64(pb.total) * 100
	filledWidth := int(float64(pb.width) * float64(pb.current) / float64(pb.total))

	// Create progress bar
	bar := strings.Repeat("█", filledWidth) + strings.Repeat("░", pb.width-filledWidth)

	// Calculate speed and ETA
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

	// Format sizes
	currentStr := util.FormatBytes(pb.current)
	totalStr := util.FormatBytes(pb.total)

	// Build status line
	status := fmt.Sprintf("\r%s [%s] %.1f%% (%s/%s)",
		pb.description, bar, percentage, currentStr, totalStr)

	if speedStr != "" {
		status += fmt.Sprintf(" %s", speedStr)
	}

	if etaStr != "" {
		status += fmt.Sprintf(" ETA: %s", etaStr)
	}

	// Pad to clear previous line
	status = fmt.Sprintf("%-80s", status)

	fmt.Print(status)
}

// MultiProgress handles multiple progress bars
type MultiProgress struct {
	bars   []*ProgressBar
	mutex  sync.Mutex
	active bool
}

func NewMultiProgress() *MultiProgress {
	return &MultiProgress{
		active: true,
	}
}

func (mp *MultiProgress) AddBar(total int64, description string) *ProgressBar {
	mp.mutex.Lock()
	defer mp.mutex.Unlock()

	bar := &ProgressBar{
		total:       total,
		current:     0,
		width:       40, // Smaller width for multiple bars
		description: description,
		startTime:   time.Now(),
		lastUpdate:  time.Now(),
	}

	mp.bars = append(mp.bars, bar)
	return bar
}

func (mp *MultiProgress) Stop() {
	mp.mutex.Lock()
	defer mp.mutex.Unlock()

	mp.active = false
	for _, bar := range mp.bars {
		bar.Finish()
	}
}
