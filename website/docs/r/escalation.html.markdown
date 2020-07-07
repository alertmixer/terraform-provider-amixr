---
layout: "amixr"
page_title: "Amixr: amixr_escalation"
sidebar_current: "docs-amixr-resource-escalation"
description: |-
  Configure an escalation policies in Amixr.
---

# amixr\_escalation

[Escalation policy](https://api-docs.amixr.io/#escalation-policies) configures what happened after incident is triggerd: who will be notified first, second, etc., and delay before notifications. 

## Example Usage

```hcl
data "amixr_user" "alex" {
  email = "alex@example.com"
}

resource "amixr_integration" "grafana-integration" {
  name = "Grafana Integration"
  type = "grafana"
}

resource "amixr_escalation" "notify_step" {
  route_id = amixr_integration.grafana-integration.default_route_id
  type = "notify_persons"
  persons_to_notify = [
    data.amixr_user.alex.id
  ]
  position = 0
}

resource "amixr_escalation" "wait_step" {
  route_id = amixr_integration.grafana-integration.default_route_id
  type = "wait"
  duration = 60
  position = 1
}
```

## Argument Reference

The following arguments are supported:

  * `route_id` - (Required) The ID of the route.
  * `position` - (Required) The position of the escalation step (starts from 0)
  * `type` - (Required) The type of escalation policy. Can be:
    - `wait` - just wait
    - `notify_persons` - notify multiple users at the same time
    - `notify_person_next_each_time` - notify one user from queue
    - `notify_on_call_from_schedule` - notify by schedule
  * `duration` - (Optional) The duration of delay for `wait` type step.
  * `persons_to_notify` - (Optional) The list of ID's of users for `notify_persons` type step.
  * `persons_to_notify_next_each_time` - (Optional) The list of ID's of users for `notify_person_next_each_time` type step.
  * `notify_on_call_from_schedule` - (Optional) ID of a Schedule for `notify_on_call_from_schedule` type step.


## Attributes Reference

The following attributes are exported:

  * `id` - The ID of the escalation policy.


## Import

Existing escalations can be imported using the escalation ID:

```sh
$ terraform import amixr_escalation.example_escalation {escalation ID}
```
