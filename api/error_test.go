package api

import (
	"errors"
	"testing"
)

func Test_IsGuestAPINotAvailable(t *testing.T) {
	testCases := []struct {
		description   string
		errorMessage  string
		expectedMatch bool
	}{
		{
			description:   "dns not ready",
			errorMessage:  "dial tcp: lookup api.5xchu.aws.gigantic.io on 10.96.0.10:53: no such host",
			expectedMatch: true,
		},
		{
			description:   "dns not ready incorrect port",
			errorMessage:  "dial tcp: lookup api.5xchu.aws.gigantic.io on 10.96.0.10:443: no such host",
			expectedMatch: false,
		},
		{
			description:   "ingress not ready get request",
			errorMessage:  "Get https://api.5xchu.aws.gigantic.io: x509: certificate is valid for ingress.local, not api.5xchu.aws.gigantic.io:",
			expectedMatch: true,
		},
		{
			description:   "API not ready get EOF request",
			errorMessage:  "Get https://api.5xchu.aws.gigantic.io/api/v1/nodes: EOF",
			expectedMatch: true,
		},
		{
			description:   "ingress not ready post request",
			errorMessage:  "Post https://api.5xchu.aws.gigantic.io: x509: certificate is valid for ingress.local, not api.5xchu.aws.gigantic.io:",
			expectedMatch: true,
		},
		{
			description:   "ingress not ready post different domain",
			errorMessage:  "Post https://api.5xchu.aws.gigantic.io: x509: certificate is valid for localhost, not api.5xchu.aws.gigantic.io:",
			expectedMatch: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			err := errors.New(tc.errorMessage)
			result := IsGuestAPINotAvailable(err)

			if result != tc.expectedMatch {
				t.Fatalf("expected %t, got %t", tc.expectedMatch, result)
			}
		})
	}
}
