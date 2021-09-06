package amixr

import (
	"fmt"
	amixr "github.com/alertmixer/amixr-go-client"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceEscalationChain() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceEscalationChainRead,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"is_default": {
				Type:     schema.TypeBool,
				Computed: true,
			},
		},
	}
}

func dataSourceEscalationChainRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*amixr.Client)
	options := &amixr.ListEscalationChainOptions{}
	nameData := d.Get("name").(string)

	options.Name = nameData

	escalationChainsResponse, _, err := client.EscalationChains.ListEscalationChains(options)

	if err != nil {
		return err
	}

	if len(escalationChainsResponse.EscalationChains) == 0 {
		return fmt.Errorf("couldn't find an escalation chain matching: %s", options.Name)
	} else if len(escalationChainsResponse.EscalationChains) != 1 {
		return fmt.Errorf("more than one escalation chain found matching: %s", options.Name)
	}

	escalationChain := escalationChainsResponse.EscalationChains[0]

	d.Set("name", escalationChain.Name)
	d.Set("is_default", escalationChain.IsDefault)

	d.SetId(escalationChain.ID)

	return nil
}
