package amixr

import (
	"github.com/alertmixer/amixr-go-client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
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
            "team_id": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
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
			"link": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"templates": &schema.Schema{
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"resolve_signal": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"grouping_key": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"slack": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"title": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"message": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"image_url": {
										Type:     schema.TypeString,
										Optional: true,
									},
								},
							},
							MaxItems: 1,
						},
					},
				},
				MaxItems: 1,
			},
		},
	}
}

func resourceIntegrationCreate(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] create amixr integration")

	client := m.(*amixr.Client)

	teamIdData := d.Get("team_id").(string)
	nameData := d.Get("name").(string)
	typeData := d.Get("type").(string)
	templatesData := d.Get("templates").([]interface{})

	createOptions := &amixr.CreateIntegrationOptions{
	    TeamId:    teamIdData,
		Name:      nameData,
		Type:      typeData,
		Templates: expandTemplates(templatesData),
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

	teamIdData := d.Get("team_id").(string)
	nameData := d.Get("name").(string)
	templateData := d.Get("templates").([]interface{})

	updateOptions := &amixr.UpdateIntegrationOptions{
	    TeamId:    teamIdData,
		Name:      nameData,
		Templates: expandTemplates(templateData),
	}

	integration, _, err := client.Integrations.UpdateIntegration(d.Id(), updateOptions)
	if err != nil {
		return err
	}

	d.SetId(integration.ID)

	return resourceIntegrationRead(d, m)
}

func resourceIntegrationRead(d *schema.ResourceData, m interface{}) error {

	client := m.(*amixr.Client)
	options := &amixr.GetIntegrationOptions{}
	integration, _, err := client.Integrations.GetIntegration(d.Id(), options)
	if err != nil {
		return err
	}

	d.Set("team_id", integration.TeamId)
	d.Set("default_route_id", integration.DefaultRouteId)
	d.Set("name", integration.Name)
	d.Set("type", integration.Type)
	d.Set("templates", flattenTemplates(integration.Templates))
	d.Set("link", integration.Link)

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

func flattenTemplates(in *amixr.Templates) []map[string]interface{} {
	templates := make([]map[string]interface{}, 0, 1)
	out := make(map[string]interface{})

	out["grouping_key"] = in.GroupingKey
	out["resolve_signal"] = in.ResolveSignal
	out["slack"] = flattenSlackTemplate(in.Slack)

	add := false

	if in.GroupingKey != nil {
		out["grouping_key"] = in.GroupingKey
		add = true
	}
	if in.ResolveSignal != nil {
		out["resolve_signal"] = in.ResolveSignal
		add = true

	}
	if in.Slack != nil {
		flattenSlackTemplate := flattenSlackTemplate(in.Slack)
		if len(flattenSlackTemplate) > 0 {
			out["resolve_signal"] = in.ResolveSignal
			add = true
		}
	}

	if add {
		templates = append(templates, out)
	}

	return templates
}

func flattenSlackTemplate(in *amixr.SlackTemplate) []map[string]interface{} {
	slackTemplates := make([]map[string]interface{}, 0, 1)

	add := false

	slackTemplate := make(map[string]interface{})

	if in.Title != nil {
		slackTemplate["title"] = in.Title
		add = true
	}
	if in.ImageURL != nil {
		slackTemplate["image_url"] = in.ImageURL
		add = true
	}
	if in.Message != nil {
		slackTemplate["message"] = in.Message
		add = true
	}

	if add {
		slackTemplates = append(slackTemplates, slackTemplate)
	}

	return slackTemplates
}

func expandTemplates(input []interface{}) *amixr.Templates {

	templates := amixr.Templates{}

	for _, r := range input {
		inputMap := r.(map[string]interface{})
		if inputMap["grouping_key"] != "" {
			gk := inputMap["grouping_key"].(string)
			templates.GroupingKey = &gk
		}
		if inputMap["resolve_signal"] != "" {
			rs := inputMap["resolve_signal"].(string)
			templates.ResolveSignal = &rs
		}
		if inputMap["slack"] == nil {
			templates.Slack = nil
		} else {
			templates.Slack = expandSlackTemplate(inputMap["slack"].([]interface{}))
		}
	}
	return &templates
}

func expandSlackTemplate(in []interface{}) *amixr.SlackTemplate {

	slackTemplate := amixr.SlackTemplate{}
	for _, r := range in {
		inputMap := r.(map[string]interface{})
		if inputMap["title"] != "" {
			t := inputMap["title"].(string)
			slackTemplate.Title = &t
		}
		if inputMap["message"] != "" {
			m := inputMap["message"].(string)
			slackTemplate.Message = &m
		}
		if inputMap["image_url"] != "" {
			iu := inputMap["image_url"].(string)
			slackTemplate.ImageURL = &iu
		}
	}
	return &slackTemplate
}
