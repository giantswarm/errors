package tenant

import (
	"errors"
	"testing"
)

func Test_IsAPINotAvailable(t *testing.T) {
	testCases := []struct {
		description   string
		errorMessage  string
		expectedMatch bool
	}{
		{
			description:   "case 1: dns not ready",
			errorMessage:  "dial tcp: lookup api.5xchu.aws.gigantic.io on 10.96.0.10:53: no such host",
			expectedMatch: true,
		},
		{
			description:   "case 2: dns not ready incorrect port",
			errorMessage:  "dial tcp: lookup api.5xchu.aws.gigantic.io on 10.96.0.10:443: no such host",
			expectedMatch: false,
		},
		{
			description:   "case 3: ingress not ready get request",
			errorMessage:  "Get https://api.5xchu.aws.gigantic.io: x509: certificate is valid for ingress.local, not api.5xchu.aws.gigantic.io:",
			expectedMatch: true,
		},
		{
			description:   "case 4: API not ready get EOF request",
			errorMessage:  "Get https://api.5xchu.aws.gigantic.io/api/v1/nodes: EOF",
			expectedMatch: true,
		},
		{
			description:   "case 5: temporary issues with the master node serving the tenant cluster API",
			errorMessage:  "Get https://api.8dnxs.g8s.gorgoth.gridscale.kvm.gigantic.io/api/v1/nodes: unexpected EOF",
			expectedMatch: true,
		},
		{
			description:   "case 6: temporary issues with the master node serving the tenant cluster API",
			errorMessage:  "Get https://api.uth29.g8s.geckon.gridscale.kvm.gigantic.io/api/v1/nodes?timeout=30s: EOF",
			expectedMatch: true,
		},
		{
			description:   "case 7: ingress not ready post request",
			errorMessage:  "Post https://api.5xchu.aws.gigantic.io: x509: certificate is valid for ingress.local, not api.5xchu.aws.gigantic.io:",
			expectedMatch: true,
		},
		{
			description:   "case 8: ingress not ready post different domain",
			errorMessage:  "Post https://api.5xchu.aws.gigantic.io: x509: certificate is valid for localhost, not api.5xchu.aws.gigantic.io:",
			expectedMatch: false,
		},
		{
			description:   "case 9: timeout getting namespace",
			errorMessage:  "Get https://api.3jwh2.k8s.aws.gigantic.io/api/v1/namespaces/giantswarm?timeout=30s: EOF",
			expectedMatch: true,
		},
		{
			description:   "case 10: timeout getting service account",
			errorMessage:  "Post https://api.3jwh2.k8s.aws.gigantic.io/api/v1/namespaces/giantswarm/serviceaccounts?timeout=30s: EOF",
			expectedMatch: true,
		},
		{
			description:   "case 11: GET timeout establishing TLS handshake",
			errorMessage:  "Get https://api.08vka.k8s.gorgoth.gridscale.kvm.gigantic.io/api/v1/nodes?timeout=30s: net/http: TLS handshake timeout",
			expectedMatch: true,
		},
		{
			description:   "case 12: server is misbehaving due to TCP lookup",
			errorMessage:  "Get https://api.ci-wip-70f9b-5e958.k8s.godsmack.westeurope.azure.gigantic.io/api/v1/nodes?timeout=30s: dial tcp: lookup api.ci-wip-70f9b-5e958.k8s.godsmack.westeurope.azure.gigantic.io on 10.96.0.10:53: server misbehaving",
			expectedMatch: true,
		},
		{
			description:   "case 13: request canceled while waiting for connection",
			errorMessage:  "Get https://api.ci-wip-2317d-c1c86.k8s.godsmack.westeurope.azure.gigantic.io/api/v1/nodes?timeout=30s: net/http: request canceled while waiting for connection (Client.Timeout exceeded while awaiting headers)",
			expectedMatch: true,
		},
		{
			description:   "case 14: request canceled due to client timeout exceeded",
			errorMessage:  "Get https://api.06bhh.k8s.geckon.gridscale.kvm.gigantic.io/api/v1/nodes?timeout=30s: net/http: request canceled (Client.Timeout exceeded while awaiting headers)",
			expectedMatch: true,
		},
		{
			description:   "case 15: certificate signed by unknown authority",
			errorMessage:  "Get https://api.ci-cur-42bc2-cba40.k8s.godsmack.westeurope.azure.gigantic.io/api/v1/nodes?timeout=30s: x509: certificate signed by unknown authority (possibly because of \"crypto/rsa: verification error\" while trying to verify candidate authority certificate \"ci-cur-42bc2-cba40.k8s.godsmack.westeurope.azure.gigantic.io\")",
			expectedMatch: true,
		},
		{
			description:   "case 16: Patch timeout establishing TLS handshake",
			errorMessage:  "Patch https://api.xca65.k8s.geckon.gridscale.kvm.gigantic.io/api/v1/nodes/worker-sruw7-689bd75b49-8gbtl?timeout=30s: net/http: TLS handshake timeout",
			expectedMatch: true,
		},
		{
			description:   "case 17: Get i/o timeout establishing TCP connection",
			errorMessage:  "Get https://api.wgrt8.k8s.godsmack.westeurope.azure.gigantic.io/api/v1/nodes?timeout=30s: dial tcp 40.113.146.2:443: i/o timeout",
			expectedMatch: true,
		},
		{
			description:   "case 18: unable to connect to broken tenant api",
			errorMessage:  "Get https://api.cl048.k8s.gauss.eu-central-1.aws.gigantic.io/api/v1/namespaces/kube-system/configmaps?labelSelector=giantswarm.io%2Fservice-type%3Dmanaged%2C+giantswarm.io%2Fmanaged-by%3Dcluster-operator: EOF",
			expectedMatch: true,
		},
		{
			description:   "case 19: unable to connect to broken tenant api with expired certs",
			errorMessage:  "Get https://api.cl048.k8s.gauss.eu-central-1.aws.gigantic.io/api/v1/nodes: x509: certificate has expired or is not yet valid",
			expectedMatch: true,
		},
		{
			description:   "case 20: dns not ready alternative error (telepresence)",
			errorMessage:  "Get https://api.72fru.k8s.godsmack.westeurope.azure.gigantic.io/api/v1/nodes: dial tcp: lookup api.72fru.k8s.godsmack.westeurope.azure.gigantic.io: no such host",
			expectedMatch: true,
		},
		{
			description:   "case 21: Get i/o timeout from awaiting header",
			errorMessage:  "Get https://api.pz8mw.k8s.geckon.gridscale.kvm.gigantic.io/api?timeout=10s: context deadline exceeded (Client.Timeout exceeded while awaiting headers)",
			expectedMatch: true,
		},
		{
			description:   "case 22: api timeout eof",
			errorMessage:  "Get https://api.jfc8o.k8s.gauss.eu-central-1.aws.gigantic.io/api?timeout=10s: EOF",
			expectedMatch: true,
		},
		{
			description:   "case 23: context deadline exceeded",
			errorMessage:  "Get https://api.pnwd0.k8s.eu-central-1.aws.cps.vodafone.com/api?timeout=10s: context deadline exceeded",
			expectedMatch: true,
		},
		{
			description:   "case 24: tenant API unavailable certificate issues",
			errorMessage:  "Get \"https://api.qh99j.k8s.gauss.eu-central-1.aws.gigantic.io/api?timeout=10s\": x509: certificate is valid for ingress.local, not api.qh99j.k8s.gauss.eu-central-1.aws.gigantic.io",
			expectedMatch: true,
		},
		{
			description:   "case 25: tenant API unavailable EOF",
			errorMessage:  "Get \"https://api.hixh7.k8s.platypus.eu-west-1.aws.gigantic.io/api?timeout=10s\": EOF",
			expectedMatch: true,
		},
		{
			description:   "case 26: tenant API unavailable certificate issues - don't match partial prefix",
			errorMessage:  "t \"https://api.qh99j.k8s.gauss.eu-central-1.aws.gigantic.io/api?timeout=10s\": x509: certificate is valid for ingress.local, not api.qh99j.k8s.gauss.eu-central-1.aws.gigantic.io",
			expectedMatch: false,
		},
		{
			description:   "case 27: ingress not ready post different domain - don't match partial prefix",
			errorMessage:  "t https://api.5xchu.aws.gigantic.io: x509: certificate is valid for localhost, not api.5xchu.aws.gigantic.io:",
			expectedMatch: false,
		},
		{
			description:   "case 28: Get i/o timeout from awaiting header - quoted URL",
			errorMessage:  "Get \"https://api.fw0ex.k8s.geckon.gridscale.kvm.gigantic.io/api?timeout=10s\": context deadline exceeded (Client.Timeout exceeded while awaiting headers)",
			expectedMatch: true,
		},
		{
			description:   "case 29: ingress not ready get request - quoted URL",
			errorMessage:  "Get \"https://api.5xchu.aws.gigantic.io\": x509: certificate is valid for ingress.local, not api.5xchu.aws.gigantic.io:",
			expectedMatch: true,
		},
		{
			description:   "case 30: API not ready get EOF request - quoted URL",
			errorMessage:  "Get \"https://api.5xchu.aws.gigantic.io/api/v1/nodes\": EOF",
			expectedMatch: true,
		},
		{
			description:   "case 31: temporary issues with the master node serving the tenant cluster API - quoted URL",
			errorMessage:  "Get \"https://api.8dnxs.g8s.gorgoth.gridscale.kvm.gigantic.io/api/v1/nodes\": unexpected EOF",
			expectedMatch: true,
		},
		{
			description:   "case 32: temporary issues with the master node serving the tenant cluster API - quoted URL",
			errorMessage:  "Get \"https://api.uth29.g8s.geckon.gridscale.kvm.gigantic.io/api/v1/nodes?timeout=30s\": EOF",
			expectedMatch: true,
		},
		{
			description:   "case 33: ingress not ready post request - quoted URL",
			errorMessage:  "Post \"https://api.5xchu.aws.gigantic.io\": x509: certificate is valid for ingress.local, not api.5xchu.aws.gigantic.io:",
			expectedMatch: true,
		},
		{
			description:   "case 34: timeout getting namespace - quoted URL",
			errorMessage:  "Get \"https://api.3jwh2.k8s.aws.gigantic.io/api/v1/namespaces/giantswarm?timeout=30s\": EOF",
			expectedMatch: true,
		},
		{
			description:   "case 35: timeout getting service account - quoted URL",
			errorMessage:  "Post \"https://api.3jwh2.k8s.aws.gigantic.io/api/v1/namespaces/giantswarm/serviceaccounts?timeout=30s\": EOF",
			expectedMatch: true,
		},
		{
			description:   "case 36: GET timeout establishing TLS handshake - quoted URL",
			errorMessage:  "Get \"https://api.08vka.k8s.gorgoth.gridscale.kvm.gigantic.io/api/v1/nodes?timeout=30s\": net/http: TLS handshake timeout",
			expectedMatch: true,
		},
		{
			description:   "case 37: server is misbehaving due to TCP lookup - quoted URL",
			errorMessage:  "Get \"https://api.ci-wip-70f9b-5e958.k8s.godsmack.westeurope.azure.gigantic.io/api/v1/nodes?timeout=30s\": dial tcp: lookup api.ci-wip-70f9b-5e958.k8s.godsmack.westeurope.azure.gigantic.io on 10.96.0.10:53: server misbehaving",
			expectedMatch: true,
		},
		{
			description:   "case 38: request canceled while waiting for connection - quoted URL",
			errorMessage:  "Get \"https://api.ci-wip-2317d-c1c86.k8s.godsmack.westeurope.azure.gigantic.io/api/v1/nodes?timeout=30s\": net/http: request canceled while waiting for connection (Client.Timeout exceeded while awaiting headers)",
			expectedMatch: true,
		},
		{
			description:   "case 39: request canceled due to client timeout exceeded - quoted URL",
			errorMessage:  "Get \"https://api.06bhh.k8s.geckon.gridscale.kvm.gigantic.io/api/v1/nodes?timeout=30s\": net/http: request canceled (Client.Timeout exceeded while awaiting headers)",
			expectedMatch: true,
		},
		{
			description:   "case 40: certificate signed by unknown authority - quoted URL",
			errorMessage:  "Get \"https://api.ci-cur-42bc2-cba40.k8s.godsmack.westeurope.azure.gigantic.io/api/v1/nodes?timeout=30s\": x509: certificate signed by unknown authority (possibly because of \"crypto/rsa: verification error\" while trying to verify candidate authority certificate \"ci-cur-42bc2-cba40.k8s.godsmack.westeurope.azure.gigantic.io\")",
			expectedMatch: true,
		},
		{
			description:   "case 41: Patch timeout establishing TLS handshake - quoted URL",
			errorMessage:  "Patch \"https://api.xca65.k8s.geckon.gridscale.kvm.gigantic.io/api/v1/nodes/worker-sruw7-689bd75b49-8gbtl?timeout=30s\": net/http: TLS handshake timeout",
			expectedMatch: true,
		},
		{
			description:   "case 42: Get i/o timeout establishing TCP connection - quoted URL",
			errorMessage:  "Get \"https://api.wgrt8.k8s.godsmack.westeurope.azure.gigantic.io/api/v1/nodes?timeout=30s\": dial tcp 40.113.146.2:443: i/o timeout",
			expectedMatch: true,
		},
		{
			description:   "case 43: unable to connect to broken tenant api - quoted URL",
			errorMessage:  "Get \"https://api.cl048.k8s.gauss.eu-central-1.aws.gigantic.io/api/v1/namespaces/kube-system/configmaps?labelSelector=giantswarm.io%2Fservice-type%3Dmanaged%2C+giantswarm.io%2Fmanaged-by%3Dcluster-operator\": EOF",
			expectedMatch: true,
		},
		{
			description:   "case 44: unable to connect to broken tenant api with expired certs - quoted URL",
			errorMessage:  "Get \"https://api.cl048.k8s.gauss.eu-central-1.aws.gigantic.io/api/v1/nodes\": x509: certificate has expired or is not yet valid",
			expectedMatch: true,
		},
		{
			description:   "case 45: dns not ready alternative error (telepresence) - quoted URL",
			errorMessage:  "Get \"https://api.72fru.k8s.godsmack.westeurope.azure.gigantic.io/api/v1/nodes\": dial tcp: lookup api.72fru.k8s.godsmack.westeurope.azure.gigantic.io: no such host",
			expectedMatch: true,
		},
		{
			description:   "case 46: api timeout eof - quoted URL",
			errorMessage:  "Get https://api.jfc8o.k8s.gauss.eu-central-1.aws.gigantic.io/api?timeout=10s: EOF",
			expectedMatch: true,
		},
		{
			description:   "case 47: context deadline exceeded - quoted URL",
			errorMessage:  "Get \"https://api.pnwd0.k8s.eu-central-1.aws.cps.vodafone.com/api?timeout=10s\": context deadline exceeded",
			expectedMatch: true,
		},
		{
			description:   "case 48: tenant API unavailable - quoted URL",
			errorMessage:  "Get \"https://api.qh99j.k8s.gauss.eu-central-1.aws.gigantic.io/api?timeout=10s\": x509: certificate is valid for ingress.local, not api.qh99j.k8s.gauss.eu-central-1.aws.gigantic.io",
			expectedMatch: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			err := errors.New(tc.errorMessage)
			result := IsAPINotAvailable(err)

			if result != tc.expectedMatch {
				t.Fatalf("expected %t, got %t", tc.expectedMatch, result)
			}
		})
	}
}
