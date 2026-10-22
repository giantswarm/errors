[![CircleCI](https://circleci.com/gh/giantswarm/errors.svg?&style=shield)](https://circleci.com/gh/giantswarm/errors)

# errors

This module contains error matchers needed by multiple projects, in
subpackages such as `tenant`.

## Usage

`tenant.IsAPINotAvailable` reports whether an error looks like the
Kubernetes API of a workload cluster is not available, for example DNS,
EOF, TLS or timeout errors while the API is coming up, or the
`tenant.APINotAvailableError` sentinel. It matches on the error message;
see `tenant/error.go` for the patterns.

```go
import "github.com/giantswarm/errors/tenant"

if tenant.IsAPINotAvailable(err) {
	// The API is not available yet; retry later.
}
```
