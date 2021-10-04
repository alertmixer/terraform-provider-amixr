package amixr

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"token": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("AMIXR_API_KEY", nil),
			},
		},
		DataSourcesMap: map[string]*schema.Resource{
			"amixr_user":             dataSourceAmixrUser(),
			"amixr_escalation_chain": dataSourceEscalationChain(),
			"amixr_schedule":         dataSourceAmixrSchedule(),
			"amixr_slack_channel":    dataSourceAmixrSlackChannel(),
			"amixr_action":           dataSourceAmixrAction(),
			"amixr_user_group":       dataSourceAmixrUserGroup(),
			"amixr_team":             dataSourceAmixrTeam(),
		},
		ResourcesMap: map[string]*schema.Resource{
			"amixr_integration":      resourceIntegration(),
			"amixr_escalation_chain": resourceEscalationChain(),
			"amixr_escalation":       resourceEscalation(),
			"amixr_route":            resourceRoute(),
			"amixr_on_call_shift":    resourceOnCallShift(),
			"amixr_schedule":         resourceSchedule(),
		},
		ConfigureContextFunc: providerConfigure,
	}
}

func providerConfigure(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
	var diags diag.Diagnostics

	config := Config{
		Token: d.Get("token").(string),
	}
	client, err := config.Client()
	if err != nil {
		return nil, diag.FromErr(err)
	}
	return client, diags
}
