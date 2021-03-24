/*
Provide a resource to manage an event webhook settings.
Example Usage
```hcl
resource "sendgrid_event_webhook" "default" {
	enabled = true
    url = "https://foo.bar/sendgrid/inbound"
    group_resubscribe = true
    delivered = true
    group_unsubscribe = true
    spam_report = true
    bounce = true
    deferred = true
    unsubscribe = true
    processed = true
    open = true
    click = true
    dropped = true
    oauth_client_id = "a-client-id"
    oauth_client_secret = "a-client-secret"
    oauth_token_url = "https://oauth.example.com/token"
}
```
*/
package sendgrid

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	sendgrid "github.com/trois-six/terraform-provider-sendgrid/sdk"
)

func resourceSendgridEventWebhook() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceSendgridEventWebhookPatch,
		ReadContext:   resourceSendgridEventWebhookRead,
		UpdateContext: resourceSendgridEventWebhookPatch,

		Schema: map[string]*schema.Schema{
			"enabled": {
				Type:        schema.TypeBool,
				Description: "Indicates if the event webhook is enabled.",
				Required:    true,
			},
			"url": {
				Type:        schema.TypeString,
				Description: "The public URL where you would like SendGrid to POST the data events from your email. Any emails sent with the given hostname provided (whose MX records have been updated to point to SendGrid) will be eventd and POSTed to this URL.",
				Required:    true,
			},
			"group_resubscribe": {
				Type:        schema.TypeBool,
				Description: "Recipient resubscribes to specific group by updating preferences. You need to enable Subscription Tracking for getting this type of event.",
				Optional:    true,
			},
			"delivered": {
				Type:        schema.TypeBool,
				Description: "Message has been successfully delivered to the receiving server.",
				Optional:    true,
			},
			"group_unsubscribe": {
				Type:        schema.TypeBool,
				Description: "Recipient unsubscribe from specific group, by either direct link or updating preferences. You need to enable Subscription Tracking for getting this type of event.",
				Optional:    true,
			},
			"spam_report": {
				Type:        schema.TypeBool,
				Description: "Recipient marked a message as spam.",
				Optional:    true,
			},
			"bounce": {
				Type:        schema.TypeBool,
				Description: "Receiving server could not or would not accept message.",
				Optional:    true,
			},
			"deferred": {
				Type:        schema.TypeBool,
				Description: "Recipient's email server temporarily rejected message.",
				Optional:    true,
			},
			"unsubscribe": {
				Type:        schema.TypeBool,
				Description: "Recipient clicked on message's subscription management link. You need to enable Subscription Tracking for getting this type of event.",
				Optional:    true,
			},
			"processed": {
				Type:        schema.TypeBool,
				Description: "Message has been received and is ready to be delivered.",
				Optional:    true,
			},
			"open": {
				Type:        schema.TypeBool,
				Description: "Recipient has opened the HTML message. You need to enable Open Tracking for getting this type of event.",
				Optional:    true,
			},
			"click": {
				Type:        schema.TypeBool,
				Description: "Recipient clicked on a link within the message. You need to enable Click Tracking for getting this type of event.",
				Optional:    true,
			},
			"dropped": {
				Type:        schema.TypeBool,
				Description: "You may see the following drop reasons: Invalid SMTPAPI header, Spam Content (if spam checker app enabled), Unsubscribed Address, Bounced Address, Spam Reporting Address, Invalid, Recipient List over Package Quota.",
				Optional:    true,
			},
			"oauth_client_id": {
				Type:        schema.TypeString,
				Description: "The client ID Twilio SendGrid sends to your OAuth server or service provider to generate an OAuth access token.",
				Optional:    true,
			},
			"oauth_client_secret": {
				Type:        schema.TypeString,
				Description: "This secret is needed only once to create an access token. SendGrid will store this secret, allowing you to update your Client ID and Token URL without passing the secret to SendGrid again. When passing data in this field, you must also include the oauth_client_id and oauth_token_url fields.",
				Optional:    true,
			},
			"oauth_token_url": {
				Type:        schema.TypeString,
				Description: "The URL where Twilio SendGrid sends the Client ID and Client Secret to generate an access token. This should be your OAuth server or service provider. When passing data in this field, you must also include the oauth_client_id field.",
				Optional:    true,
			},
		},
	}
}

func resourceSendgridEventWebhookPatch(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*sendgrid.Client)

	enabled := d.Get("enabled").(bool)
	url := d.Get("url").(string)
	groupResubscribe := d.Get("group_resubscribe").(bool)
	delivered := d.Get("delivered").(bool)
	groupUnsubscribe := d.Get("group_unsubscribe").(bool)
	spamReport := d.Get("spam_report").(bool)
	bounce := d.Get("bounce").(bool)
	deferred := d.Get("deferred").(bool)
	unsubscribe := d.Get("unsubscribe").(bool)
	processed := d.Get("processed").(bool)
	open := d.Get("open").(bool)
	click := d.Get("click").(bool)
	dropped := d.Get("dropped").(bool)
	oauthClientId := d.Get("oauth_client_id").(string)
	oauthClientSecret := d.Get("oauth_client_secret").(string)
	oauthTokenUrl := d.Get("oauth_token_url").(string)

	_, err := sendgrid.RetryOnRateLimit(ctx, d, func() (interface{}, sendgrid.RequestError) {
		return c.PatchEventWebhook(enabled, url, groupResubscribe, delivered, groupUnsubscribe, spamReport, bounce, deferred, unsubscribe, processed, open, click, dropped, oauthClientId, oauthClientSecret, oauthTokenUrl)
	})

	if err != nil {
		return diag.FromErr(err)
	}

	return resourceSendgridEventWebhookRead(ctx, d, m)
}

func resourceSendgridEventWebhookRead(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*sendgrid.Client)

	webhook, err := c.ReadEventWebhook()
	if err.Err != nil {
		return diag.FromErr(err.Err)
	}

	//nolint:errcheck
	d.Set("enabled", webhook.Enabled)
	//nolint:errcheck
	d.Set("url", webhook.Url)
	//nolint:errcheck
	d.Set("group_resubscribe", webhook.GroupResubscribe)
	//nolint:errcheck
	d.Set("delivered", webhook.Delivered)
	//nolint:errcheck
	d.Set("group_unsubscribe", webhook.GroupUnsubscribe)
	//nolint:errcheck
	d.Set("spam_report", webhook.SpamReport)
	//nolint:errcheck
	d.Set("bounce", webhook.Bounce)
	//nolint:errcheck
	d.Set("deferred", webhook.Delivered)
	//nolint:errcheck
	d.Set("unsubscribe", webhook.Unsubscribe)
	//nolint:errcheck
	d.Set("processed", webhook.Processed)
	//nolint:errcheck
	d.Set("open", webhook.Open)
	//nolint:errcheck
	d.Set("click", webhook.Click)
	//nolint:errcheck
	d.Set("dropped", webhook.Dropped)
	//nolint:errcheck
	d.Set("oauth_client_id", webhook.OAuthClientId)
	//nolint:errcheck
	d.Set("oauth_token_url", webhook.OAuthTokenUrl)

	return nil
}
