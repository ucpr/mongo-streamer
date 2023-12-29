package persistent

import (
	"context"
	"sync"
	"time"

	"github.com/ucpr/mongo-streamer/pkg/log"
)

const (
	// DefaultBufferCap is the default capacity of the buffer
	defaultBufferCap = 50
)

type Storage interface {
	Write(s string) error
	Clear() error
	Read() (string, error)
	Close(ctx context.Context) error
}

type StorageBuffer interface {
	Watch(ctx context.Context)
	Set(s string) error
	Get() (string, error)
	Flush() error
	Clear() error
	Close(ctx context.Context) error
}

type Buffer struct {
	sync.Mutex

	// data is the slice of strings that the buffer holds
	data []string
	// cap is the maximum number of elements in the buffer
	cap int
	// interval is the interval at which the buffer is flushed
	interval time.Duration
	// storage is the storage that writes the buffer to the persistent storage
	storage Storage
}

// NewBuffer creates a new buffer with the given capacity
func NewBuffer(cap int, interval time.Duration, writer Storage) (*Buffer, error) {
	buf := &Buffer{
		data:     make([]string, 0, cap),
		cap:      cap,
		interval: interval,
		storage:  writer,
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
			return
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

// Get returns data from persistent storage
func (b *Buffer) Get() (string, error) {
	b.Lock()
	defer b.Unlock()

	s, err := b.storage.Read()
	if err != nil {
		return "", err
	}

	return s, nil
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
	if err := b.storage.Write(data); err != nil {
		return err
	}

	// initialize
	b.data = make([]string, 0, b.cap)

	return nil
}

// Clear clears the buffer and saved data
func (b *Buffer) Clear() error {
	b.Lock()
	defer b.Unlock()

	if err := b.storage.Clear(); err != nil {
		return err
	}

	b.data = make([]string, 0, b.cap)

	return nil
}

// Close closes the Buffer.
// flush the data in the buffer before close
func (b *Buffer) Close(ctx context.Context) error {
	defer func() {
		if err := b.storage.Close(ctx); err != nil {
			log.Error("failed to close file", err)
		}
	}()
	return b.Flush()
}
