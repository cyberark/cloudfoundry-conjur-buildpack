package main

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"math/big"
	"net"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/cyberark/conjur-api-go/conjurapi"
	"github.com/stretchr/testify/assert"
)

type testCertBundle struct {
	caCert        *x509.Certificate
	caCertPEM     string
	clientCertPEM string
	clientKeyPEM  string
	serverTLSCert tls.Certificate
}

func generateTestCertBundle(t *testing.T) testCertBundle {
	t.Helper()

	// ── CA ────────────────────────────────────────────────────────────────────
	caKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	assert.NoError(t, err)

	caTemplate := &x509.Certificate{
		SerialNumber:          big.NewInt(1),
		Subject:               pkix.Name{CommonName: "Test CA"},
		NotBefore:             time.Now().Add(-time.Hour),
		NotAfter:              time.Now().Add(24 * time.Hour),
		IsCA:                  true,
		BasicConstraintsValid: true,
		KeyUsage:              x509.KeyUsageCertSign | x509.KeyUsageCRLSign,
	}
	caDER, err := x509.CreateCertificate(rand.Reader, caTemplate, caTemplate, &caKey.PublicKey, caKey)
	assert.NoError(t, err)

	caCert, err := x509.ParseCertificate(caDER)
	assert.NoError(t, err)
	caCertPEM := string(pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: caDER}))

	clientKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	assert.NoError(t, err)

	clientTemplate := &x509.Certificate{
		SerialNumber: big.NewInt(2),
		Subject:      pkix.Name{CommonName: "test-workload"},
		NotBefore:    time.Now().Add(-time.Hour),
		NotAfter:     time.Now().Add(24 * time.Hour),
		ExtKeyUsage:  []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth},
	}
	clientDER, err := x509.CreateCertificate(rand.Reader, clientTemplate, caCert, &clientKey.PublicKey, caKey)
	assert.NoError(t, err)

	clientKeyBytes, err := x509.MarshalECPrivateKey(clientKey)
	assert.NoError(t, err)

	clientCertPEM := string(pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: clientDER}))
	clientKeyPEM := string(pem.EncodeToMemory(&pem.Block{Type: "EC PRIVATE KEY", Bytes: clientKeyBytes}))

	serverKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	assert.NoError(t, err)

	serverTemplate := &x509.Certificate{
		SerialNumber: big.NewInt(3),
		Subject:      pkix.Name{CommonName: "127.0.0.1"},
		NotBefore:    time.Now().Add(-time.Hour),
		NotAfter:     time.Now().Add(24 * time.Hour),
		IPAddresses:  []net.IP{net.ParseIP("127.0.0.1")},
		ExtKeyUsage:  []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
	}
	serverDER, err := x509.CreateCertificate(rand.Reader, serverTemplate, caCert, &serverKey.PublicKey, caKey)
	assert.NoError(t, err)

	serverKeyBytes, err := x509.MarshalECPrivateKey(serverKey)
	assert.NoError(t, err)

	serverTLSCert, err := tls.X509KeyPair(
		pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: serverDER}),
		pem.EncodeToMemory(&pem.Block{Type: "EC PRIVATE KEY", Bytes: serverKeyBytes}),
	)
	assert.NoError(t, err)

	return testCertBundle{
		caCert:        caCert,
		caCertPEM:     caCertPEM,
		clientCertPEM: clientCertPEM,
		clientKeyPEM:  clientKeyPEM,
		serverTLSCert: serverTLSCert,
	}
}

func startMTLSServer(t *testing.T, bundle testCertBundle, handler http.Handler) *httptest.Server {
	t.Helper()

	caPool := x509.NewCertPool()
	caPool.AddCert(bundle.caCert)

	srv := httptest.NewUnstartedServer(handler)
	srv.TLS = &tls.Config{
		Certificates: []tls.Certificate{bundle.serverTLSCert},
		ClientAuth:   tls.RequireAndVerifyClientCert,
		ClientCAs:    caPool,
	}
	srv.StartTLS()
	t.Cleanup(srv.Close)
	return srv
}

const sampleTokenJSON = `{"protected":"eyJhbGciOiJjb25qdXIub3JnL3Nsb3NpbG8vdjIiLCJraWQiOiJ0ZXN0In0=","payload":"eyJzdWIiOiJ0ZXN0IiwiaWF0IjoxNjAwMDAwMDAwfQ==","signature":"dGVzdA=="}`

func TestNormalizeApplianceURL(t *testing.T) {
	testCases := []struct {
		description string
		input       string
		expected    string
	}{
		{
			description: "SaaS URL without /api gets /api appended",
			input:       "https://tenant.secretsmgr.cyberark.cloud",
			expected:    "https://tenant.secretsmgr.cyberark.cloud/api",
		},
		{
			description: "SaaS URL that already contains /api is unchanged",
			input:       "https://tenant.secretsmgr.cyberark.cloud/api",
			expected:    "https://tenant.secretsmgr.cyberark.cloud/api",
		},
		{
			description: "SaaS integration URL gets /api appended",
			input:       "https://tenant.secretsmgr.integration-cyberark.cloud",
			expected:    "https://tenant.secretsmgr.integration-cyberark.cloud/api",
		},
		{
			description: "Non-SaaS URL is returned unchanged",
			input:       "https://conjur.example.com",
			expected:    "https://conjur.example.com",
		},
		{
			description: "Non-SaaS URL with trailing slash has slash removed",
			input:       "https://conjur.example.com/",
			expected:    "https://conjur.example.com",
		},
		{
			description: "Edge URL (non-SaaS) is returned unchanged",
			input:       "https://us-edge.acme.dev",
			expected:    "https://us-edge.acme.dev",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			result := normalizeApplianceURL(tc.input)
			assert.Equal(t, tc.expected, result)
		})
	}
}

func TestNewCertProviderFromEnv_MissingServiceID(t *testing.T) {
	os.Unsetenv("CONJUR_AUTHN_SERVICE_ID")
	os.Setenv("CONJUR_AUTHN_CERT", "cert-pem")
	os.Setenv("CONJUR_AUTHN_CERT_KEY", "key-pem")
	defer func() {
		os.Unsetenv("CONJUR_AUTHN_CERT")
		os.Unsetenv("CONJUR_AUTHN_CERT_KEY")
	}()

	_, err := newCertProviderFromEnv()
	assert.ErrorContains(t, err, "CONJUR_AUTHN_SERVICE_ID")
}

func TestNewCertProviderFromEnv_MissingCert(t *testing.T) {
	os.Setenv("CONJUR_AUTHN_SERVICE_ID", "acme-vm")
	os.Unsetenv("CONJUR_AUTHN_CERT")
	os.Setenv("CONJUR_AUTHN_CERT_KEY", "key-pem")
	defer func() {
		os.Unsetenv("CONJUR_AUTHN_SERVICE_ID")
		os.Unsetenv("CONJUR_AUTHN_CERT_KEY")
	}()

	_, err := newCertProviderFromEnv()
	assert.ErrorContains(t, err, "CONJUR_AUTHN_CERT")
}

func TestNewCertProviderFromEnv_MissingCertKey(t *testing.T) {
	os.Setenv("CONJUR_AUTHN_SERVICE_ID", "acme-vm")
	os.Setenv("CONJUR_AUTHN_CERT", "cert-pem")
	os.Unsetenv("CONJUR_AUTHN_CERT_KEY")
	defer func() {
		os.Unsetenv("CONJUR_AUTHN_SERVICE_ID")
		os.Unsetenv("CONJUR_AUTHN_CERT")
	}()

	_, err := newCertProviderFromEnv()
	assert.ErrorContains(t, err, "CONJUR_AUTHN_CERT_KEY")
}

func TestAuthenticateWithCert_InvalidClientCert(t *testing.T) {
	_, err := authenticateWithCert(
		"https://edge.example.com", "conjur", "acme-vm",
		"host/data/vm-01",
		"not-a-valid-cert", "not-a-valid-key",
		"",
	)
	assert.Error(t, err)
	assert.ErrorContains(t, err, "failed to load client certificate and key")
}

func TestAuthenticateWithCert_InvalidCACert(t *testing.T) {
	bundle := generateTestCertBundle(t)

	_, err := authenticateWithCert(
		"https://edge.example.com", "conjur", "acme-vm",
		"host/data/vm-01",
		bundle.clientCertPEM, bundle.clientKeyPEM,
		"not-valid-ca-pem",
	)
	assert.Error(t, err)
	assert.ErrorContains(t, err, "failed to parse CA certificate")
}

func TestAuthenticateWithCert_ServerReturnsUnauthorized(t *testing.T) {
	bundle := generateTestCertBundle(t)

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("401 UNAUTHORIZED"))
	})
	srv := startMTLSServer(t, bundle, handler)

	_, err := authenticateWithCert(
		srv.URL, "conjur", "acme-vm",
		"host/data/vm-workloads/vm-01",
		bundle.clientCertPEM, bundle.clientKeyPEM,
		bundle.caCertPEM,
	)
	assert.Error(t, err)
	assert.ErrorContains(t, err, "authentication failed with status 401")
}

func TestAuthenticateWithCert_ServerReturnsInvalidBase64(t *testing.T) {
	bundle := generateTestCertBundle(t)

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("!!!not-valid-base64!!!"))
	})
	srv := startMTLSServer(t, bundle, handler)

	_, err := authenticateWithCert(
		srv.URL, "conjur", "acme-vm",
		"host/data/vm-workloads/vm-01",
		bundle.clientCertPEM, bundle.clientKeyPEM,
		bundle.caCertPEM,
	)
	assert.Error(t, err)
	assert.ErrorContains(t, err, "failed to decode authentication token")
}

func TestAuthenticateWithCert_Success(t *testing.T) {
	bundle := generateTestCertBundle(t)

	encodedToken := base64.StdEncoding.EncodeToString([]byte(sampleTokenJSON))

	var capturedPath string
	var capturedHeader string
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		capturedPath = r.RequestURI
		capturedHeader = r.Header.Get("Accept-Encoding")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(encodedToken))
	})
	srv := startMTLSServer(t, bundle, handler)

	login := "host/data/vm-workloads/vm-01"
	token, err := authenticateWithCert(
		srv.URL, "conjur", "acme-vm",
		login,
		bundle.clientCertPEM, bundle.clientKeyPEM,
		bundle.caCertPEM,
	)

	assert.NoError(t, err)
	assert.Equal(t, sampleTokenJSON, string(token), "decoded token should match the original JSON")

	assert.Equal(t, "/authn-cert/acme-vm/conjur/host%2Fdata%2Fvm-workloads%2Fvm-01/authenticate", capturedPath)
	assert.Equal(t, "base64", capturedHeader)
}

func TestAuthenticateWithCert_LoginContainingSpecialChars(t *testing.T) {
	bundle := generateTestCertBundle(t)

	encodedToken := base64.StdEncoding.EncodeToString([]byte(sampleTokenJSON))

	var capturedPath string
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		capturedPath = r.RequestURI
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(encodedToken))
	})
	srv := startMTLSServer(t, bundle, handler)

	_, err := authenticateWithCert(
		srv.URL, "conjur", "acme-vm",
		"host/data/org/vm-01",
		bundle.clientCertPEM, bundle.clientKeyPEM,
		bundle.caCertPEM,
	)

	assert.NoError(t, err)
	// "/" in the login must be percent-encoded as "%2F"
	assert.Contains(t, capturedPath, "%2F")
	assert.NotContains(t, capturedPath, "host/data")
}

func TestAuthenticateWithCert_SPIFFEMode(t *testing.T) {
	bundle := generateTestCertBundle(t)

	encodedToken := base64.StdEncoding.EncodeToString([]byte(sampleTokenJSON))

	var capturedPath string
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		capturedPath = r.RequestURI
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(encodedToken))
	})
	srv := startMTLSServer(t, bundle, handler)

	_, err := authenticateWithCert(
		srv.URL, "conjur", "acme-vm",
		"", // empty login → SPIFFE mode
		bundle.clientCertPEM, bundle.clientKeyPEM,
		bundle.caCertPEM,
	)

	assert.NoError(t, err)
	// URL must NOT contain a workload-id segment:
	// expected: /authn-cert/acme-vm/conjur/authenticate
	// must NOT be: /authn-cert/acme-vm/conjur//authenticate
	assert.Equal(t, "/authn-cert/acme-vm/conjur/authenticate", capturedPath)
}

func TestNewAutoProvider_SelectsAPIKeyByDefault(t *testing.T) {
	vcapServices := `{"cyberark-conjur":[{"credentials":{"appliance_url":"https://conjur.test","account":"conjur","authn_login":"host/app","authn_api_key":""}}]}`
	os.Setenv("VCAP_SERVICES", vcapServices)
	defer os.Unsetenv("VCAP_SERVICES")

	_, err := NewAutoProvider()
	if err != nil {
		assert.NotContains(t, err.Error(), "CONJUR_AUTHN_SERVICE_ID")
		assert.NotContains(t, err.Error(), "CONJUR_AUTHN_CERT")
	}
}

func TestNewAutoProvider_SelectsCertProviderWhenConfigured(t *testing.T) {
	vcapServices := `{"cyberark-conjur":[{"credentials":{"appliance_url":"https://conjur.test","account":"conjur","authn_login":"host/app","authn_type":"cert"}}]}`
	os.Setenv("VCAP_SERVICES", vcapServices)
	defer os.Unsetenv("VCAP_SERVICES")

	_, err := NewAutoProvider()
	assert.Error(t, err)
	assert.ErrorContains(t, err, "CONJUR_AUTHN_SERVICE_ID")
}

func TestApplyTelemetry_SetsCorrectFields(t *testing.T) {
	config := conjurapi.Config{}
	applyTelemetry(&config)

	assert.Equal(t, "CloudFoundry Conjur Buildpack", config.IntegrationName)
	assert.Equal(t, "cybr-secretsmanager", config.IntegrationType)
	assert.Equal(t, "CyberArk", config.VendorName)
	assert.Equal(t, Version, config.IntegrationVersion)
}

func TestVersion_DefaultIsDev(t *testing.T) {
	assert.NotEmpty(t, Version)
}

func TestConjurCredentials_CertFieldsUnmarshalled(t *testing.T) {
	vcapJSON := `{
		"cyberark-conjur": [{
			"credentials": {
				"appliance_url":    "https://edge.example.com",
				"account":          "conjur",
				"authn_login":      "host/data/vms/vm-01",
				"ssl_certificate":  "ca-cert-pem",
				"authn_type":       "cert",
				"authn_service_id": "acme-vm",
				"authn_cert":       "client-cert-pem",
				"authn_cert_key":   "client-key-pem"
			}
		}]
	}`

	services := VcapServices{}
	err := json.Unmarshal([]byte(vcapJSON), &services)
	assert.NoError(t, err)

	creds := services.ConjurInfo.Credentials
	assert.Equal(t, "https://edge.example.com", creds.ApplianceURL)
	assert.Equal(t, "conjur", creds.Account)
	assert.Equal(t, "host/data/vms/vm-01", creds.Login)
	assert.Equal(t, "cert", creds.AuthnType)
	assert.Equal(t, "acme-vm", creds.ServiceID)
	assert.Equal(t, "client-cert-pem", creds.ClientCert)
	assert.Equal(t, "client-key-pem", creds.ClientKey)
}
