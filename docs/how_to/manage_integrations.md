---
description: Manage integrations with `goops` charms.
---

# Manage integrations

## Directly

`goops` provides functions to manage relations, allowing you to get relation IDs, list relation units, and set or get relation data. Those are the same functions that Juju exposes through hook commands.

### 1. Declare the relation endpoint

To integrate with another charm, declare the relations in your charm’s `charmcraft.yaml` file.

To exchange data with another charm, define a provides or requires endpoint including an interface name. By convention, the interface name should be unique in the ecosystem. Each relation must have an endpoint, which your charm will use to refer to the relation.

For example, to declare a relation with a PostgreSQL database, you can add the following to your `charmcraft.yaml`:

```yaml
requires:
  db:
    interface: postgresql_client
    limit: 1
```

### 2. Implement the relation logic in your charm

Use `goops` to read the relation data:

```go
package charm

import (
	"fmt"

	"github.com/gruyaume/goops"
)

func GetDatabaseURL(relationName string) (string, error) {
	relationIDs, err := goops.GetRelationIDs(relationName)
	if err != nil {
		return "", fmt.Errorf("could not get relation IDs: %w", err)
	}

	if len(relationIDs) == 0 {
		return "", fmt.Errorf("no relation IDs found for %s", relationName)
	}

	relationID := relationIDs[0]

	relationUnits, err := goops.ListRelationUnits(relationID)
	if err != nil {
		return "", fmt.Errorf("could not get relation list: %w", err)
	}

	if len(relationUnits) == 0 {
		return "", fmt.Errorf("no relation units found for ID: %s", relationID)
	}

	relationData, err := goops.GetAppRelationData(relationID, relationUnits[0])
	if err != nil {
		goops.LogDebugf("Could not get relation data: %s", err.Error())
		return "", fmt.Errorf("could not get relation data for ID %s: %w", relationID, err)
	}

	endpoints, ok := relationData["endpoints"]
	if !ok {
		return "", fmt.Errorf("no endpoints found in relation data for ID %s", relationID)
	}

	return endpoints, nil
}
```

!!! info
    Learn more about relation management in charms:

    - [Juju Hook commands :octicons-link-external-24:](https://documentation.ubuntu.com/juju/3.6/reference/hook-command/list-of-hook-commands/)
    - [goops API reference :octicons-link-external-24:](https://pkg.go.dev/github.com/gruyaume/goops)

## Using Charm Libraries

In most cases, charms should not directly read and write to relation data. Instead, they should do so indirectly using [Charm Libraries](../reference/charm_libraries.md), which encapsulate the relation logic.

```go
package charm

import (
	"github.com/gruyaume/charm-libraries/postgresql"
)

func GetDatabaseURL(relationName string) (string, error) {
	i := &postgresql.Integration{
		RelationName: relationName,
	}

	return i.GetDatabaseURL()
}
```
