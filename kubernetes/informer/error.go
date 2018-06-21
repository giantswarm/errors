package informer

import (
	"github.com/giantswarm/microerror"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var StatusForbiddenError = microerror.New("status forbidden")

func IsStatusForbidden(err error) bool {
	if err == nil {
		return false
	}

	c := microerror.Cause(err)
	if c == StatusForbiddenError {
		return true
	}

	statusErr, ok := c.(*errors.StatusError)
	if !ok {
		return false
	}

	if statusErr.Status().Reason == metav1.StatusReasonForbidden {
		return true
	}

	return false
}
