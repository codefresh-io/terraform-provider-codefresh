---
layout: "codefresh"
page_title: "Provider: Codefresh"
sidebar_current: "docs-codefresh-index"
description: |-
  The Codefresh provider is used to manage Codefresh resources.
---

# Codefresh Provider

The Codefresh Provider can be used to configure [Codefresh](https://codefresh.io/) resources - pipelines, projects, accounts, etc using the [Codefresh API](https://codefresh.io/docs/docs/integrations/codefresh-api/).

## Authenticating to Codefresh

The Codefresh API requires the [authentication key](https://codefresh.io/docs/docs/integrations/codefresh-api/#authentication-instructions) to authenticate.
The key can be passed either as provider's attribute or as environment variable - `CODEFRESH_API_KEY`.

## Argument Reference

The following arguments are supported:

- `token` - (Optional) The client API token. This can also be sourced from the `CODEFRESH_API_KEY` environment variable.
- `api_url` -(Optional) Default value - https://g.codefresh.io/api.

## Recommendation for creation Accounts, Users, Teams, Permissions
* create users and accounts using [accounts_users module](modules/accounts_users.md) and Codefresh Admin token 
* Create and save in tf state api_keys using [accounts_token module](modules/accounts_token.md)
* Create teams using [teams module](modules/teams.md)
* Create permissions - [see example](../examples/permissions)

