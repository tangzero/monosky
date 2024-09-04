package monosky

import (
	"context"
	"encoding/json"

	"github.com/bluesky-social/indigo/api/atproto"
	"github.com/bluesky-social/indigo/xrpc"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/keybase/go-keychain"
)

var DefaultClient = &Client{xrpc: &xrpc.Client{Host: "https://bsky.social"}}

type Client struct {
	xrpc *xrpc.Client
}

// LoadAuthInfo loads the authentication info from the keychain
// TODO: add error handling
func (client *Client) LoadAuthInfo() tea.Cmd {
	queryItem := keychain.NewItem()
	queryItem.SetSecClass(keychain.SecClassGenericPassword)
	queryItem.SetService("Monosky")
	queryItem.SetMatchLimit(keychain.MatchLimitOne)
	queryItem.SetReturnAttributes(true)
	queryItem.SetReturnData(true)

	results, err := keychain.QueryItem(queryItem)
	if err != nil || len(results) == 0 {
		return AskToLogin
	}

	var auth xrpc.AuthInfo
	if err := json.Unmarshal(results[0].Data, &auth); err != nil {
		return AskToLogin
	}

	client.xrpc.Auth = &auth
	client.xrpc.Auth.AccessJwt = auth.RefreshJwt // TODO: send a PR to fix this issue when refreshing the session
	output, err := atproto.ServerRefreshSession(context.Background(), client.xrpc)
	if err != nil {
		client.xrpc.Auth = nil
		return AskToLogin
	}

	client.SaveAuthInfo(&xrpc.AuthInfo{
		AccessJwt:  output.AccessJwt,
		RefreshJwt: output.RefreshJwt,
		Handle:     output.Handle,
		Did:        output.Did,
	})
	return nil
}

// SaveAuthInfo saves the authentication info to the keychain
// TODO: add error handling
func (client *Client) SaveAuthInfo(authInfo *xrpc.AuthInfo) {
	client.xrpc.Auth = authInfo
	data, _ := json.Marshal(authInfo)

	authItem := keychain.NewItem()
	authItem.SetSecClass(keychain.SecClassGenericPassword)
	authItem.SetService("Monosky")
	authItem.SetAccount(client.xrpc.Auth.Handle)
	authItem.SetAuthenticationType(keychain.AuthenticationTypeKey)
	authItem.SetData(data)
	authItem.SetAccessible(keychain.AccessibleWhenUnlocked)

	keychain.DeleteItem(authItem)
	keychain.AddItem(authItem)
}

// LoggedIn returns true if we're able to load the authentication info
func (client *Client) LoggedIn() bool {
	return client.xrpc.Auth != nil
}
