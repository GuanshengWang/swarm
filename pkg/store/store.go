package store

import (
	"crypto/tls"
	"errors"
	"time"

	log "github.com/Sirupsen/logrus"
)

// Backend represents a KV Store Backend
type Backend string

const (
	// CONSUL backend
	CONSUL Backend = "consul"
	// ETCD backend
	ETCD = "etcd"
	// ZK backend
	ZK = "zk"
)

var (
	// ErrNotSupported is exported
	ErrNotSupported = errors.New("Backend storage not supported yet, please choose another one")
	// ErrNotImplemented is exported
	ErrNotImplemented = errors.New("Call not implemented in current backend")
	// ErrNotReachable is exported
	ErrNotReachable = errors.New("Api not reachable")
	// ErrCannotLock is exported
	ErrCannotLock = errors.New("Error acquiring the lock")
	// ErrWatchDoesNotExist is exported
	ErrWatchDoesNotExist = errors.New("No watch found for specified key")
	// ErrKeyModified is exported
	ErrKeyModified = errors.New("Unable to complete atomic operation, key modified")
	// ErrKeyNotFound is exported
	ErrKeyNotFound = errors.New("Key not found in store")
)

// Config contains the options for a storage client
type Config struct {
	TLS     *tls.Config
	Timeout time.Duration
}

// Store represents the backend K/V storage
// Each store should support every call listed
// here. Or it couldn't be implemented as a K/V
// backend for libkv
type Store interface {
	// Put a value at the specified key
	Put(key string, value []byte) error

	// Get a value given its key
	Get(key string) (*KVPair, error)

	// Delete the value at the specified key
	Delete(key string) error

	// Verify if a Key exists in the store
	Exists(key string) (bool, error)

	// Watch changes on a key.
	// Returns a channel that will receive changes or an error.
	// Upon creating a watch, the current value will be sent to the channel.
	// Providing a non-nil stopCh can be used to stop watching.
	Watch(key string, stopCh <-chan struct{}) (<-chan *KVPair, error)

	// WatchTree watches changes on a "directory"
	// Returns a channel that will receive changes or an error.
	// Upon creating a watch, the current value will be sent to the channel.
	// Providing a non-nil stopCh can be used to stop watching.
	WatchTree(prefix string, stopCh <-chan struct{}) (<-chan []*KVPair, error)

	// CreateLock for a given key.
	// The returned Locker is not held and must be acquired with `.Lock`.
	// value is optional.
	CreateLock(key string, value []byte) (Locker, error)

	// List the content of a given prefix
	List(prefix string) ([]*KVPair, error)

	// DeleteTree deletes a range of keys based on prefix
	DeleteTree(prefix string) error

	// Atomic operation on a single value
	AtomicPut(key string, value []byte, previous *KVPair) (bool, error)

	// Atomic delete of a single value
	AtomicDelete(key string, previous *KVPair) (bool, error)
}

// KVPair represents {Key, Value, Lastindex} tuple
type KVPair struct {
	Key       string
	Value     []byte
	LastIndex uint64
}

// WatchCallback is used for watch methods on keys
// and is triggered on key change
type WatchCallback func(entries ...*KVPair)

// Locker provides locking mechanism on top of the store.
// Similar to `sync.Lock` except it may return errors.
type Locker interface {
	Lock() (<-chan struct{}, error)
	Unlock() error
}

// Initialize creates a new Store object, initializing the client
type Initialize func(addrs []string, options *Config) (Store, error)

var (
	// Backend initializers
	initializers = map[Backend]Initialize{
		CONSUL: InitializeConsul,
		ETCD:   InitializeEtcd,
		ZK:     InitializeZookeeper,
	}
)

// CreateStore creates a an instance of store
func CreateStore(backend Backend, addrs []string, options *Config) (Store, error) {
	if init, exists := initializers[backend]; exists {
		log.WithFields(log.Fields{"backend": backend}).Debug("Initializing store service")
		return init(addrs, options)
	}

	return nil, ErrNotSupported
}
