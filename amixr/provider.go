package amixr

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func Provider() terraform.ResourceProvider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"token": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("AMIXR_API_KEY", nil),
			},
		},
		DataSourcesMap: map[string]*schema.Resource{
			"amixr_user":          dataSourceAmixrUser(),
			"amixr_schedule":      dataSourceAmixrSchedule(),
			"amixr_slack_channel": dataSourceAmixrSlackChannel(),
			"amixr_action":        dataSourceAmixrAction(),
			"amixr_user_group":    dataSourceAmixrUserGroup(),
		},
		ResourcesMap: map[string]*schema.Resource{
			"amixr_integration":   resourceIntegration(),
			"amixr_escalation":    resourceEscalation(),
			"amixr_route":         resourceRoute(),
			"amixr_on_call_shift": resourceOnCallShift(),
			"amixr_schedule":      resourceSchedule(),
		},
		ConfigureFunc: providerConfigure,
	}
}

func providerConfigure(d *schema.ResourceData) (interface{}, error) {
	config := Config{
		Token: d.Get("token").(string),
	}

	return config.Client()
}
