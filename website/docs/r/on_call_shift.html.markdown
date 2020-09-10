---
layout: "amixr"
page_title: "Amixr: amixr_on_call_shift"
sidebar_current: "docs-amixr-resource-on_call_shift"
description: |-
  Creates and manages schedule on-call shifts in Amixr.
---

# amixr\_on_call_shifts

Manage an [on-call shift](https://api-docs.amixr.io/#on-call-shifts).

## Example Usage

```hcl
data "amixr_user" "alex" {
  email = "alex@example.com"
}

resource "amixr_schedule" "example_schedule" {
  name      = "Example Schedule"
  time_zone = "America/New_York"
}

resource "amixr_on_call_shift" "example_shift" {
  name = "Example Shift"
  schedule_id = amixr_schedule.example_schedule.id
  type = "recurrent_event"
  start = "2020-09-07T14:00:00"
  duration = 18000
  frequency = "weekly"
  interval = 2
  by_day = ["MO", "FR"]
  week_start = "MO"
  users = [
    data.amixr_user.alex.id
  ]
}

```

## Argument Reference

The following arguments are supported:

  * `name` - (Required) The shift's name
  * `schedule_id` - (Required) The ID of the schedule
  * `type` - (Required) The shift's type. Can be:
    - `single_event` - the event will be triggered once and does not repeat
    - `recurrent_event` - the event will be repeated in accordance with the recurrence rules
  * `start` - (Required) The start time of the on-call shift. This parameter takes a date format as `yyyy-MM-dd'T'HH:mm:ss` (for example "2020-09-05T08:00:00")
  * `duration` - (Required) The duration of the event
  * `frequency` - (Required for recurrent events) The frequency of the event. Can be: `daily`, `weekly`, `monthly`
  * `interval` - (Optional) This parameter takes a positive integer representing at which intervals the recurrence rule repeats
  * `week_start` - (Optional) Start day of the week in iCal format. Can be: `SU` (Sunday), `MO` (Monday), `TU` (Tuesday), `WE` (Wednesday), `TH` (Thursday), `FR` (Friday), `SA` (Saturday). Default: `SU`
  * `by_day` - (Optional) This parameter takes a list of days in iCal format (`SU`, `MO`, `TU`, `WE`, `TH`, `FR`, `SA`)
  * `by_month` - (Optional) This parameter takes a list of months. Valid values are 1 to 12
  * `by_monthday` - (Optional) This parameter takes a list of days of the month. Valid values are 1 to 31 or -31 to -1
  * `users` - (Optional) The list of on-call users
  
Please look [RFC 5545](https://tools.ietf.org/html/rfc5545#section-3.3.10) for more information about the recurrence rules.

## Attributes Reference

The following attributes are exported:

  * `id` - The ID of the on-call shift.
  

## Import

Existing on-call shift can be imported using the shift ID:

```sh
$ terraform import amixr_on_call_shift.example_shift {on-call shift ID}
```

