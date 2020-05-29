---
layout: "amixr"
page_title: "Amixr: amixr_integration"
sidebar_current: "docs-amixr-resource-integration"
description: |-
  Creates and manages an integration in Amixr.
---

# amixr\_integration

[Integrations](https://api-docs.amixr.io/#integrations) are sources of alerts and incidents for Amixr.

## Example Usage

```hcl
resource "amixr_integration" "example" {
  name    = "Grafana Integration"
  type    = "grafana"
}
```

## Argument Reference

The following arguments are supported:

  * `name` - (Optional) The name of the service integration.
  * `type` - (Required) The type of integration. Can be:
  `grafana`, `webhook`, `alertmanager`, `kapacitor` , `fabric`,
  `newrelic`, `datadog`, `pagerduty`, `pingdom`,
  `elastalert`, `amazon_sns`, `curler`, `sentry`,
  `formatted_webhook`, `heartbeat`, `demo`, `stackdriver`,
  `uptimerobot`, `sentry_platform`, `zabbix`, `prtg`
   or `inbound_email`.

## Attributes Reference

The following attributes are exported:

  * `id` - The ID of the integration.
