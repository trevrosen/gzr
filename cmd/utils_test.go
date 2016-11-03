package cmd

import "testing"

var testGozerConfigYaml = `---
apiVersion: v1
deployments:
  -
    build-steps:
      - "make build"
      - "make tag"
      - "make push"
    default-image: "bypass/bypass:latest"
    name: bypass-gozer
    temple-path: bypass/deployment/bypass-api.yaml
services:
  -
    externalPort: 443
    internalPort: 80
    name: bypass-gozer
    params:
      service.beta.kubernetes.io/aws-load-balancer-backend-protocol: http
      service.beta.kubernetes.io/aws-load-balancer-ssl-cert: "arn:aws:acm:us-east-1:562983362877:certificate/e33cc7bd-afe2-4791-ac96-e2a6412e2330"
      service.bypass.hostname: gozer.bypassmobile.com
      service.bypass.roledefs: false
    selector:
      name: bypass-gozer
    type: LoadBalancer`

func TestGozerConfig(t *testing.T) {
	config, err := gozerConfigFromBytes([]byte(testGozerConfigYaml))

	versionOK := config.ApiVersion == "v1"
	deploymentsNameOK := config.Deployments[0].Name == "bypass-gozer"
	serviceExternalPortOK := config.Services[0].ExternalPort == 443

	if err != nil {
		t.Errorf("Parse error: %s", err.Error())
	}

	if !versionOK || !deploymentsNameOK || !serviceExternalPortOK {
		t.Errorf("Expected sanity check values do not match fixture")
	}
}
