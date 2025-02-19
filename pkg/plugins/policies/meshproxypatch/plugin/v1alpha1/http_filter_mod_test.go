package v1alpha1_test

import (
	envoy_listener "github.com/envoyproxy/go-control-plane/envoy/config/listener/v3"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"sigs.k8s.io/yaml"

	core_xds "github.com/kumahq/kuma/pkg/core/xds"
	api "github.com/kumahq/kuma/pkg/plugins/policies/meshproxypatch/api/v1alpha1"
	plugin "github.com/kumahq/kuma/pkg/plugins/policies/meshproxypatch/plugin/v1alpha1"
	util_proto "github.com/kumahq/kuma/pkg/util/proto"
	"github.com/kumahq/kuma/pkg/xds/generator"
)

var _ = Describe("HTTP Filter modifications", func() {

	type testCase struct {
		listeners     []string
		modifications []string
		expected      string
	}

	DescribeTable("should apply modifications",
		func(given testCase) {
			// given
			set := core_xds.NewResourceSet()
			for _, listenerYAML := range given.listeners {
				listener := &envoy_listener.Listener{}
				err := util_proto.FromYAML([]byte(listenerYAML), listener)
				Expect(err).ToNot(HaveOccurred())
				set.Add(&core_xds.Resource{
					Name:     listener.Name,
					Origin:   generator.OriginInbound,
					Resource: listener,
				})
			}

			var mods []api.Modification
			for _, modificationYAML := range given.modifications {
				modification := api.Modification{}
				err := yaml.Unmarshal([]byte(modificationYAML), &modification)
				Expect(err).ToNot(HaveOccurred())
				mods = append(mods, modification)
			}

			// when
			err := plugin.ApplyMods(set, mods)

			// then
			Expect(err).ToNot(HaveOccurred())
			resp, err := set.List().ToDeltaDiscoveryResponse()
			Expect(err).ToNot(HaveOccurred())
			actual, err := util_proto.ToYAML(resp)
			Expect(err).ToNot(HaveOccurred())
			Expect(actual).To(MatchYAML(given.expected))
		},
		Entry("should add filter as the last", testCase{
			listeners: []string{
				`
                name: inbound:192.168.0.1:8080
                trafficDirection: INBOUND
                address:
                  socketAddress:
                    address: 192.168.0.1
                    portValue: 8080
                filterChains:
                - filters:
                  - name: envoy.filters.network.http_connection_manager
                    typedConfig:
                      '@type': type.googleapis.com/envoy.extensions.filters.network.http_connection_manager.v3.HttpConnectionManager
                      statPrefix: localhost_8080
                      httpFilters:
                      - name: envoy.filters.http.router`,
			},
			modifications: []string{`
                httpFilter:
                   operation: AddLast
                   value: |
                     name: envoy.filters.http.cors
`,
			},
			expected: `
            resources:
            - name: inbound:192.168.0.1:8080
              resource:
                '@type': type.googleapis.com/envoy.config.listener.v3.Listener
                address:
                  socketAddress:
                    address: 192.168.0.1
                    portValue: 8080
                filterChains:
                - filters:
                  - name: envoy.filters.network.http_connection_manager
                    typedConfig:
                      '@type': type.googleapis.com/envoy.extensions.filters.network.http_connection_manager.v3.HttpConnectionManager
                      httpFilters:
                      - name: envoy.filters.http.router
                      - name: envoy.filters.http.cors
                      statPrefix: localhost_8080
                name: inbound:192.168.0.1:8080
                trafficDirection: INBOUND`,
		}),
		Entry("should add filter as the first", testCase{
			listeners: []string{
				`
                name: inbound:192.168.0.1:8080
                trafficDirection: INBOUND
                address:
                  socketAddress:
                    address: 192.168.0.1
                    portValue: 8080
                filterChains:
                - filters:
                  - name: envoy.filters.network.http_connection_manager
                    typedConfig:
                      '@type': type.googleapis.com/envoy.extensions.filters.network.http_connection_manager.v3.HttpConnectionManager
                      statPrefix: localhost_8080
                      httpFilters:
                      - name: envoy.filters.http.router`,
			},
			modifications: []string{`
                httpFilter:
                   operation: AddFirst
                   value: |
                     name: envoy.filters.http.cors
`,
			},
			expected: `
            resources:
            - name: inbound:192.168.0.1:8080
              resource:
                '@type': type.googleapis.com/envoy.config.listener.v3.Listener
                address:
                  socketAddress:
                    address: 192.168.0.1
                    portValue: 8080
                filterChains:
                - filters:
                  - name: envoy.filters.network.http_connection_manager
                    typedConfig:
                      '@type': type.googleapis.com/envoy.extensions.filters.network.http_connection_manager.v3.HttpConnectionManager
                      httpFilters:
                      - name: envoy.filters.http.cors
                      - name: envoy.filters.http.router
                      statPrefix: localhost_8080
                name: inbound:192.168.0.1:8080
                trafficDirection: INBOUND`,
		}),
		Entry("should remove all filters from all listeners when there is no match section", testCase{
			listeners: []string{
				`
                address:
                  socketAddress:
                    address: 192.168.0.1
                    portValue: 8080
                filterChains:
                - filters:
                  - name: envoy.filters.network.http_connection_manager
                    typedConfig:
                      '@type': type.googleapis.com/envoy.extensions.filters.network.http_connection_manager.v3.HttpConnectionManager
                      httpFilters:
                      - name: envoy.filters.http.router
                      - name: envoy.filters.http.cors
                      statPrefix: localhost_8080
                name: inbound:192.168.0.1:8080
                trafficDirection: INBOUND`,
			},
			modifications: []string{`
                httpFilter:
                   operation: Remove
`,
			},
			expected: `
            resources:
            - name: inbound:192.168.0.1:8080
              resource:
                '@type': type.googleapis.com/envoy.config.listener.v3.Listener
                address:
                  socketAddress:
                    address: 192.168.0.1
                    portValue: 8080
                filterChains:
                - filters:
                  - name: envoy.filters.network.http_connection_manager
                    typedConfig:
                      '@type': type.googleapis.com/envoy.extensions.filters.network.http_connection_manager.v3.HttpConnectionManager
                      statPrefix: localhost_8080
                name: inbound:192.168.0.1:8080
                trafficDirection: INBOUND`,
		}),
		Entry("should remove all filters from all listeners when there is inbound match section", testCase{
			listeners: []string{
				`
                address:
                  socketAddress:
                    address: 192.168.0.1
                    portValue: 8080
                filterChains:
                - filters:
                  - name: envoy.filters.network.http_connection_manager
                    typedConfig:
                      '@type': type.googleapis.com/envoy.extensions.filters.network.http_connection_manager.v3.HttpConnectionManager
                      httpFilters:
                      - name: envoy.filters.http.router
                      - name: envoy.filters.http.cors
                      statPrefix: localhost_8080
                name: inbound:192.168.0.1:8080
                trafficDirection: INBOUND`,
			},
			modifications: []string{`
                httpFilter:
                   operation: Remove
                   match:
                     origin: inbound
`,
			},
			expected: `
            resources:
            - name: inbound:192.168.0.1:8080
              resource:
                '@type': type.googleapis.com/envoy.config.listener.v3.Listener
                address:
                  socketAddress:
                    address: 192.168.0.1
                    portValue: 8080
                filterChains:
                - filters:
                  - name: envoy.filters.network.http_connection_manager
                    typedConfig:
                      '@type': type.googleapis.com/envoy.extensions.filters.network.http_connection_manager.v3.HttpConnectionManager
                      statPrefix: localhost_8080
                name: inbound:192.168.0.1:8080
                trafficDirection: INBOUND`,
		}),
		Entry("should remove all filters from picked listener", testCase{
			listeners: []string{
				`
                address:
                  socketAddress:
                    address: 192.168.0.1
                    portValue: 8080
                filterChains:
                - filters:
                  - name: envoy.filters.network.http_connection_manager
                    typedConfig:
                      '@type': type.googleapis.com/envoy.extensions.filters.network.http_connection_manager.v3.HttpConnectionManager
                      httpFilters:
                      - name: envoy.filters.http.router
                      - name: envoy.filters.http.cors
                      statPrefix: localhost_8080
                name: inbound:192.168.0.1:8080
                trafficDirection: INBOUND`,
				`
                address:
                  socketAddress:
                    address: 192.168.0.1
                    portValue: 8081
                filterChains:
                - filters:
                  - name: envoy.filters.network.http_connection_manager
                    typedConfig:
                      '@type': type.googleapis.com/envoy.extensions.filters.network.http_connection_manager.v3.HttpConnectionManager
                      httpFilters:
                      - name: envoy.filters.http.router
                      - name: envoy.filters.http.cors
                      statPrefix: localhost_8081
                name: inbound:192.168.0.1:8081
                trafficDirection: INBOUND`,
			},
			modifications: []string{`
                httpFilter:
                   operation: Remove
                   match:
                     listenerName: inbound:192.168.0.1:8080
`,
			},
			expected: `
            resources:
            - name: inbound:192.168.0.1:8080
              resource:
                '@type': type.googleapis.com/envoy.config.listener.v3.Listener
                address:
                  socketAddress:
                    address: 192.168.0.1
                    portValue: 8080
                filterChains:
                - filters:
                  - name: envoy.filters.network.http_connection_manager
                    typedConfig:
                      '@type': type.googleapis.com/envoy.extensions.filters.network.http_connection_manager.v3.HttpConnectionManager
                      statPrefix: localhost_8080
                name: inbound:192.168.0.1:8080
                trafficDirection: INBOUND
            - name: inbound:192.168.0.1:8081
              resource:
                '@type': type.googleapis.com/envoy.config.listener.v3.Listener
                address:
                  socketAddress:
                    address: 192.168.0.1
                    portValue: 8081
                filterChains:
                - filters:
                  - name: envoy.filters.network.http_connection_manager
                    typedConfig:
                      '@type': type.googleapis.com/envoy.extensions.filters.network.http_connection_manager.v3.HttpConnectionManager
                      httpFilters:
                      - name: envoy.filters.http.router
                      - name: envoy.filters.http.cors
                      statPrefix: localhost_8081
                name: inbound:192.168.0.1:8081
                trafficDirection: INBOUND`,
		}),
		Entry("should remove all filters of given name from all listeners", testCase{
			listeners: []string{
				`
                address:
                  socketAddress:
                    address: 192.168.0.1
                    portValue: 8080
                filterChains:
                - filters:
                  - name: envoy.filters.network.http_connection_manager
                    typedConfig:
                      '@type': type.googleapis.com/envoy.extensions.filters.network.http_connection_manager.v3.HttpConnectionManager
                      httpFilters:
                      - name: envoy.filters.http.router
                      - name: envoy.filters.http.cors
                      statPrefix: localhost_8080
                name: inbound:192.168.0.1:8080
                trafficDirection: INBOUND`,
				`
                address:
                  socketAddress:
                    address: 192.168.0.1
                    portValue: 8081
                filterChains:
                - filters:
                  - name: envoy.filters.network.http_connection_manager
                    typedConfig:
                      '@type': type.googleapis.com/envoy.extensions.filters.network.http_connection_manager.v3.HttpConnectionManager
                      httpFilters:
                      - name: envoy.filters.http.router
                      - name: envoy.filters.http.cors
                      statPrefix: localhost_8081
                name: inbound:192.168.0.1:8081
                trafficDirection: INBOUND`,
			},
			modifications: []string{`
                httpFilter:
                   operation: Remove
                   match:
                     name: envoy.filters.http.cors
`,
			},
			expected: `
            resources:
            - name: inbound:192.168.0.1:8080
              resource:
                '@type': type.googleapis.com/envoy.config.listener.v3.Listener
                address:
                  socketAddress:
                    address: 192.168.0.1
                    portValue: 8080
                filterChains:
                - filters:
                  - name: envoy.filters.network.http_connection_manager
                    typedConfig:
                      '@type': type.googleapis.com/envoy.extensions.filters.network.http_connection_manager.v3.HttpConnectionManager
                      httpFilters:
                      - name: envoy.filters.http.router
                      statPrefix: localhost_8080
                name: inbound:192.168.0.1:8080
                trafficDirection: INBOUND
            - name: inbound:192.168.0.1:8081
              resource:
                '@type': type.googleapis.com/envoy.config.listener.v3.Listener
                address:
                  socketAddress:
                    address: 192.168.0.1
                    portValue: 8081
                filterChains:
                - filters:
                  - name: envoy.filters.network.http_connection_manager
                    typedConfig:
                      '@type': type.googleapis.com/envoy.extensions.filters.network.http_connection_manager.v3.HttpConnectionManager
                      httpFilters:
                      - name: envoy.filters.http.router
                      statPrefix: localhost_8081
                name: inbound:192.168.0.1:8081
                trafficDirection: INBOUND`,
		}),
		Entry("should add filter after already defined", testCase{
			listeners: []string{
				`
                address:
                  socketAddress:
                    address: 192.168.0.1
                    portValue: 8080
                filterChains:
                - filters:
                  - name: envoy.filters.network.http_connection_manager
                    typedConfig:
                      '@type': type.googleapis.com/envoy.extensions.filters.network.http_connection_manager.v3.HttpConnectionManager
                      httpFilters:
                      - name: envoy.filters.http.router
                      - name: envoy.filters.http.gzip
                      statPrefix: localhost_8080
                name: inbound:192.168.0.1:8080
                trafficDirection: INBOUND`,
			},
			modifications: []string{`
                httpFilter:
                   operation: AddAfter
                   match:
                     name: envoy.filters.http.router
                   value: |
                     name: envoy.filters.http.cors
`,
			},
			expected: `
            resources:
            - name: inbound:192.168.0.1:8080
              resource:
                '@type': type.googleapis.com/envoy.config.listener.v3.Listener
                address:
                  socketAddress:
                    address: 192.168.0.1
                    portValue: 8080
                filterChains:
                - filters:
                  - name: envoy.filters.network.http_connection_manager
                    typedConfig:
                      '@type': type.googleapis.com/envoy.extensions.filters.network.http_connection_manager.v3.HttpConnectionManager
                      httpFilters:
                      - name: envoy.filters.http.router
                      - name: envoy.filters.http.cors
                      - name: envoy.filters.http.gzip
                      statPrefix: localhost_8080
                name: inbound:192.168.0.1:8080
                trafficDirection: INBOUND`,
		}),
		Entry("should not add filter when name is not matched", testCase{
			listeners: []string{
				`
                address:
                  socketAddress:
                    address: 192.168.0.1
                    portValue: 8080
                filterChains:
                - filters:
                  - name: envoy.filters.network.http_connection_manager
                    typedConfig:
                      '@type': type.googleapis.com/envoy.extensions.filters.network.http_connection_manager.v3.HttpConnectionManager
                      httpFilters:
                      - name: envoy.filters.http.router
                      statPrefix: localhost_8080
                name: inbound:192.168.0.1:8080
                trafficDirection: INBOUND`,
			},
			modifications: []string{`
                httpFilter:
                   operation: AddAfter
                   match:
                     name: envoy.filters.http.gzip
                   value: |
                     name: envoy.filters.http.cors
`,
			},
			expected: `
            resources:
            - name: inbound:192.168.0.1:8080
              resource:
                '@type': type.googleapis.com/envoy.config.listener.v3.Listener
                address:
                  socketAddress:
                    address: 192.168.0.1
                    portValue: 8080
                filterChains:
                - filters:
                  - name: envoy.filters.network.http_connection_manager
                    typedConfig:
                      '@type': type.googleapis.com/envoy.extensions.filters.network.http_connection_manager.v3.HttpConnectionManager
                      httpFilters:
                      - name: envoy.filters.http.router
                      statPrefix: localhost_8080
                name: inbound:192.168.0.1:8080
                trafficDirection: INBOUND`,
		}),
		Entry("should add filter before already defined", testCase{
			listeners: []string{
				`
                address:
                  socketAddress:
                    address: 192.168.0.1
                    portValue: 8080
                filterChains:
                - filters:
                  - name: envoy.filters.network.http_connection_manager
                    typedConfig:
                      '@type': type.googleapis.com/envoy.extensions.filters.network.http_connection_manager.v3.HttpConnectionManager
                      httpFilters:
                      - name: envoy.filters.http.router
                      - name: envoy.filters.http.gzip
                      statPrefix: localhost_8080
                name: inbound:192.168.0.1:8080
                trafficDirection: INBOUND`,
			},
			modifications: []string{`
                httpFilter:
                   operation: AddBefore
                   match:
                     name: envoy.filters.http.gzip
                   value: |
                     name: envoy.filters.http.cors
`,
			},
			expected: `
            resources:
            - name: inbound:192.168.0.1:8080
              resource:
                '@type': type.googleapis.com/envoy.config.listener.v3.Listener
                address:
                  socketAddress:
                    address: 192.168.0.1
                    portValue: 8080
                filterChains:
                - filters:
                  - name: envoy.filters.network.http_connection_manager
                    typedConfig:
                      '@type': type.googleapis.com/envoy.extensions.filters.network.http_connection_manager.v3.HttpConnectionManager
                      httpFilters:
                      - name: envoy.filters.http.router
                      - name: envoy.filters.http.cors
                      - name: envoy.filters.http.gzip
                      statPrefix: localhost_8080
                name: inbound:192.168.0.1:8080
                trafficDirection: INBOUND`,
		}),
		Entry("should patch resource matching filter name", testCase{
			listeners: []string{
				`
                address:
                  socketAddress:
                    address: 192.168.0.1
                    portValue: 8080
                filterChains:
                - filters:
                  - name: envoy.filters.network.http_connection_manager
                    typedConfig:
                      '@type': type.googleapis.com/envoy.extensions.filters.network.http_connection_manager.v3.HttpConnectionManager
                      httpFilters:
                      - name: envoy.filters.http.router
                        typedConfig:
                          '@type': type.googleapis.com/envoy.extensions.filters.http.router.v3.Router
                          startChildSpan: true
                      - name: envoy.filters.http.gzip
                      statPrefix: localhost_8080
                name: inbound:192.168.0.1:8080
                trafficDirection: INBOUND`,
			},
			modifications: []string{`
                httpFilter:
                   operation: Patch
                   match:
                     name: envoy.filters.http.router
                   value: |
                     typedConfig:
                       '@type': type.googleapis.com/envoy.extensions.filters.http.router.v3.Router
                       dynamicStats: false
`,
			},
			expected: `
            resources:
            - name: inbound:192.168.0.1:8080
              resource:
                '@type': type.googleapis.com/envoy.config.listener.v3.Listener
                address:
                  socketAddress:
                    address: 192.168.0.1
                    portValue: 8080
                filterChains:
                - filters:
                  - name: envoy.filters.network.http_connection_manager
                    typedConfig:
                      '@type': type.googleapis.com/envoy.extensions.filters.network.http_connection_manager.v3.HttpConnectionManager
                      httpFilters:
                      - name: envoy.filters.http.router
                        typedConfig:
                          '@type': type.googleapis.com/envoy.extensions.filters.http.router.v3.Router
                          startChildSpan: true
                          dynamicStats: false
                      - name: envoy.filters.http.gzip
                      statPrefix: localhost_8080
                name: inbound:192.168.0.1:8080
                trafficDirection: INBOUND`,
		}),
		Entry("should patch resource matching listener tags", testCase{
			listeners: []string{
				`
                address:
                  socketAddress:
                    address: 192.168.0.1
                    portValue: 8080
                metadata:
                  filterMetadata:
                    io.kuma.tags:
                      kuma.io/service: backend
                filterChains:
                - filters:
                  - name: envoy.filters.network.http_connection_manager
                    typedConfig:
                      '@type': type.googleapis.com/envoy.extensions.filters.network.http_connection_manager.v3.HttpConnectionManager
                      httpFilters:
                      - name: envoy.filters.http.router
                        typedConfig:
                          '@type': type.googleapis.com/envoy.extensions.filters.http.router.v3.Router
                          startChildSpan: true
                      - name: envoy.filters.http.gzip
                      statPrefix: localhost_8080
                name: inbound:192.168.0.1:8080
                trafficDirection: INBOUND`,
			},
			modifications: []string{`
                httpFilter:
                   operation: Patch
                   match:
                     name: envoy.filters.http.router
                     listenerTags:
                       kuma.io/service: backend
                   value: |
                     typedConfig:
                       '@type': type.googleapis.com/envoy.extensions.filters.http.router.v3.Router
                       dynamicStats: false
`,
			},
			expected: `
            resources:
            - name: inbound:192.168.0.1:8080
              resource:
                '@type': type.googleapis.com/envoy.config.listener.v3.Listener
                address:
                  socketAddress:
                    address: 192.168.0.1
                    portValue: 8080
                metadata:
                  filterMetadata:
                    io.kuma.tags:
                      kuma.io/service: backend
                filterChains:
                - filters:
                  - name: envoy.filters.network.http_connection_manager
                    typedConfig:
                      '@type': type.googleapis.com/envoy.extensions.filters.network.http_connection_manager.v3.HttpConnectionManager
                      httpFilters:
                      - name: envoy.filters.http.router
                        typedConfig:
                          '@type': type.googleapis.com/envoy.extensions.filters.http.router.v3.Router
                          startChildSpan: true
                          dynamicStats: false
                      - name: envoy.filters.http.gzip
                      statPrefix: localhost_8080
                name: inbound:192.168.0.1:8080
                trafficDirection: INBOUND`,
		}),
		Entry("should not patch resource not matching listener tags", testCase{
			listeners: []string{
				`
                address:
                  socketAddress:
                    address: 192.168.0.1
                    portValue: 8080
                metadata:
                  filterMetadata:
                    io.kuma.tags:
                      kuma.io/service: backend
                filterChains:
                - filters:
                  - name: envoy.filters.network.http_connection_manager
                    typedConfig:
                      '@type': type.googleapis.com/envoy.extensions.filters.network.http_connection_manager.v3.HttpConnectionManager
                      httpFilters:
                      - name: envoy.filters.http.router
                        typedConfig:
                          '@type': type.googleapis.com/envoy.extensions.filters.http.router.v3.Router
                          startChildSpan: true
                      - name: envoy.filters.http.gzip
                      statPrefix: localhost_8080
                name: inbound:192.168.0.1:8080
                trafficDirection: INBOUND`,
			},
			modifications: []string{`
                httpFilter:
                   operation: Patch
                   match:
                     name: envoy.filters.http.router
                     listenerTags:
                       kuma.io/service: web
                   value: |
                     typedConfig:
                       '@type': type.googleapis.com/envoy.extensions.filters.http.router.v3.Router
                       dynamicStats: false
`,
			},
			expected: `
            resources:
            - name: inbound:192.168.0.1:8080
              resource:
                '@type': type.googleapis.com/envoy.config.listener.v3.Listener
                address:
                  socketAddress:
                    address: 192.168.0.1
                    portValue: 8080
                metadata:
                  filterMetadata:
                    io.kuma.tags:
                      kuma.io/service: backend
                filterChains:
                - filters:
                  - name: envoy.filters.network.http_connection_manager
                    typedConfig:
                      '@type': type.googleapis.com/envoy.extensions.filters.network.http_connection_manager.v3.HttpConnectionManager
                      httpFilters:
                      - name: envoy.filters.http.router
                        typedConfig:
                          '@type': type.googleapis.com/envoy.extensions.filters.http.router.v3.Router
                          startChildSpan: true
                      - name: envoy.filters.http.gzip
                      statPrefix: localhost_8080
                name: inbound:192.168.0.1:8080
                trafficDirection: INBOUND`,
		}),
	)
})
