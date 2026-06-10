package auth

import (
	"golang.org/x/oauth2"
)

// staticTokenSuffix is appended to the instance name when storing a static
// API token. This namespaces static tokens away from OAuth tokens so that
// loading an OAuth token for an instance does not accidentally return a
// static-token entry (which has no expiry or refresh information).
const staticTokenSuffix = ":static"

// staticInstanceKey returns the storage key used for the static API token of
// the named instance.
func staticInstanceKey(instanceName string) string {
	return instanceName + staticTokenSuffix
}

// SaveStaticToken stores a plain API token (not an OAuth token) in the
// provided TokenStore. It is wrapped in an oauth2.Token so the existing
// encrypted-file / keyring storage can be reused without modification.
func SaveStaticToken(store TokenStore, instanceName, apiToken string) error {
	t := &oauth2.Token{
		AccessToken: apiToken,
		TokenType:   "Bearer",
		// Deliberately zero Expiry: static tokens do not expire through OAuth.
	}
	return store.Save(staticInstanceKey(instanceName), t)
}

// LoadStaticToken retrieves a plain API token from the TokenStore.
// It returns ("", nil) when no static token has been stored for the instance.
func LoadStaticToken(store TokenStore, instanceName string) (string, error) {
	key := staticInstanceKey(instanceName)
	if !store.Exists(key) {
		return "", nil
	}
	t, err := store.Load(key)
	if err != nil {
		return "", err
	}
	return t.AccessToken, nil
}

// DeleteStaticToken removes the static API token for an instance from the
// TokenStore. It is a no-op if no token was stored.
func DeleteStaticToken(store TokenStore, instanceName string) error {
	return store.Delete(staticInstanceKey(instanceName))
}

// StaticTokenExists reports whether a static API token is stored for the
// given instance.
func StaticTokenExists(store TokenStore, instanceName string) bool {
	return store.Exists(staticInstanceKey(instanceName))
}
