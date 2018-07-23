package guest

import (
	"regexp"

	"github.com/giantswarm/microerror"
)

var (
	APINotAvailablePatterns = []*regexp.Regexp{
		regexp.MustCompile(`dial tcp: lookup .* on .*:53: (no such host|server misbehaving)`),
		regexp.MustCompile(`Get https://api\..*/api/v1/nodes.* (unexpected )?EOF`),
		regexp.MustCompile(`[Get|Post] https://api\..*/api/v1/namespaces/*/.* (unexpected )?EOF`),
		regexp.MustCompile(`Get https://api\..*/api/v1/nodes.* net/http: (TLS handshake timeout|request canceled while waiting for connection).*?`),
		regexp.MustCompile(`[Get|Post] https://api\..*: x509: certificate is valid for ingress.local, not api\..*`),
	}
)

// APINotAvailableError is returned when the guest Kubernetes API is not
// available.
var APINotAvailableError = microerror.New("API not available")

// IsAPINotAvailable asserts APINotAvailableError.
func IsAPINotAvailable(err error) bool {
	if err == nil {
		return false
	}

	c := microerror.Cause(err)

	for _, re := range APINotAvailablePatterns {
		matched := re.MatchString(c.Error())

		if matched {
			return true
		}
	}

	if c == APINotAvailableError {
		return true
	}

	return false
}
