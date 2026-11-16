module github.com/cyberark/cloudfoundry-conjur-buildpack/conjur-env

go 1.26.10

require (
	github.com/cyberark/conjur-api-go v0.17.1
	github.com/cyberark/summon v0.13.1
	github.com/stretchr/testify v1.12.1
)

require (
	github.com/Masterminds/semver/v3 v3.5.0 // indirect
	github.com/Microsoft/go-winio v0.6.2 // indirect
	github.com/apapsch/go-jsonmerge/v2 v2.0.0 // indirect
	github.com/aws/aws-sdk-go-v2 v1.47.1 // indirect
	github.com/aws/aws-sdk-go-v2/config v1.33.6 // indirect
	github.com/aws/aws-sdk-go-v2/credentials v1.20.6 // indirect
	github.com/aws/aws-sdk-go-v2/feature/ec2/imds v1.20.1 // indirect
	github.com/aws/aws-sdk-go-v2/internal/configsources v1.5.4 // indirect
	github.com/aws/aws-sdk-go-v2/internal/endpoints/v2 v2.8.4 // indirect
	github.com/aws/aws-sdk-go-v2/internal/v4a v1.5.4 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/accept-encoding v1.13.19 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/presigned-url v1.14.4 // indirect
	github.com/aws/aws-sdk-go-v2/service/signin v1.10.1 // indirect
	github.com/aws/aws-sdk-go-v2/service/sso v1.38.1 // indirect
	github.com/aws/aws-sdk-go-v2/service/ssooidc v1.43.1 // indirect
	github.com/aws/aws-sdk-go-v2/service/sts v1.51.1 // indirect
	github.com/aws/smithy-go v1.28.2 // indirect
	github.com/bgentry/go-netrc v0.0.0-20140422174119-9fd32a8b3d3d // indirect
	github.com/danieljoos/wincred v1.2.3 // indirect
	github.com/go-jose/go-jose/v4 v4.1.6 // indirect
	github.com/godbus/dbus/v5 v5.2.2 // indirect
	github.com/google/uuid v1.6.0 // indirect
	github.com/oapi-codegen/runtime v1.9.0 // indirect
	github.com/sirupsen/logrus v1.10.2 // indirect
	github.com/spiffe/go-spiffe/v2 v2.8.2 // indirect
	github.com/zalando/go-keyring v0.2.8 // indirect
	go.yaml.in/yaml/v3 v3.0.5 // indirect
	golang.org/x/net v0.62.0 // indirect
	golang.org/x/sys v0.49.0 // indirect
	golang.org/x/text v0.43.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20260928230214-8a89bd6388cc // indirect
	google.golang.org/grpc v1.84.0 // indirect
	google.golang.org/protobuf v1.36.12 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

// Vulnerable go-jose (fixed in 4.1.6) and x/net (fixed in 0.62.0) versions still required by go-spiffe, conjur-api-go and grpc
exclude (
	github.com/go-jose/go-jose/v4 v4.1.4
	github.com/go-jose/go-jose/v4 v4.1.5
	golang.org/x/net v0.48.0
	golang.org/x/net v0.57.0
	golang.org/x/net v0.60.0
)
