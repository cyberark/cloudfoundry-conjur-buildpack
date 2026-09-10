package main

import (
	"github.com/cyberark/conjur-api-go/conjurapi"
)

// NewAPIProvider returns a Conjur API client that authenticates using an API
// key. Credentials are read from VCAP_SERVICES via setConjurCredentialsEnv.
//
// Required environment variables (populated from VCAP_SERVICES credentials):
//   - CONJUR_APPLIANCE_URL   – Conjur server URL
//   - CONJUR_ACCOUNT         – Conjur account name
//   - CONJUR_AUTHN_LOGIN     – login / role ID
//   - CONJUR_AUTHN_API_KEY   – API key for the role
//
// Optional:
//   - CONJUR_SSL_CERTIFICATE – CA certificate bundle (PEM) for TLS verification
func NewAPIProvider() (Provider, error) {
	err := setConjurCredentialsEnv()
	if err != nil {
		return nil, err
	}
	defer unsetEnv()
	return newAPIProviderFromEnv()
}

func newAPIProviderFromEnv() (Provider, error) {
	config, err := conjurapi.LoadConfig()
	if err != nil {
		return nil, err
	}
	applyTelemetry(&config)
	return conjurapi.NewClientFromEnvironment(config)
}
