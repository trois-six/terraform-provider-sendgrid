package sendgrid

import (
	"context"
	"errors"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	sendgrid "github.com/trois-six/terraform-provider-sendgrid/sdk"
)

func dataSendgridTemplateVersion() *schema.Resource {
	s := resourceSendgridTemplateVersion().Schema

	for key, val := range s {
		if key != "template_id" {
			val.Computed = true
			val.Optional = false
			val.Required = false
			val.Default = nil
			val.ValidateFunc = nil
		}
	}
	return &schema.Resource{
		ReadContext: dataSendgridTemplateVersionRead,
		Schema:      s,
	}
}

func dataSendgridTemplateVersionRead(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	templateId := d.Get("template_id").(string)
	c := m.(*sendgrid.Client)

	template, err := c.ReadTemplate(templateId)
	if err != nil {
		return diag.FromErr(err)
	}
	var activeVersion *sendgrid.TemplateVersion
	for _, version := range template.Versions {
		if version.Active == 1 {
			activeVersion = &version
			break
		}
	}

	if activeVersion == nil {
		return diag.FromErr(errors.New("no recent version found for template_id"))
	}

	d.SetId(activeVersion.ID)

	if err := parseTemplateVersion(d, activeVersion); err != nil {
		return diag.FromErr(err)
	}
	return nil
}
