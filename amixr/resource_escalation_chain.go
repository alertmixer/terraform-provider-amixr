package amixr

import (
	amixr "github.com/alertmixer/amixr-go-client"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
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
			"is_default": &schema.Schema{
				Type:     schema.TypeBool,
				Computed: true,
			},
		},
	}
}

func resourceEscalationChainCreate(d *schema.ResourceData, m interface{}) error {
	client := m.(*amixr.Client)

	nameData := d.Get("name").(string)
	createOptions := &amixr.CreateEscalationChainOptions{
		Name: nameData,
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
	d.Set("is_default", escalationChain.IsDefault)

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
