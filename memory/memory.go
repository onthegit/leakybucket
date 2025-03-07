package memory

import (
	"sync"
	"time"

	"github.com/onthegit/leakybucket"
)

type bucket struct {
	capacity  uint
	remaining uint
	reset     time.Time
	rate      time.Duration
	mutex     *sync.Mutex
}

func (b *bucket) Capacity() uint {
	return b.capacity
}

// Remaining space in the bucket.
func (b *bucket) Remaining() uint {
	return b.remaining
}

// Reset returns when the bucket will be drained.
func (b *bucket) Reset() time.Time {
	return b.reset
}

// Add to the bucket.
func (b *bucket) Add(amount uint) (leakybucket.BucketState, error) {
	b.mutex.Lock()
	defer b.mutex.Unlock()
	if time.Now().After(b.reset) {
		b.reset = time.Now().Add(b.rate)
		b.remaining = b.capacity
	}
	if amount > b.remaining {
		return leakybucket.BucketState{Capacity: b.capacity, Remaining: b.remaining, Reset: b.reset}, leakybucket.ErrorFull
	}
	b.remaining -= amount
	return leakybucket.BucketState{Capacity: b.capacity, Remaining: b.remaining, Reset: b.reset}, nil
}

// Storage is a non thread-safe in-memory leaky bucket factory.
type Storage struct {
	buckets map[string]*bucket
	mu      *sync.Mutex
}

// New initializes the in-memory bucket store.
func New() *Storage {
	return &Storage{
		buckets: make(map[string]*bucket),
		mu:      new(sync.Mutex),
	}
}

// Create a bucket.
func (s *Storage) Create(name string, capacity uint, rate time.Duration) (leakybucket.Bucket, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	b, ok := s.buckets[name]
	if ok {
		return b, nil
	}
	b = &bucket{
		mutex:     new(sync.Mutex),
		capacity:  capacity,
		remaining: capacity,
		reset:     time.Now().Add(rate),
		rate:      rate,
	}
	s.buckets[name] = b
	return b, nil
}

func (s *Storage) Remove(name string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.buckets, name)
}
