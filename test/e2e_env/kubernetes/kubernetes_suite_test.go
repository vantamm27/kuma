package kubernetes_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"

	"github.com/kumahq/kuma/pkg/test"
	"github.com/kumahq/kuma/test/e2e_env/kubernetes/container_patch"
	"github.com/kumahq/kuma/test/e2e_env/kubernetes/defaults"
	"github.com/kumahq/kuma/test/e2e_env/kubernetes/externalservices"
	"github.com/kumahq/kuma/test/e2e_env/kubernetes/gateway"
	"github.com/kumahq/kuma/test/e2e_env/kubernetes/graceful"
	"github.com/kumahq/kuma/test/e2e_env/kubernetes/healthcheck"
	"github.com/kumahq/kuma/test/e2e_env/kubernetes/inspect"
	"github.com/kumahq/kuma/test/e2e_env/kubernetes/jobs"
	"github.com/kumahq/kuma/test/e2e_env/kubernetes/k8s_api_bypass"
	"github.com/kumahq/kuma/test/e2e_env/kubernetes/kic"
	"github.com/kumahq/kuma/test/e2e_env/kubernetes/membership"
	"github.com/kumahq/kuma/test/e2e_env/kubernetes/meshcircuitbreaker"
	"github.com/kumahq/kuma/test/e2e_env/kubernetes/meshfaultinjection"
	"github.com/kumahq/kuma/test/e2e_env/kubernetes/meshhealthcheck"
	"github.com/kumahq/kuma/test/e2e_env/kubernetes/meshhttproute"
	"github.com/kumahq/kuma/test/e2e_env/kubernetes/meshproxypatch"
	"github.com/kumahq/kuma/test/e2e_env/kubernetes/meshratelimit"
	"github.com/kumahq/kuma/test/e2e_env/kubernetes/meshretry"
	"github.com/kumahq/kuma/test/e2e_env/kubernetes/meshtimeout"
	"github.com/kumahq/kuma/test/e2e_env/kubernetes/meshtrafficpermission"
	"github.com/kumahq/kuma/test/e2e_env/kubernetes/observability"
	"github.com/kumahq/kuma/test/e2e_env/kubernetes/reachableservices"
	"github.com/kumahq/kuma/test/e2e_env/kubernetes/trafficlog"
	"github.com/kumahq/kuma/test/e2e_env/kubernetes/virtualoutbound"
	"github.com/kumahq/kuma/test/framework/envs/kubernetes"
)

func TestE2E(t *testing.T) {
	test.RunE2ESpecs(t, "E2E Kubernetes Suite")
}

var _ = SynchronizedBeforeSuite(kubernetes.SetupAndGetState, kubernetes.RestoreState)

// SynchronizedAfterSuite keeps the main process alive until all other processes finish.
// Otherwise, we would close port-forward to the CP and remaining tests executed in different processes may fail.
var _ = SynchronizedAfterSuite(func() {}, func() {})

var _ = Describe("Virtual Probes", healthcheck.VirtualProbes, Ordered)
var _ = Describe("Gateway", gateway.Gateway, Ordered)
var _ = Describe("Gateway - Cross-mesh", gateway.CrossMeshGatewayOnKubernetes, Ordered)
var _ = Describe("Gateway - Gateway API", Label("arm-not-supported"), gateway.GatewayAPI, Ordered)
var _ = Describe("Gateway - mTLS", gateway.Mtls, Ordered)
var _ = Describe("Gateway - Resources", gateway.Resources, Ordered)
var _ = Describe("Graceful", graceful.Graceful, Ordered)
var _ = Describe("Eviction", graceful.Eviction, Ordered)
var _ = Describe("Jobs", jobs.Jobs)
var _ = Describe("Membership", membership.Membership, Ordered)
var _ = Describe("Container Patch", container_patch.ContainerPatch, Ordered)
var _ = Describe("Metrics", observability.ApplicationsMetrics, Ordered)
var _ = Describe("Tracing", observability.Tracing, Ordered)
var _ = Describe("MeshTrace", observability.PluginTest, Ordered)
var _ = Describe("Traffic Log", trafficlog.TCPLogging, Ordered)
var _ = Describe("Inspect", inspect.Inspect, Ordered)
var _ = Describe("K8S API Bypass", k8s_api_bypass.K8sApiBypass, Ordered)
var _ = Describe("Reachable Services", reachableservices.ReachableServices, Ordered)
var _ = Describe("Defaults", defaults.Defaults, Ordered)
var _ = Describe("External Services", externalservices.ExternalServices, Ordered)
var _ = Describe("Virtual Outbound", virtualoutbound.VirtualOutbound, Ordered)
var _ = Describe("Kong Ingress Controller", Label("arm-not-supported"), kic.KICKubernetes, Ordered)
var _ = Describe("MeshTrafficPermission API", meshtrafficpermission.API, Ordered)
var _ = Describe("MeshRateLimit API", meshratelimit.API, Ordered)
var _ = Describe("MeshTimeout API", meshtimeout.MeshTimeout, Ordered)
var _ = Describe("MeshHealthCheck API", meshhealthcheck.API, Ordered)
var _ = Describe("MeshCircuitBreaker API", meshcircuitbreaker.API, Ordered)
var _ = Describe("MeshCircuitBreaker", meshcircuitbreaker.MeshCircuitBreaker, Ordered)
var _ = Describe("MeshRetry", meshretry.API, Ordered)
var _ = Describe("MeshProxyPatch", meshproxypatch.MeshProxyPatch, Ordered)
var _ = Describe("MeshFaultInjection", meshfaultinjection.API, Ordered)
var _ = Describe("MeshHTTPRoute", meshhttproute.Test, Ordered)
