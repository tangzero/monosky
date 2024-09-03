package components

import (
	"encoding/json"

	"github.com/bluesky-social/indigo/xrpc"
	"github.com/keybase/go-keychain"
)

var DefaultClient = &Client{&xrpc.Client{Host: "https://bsky.social"}}

type Client struct {
	*xrpc.Client
}

// LoadAuthInfo loads the authentication info from the keychain
// TODO: add error handling
func (client *Client) LoadAuthInfo() {
	queryItem := keychain.NewItem()
	queryItem.SetSecClass(keychain.SecClassGenericPassword)
	queryItem.SetService("Monosky")
	queryItem.SetMatchLimit(keychain.MatchLimitOne)
	queryItem.SetReturnAttributes(true)
	queryItem.SetReturnData(true)

	results, err := keychain.QueryItem(queryItem)
	if err != nil || len(results) == 0 {
		return
	}

	var auth xrpc.AuthInfo
	if err := json.Unmarshal(results[0].Data, &auth); err != nil {
		return
	}

	client.Auth = &auth
}

// SaveAuthInfo saves the authentication info to the keychain
// TODO: add error handling
func (client *Client) SaveAuthInfo() {
	data, _ := json.Marshal(DefaultClient.Auth)

	authItem := keychain.NewItem()
	authItem.SetSecClass(keychain.SecClassGenericPassword)
	authItem.SetService("Monosky")
	authItem.SetAccount(client.Auth.Handle)
	authItem.SetAuthenticationType(keychain.AuthenticationTypeKey)
	authItem.SetData(data)
	authItem.SetAccessible(keychain.AccessibleWhenUnlocked)

	keychain.DeleteItem(authItem)
	keychain.AddItem(authItem)
}

// LoggedIn returns true if we're able to load the authentication info
func (client *Client) LoggedIn() bool {
	return client.Auth != nil
}
