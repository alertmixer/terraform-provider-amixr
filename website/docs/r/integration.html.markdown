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
  name      = "Grafana Integration"
  type      = "grafana"
  templates {
      grouping_key = "custom uuid"
      resolve_signal = "{{ 1 if payload.resolved == 'ok' else 0 }}"
      slack {
          title = "Custom title {{ payload.title }}"
          message = "Custom message {{ payload.message }}"
          image_url = "http://example.com/custom_image.jpg"
    }
  }
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
  * `templates` - (Optional) Jinja2 templates for Alert payload. Includes:
    - `grouping_key`- (Optional) Template for the key by which alerts are grouped.
    - `resolve_signal`- (Optional) Template for sending a signal to resolve the Incident. This template should output one of the following values: ok, true, 1 (case insensitive)
    - `slack`- (Optional) Templates for Slack:
        - `title`- (Optional) Template for Alert title.
        - `message`- (Optional) Template for Alert message.
        - `image_url`- (Optional) Template for Alert image url.

    To set a parameter to default value just remove it or set the value to `null`.

## Attributes Reference

The following attributes are exported:

  * `id` - The ID of the integration.
  

## Import

Existing integrations can be imported using the integration ID:

```sh
$ terraform import amixr_integration.example_integration {integration ID}
```

