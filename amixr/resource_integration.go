package amixr

import (
	amixr "github.com/alertmixer/amixr-go-client"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"log"
)

var integrationTypes = []string{
	"grafana",
	"webhook",
	"alertmanager",
	"kapacitor",
	"fabric",
	"newrelic",
	"datadog",
	"pagerduty",
	"pingdom",
	"elastalert",
	"amazon_sns",
	"curler",
	"sentry",
	"formatted_webhook",
	"heartbeat",
	"demo",
	"manual",
	"stackdriver",
	"uptimerobot",
	"sentry_platform",
	"zabbix",
	"prtg",
	"slack_channel",
	"inbound_email",
}

func resourceIntegration() *schema.Resource {
	return &schema.Resource{
		Create: resourceIntegrationCreate,
		Read:   resourceIntegrationRead,
		Update: resourceIntegrationUpdate,
		Delete: resourceIntegrationDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},
			"type": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringInSlice(integrationTypes, false),
				ForceNew:     true,
			},
			"default_route_id": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceIntegrationCreate(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] create amixr integration")

	client := m.(*amixr.Client)

	nameData := d.Get("name").(string)
	typeData := d.Get("type").(string)
	createOptions := &amixr.CreateIntegrationOptions{
		Name: nameData,
		Type: typeData,
	}

	integration, _, err := client.Integrations.CreateIntegration(createOptions)
	if err != nil {
		return err
	}

	d.SetId(integration.ID)

	return resourceIntegrationRead(d, m)
}

func resourceIntegrationUpdate(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] update amixr integration")

	client := m.(*amixr.Client)

	nameData := d.Get("name").(string)
	updateOptions := &amixr.UpdateIntegrationOptions{
		Name: nameData,
	}

	integration, _, err := client.Integrations.UpdateIntegration(d.Id(), updateOptions)
	if err != nil {
		return err
	}

	d.SetId(integration.ID)

	return resourceIntegrationRead(d, m)
}

func resourceIntegrationRead(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] read amixr integration")

	client := m.(*amixr.Client)
	options := &amixr.GetIntegrationOptions{}
	integration, _, err := client.Integrations.GetIntegration(d.Id(), options)
	if err != nil {
		return err
	}

	d.Set("default_route_id", integration.DefaultRouteId)
	d.Set("name", integration.Name)
	d.Set("type", integration.Type)

	return nil
}

func resourceIntegrationDelete(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] delete amixr integration")

	client := m.(*amixr.Client)
	options := &amixr.DeleteIntegrationOptions{}
	_, err := client.Integrations.DeleteIntegration(d.Id(), options)
	if err != nil {
		return err
	}

	d.SetId("")

	return nil
}
