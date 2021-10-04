package amixr

import (
	"github.com/alertmixer/amixr-go-client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceEscalationChain() *schema.Resource {
	return &schema.Resource{
		Create: resourceEscalationChainCreate,
		Read:   handleNonExistentResource(resourceEscalationChainRead),
		Update: resourceEscalationChainUpdate,
		Delete: resourceEscalationChainDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"team_id": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func resourceEscalationChainCreate(d *schema.ResourceData, m interface{}) error {
	client := m.(*amixr.Client)

	nameData := d.Get("name").(string)
	teamIdData := d.Get("team_id").(string)
	
	createOptions := &amixr.CreateEscalationChainOptions{
		Name: nameData,
		TeamId: teamIdData,
	}

	escalationChain, _, err := client.EscalationChains.CreateEscalationChain(createOptions)
	if err != nil {
		return err
	}

	d.SetId(escalationChain.ID)

	return resourceEscalationChainRead(d, m)
}

func resourceEscalationChainRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*amixr.Client)

	escalationChain, _, err := client.EscalationChains.GetEscalationChain(d.Id(), &amixr.GetEscalationChainOptions{})
	if err != nil {
		return err
	}

	d.Set("name", escalationChain.Name)
	d.Set("team_id", escalationChain.TeamId)

	return nil
}

func resourceEscalationChainUpdate(d *schema.ResourceData, m interface{}) error {
	client := m.(*amixr.Client)

	nameData := d.Get("name").(string)

	updateOptions := &amixr.UpdateEscalationChainOptions{
		Name: nameData,
	}

	escalationChain, _, err := client.EscalationChains.UpdateEscalationChain(d.Id(), updateOptions)
	if err != nil {
		return err
	}

	d.SetId(escalationChain.ID)
	return resourceEscalationChainRead(d, m)
}

func resourceEscalationChainDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(*amixr.Client)

	_, err := client.EscalationChains.DeleteEscalationChain(d.Id(), &amixr.DeleteEscalationChainOptions{})
	if err != nil {
		return err
	}

	d.SetId("")

	return nil
}
