package sendgrid

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	sendgrid "github.com/trois-six/terraform-provider-sendgrid/sdk"
)

func dataSendgridUnsubscribeGroup() *schema.Resource {
	s := resourceSendgridUnsubscribeGroup().Schema

	for _, val := range s {
		val.Computed = true
		val.Optional = false
		val.Required = false
		val.Default = nil
		val.ValidateFunc = nil
		val.ForceNew = true
	}

	s["group_id"] = &schema.Schema{
		Type:        schema.TypeString,
		Optional:    true,
		Description: "The id of the unsubscribe group to retrieve",
	}
	s["name"] = &schema.Schema{
		Type:        schema.TypeString,
		Optional:    true,
		Description: "The name of the unsubscribe group to retrieve",
	}
	return &schema.Resource{
		ReadContext: dataSendgridUnsubscribeGroupRead,
		Schema:      s,
	}
}

func dataSendgridUnsubscribeGroupRead(context context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	groupId := d.Get("group_id").(string)
	name := d.Get("name").(string)
	c := m.(*sendgrid.Client)

	if groupId != "" {
		d.SetId(groupId)
		return resourceSendgridUnsubscribeGroupRead(context, d, m)
	} else if name != "" {
		groups, err := c.ReadUnsubscribeGroups()
		if err.Err != nil {
			return diag.FromErr(err.Err)
		}
		for _, group := range groups {
			if group.Name == name {
				d.SetId(fmt.Sprint(group.ID))
				if err := sendgridUnsubscribeGroupParse(&group, d); err != nil {
					return diag.FromErr(err)
				}
				return nil
			}
		}
		return diag.Errorf("unable to find a unsubscribe group with name '%s'", name)
	} else {
		return diag.Errorf("either 'group_id' or 'name' must be specified for data.sendgrid_unsubscribe_group")
	}

}
