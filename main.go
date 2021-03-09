package main

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/plugin"
	"github.com/trois-six/terraform-provider-sendgrid/sendgrid"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: sendgrid.Provider,
	})
}
