package amixr

import (
	amixr "github.com/alertmixer/amixr-go-client"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"log"
)

func resourceRoute() *schema.Resource {
	return &schema.Resource{
		Create: resourceRouteCreate,
		Read:   handleNonExistentResource(resourceRouteRead),
		Update: resourceRouteUpdate,
		Delete: resourceRouteDelete,

		Schema: map[string]*schema.Schema{
			"integration_id": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
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
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"channel_id": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
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
	routingRegexData := d.Get("routing_regex").(string)
	positionData := d.Get("position").(int)
	slackData := d.Get("slack").([]interface{})

	createOptions := &amixr.CreateRouteOptions{
		IntegrationId: integrationIdData,
		RoutingRegex:  routingRegexData,
		Position:      &positionData,
		ManualOrder:   true,
		Slack:         expandRouteSlack(slackData),
	}

	log.Printf("[DEBUG] ,%v", expandRouteSlack(slackData))

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
	d.Set("routing_regex", route.RoutingRegex)
	d.Set("position", route.Position)
	d.Set("slack", flattenRouteSlack(route.SlackRoute))

	return nil
}

func resourceRouteUpdate(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] update amixr route")
	client := m.(*amixr.Client)

	routingRegexData := d.Get("routing_regex").(string)
	positionData := d.Get("position").(int)
	slackData := d.Get("slack").([]interface{})

	updateOptions := &amixr.UpdateRouteOptions{
		RoutingRegex: routingRegexData,
		Position:     &positionData,
		ManualOrder:  true,
		Slack:        expandRouteSlack(slackData),
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

func flattenRouteSlack(input *amixr.SlackRoute) []map[string]interface{} {
	slack := make([]map[string]interface{}, 0, 1)
	out := make(map[string]interface{})
	if input.ChannelId != nil {
		out["channel_id"] = *input.ChannelId
	} else {
		out["channel_id"] = nil
	}

	slack = append(slack, out)
	return slack
}

func expandRouteSlack(input []interface{}) *amixr.SlackRoute {
	log.Printf("[DEBUG] expand slack route")
	log.Printf("[DEBUG] input %v", input)

	slackRoute := amixr.SlackRoute{}
	for _, r := range input {
		inputMap := r.(map[string]interface{})

		if inputMap["channel_id"] != nil {
			channelId := inputMap["channel_id"].(string)
			slackRoute.ChannelId = &channelId
		} else {
			slackRoute.ChannelId = nil
		}

	}

	log.Printf("[DEBUG] calculated from input %v", slackRoute)

	return &slackRoute

}
