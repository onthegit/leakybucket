package leakybucket

import (
	"errors"
	"time"
)

var (
	// ErrorFull is returned when the amount requested to add exceeds the remaining space in the bucket.
	ErrorFull = errors.New("add exceeds free capacity")
)

// Bucket interface for interacting with leaky buckets: https://en.wikipedia.org/wiki/Leaky_bucket
type Bucket interface {
	// Capacity of the bucket.
	Capacity() uint

	// Remaining space in the bucket.
	Remaining() uint

	// Reset returns when the bucket will be drained.
	Reset() time.Time

	// Add to the bucket. MUST return bucket state after adding, even if an error was encountered
	Add(uint) (BucketState, error)
}

// BucketState is a snapshot of a bucket's properties.
type BucketState struct {
	Capacity  uint
	Remaining uint
	Reset     time.Time
}

// Storage interface for generating buckets keyed by a string.
type Storage interface {
	// Create a bucket with a name, capacity, and rate.
	// rate is how long it takes for full capacity to drain.
	Create(name string, capacity uint, rate time.Duration) (Bucket, error)

	//removes the bucket
	Remove(name string)
}
