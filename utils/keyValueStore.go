package utils

import (
	"oauth2-provider/models"
	"sync"
	"time"

	"github.com/google/uuid"
)

type KeyValueStore interface {
	Code(ar *models.AuthorizationRequest) string
	Revoke(code string) (*models.AuthorizationRequest, bool)
}

//var KVS *DefaultKeyValueStore

/*var (
	kvsSetOnce sync.Once
	kvs        KeyValueStore = NewDefaultKeyValueStore()
)*/

/*func SetCustomKeyValueStore(newKvs KeyValueStore) {
	kvsSetOnce.Do(func() {
		kvs = newKvs
	})
}

func getKeyValueStore() KeyValueStore {
	return kvs
}*/

/***********************************************************************/
/*                                                                     */
/*                   PROVIDE DEFAULT KEY VALUE STORE                   */
/*                  SHOULD NEVER BE USED IN PRODUCTION                 */
/*             PROVIDE ONE CALLING SetCustomKeyValueStore              */
/*                                                                     */
/***********************************************************************/

type AuthorizationRequestWithTimer struct {
	models.AuthorizationRequest
	timer *time.Timer
}

type DefaultKeyValueStore struct {
	codes map[string]*AuthorizationRequestWithTimer
	sync.Mutex
}

func (kvs *DefaultKeyValueStore) Code(ar *models.AuthorizationRequest) string {
	code := uuid.New().String()
	tokenTimer := time.AfterFunc(1*time.Minute, func() {
		kvs.Lock()
		defer kvs.Unlock()
		delete(kvs.codes, code)
	})

	kvs.Lock()
	defer kvs.Unlock()
	kvs.codes[code] = &AuthorizationRequestWithTimer{
		AuthorizationRequest: *ar,
		timer:                tokenTimer,
	}

	return code
}

func (kvs *DefaultKeyValueStore) Revoke(code string) (*models.AuthorizationRequest, bool) {
	kvs.Lock()
	defer kvs.Unlock()
	arwt, ok := kvs.codes[code]
	if ok {
		arwt.timer.Stop()
		delete(kvs.codes, code)
		return &arwt.AuthorizationRequest, ok
	}
	return nil, false
}

func NewDefaultKeyValueStore() *DefaultKeyValueStore {
	return &DefaultKeyValueStore{
		codes: make(map[string]*AuthorizationRequestWithTimer),
	}
}
