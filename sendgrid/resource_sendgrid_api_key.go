/*
Provide a resource to manage an API key.
Example Usage
```hcl
resource "sendgrid_api_key" "api_key" {
	name   = "my-api-key"
	scopes = [
		"mail.send",
		"sender_verification_eligible",
	]
}
```
Import
An API key can be imported, e.g.
```hcl
$ terraform import sendgrid_api_key.api_key apiKeyID
```
*/
package sendgrid

import (
	"context"
	"reflect"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	sendgrid "github.com/trois-six/terraform-provider-sendgrid/sdk"
)

func resourceSendgridAPIKey() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceSendgridAPIKeyCreate,
		ReadContext:   resourceSendgridAPIKeyRead,
		UpdateContext: resourceSendgridAPIKeyUpdate,
		DeleteContext: resourceSendgridAPIKeyDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Description: "The name you will use to describe this API Key.",
				Required:    true,
			},
			"sub_user_on_behalf_of": {
				Type:        schema.TypeString,
				Description: "The subuser's username. Generates the API call as if the subuser account was making the call",
				Optional:    true,
			},
			"scopes": {
				Type:        schema.TypeSet,
				Description: "The individual permissions that you are giving to this API Key.",
				Optional:    true,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"api_key": {
				Type:        schema.TypeString,
				Description: "The API key created by the API.",
				Computed:    true,
			},
		},
	}
}

func scopeInScopes(scopes []string, scope string) bool {
	for _, v := range scopes {
		if v == scope {
			return true
		}
	}

	return false
}

func resourceSendgridAPIKeyCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var scopes []string

	c := m.(*sendgrid.Client)
	name := d.Get("name").(string)
	c.OnBehalfOf = d.Get("sub_user_on_behalf_of").(string)

	for _, scope := range d.Get("scopes").(*schema.Set).List() {
		scopes = append(scopes, scope.(string))
	}

	if ok := scopeInScopes(scopes, "sender_verification_eligible"); !ok {
		scopes = append(scopes, "sender_verification_eligible")
	}

	if ok := scopeInScopes(scopes, "2fa_required"); !ok {
		scopes = append(scopes, "2fa_required")
	}

	apiKey, err := c.CreateAPIKey(name, scopes)
	if err.Err != nil {
		return diag.FromErr(err.Err)
	}

	d.SetId(apiKey.ID)
	//nolint:errcheck
	d.Set("api_key", apiKey.APIKey)

	return resourceSendgridAPIKeyRead(ctx, d, m)
}

func resourceSendgridAPIKeyRead(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*sendgrid.Client)

	c.OnBehalfOf = d.Get("sub_user_on_behalf_of").(string)

	apiKey, err := c.ReadAPIKey(d.Id())
	if err.Err != nil {
		return diag.FromErr(err.Err)
	}

	//nolint:errcheck
	d.Set("name", apiKey.Name)
	//nolint:errcheck
	d.Set("scopes", apiKey.Scopes)

	return nil
}

func hasDiff(o, n interface{}) bool {
	if eq, ok := o.(schema.Equal); ok {
		return !eq.Equal(n)
	}

	return !reflect.DeepEqual(o, n)
}

func resourceSendgridAPIKeyUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*sendgrid.Client)

	c.OnBehalfOf = d.Get("sub_user_on_behalf_of").(string)

	a := sendgrid.APIKey{
		ID:   d.Id(),
		Name: d.Get("name").(string),
	}

	o, n := d.GetChange("scopes")
	n.(*schema.Set).Add("sender_verification_eligible")
	n.(*schema.Set).Add("2fa_required")

	if ok := hasDiff(o, n); ok {
		var scopes []string
		for _, scope := range d.Get("scopes").(*schema.Set).List() {
			scopes = append(scopes, scope.(string))
		}

		a.Scopes = scopes
	}

	if _, err := c.UpdateAPIKey(d.Id(), a.Name, a.Scopes); err.Err != nil {
		return diag.FromErr(err.Err)
	}

	return resourceSendgridAPIKeyRead(ctx, d, m)
}

func resourceSendgridAPIKeyDelete(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*sendgrid.Client)

	c.OnBehalfOf = d.Get("sub_user_on_behalf_of").(string)

	_, err := c.DeleteAPIKey(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}
