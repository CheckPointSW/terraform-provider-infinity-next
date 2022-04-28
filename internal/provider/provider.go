package provider

import (
	"context"

	"github.com/CheckPointSW/terraform-provider-infinity-next/internal/api"
	"github.com/CheckPointSW/terraform-provider-infinity-next/internal/resources"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func init() {
	// Set descriptions to support markdown syntax, this will be used in document generation
	// and the language server.
	schema.DescriptionKind = schema.StringMarkdown
}

func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"region": {
				Description:      "The region where Infinity Policy operations will take place. Options are: us, eu",
				Type:             schema.TypeString,
				Optional:         true,
				ForceNew:         true,
				ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{"eu", "us"}, false)),
				DefaultFunc:      schema.EnvDefaultFunc("INEXT_REGION", "eu"),
			},
			"client_id": {
				Description: "The client id for API operations, You can retrieve this\n" +
					"from the 'Global Settings -> API Keys' section of the Infinity Next portal",
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("INEXT_CLIENT_ID", ""),
			},
			"access_key": {
				Description: "The access key for API operations. You can retrieve this\n" +
					"from the 'Global Settings -> API Keys' section of the Infinity Next portal",
				Type:        schema.TypeString,
				Required:    true,
				Sensitive:   true,
				DefaultFunc: schema.EnvDefaultFunc("INEXT_ACCESS_KEY", ""),
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"inext_log_trigger":            resources.ResourceLogTrigger(),
			"inext_appsec_gateway_profile": resources.ResourceAppSecGatewayProfile(),
			"inext_embedded_profile":       resources.ResourceEmbeddedProfile(),
			"inext_docker_profile":         resources.ResourceDockerProfile(),
			"inext_kubernetes_profile":     resources.ResourceKubernetesProfile(),
			"inext_web_app_asset":          resources.ResourceWebAppAsset(),
			"inext_web_api_asset":          resources.ResourceWebAPIAsset(),
			"inext_web_app_practice":       resources.ResourceWebAppPractice(),
			"inext_web_api_practice":       resources.ResourceWebAPIPractice(),
			"inext_trusted_sources":        resources.ResourceTrustedSources(),
			"inext_exceptions":             resources.ResourceExceptions(),
			"inext_access_token":           resources.ResourceAccessToken(),
		},
		ConfigureContextFunc: providerConfigure,
	}
}

func providerConfigure(ctx context.Context, d *schema.ResourceData) (any, diag.Diagnostics) {
	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	region := d.Get("region").(string)
	client_id := d.Get("client_id").(string)
	access_key := d.Get("access_key").(string)

	if len(client_id) == 0 || len(access_key) == 0 {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Missing credentials",
			Detail:   "Must define client_id and access_key",
		})
		return nil, diags
	}

	client := api.NewClient()

	switch region {
	case "eu":
		client.SetHost("https://cloudinfra-gw.portal.checkpoint.com")
		client.SetEndpoint("/app/i2/graphql/V1")
		// client.SetHost("https://dev-cloudinfra-gw.kube1.iaas.checkpoint.com")
		// client.SetEndpoint("/app/infinity2gem/graphql/V1")
	case "us":
		client.SetHost("https://cloudinfra-gw-us.portal.checkpoint.com")
		client.SetEndpoint("/app/i2/graphql/V1")
	case "dev":
		client.SetHost("https://dev-cloudinfra-gw.kube1.iaas.checkpoint.com")
		client.SetEndpoint("/app/infinity2gem/graphql/V1")
	default:
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Invalid region",
			Detail:   "region must be eu or us",
		})
		return nil, diags
	}

	if err := client.InfinityPortalAuthentication(client_id, access_key); err != nil {
		return nil, diag.FromErr(err)
	}

	return client, nil
}
