package amixr

import (
	"fmt"
	amixr "github.com/alertmixer/amixr-go-client"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"log"
)

var onCallShiftTypeOptions = []string{
	"rolling_users",
	"recurrent_event",
	"single_event",
}

var onCallShiftFrequencyOptions = []string{
	"daily",
	"weekly",
	"monthly",
}

var onCallShiftWeekDayOptions = []string{
	"MO",
	"TU",
	"WE",
	"TH",
	"FR",
	"ST",
	"SU",
}

var sourceTerraform = 3

func resourceOnCallShift() *schema.Resource {
	return &schema.Resource{
		Create: resourceOnCallShiftCreate,
		Read:   resourceOnCallShiftRead,
		Update: resourceOnCallShiftUpdate,
		Delete: resourceOnCallShiftDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"schedule_id": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},
			"name": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},
			"level": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
			},
			"type": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringInSlice(onCallShiftTypeOptions, false),
			},
			"start": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},
			"duration": &schema.Schema{
				Type:         schema.TypeInt,
				Required:     true,
				ValidateFunc: validation.IntAtLeast(0),
			},
			"frequency": &schema.Schema{
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice(onCallShiftFrequencyOptions, false),
			},
			"users": {
				Type: schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Optional: true,
			},
			"rolling_users": {
				Type: schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeSet,
					Elem: &schema.Schema{
						Type: schema.TypeString,
					},
				},
				Optional: true,
			},
			"interval": &schema.Schema{
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: validation.IntAtLeast(1),
			},
			"week_start": &schema.Schema{
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice(onCallShiftWeekDayOptions, false),
			},
			"by_day": {
				Type: schema.TypeSet,
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: validation.StringInSlice(onCallShiftWeekDayOptions, false),
				},
				Optional: true,
			},
			"by_month": {
				Type: schema.TypeSet,
				Elem: &schema.Schema{
					Type:         schema.TypeInt,
					ValidateFunc: validation.IntBetween(1, 12),
				},
				Optional: true,
			},
			"by_monthday": {
				Type: schema.TypeSet,
				Elem: &schema.Schema{
					Type:         schema.TypeInt,
					ValidateFunc: validation.IntBetween(-31, 31),
				},
				Optional: true,
			},
		},
	}
}

func resourceOnCallShiftCreate(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] create amixr on-call shift")

	client := m.(*amixr.Client)

	scheduleIdData := d.Get("schedule_id").(string)
	typeData := d.Get("type").(string)
	nameData := d.Get("name").(string)
	startData := d.Get("start").(string)
	durationData := d.Get("duration").(int)

	createOptions := &amixr.CreateOnCallShiftOptions{
		ScheduleId: scheduleIdData,
		Type:       typeData,
		Name:       nameData,
		Start:      startData,
		Duration:   durationData,
		Source:     sourceTerraform,
	}

	levelData := d.Get("level").(int)
	createOptions.Level = &levelData

	frequencyData, frequencyOk := d.GetOk("frequency")
	if frequencyOk {
		if typeData != "single_event" {
			f := frequencyData.(string)
			createOptions.Frequency = &f
		} else {
			return fmt.Errorf("frequency can not be set with type: %s", typeData)
		}
	}

	usersData, usersDataOk := d.GetOk("users")
	if usersDataOk {
		createOptions.Users = stringSetToStringSlice(usersData.(*schema.Set))
	}

	intervalData, intervalOk := d.GetOk("interval")
	if intervalOk {
		if typeData != "single_event" {
			i := intervalData.(int)
			createOptions.Interval = &i
		} else {
			return fmt.Errorf("interval can not be set with type: %s", typeData)
		}
	}

	weekStartData, weekStartOk := d.GetOk("week_start")
	if weekStartOk {
		if typeData != "single_event" {
			w := weekStartData.(string)
			createOptions.WeekStart = &w
		} else {
			return fmt.Errorf("week_start can not be set with type: %s", typeData)
		}
	}

	byDayData, byDayOk := d.GetOk("by_day")
	if byDayOk {
		if typeData != "single_event" {
			createOptions.ByDay = stringSetToStringSlice(byDayData.(*schema.Set))
		} else {
			return fmt.Errorf("by_day can not be set with type: %s", typeData)
		}
	}

	byMonthData, byMonthOk := d.GetOk("by_month")
	if byMonthOk {
		if typeData != "single_event" {
			createOptions.ByMonth = intSetToIntSlice(byMonthData.(*schema.Set))
		} else {
			return fmt.Errorf("by_month can not be set with type: %s", typeData)
		}
	}

	byMonthdayData, byMonthdayOk := d.GetOk("by_monthday")
	if byMonthdayOk {
		if typeData != "single_event" {
			createOptions.ByMonthday = intSetToIntSlice(byMonthdayData.(*schema.Set))
		} else {
			return fmt.Errorf("by_monthday can not be set with type: %s", typeData)
		}
	}

	onCallShift, _, err := client.OnCallShifts.CreateOnCallShift(createOptions)
	if err != nil {
		return err
	}

	d.SetId(onCallShift.ID)

	return resourceOnCallShiftRead(d, m)
}

func resourceOnCallShiftUpdate(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] update amixr on-call shift")

	client := m.(*amixr.Client)

	typeData := d.Get("type").(string)
	nameData := d.Get("name").(string)
	startData := d.Get("start").(string)
	durationData := d.Get("duration").(int)

	updateOptions := &amixr.UpdateOnCallShiftOptions{
		Type:     typeData,
		Name:     nameData,
		Start:    startData,
		Duration: durationData,
		Source:   sourceTerraform,
	}

	levelData := d.Get("level").(int)
	updateOptions.Level = &levelData

	frequencyData, frequencyOk := d.GetOk("frequency")
	if frequencyOk {
		if typeData != "single_event" {
			f := frequencyData.(string)
			updateOptions.Frequency = &f
		} else {
			return fmt.Errorf("frequency can not be set with type: %s", typeData)
		}
	}

	usersData, usersDataOk := d.GetOk("users")
	if usersDataOk {
		updateOptions.Users = stringSetToStringSlice(usersData.(*schema.Set))
	}

	intervalData, intervalOk := d.GetOk("interval")
	if intervalOk {
		if typeData != "single_event" {
			i := intervalData.(int)
			updateOptions.Interval = &i
		} else {
			return fmt.Errorf("interval can not be set with type: %s", typeData)
		}
	}

	weekStartData, weekStartOk := d.GetOk("week_start")
	if weekStartOk {
		if typeData != "single_event" {
			w := weekStartData.(string)
			updateOptions.WeekStart = &w
		} else {
			return fmt.Errorf("week_start can not be set with type: %s", typeData)
		}
	}

	byDayData, byDayOk := d.GetOk("by_day")
	if byDayOk {
		if typeData != "single_event" {
			updateOptions.ByDay = stringSetToStringSlice(byDayData.(*schema.Set))
		} else {
			return fmt.Errorf("by_day can not be set with type: %s", typeData)
		}
	}

	byMonthData, byMonthOk := d.GetOk("by_month")
	if byMonthOk {
		if typeData != "single_event" {
			updateOptions.ByMonth = intSetToIntSlice(byMonthData.(*schema.Set))
		} else {
			return fmt.Errorf("by_month can not be set with type: %s", typeData)
		}
	}

	byMonthDayData, byMonthDayOk := d.GetOk("by_monthday")
	if byMonthDayOk {
		if typeData != "single_event" {
			updateOptions.ByMonthday = intSetToIntSlice(byMonthDayData.(*schema.Set))
		} else {
			return fmt.Errorf("by_monthday can not be set with type: %s", typeData)
		}
	}

	onCallShift, _, err := client.OnCallShifts.UpdateOnCallShift(d.Id(), updateOptions)
	if err != nil {
		return err
	}

	d.SetId(onCallShift.ID)

	return resourceOnCallShiftRead(d, m)
}

func resourceOnCallShiftRead(d *schema.ResourceData, m interface{}) error {

	client := m.(*amixr.Client)
	options := &amixr.GetOnCallShiftOptions{}
	onCallShift, _, err := client.OnCallShifts.GetOnCallShift(d.Id(), options)

	if err != nil {
		return err
	}

	d.Set("schedule_id", onCallShift.ScheduleId)
	d.Set("name", onCallShift.Name)
	d.Set("type", onCallShift.Type)
	d.Set("level", onCallShift.Level)
	d.Set("start", onCallShift.Start)
	d.Set("duration", onCallShift.Duration)
	d.Set("frequency", onCallShift.Frequency)
	d.Set("users", onCallShift.Users)
	d.Set("start", onCallShift.Start)

	return nil
}

func resourceOnCallShiftDelete(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] delete amixr on-call shift")

	client := m.(*amixr.Client)
	options := &amixr.DeleteOnCallShiftOptions{}
	_, err := client.OnCallShifts.DeleteOnCallShift(d.Id(), options)
	if err != nil {
		return err
	}

	d.SetId("")

	return nil
}
