---
layout: "amixr"
page_title: "Amixr: amixr_action"
sidebar_current: "docs-amixr-datasource-action"
description: |-
  Get information about an action (outgoing webhook).
---

# amixr\_action

Get information about a specific [action](https://api-docs.amixr.io/#actions-outgoing-webhooks) (outgoing webhook).

## Example Usage

```hcl
resource "amixr_integration" "example_integration" {
  name      = "Grafana Integration"
  type      = "grafana"
}

data "amixr_action" "example_action" {
  name = "test"
  integration_id = amixr_integration.example_integration.id
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The action name.
* `integration_id` - (Required) The ID of the integration.

## Attributes Reference

* `id` - The ID of the found action.
