---
description: Manage secrets with `goops` charms.
---

# How-to manage secrets

Both Juju users and charms can manage secrets. Here we cover how charms can read and write secrets using `goops`.

`goops` provides functions to manage secrets, allowing you to read, create, grant and revoke secrets. Those are the same functions that Juju exposes through hook commands.

```go
package charm

import (
	"fmt"
	"time"

	"github.com/gruyaume/goops"
)

const (
	CaCertificateSecretLabel = "active-ca-certificates"
)

func Configure() error {
	err := generateAndStoreRootCertificate()
	if err != nil {
		return fmt.Errorf("could not generate and store root certificate: %w", err)
	}

	return nil
}

func generateAndStoreRootCertificate() error {
	_, err := goops.GetSecretByLabel(CaCertificateSecretLabel, false, true)
	if err != nil {
		goops.LogInfof("could not get secret: %s", err.Error())

		secretContent := map[string]string{
			"private-key":    "Example private key",
			"ca-certificate": "Example CA certificate",
		}

		expiry := time.Now().AddDate(1, 0, 0)

		output, err := goops.AddSecret(&goops.AddSecretOptions{
			Content:     secretContent,
			Description: "ca certificate and private key for the certificates charm",
			Expire:      expiry,
			Label:       CaCertificateSecretLabel,
			Rotate:      goops.RotateNever,
			Owner:       goops.OwnerApplication,
		})
		if err != nil {
			return fmt.Errorf("could not add secret: %w", err)
		}

		goops.LogInfof("Created new secret with ID: %s", output)

		return nil
	}

	secretInfo, err := goops.GetSecretInfoByLabel(CaCertificateSecretLabel)
	if err != nil {
		return fmt.Errorf("could not get secret info: %w", err)
	}

	if secretInfo == nil {
		return fmt.Errorf("secret info is nil")
	}

	return nil
}
```

!!! info
    Learn more about secret management in charms:

    - [Juju Hook commands :octicons-link-external-24:](https://documentation.ubuntu.com/juju/3.6/reference/hook-command/list-of-hook-commands/)
    - [goops API reference :octicons-link-external-24:](https://pkg.go.dev/github.com/gruyaume/goops)
