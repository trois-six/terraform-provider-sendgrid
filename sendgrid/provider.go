/*
Resources List

API key Resource
  sendgrid_api_key

Subuser resource
  sendgrid_subuser

Template Resources
  sendgrid_template
  sendgrid_template_version
*/
package sendgrid

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	sendgrid "github.com/trois-six/terraform-provider-sendgrid/sdk"
)

// Provider terraform.ResourceProvider.
func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"api_key": {
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   true,
				DefaultFunc: schema.EnvDefaultFunc("SENDGRID_API_KEY", nil),
			},
			"host": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("SENDGRID_HOST", nil),
			},
			"subuser": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("SENDGRID_SUBUSER", nil),
			},
		},

		ResourcesMap: map[string]*schema.Resource{
			"sendgrid_api_key":          resourceSendgridAPIKey(),
			"sendgrid_subuser":          resourceSendgridSubuser(),
			"sendgrid_template":         resourceSendgridTemplate(),
			"sendgrid_template_version": resourceSendgridTemplateVersion(),
		},

		ConfigureContextFunc: providerConfigure,
	}
}

func providerConfigure(_ context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
	var diags diag.Diagnostics

	apiKey, ok := d.Get("api_key").(string)
	if !ok || apiKey == "" {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Sendgrid API key wasn't provided",
			Detail: "Unable to retrieve the API key, " +
				"either from the configuration of the provider, " +
				"nor the env variable SENDGRID_API_KEY",
		})

		return nil, diags
	}

	host := d.Get("host").(string)
	subuser := d.Get("subuser").(string)

	return sendgrid.NewClient(apiKey, host, subuser), diags
}
