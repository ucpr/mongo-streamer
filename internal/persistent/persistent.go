package persistent

import (
	"context"
	"os"
	"sync"
	"time"

	"github.com/ucpr/mongo-streamer/internal/log"
)

const (
	// DefaultBufferCap is the default capacity of the buffer
	defaultBufferCap = 50
)

type Buffer struct {
	sync.Mutex
	// data is the slice of strings that the buffer holds
	data []string
	// cap is the maximum number of elements in the buffer
	cap int
	// interval is the interval at which the buffer is flushed
	interval time.Duration
	// file is the file to which the buffer is flushed
	file *os.File
}

// NewBuffer creates a new buffer with the given capacity
func NewBuffer(cap int, interval time.Duration, filepath string) (*Buffer, error) {
	// TODO: check exists dir
	f, err := os.OpenFile(filepath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0o644)
	if err != nil {
		return nil, err
	}

	buf := &Buffer{
		data:     make([]string, 0, cap),
		cap:      cap,
		interval: interval,
		file:     f,
	}
	if cap < 0 {
		buf.cap = defaultBufferCap
	}
	return buf, nil
}

// Watch starts a loop that flushes the buffer at the given interval
func (b *Buffer) Watch(ctx context.Context) {
	ticker := time.NewTicker(b.interval)
	for {
		select {
		case <-ticker.C:
			if err := b.Flush(); err != nil {
				log.Info("failed to flush buffer", err)
			}
		case <-ctx.Done():
			log.Info("buffer watcher stopped")
			break
		}
	}
}

// Set adds a string to the buffer
func (b *Buffer) Set(s string) error {
	b.Lock()
	defer b.Unlock()

	b.data = append(b.data, s)

	// flush when capacity is reached
	if len(b.data) >= b.cap {
		b.Flush()
	}
	return nil
}

// Get returns the last element of the buffer
func (b *Buffer) Get() string {
	b.Lock()
	defer b.Unlock()

	return b.data[len(b.data)-1]
}

// Flush writes the buffer to the persistent storage.
// Flush method is groutine safe and can be called concurrently.
func (b *Buffer) Flush() error {
	b.Lock()
	defer b.Unlock()

	// do nothing if buffer is empty
	if len(b.data) < 1 {
		return nil
	}

	// get the last element
	data := b.data[len(b.data)-1]

	// save bufferd data to file
	if _, err := b.file.WriteString(data); err != nil {
		return err
	}

	// initialize
	b.data = make([]string, 0, b.cap)

	return nil
}

// Close closes the Buffer
func (b *Buffer) Close() error {
	defer func() {
		if err := b.file.Close(); err != nil {
			log.Error("failed to close file", err)
		}
	}()
	return b.Flush()
}
