package main

import (
	"os"

	"github.com/cyberark/conjur-api-go/conjurapi"
)

func NewAutoProvider() (Provider, error) {
	err := setConjurCredentialsEnv()
	if err != nil {
		return nil, err
	}
	defer unsetEnv()

	if os.Getenv("CONJUR_BUILDPACK_AUTHN_TYPE") == certAuthnType {
		return newCertProviderFromEnv()
	}
	return newAPIProviderFromEnv()
}

func applyTelemetry(config *conjurapi.Config) {
	config.SetIntegrationName("CloudFoundry Conjur Buildpack")
	config.SetIntegrationVersion(Version)
	config.SetIntegrationType("cybr-secretsmanager")
	config.SetVendorName("CyberArk")
}
