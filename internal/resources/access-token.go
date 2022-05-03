package resources

import (
	"context"
	"strconv"

	"github.com/CheckPointSW/terraform-provider-infinity-next/internal/api"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func ResourceAccessToken() *schema.Resource {
	return &schema.Resource{
		Description: "Access Token",

		CreateContext: resourceAccessTokenCreate,
		ReadContext:   resourceAccessTokenRead,
		UpdateContext: resourceAccessTokenUpdate,
		DeleteContext: resourceAccessTokenDelete,

		Schema: map[string]*schema.Schema{
			"value": {
				Type:      schema.TypeString,
				Optional:  true,
				Computed:  true,
				Sensitive: true,
			},
		},
	}
}

func resourceAccessTokenCreate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	var diags diag.Diagnostics
	c := meta.(*api.Client)
	token := c.GetToken()
	d.SetId(strconv.Itoa(schema.HashString(token)))
	d.Set("value", token)

	return diags
}

func resourceAccessTokenRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	return nil
}

func resourceAccessTokenUpdate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	return nil
}

func resourceAccessTokenDelete(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	var diags diag.Diagnostics
	d.SetId("")

	return diags
}
