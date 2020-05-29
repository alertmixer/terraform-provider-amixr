---
layout: "amixr"
page_title: "Provider: Amixr"
sidebar_current: "docs-amixr-index"
description: |-
  Amixr is an incident management platform
---

# Amixr Provider

[Amixr](https://amixr.io/) is next-gen incident management platform for DevOps and SRE. The platform allows you to optimize channels, recipients, and content in order to increase the speed of solving IT problems.

Use the navigation to the left to read about the available resources.

## Example Usage

```hcl
# Configure the Amixr provider
provider "amixr" {
  token = var.amixr_token
}

# Create an Amixr integration
resource "amixr_integration" "grafana-integration" {
  name = "Grafana Integration"
  type = "grafana"
}

# Create an Amixr escalation policy step
resource "amixr_escalation" "wait_step" {
  route_id = amixr_integration.grafana-integration.default_route_id
  type = "wait"
  duration = 60
  position = 0
}
```

## Argument Reference

The following arguments are supported in the `provider` block:

* `token` - (Required) The authorization token. It can also be sourced from the `AMIXR_API_KEY` environment variable. See [API Documentation](https://api-docs.amixr.io/#authentication) for more information.
