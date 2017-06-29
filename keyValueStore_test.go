package oauth2Provider

import (
	"bytes"
	"fmt"
	"testing"
	"time"
)

type FakeKeyValueStore struct {
	code           map[string][]byte
	codeExpiration map[string]time.Time
}

func (t *FakeKeyValueStore) Set(key, value []byte, d time.Duration) error {
	expiration := time.Now().Add(d)
	sKey := string(key)
	t.code[sKey] = value
	t.codeExpiration[sKey] = expiration
	return nil
}

func (t *FakeKeyValueStore) Get(key []byte) ([]byte, error) {
	sKey := string(key)
	expired := t.codeExpiration[string(key)]

	if time.Now().Before(expired) {
		return t.code[sKey], nil
	}

	return nil, nil
}

func (t *FakeKeyValueStore) Del(key []byte) ([]byte, error) {
	sKey := string(key)
	res := t.code[sKey]
	delete(t.code, sKey)
	return res, nil
}

func (t *FakeKeyValueStore) Len() int {
	return len(t.code)
}

func (t *FakeKeyValueStore) Log() {
	fmt.Println("logging kvs values...")
	for k, v := range t.code {
		fmt.Printf("key/valye %v/%v\n", k, v)
	}
}

func NewFakeKeyValueStore() *FakeKeyValueStore {
	return &FakeKeyValueStore{
		code:           make(map[string][]byte),
		codeExpiration: make(map[string]time.Time),
	}
}

func TestKeyValueStore(t *testing.T) {
	setKeyValueStore(NewFakeKeyValueStore())

	key := []byte("1-2-3")
	val := []byte{4, 5, 6}

	kvs.Set(key, val, 30*time.Second)

	if got, _ := kvs.Get(key); !bytes.Equal(val, got) {
		t.Errorf("Did not got expected value, got [%v] while expecting [%v]", got, val)
	}
}
