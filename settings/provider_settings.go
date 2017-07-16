package settings

import (
	"oauth2-provider/client"
	"oauth2-provider/user"
	"oauth2-provider/utils"
)

type Oauth2ProviderSettings struct {
	utils.KeyValueStore
	client.ClientManager
	user.UserManager
}

func DefaultOauth2ProviderSettings() *Oauth2ProviderSettings {
	return &Oauth2ProviderSettings{
		KeyValueStore: utils.NewDefaultKeyValueStore(),
		ClientManager: &client.DefaultClientManager{},
		UserManager:   &user.DefaultUserManager{},
	}
}
