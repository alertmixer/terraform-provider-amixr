package amixr

import (
	"github.com/alertmixer/amixr-go-client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"log"
)

func resourceRoute() *schema.Resource {
	return &schema.Resource{
		Create: resourceRouteCreate,
		Read:   handleNonExistentResource(resourceRouteRead),
		Update: resourceRouteUpdate,
		Delete: resourceRouteDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"integration_id": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"escalation_chain_id": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"position": &schema.Schema{
				Type:     schema.TypeInt,
				Required: true,
			},
			"routing_regex": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"slack": &schema.Schema{
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"channel_id": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
				MaxItems: 1,
			},
		},
	}
}

func resourceRouteCreate(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] create amixr route")

	client := m.(*amixr.Client)

	integrationIdData := d.Get("integration_id").(string)
	escalationChainIdData := d.Get("escalation_chain_id").(string)
	routingRegexData := d.Get("routing_regex").(string)
	positionData := d.Get("position").(int)
	slackData := d.Get("slack").([]interface{})

	createOptions := &amixr.CreateRouteOptions{
		IntegrationId:     integrationIdData,
		EscalationChainId: escalationChainIdData,
		RoutingRegex:      routingRegexData,
		Position:          &positionData,
		ManualOrder:       true,
		Slack:             expandRouteSlack(slackData),
	}

	route, _, err := client.Routes.CreateRoute(createOptions)
	if err != nil {
		return err
	}

	d.SetId(route.ID)

	return resourceRouteRead(d, m)
}

func resourceRouteRead(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] read amixr route")

	client := m.(*amixr.Client)

	route, _, err := client.Routes.GetRoute(d.Id(), &amixr.GetRouteOptions{})
	if err != nil {
		return err
	}

	d.Set("integration_id", route.IntegrationId)
	d.Set("escalation_chain_id", route.EscalationChainId)
	d.Set("routing_regex", route.RoutingRegex)
	d.Set("position", route.Position)
	d.Set("slack", flattenRouteSlack(route.SlackRoute))

	return nil
}

func resourceRouteUpdate(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] update amixr route")
	client := m.(*amixr.Client)

	escalationChainIdData := d.Get("escalation_chain_id").(string)
	routingRegexData := d.Get("routing_regex").(string)
	positionData := d.Get("position").(int)
	slackData := d.Get("slack").([]interface{})

	updateOptions := &amixr.UpdateRouteOptions{
		EscalationChainId: escalationChainIdData,
		RoutingRegex:      routingRegexData,
		Position:          &positionData,
		ManualOrder:       true,
		Slack:             expandRouteSlack(slackData),
	}

	route, _, err := client.Routes.UpdateRoute(d.Id(), updateOptions)
	if err != nil {
		return err
	}

	d.SetId(route.ID)
	return resourceRouteRead(d, m)
}

func resourceRouteDelete(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] delete amixr route")

	client := m.(*amixr.Client)

	_, err := client.Routes.DeleteRoute(d.Id(), &amixr.DeleteRouteOptions{})
	if err != nil {
		return err
	}

	d.SetId("")

	return nil
}

func flattenRouteSlack(in *amixr.SlackRoute) []map[string]interface{} {
	slack := make([]map[string]interface{}, 0, 1)

	out := make(map[string]interface{})

	if in.ChannelId != nil {
		out["channel_id"] = in.ChannelId
		slack = append(slack, out)
	}
	return slack
}

func expandRouteSlack(in []interface{}) *amixr.SlackRoute {
	slackRoute := amixr.SlackRoute{}

	for _, r := range in {
		inputMap := r.(map[string]interface{})
		if inputMap["channel_id"] != "" {
			channelId := inputMap["channel_id"].(string)
			slackRoute.ChannelId = &channelId
		}
	}

	return &slackRoute

}
