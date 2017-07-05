package oauth2Provider

import (
	"strings"
	"testing"
	"time"
)

type MockKeyValueStore struct {
	code map[string][]byte
}

func (t *MockKeyValueStore) Set(key string, value string, d time.Duration) error {
	return nil
}

func (t *MockKeyValueStore) Get(key string) (string, error) {
	return "7-8-9", nil
}

func (t *MockKeyValueStore) Del(key string) (string, error) {
	return "7-8-9", nil
}

func (t *MockKeyValueStore) Len() int {
	return 1
}

func (t *MockKeyValueStore) String() string {
	return "Calling mock kvs log method..."
}

func NewMockKeyValueStore() *MockKeyValueStore {
	return &MockKeyValueStore{
		code: make(map[string][]byte),
	}
}

func TestDefaultKeyValueStore(t *testing.T) {

	//As some previous test may have initialed it, we must take care of the current kvs size...
	cuurentKvsLength := getKeyValueStore().Len()

	key := "1-2-3"
	val := "4-5-6"

	//Set the value
	getKeyValueStore().Set(key, val, 1*time.Second)

	//Check value is well created
	if got, _ := getKeyValueStore().Get(key); val != got {
		t.Errorf("Did not got expected value on get method : [%v] while expecting [%v]", got, val)
	}

	//Check length growth 1
	if len := getKeyValueStore().Len(); len != cuurentKvsLength+1 {
		t.Errorf("Did not got expected length value : got [%v] while expecting [%v]", len, cuurentKvsLength+1)
	}

	//Check Stringer method
	if log := getKeyValueStore().String(); !strings.Contains(log, "1-2-3 --> [4 5 6]\n") {
		t.Errorf("Did not get expected kvs.String() result, got %v", log)
	}

	//Checking data expiration
	time.Sleep(1 * time.Second)
	if got, _ := getKeyValueStore().Get(key); got != "" {
		t.Errorf("Was expecting a nil value due to expiration on get method : [%v] while expecting nil", got)
	}

	//Check Del method
	if got, _ := getKeyValueStore().Del(key); val != got {
		t.Errorf("Did not got expected value on delete method : got [%v] while expecting [%v]", got, val)
	}

	//Check new size value after deleting entry
	if len := getKeyValueStore().Len(); len != cuurentKvsLength {
		t.Errorf("Did not got expected length value : got [%v] while expecting [%v]", len, cuurentKvsLength)
	}
}

func TestCustomKeyValueStore(t *testing.T) {

	SetCustomKeyValueStore(NewMockKeyValueStore())

	key := "1-2-3"
	val := "4-5-6"

	getKeyValueStore().Set(key, val, 30*time.Second)

	if got, _ := getKeyValueStore().Get(key); "4-5-6" != got {
		t.Errorf("Did not got expected value, got [%v] while expecting [%v]", got, val)
	}
}
