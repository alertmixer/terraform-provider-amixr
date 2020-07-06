package amixr

import (
	"fmt"
	amixr "github.com/alertmixer/amixr-go-client"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"log"
)

func dataSourceAmixrSlackChannel() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAmixrSlackChannelRead,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"slack_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceAmixrSlackChannelRead(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] read amixr slack_channel")

	client := m.(*amixr.Client)
	options := &amixr.ListSlackChannelOptions{}
	nameData := d.Get("name").(string)

	options.ChannelName = nameData

	slackChannelsResponse, _, err := client.SlackChannels.ListSlackChannels(options)

	if err != nil {
		return err
	}

	if len(slackChannelsResponse.SlackChannels) == 0 {
		return fmt.Errorf("couldn't find a slack_channel matching: %s", options.ChannelName)
	} else if len(slackChannelsResponse.SlackChannels) != 1 {
		return fmt.Errorf("more than one slack_channel found matching: %s", options.ChannelName)
	}

	slack_channel := slackChannelsResponse.SlackChannels[0]

	d.SetId(slack_channel.SlackId)
	d.Set("name", slack_channel.Name)
	d.Set("slack_id", slack_channel.SlackId)

	return nil
}
