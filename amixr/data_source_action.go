package amixr

import (
	"fmt"
	amixr "github.com/alertmixer/amixr-go-client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"log"
)

func dataSourceAmixrAction() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAmixrActionRead,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func dataSourceAmixrActionRead(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] read amixr action")

	client := m.(*amixr.Client)
	options := &amixr.ListCustomActionOptions{}
	nameData := d.Get("name").(string)

	options.Name = nameData

	customActionsResponse, _, err := client.CustomActions.ListCustomActions(options)

	if err != nil {
		return err
	}

	if len(customActionsResponse.CustomActions) == 0 {
		return fmt.Errorf("couldn't find an action matching: %s", options.Name)
	} else if len(customActionsResponse.CustomActions) != 1 {
		return fmt.Errorf("more than one action found matching: %s", options.Name)
	}

	custom_action := customActionsResponse.CustomActions[0]

	d.SetId(custom_action.ID)
	d.Set("name", custom_action.Name)

	return nil
}
