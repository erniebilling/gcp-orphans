package gcp

import (
	"encoding/json"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	gcpCompute "google.golang.org/api/compute/v1"
)

type GCPClientInterface interface {
	Connect(creds *GCPCredentials) error
	GetFirewallRules() (*[]FirewallRule, error)
	GetVMInstances() (*[]VMInstance, error)
}

/*
	Matches a Google Developers service account JSON key file
*/
type GCPCredentials struct {
	Type                    string `json:"type"`
	ProjectID               string `json:"project_id"`
	PrivateKeyID            string `json:"private_key_id"`
	PrivateKey              string `json:"private_key"`
	ClientEmail             string `json:"client_email"`
	ClientID                string `json:"client_id"`
	AuthURI                 string `json:"auth_uri"`
	TokenURI                string `json:"token_uri"`
	AuthProviderX509CertURL string `json:"auth_provider_x509_cert_url"`
	ClientX509CertURL       string `json:"client_x509_cert_url"`
}

type GCPClient struct {
	ComputeService *gcpCompute.Service
	Creds          *GCPCredentials
}

func (c *GCPClient) Connect(creds *GCPCredentials) error {
	c.Creds = creds

	credentialsJson, err := json.Marshal(creds)
	if err != nil {
		return err
	}

	conf, err := google.JWTConfigFromJSON(credentialsJson, gcpCompute.ComputeScope)
	if err != nil {
		return err
	}

	client := conf.Client(oauth2.NoContext)

	c.ComputeService, err = gcpCompute.New(client)

	return err
}

type FirewallRule struct {
	Name string
	URL  string
	Tags []string
}

func (c *GCPClient) GetFirewallRules() (*[]FirewallRule, error) {
	firewalls, err := c.ComputeService.Firewalls.List(c.Creds.ProjectID).Filter(`direction="INGRESS"`).Do()
	if err != nil {
		return nil, err
	}

	rules := make([]FirewallRule, len(firewalls.Items))

	for n, firewall := range firewalls.Items {
		rules[n] = FirewallRule{Name: firewall.Name, URL: firewall.SelfLink, Tags: firewall.TargetTags}
	}
	return &rules, nil
}

type VMInstance struct {
	Name        string
	URL         string
	NetworkTags []string
}

func (c *GCPClient) GetVMInstances() (*[]VMInstance, error) {
	vmAggInstances, err := c.ComputeService.Instances.AggregatedList(c.Creds.ProjectID).Do()
	if err != nil {
		return nil, err
	}

	var instances []VMInstance

	for _, vmInstances := range vmAggInstances.Items {
		for _, vmInstance := range vmInstances.Instances {
			instances = append(instances, VMInstance{Name: vmInstance.Name, URL: vmInstance.SelfLink, NetworkTags: vmInstance.Tags.Items})
		}
	}

	return &instances, nil
}
