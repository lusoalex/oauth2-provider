package oauth2Provider

import (
	"bytes"
	"fmt"
	"sync"
	"time"
)

type KeyValueStore interface {
	Set(key, value []byte, d time.Duration) error
	Get(key []byte) ([]byte, error)
	Del(key []byte) ([]byte, error)
	Len() int
	String() string
}

var (
	kvsSetOnce sync.Once
	kvs        KeyValueStore = NewDefaultKeyValueStore()
)

func SetCustomKeyValueStore(newKvs KeyValueStore) {
	kvsSetOnce.Do(func() {
		kvs = newKvs
	})
}

func getKeyValueStore() KeyValueStore {
	return kvs
}

/***********************************************************************/
/*                                                                     */
/*                   PROVIDE DEFAULT KEY VALUE STORE                   */
/*                  SHOULD NEVER BE USED IN PRODUCTION                 */
/*             PROVIDE ONE CALLING SetCustomKeyValueStore              */
/*                                                                     */
/***********************************************************************/
type DefaultKeyValueStore struct {
	code           map[string][]byte
	codeExpiration map[string]time.Time
}

func (t *DefaultKeyValueStore) Set(key, value []byte, d time.Duration) error {
	expiration := time.Now().Add(d)
	sKey := string(key)
	t.code[sKey] = value
	t.codeExpiration[sKey] = expiration
	return nil
}

func (t *DefaultKeyValueStore) Get(key []byte) ([]byte, error) {
	sKey := string(key)
	expired := t.codeExpiration[string(key)]

	if time.Now().Before(expired) {
		return t.code[sKey], nil
	}

	return nil, nil
}

func (t *DefaultKeyValueStore) Del(key []byte) ([]byte, error) {
	sKey := string(key)
	res := t.code[sKey]
	delete(t.code, sKey)
	return res, nil
}

func (t *DefaultKeyValueStore) Len() int {
	return len(t.code)
}

func (t *DefaultKeyValueStore) String() string {
	var buffer bytes.Buffer
	for k, v := range t.code {
		buffer.WriteString(fmt.Sprintf("%v --> %v\n", k, v))
	}
	return buffer.String()
}

func NewDefaultKeyValueStore() *DefaultKeyValueStore {
	return &DefaultKeyValueStore{
		code:           make(map[string][]byte),
		codeExpiration: make(map[string]time.Time),
	}
}
