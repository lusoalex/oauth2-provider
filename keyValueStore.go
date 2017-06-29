package oauth2Provider

import "time"

type KeyValueStore interface {
	Set(key, value []byte, d time.Duration) error
	Get(key []byte) ([]byte, error)
	Del(key []byte) ([]byte, error)
	Len() int
	Log()
}
