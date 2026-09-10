package main

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strings"

	"github.com/cyberark/conjur-api-go/conjurapi"
)

const certAuthnType = "cert"

// conjurCloudRegexp matches Secrets Manager SaaS appliance URLs.
// This mirrors the pattern used internally by the conjur-api-go library.
var conjurCloudRegexp = regexp.MustCompile(
	`(\.secretsmgr|-secretsmanager)` +
		`(\.cyberark\.cloud|` +
		`\.integration-cyberark\.cloud|` +
		`\.test-cyberark\.cloud|` +
		`\.dev-cyberark\.cloud|` +
		`\.sandbox-cyberark\.cloud|` +
		`\.pt-cyberark\.cloud)`,
)

// NewCertProvider returns a Conjur API client that authenticates using an
// X.509 client certificate via the authn-cert authenticator. Credentials are
// read from VCAP_SERVICES (via setConjurCredentialsEnv).
//
// Required environment variables (populated from VCAP_SERVICES credentials):
//   - CONJUR_APPLIANCE_URL      – Conjur / Edge base URL
//   - CONJUR_ACCOUNT            – Conjur account (typically "conjur")
//   - CONJUR_AUTHN_SERVICE_ID   – certificate authenticator name, e.g. "acme-vm"
//   - CONJUR_AUTHN_CERT         – client certificate in PEM format
//   - CONJUR_AUTHN_CERT_KEY     – client private key in PEM format
//
// Optional:
//   - CONJUR_AUTHN_LOGIN        – workload identity path, e.g. "host/data/vms/vm-01".
//     Required for request host-mode (default). Omit for SPIFFE host-mode, where
//     the server derives the workload from the certificate's SPIFFE ID.
//   - CONJUR_SSL_CERTIFICATE    – CA certificate bundle (PEM) for server verification
func NewCertProvider() (Provider, error) {
	err := setConjurCredentialsEnv()
	if err != nil {
		return nil, err
	}
	defer unsetEnv()
	return newCertProviderFromEnv()
}

func newCertProviderFromEnv() (Provider, error) {
	serviceID := os.Getenv("CONJUR_AUTHN_SERVICE_ID")
	certPEM := os.Getenv("CONJUR_AUTHN_CERT")
	certKeyPEM := os.Getenv("CONJUR_AUTHN_CERT_KEY")

	if serviceID == "" {
		return nil, fmt.Errorf("CONJUR_AUTHN_SERVICE_ID is required for certificate authentication")
	}
	if certPEM == "" {
		return nil, fmt.Errorf("CONJUR_AUTHN_CERT is required for certificate authentication")
	}
	if certKeyPEM == "" {
		return nil, fmt.Errorf("CONJUR_AUTHN_CERT_KEY is required for certificate authentication")
	}

	applianceURL := os.Getenv("CONJUR_APPLIANCE_URL")
	account := os.Getenv("CONJUR_ACCOUNT")
	login := os.Getenv("CONJUR_AUTHN_LOGIN")
	sslCert := os.Getenv("CONJUR_SSL_CERTIFICATE")

	token, err := authenticateWithCert(
		applianceURL, account, serviceID, login, certPEM, certKeyPEM, sslCert,
	)
	if err != nil {
		return nil, fmt.Errorf("certificate authentication failed: %v", err)
	}

	config, err := conjurapi.LoadConfig()
	if err != nil {
		return nil, err
	}
	applyTelemetry(&config)

	return conjurapi.NewClientFromToken(config, string(token))
}

func authenticateWithCert(
	applianceURL, account, serviceID, login, certPEM, certKeyPEM, caCertPEM string,
) ([]byte, error) {
	clientCert, err := tls.X509KeyPair([]byte(certPEM), []byte(certKeyPEM))
	if err != nil {
		return nil, fmt.Errorf("failed to load client certificate and key: %v", err)
	}

	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{clientCert},
	}

	if caCertPEM != "" {
		caPool := x509.NewCertPool()
		if !caPool.AppendCertsFromPEM([]byte(caCertPEM)) {
			return nil, fmt.Errorf("failed to parse CA certificate")
		}
		tlsConfig.RootCAs = caPool
	}
	normalizedBase := normalizeApplianceURL(applianceURL)
	var authnURL string
	if login != "" {
		encodedLogin := url.PathEscape(login)
		authnURL = fmt.Sprintf(
			"%s/authn-cert/%s/%s/%s/authenticate",
			normalizedBase, serviceID, account, encodedLogin,
		)
	} else {
		// SPIFFE mode: the server derives the workload identity from the
		// client certificate's SPIFFE ID; no login segment in the URL.
		authnURL = fmt.Sprintf(
			"%s/authn-cert/%s/%s/authenticate",
			normalizedBase, serviceID, account,
		)
	}

	transport := &http.Transport{TLSClientConfig: tlsConfig}
	httpClient := &http.Client{Transport: transport}

	req, err := http.NewRequest(http.MethodPost, authnURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create authentication request: %v", err)
	}
	req.Header.Set("Accept-Encoding", "base64")

	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("authentication request failed: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read authentication response: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("authentication failed with status %d: %s",
			resp.StatusCode, string(body))
	}

	tokenJSON, err := base64.StdEncoding.DecodeString(strings.TrimSpace(string(body)))
	if err != nil {
		return nil, fmt.Errorf("failed to decode authentication token: %v", err)
	}

	return tokenJSON, nil
}

// normalizeApplianceURL ensures the appliance URL has the correct base path.
// For Secrets Manager SaaS deployments (*.cyberark.cloud), it appends "/api"
// when the URL does not already include it. This mirrors the behaviour of the
// conjur-api-go library's internal URL normalization (normalizeBaseURL).
func normalizeApplianceURL(applianceURL string) string {
	base := strings.TrimSuffix(applianceURL, "/")
	if conjurCloudRegexp.MatchString(base) && !strings.Contains(base, "/api") {
		return base + "/api"
	}
	return base
}
