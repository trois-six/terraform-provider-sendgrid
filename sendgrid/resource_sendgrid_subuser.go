/*
Provide a resource to manage a subuser.
Example Usage
```hcl
resource "sendgrid_subuser" "subuser" {
	username = "my-subuser"
	email    = "subuser@example.org"
	password = "Passw0rd!"
	ips      = [
		"127.0.0.1"
	]
}
```
Import
A subuser can be imported, e.g.
```hcl
$ terraform import sendgrid_subuser.subuser userName
```
*/
package sendgrid

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	sendgrid "github.com/trois-six/terraform-provider-sendgrid/sdk"
)

func resourceSendgridSubuser() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceSendgridSubuserCreate,
		ReadContext:   resourceSendgridSubuserRead,
		UpdateContext: resourceSendgridSubuserUpdate,
		DeleteContext: resourceSendgridSubuserDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"username": {
				Type:        schema.TypeString,
				Description: "The name of the subuser.",
				Required:    true,
			},
			"password": {
				Type:        schema.TypeString,
				Description: "The password the subuser will use when logging into SendGrid.",
				Sensitive:   true,
				Required:    true,
			},
			"email": {
				Type:        schema.TypeString,
				Description: "The email of the subuser.",
				Required:    true,
			},
			"ips": {
				Type:        schema.TypeSet,
				Description: "The IP addresses that should be assigned to this subuser.",
				Required:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				//ValidateDiagFunc: scopeInScopes("sender_verification_eligible"),
			},
			"user_id": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"disabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"signup_session_token": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"authorization_token": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"credit_allocation_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceSendgridSubuserCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*sendgrid.Client)

	username := d.Get("username").(string)
	password := d.Get("password").(string)
	email := d.Get("email").(string)
	ips := d.Get("ips").([]string)

	Subuser, err := c.CreateSubuser(username, email, password, ips)
	if err != nil {
		return diag.FromErr(err)
	}

	//nolint:errcheck
	d.Set("user_id", Subuser.UserID)
	//nolint:errcheck
	d.Set("signup_session_token", Subuser.SignupSessionToken)
	//nolint:errcheck
	d.Set("authorization_token", Subuser.AuthorizationToken)
	//nolint:errcheck
	d.Set("credit_allocation_type", Subuser.CreditAllocation.Type)

	if d.Get("disabled").(bool) {
		d.SetId(Subuser.UserName)
		return resourceSendgridSubuserUpdate(ctx, d, m)
	}

	d.SetId(Subuser.UserName)

	return nil
}

func resourceSendgridSubuserRead(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*sendgrid.Client)

	Subuser, err := c.ReadSubuser(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	//nolint:errcheck
	d.Set("user_id", Subuser.ID)
	//nolint:errcheck
	d.Set("disabled", Subuser.Disabled)
	//nolint:errcheck
	d.Set("email", Subuser.Email)

	return nil
}

func resourceSendgridSubuserUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*sendgrid.Client)

	if d.HasChange("disabled") {
		_, err := c.UpdateSubuser(d.Id(), d.Get("disabled").(bool))
		if err != nil {
			return diag.FromErr(err)
		}
	}

	return resourceSendgridSubuserRead(ctx, d, m)
}

func resourceSendgridSubuserDelete(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*sendgrid.Client)

	_, err := c.DeleteSubuser(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceSendgridSubuserImport(ctx context.Context, d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, diag.Diagnostics) {
	if diags := resourceSendgridSubuserRead(ctx, d, m); diags != nil {
		return nil, diags
	}
	return []*schema.ResourceData{d}, nil
}
