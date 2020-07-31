package amixr

import (
	"fmt"
	amixr "github.com/alertmixer/amixr-go-client"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"log"
)

func dataSourceAmixrCustomAction() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAmixrCustomActionRead,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"integration_id": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func dataSourceAmixrCustomActionRead(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] read amixr custom action")

	client := m.(*amixr.Client)
	options := &amixr.ListCustomActionOptions{}
	nameData := d.Get("name").(string)
	integrationIdData := d.Get("integration_id").(string)

	options.Name = nameData
	options.IntegrationId = integrationIdData

	customActionsResponse, _, err := client.CustomActions.ListCustomActions(options)

	if err != nil {
		return err
	}

	if len(customActionsResponse.CustomActions) == 0 {
		return fmt.Errorf("couldn't find a custom action matching: %s %s", options.Name, options.IntegrationId)
	} else if len(customActionsResponse.CustomActions) != 1 {
		return fmt.Errorf("more than one custom found matching: %s %s", options.Name, options.IntegrationId)
	}

	custom_action := customActionsResponse.CustomActions[0]

	d.SetId(custom_action.ID)
	d.Set("name", custom_action.Name)
	d.Set("integration_id", custom_action.IntegrationId)

	return nil
}
