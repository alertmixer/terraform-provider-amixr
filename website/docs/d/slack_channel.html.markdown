---
layout: "amixr"
page_title: "Amixr: amixr_slack_channel"
sidebar_current: "docs-amixr-datasource-slack_channel"
description: |-
  Get information about a Slack channel.
---

# amixr\_slack_channel

Get information about a specific [Slack channel](https://api-docs.amixr.io/#slack-channels).

## Example Usage

```hcl
data "amixr_slack_channel" "example_slack_channel" {
  name = "example_slack_channel"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The Slack channel name.

## Attributes Reference

* `slack_id` - The Slack ID of the found Slack channel.
