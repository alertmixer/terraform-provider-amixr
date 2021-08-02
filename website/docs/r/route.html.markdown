---
layout: "amixr"
page_title: "Amixr: amixr_route"
sidebar_current: "docs-amixr-resource-route"
description: |-
  Creates and manages integration routes in Amixr.
---

# amixr\_route

[Routes](https://api-docs.amixr.io/#routes) allow to direct different alerts to different messenger channels and escalation chains.

## Example Usage

```hcl
data "amixr_slack_channel" "example_slack_channel" {
  name = "example_slack_channel"
}

data "amixr_escalation_chain" "default" {
    name = "default"
}

resource "amixr_integration" "example_integration" {
  name    = "Grafana Integration"
  type    = "grafana"
}

resource "amixr_route" "example_route"{ 
  integration_id = amixr_integration.example_integration.id
  escalation_chain_id = data.amixr_escalation_chain.default.id
  routing_regex  = "us-(east|west)"
  position       = 0
  slack {
      channel_id = data.amixr_slack_channel.example_slack_channel.slack_id
  }
}
```

## Argument Reference

The following arguments are supported:

  * `integration_id` - (Required) The ID of the integration.
  * `escalation_chain_id` - (Required) The ID of the escalation chain.
  * `routing_regex` - (Required) Python Regex query. Amixr choose the route for an alert in case there is a match inside the whole alert payload.
  * `position` - (Required) The position of the route (starts from 0)
  * `slack` - (Optional) Dictionary with slack-specific settings for a route. Includes:
    - `channel_id` - Slack channel id. Alerts will be directed to this channel in Slack.


## Attributes Reference

The following attributes are exported:

  * `id` - The ID of the route.
  

## Import

Existing routes can be imported using the route ID:

```sh
$ terraform import amixr_route.example_route {route ID}
```

