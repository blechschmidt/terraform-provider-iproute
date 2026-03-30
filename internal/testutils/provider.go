package testutils

import (
	"github.com/example/terraform-provider-iproute/internal/provider"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
)

// TestAccProtoV6ProviderFactories returns provider factories for acceptance tests.
func TestAccProtoV6ProviderFactories() map[string]func() (tfprotov6.ProviderServer, error) {
	return map[string]func() (tfprotov6.ProviderServer, error){
		"iproute": providerserver.NewProtocol6WithError(provider.New("test")()),
	}
}
