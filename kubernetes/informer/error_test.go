package informer

import (
	"testing"

	"github.com/giantswarm/microerror"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

func Test_IsStatusForbidden(t *testing.T) {
	testCases := []struct {
		name          string
		err           error
		expectedMatch bool
	}{
		{
			name:          "case 0: match StatusForbiddenError",
			err:           StatusForbiddenError,
			expectedMatch: true,
		},
		{
			name:          "case 1: match apimachinery StatusError with Forbidden reason",
			err:           errors.NewForbidden(schema.GroupResource{}, "unittest", nil),
			expectedMatch: true,
		},
		{
			name:          "case 2: don't match nil",
			err:           nil,
			expectedMatch: false,
		},
		{
			name:          "case 3: don't match unknown error",
			err:           microerror.New("unknown error"),
			expectedMatch: false,
		},
		{
			name:          "case 4: don't match apimachinery StatusError with Unauthorized reason",
			err:           errors.NewUnauthorized(""),
			expectedMatch: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			match := IsStatusForbidden(tc.err)
			if match != tc.expectedMatch {
				t.Fatalf("match == %v, expected %v", match, tc.expectedMatch)
			}
		})
	}
}
