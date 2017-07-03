package oauth2Provider

import (
	"bytes"
	"fmt"
	"sync"
	"time"
)

type KeyValueStore interface {
	Set(key string, value string, d time.Duration) error
	Get(key string) (string, error)
	Del(key string) (string, error)
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

func (t *DefaultKeyValueStore) Set(key string, value string, d time.Duration) error {
	expiration := time.Now().Add(d)
	t.code[key] = value
	t.codeExpiration[key] = expiration
	return nil
}

func (t *DefaultKeyValueStore) Get(key string) (string, error) {
	expired := t.codeExpiration[key]

	if time.Now().Before(expired) {
		return t.code[key], nil
	}

	return nil, nil
}

func (t *DefaultKeyValueStore) Del(key string) (string, error) {
	res := t.code[key]
	delete(t.code, key)
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
