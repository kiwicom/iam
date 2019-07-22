package storage

import (
	"time"

	"github.com/getsentry/raven-go"
	"github.com/pkg/errors"
)

type cache interface {
	Get(key string, value interface{}) error
	Set(key string, value interface{}, ttl time.Duration) error
	Del(key string) error
}

// LockManager manage locks using Cache, to prevent multiple expensive actions to
// be run at the same time.
type LockManager struct {
	cache      cache
	retryDelay time.Duration
	expiration time.Duration
}

// NewLockManager initializes and returns a LockManager for Redis
func NewLockManager(cache cache, retryDelay, expiration time.Duration) *LockManager {
	return &LockManager{
		cache:      cache,
		retryDelay: retryDelay,
		expiration: expiration,
	}
}

// ErrLockExists indicates that a lock was not created because one was already
// present. This error is returned after the old lock is deleted or expired.
var ErrLockExists = errors.New("lock was not created because one was already present")

// Create creates a lock to prevent having multiple instances of this service
// doing an expensive action at the same time. If a lock already exists,
// the function will not create one, it will wait until the existing one is
// deleted or expired before returning ErrLockExists.
func (l *LockManager) Create(name string) error {
	var lock time.Time
	var exists bool
	key := "lock:" + name

	// Check if a lock already exists.
	err := l.cache.Get(key, &lock)
	for err == nil {
		// If it does, wait for it to expire or be deleted.
		exists = true
		time.Sleep(l.retryDelay)
		err = l.cache.Get(key, &lock)
	}
	if exists {
		return ErrLockExists
	}
	if err != ErrNotFound {
		err = errors.Wrap(err, "error checking if a lock exists")
		raven.CaptureError(err, nil)
	}

	lock = time.Now()
	err = l.cache.Set(key, lock, l.expiration)
	if err != nil {
		err = errors.Wrap(err, "error creating lock")
		raven.CaptureError(err, nil)
	}
	return nil
}

// Delete removes a lock for the provided name.
func (l *LockManager) Delete(name string) {
	key := "lock:" + name
	err := l.cache.Del(key)
	if err != nil {
		raven.CaptureError(err, nil)
	}
}
