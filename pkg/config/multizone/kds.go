package multizone

import (
	"github.com/pkg/errors"
	"go.uber.org/multierr"

	"github.com/kumahq/kuma/pkg/config"
	config_types "github.com/kumahq/kuma/pkg/config/types"
)

type KdsServerConfig struct {
	// Port of a gRPC server that serves Kuma Discovery Service (KDS).
	GrpcPort uint32 `json:"grpcPort" envconfig:"kuma_multizone_global_kds_grpc_port"`
	// Interval for refreshing state of the world
	RefreshInterval config_types.Duration `json:"refreshInterval" envconfig:"kuma_multizone_global_kds_refresh_interval"`
	// Interval for flushing Zone Insights (stats of multi-zone communication)
	ZoneInsightFlushInterval config_types.Duration `json:"zoneInsightFlushInterval" envconfig:"kuma_multizone_global_kds_zone_insight_flush_interval"`
	// TlsCertFile defines a path to a file with PEM-encoded TLS cert.
	TlsCertFile string `json:"tlsCertFile" envconfig:"kuma_multizone_global_kds_tls_cert_file"`
	// TlsKeyFile defines a path to a file with PEM-encoded TLS key.
	TlsKeyFile string `json:"tlsKeyFile" envconfig:"kuma_multizone_global_kds_tls_key_file"`
	// TlsMinVersion defines the minimum TLS version to be used
	TlsMinVersion string `json:"tlsMinVersion" envconfig:"kuma_multizone_global_kds_tls_min_version"`
	// TlsMaxVersion defines the maximum TLS version to be used
	TlsMaxVersion string `json:"tlsMaxVersion" envconfig:"kuma_multizone_global_kds_tls_max_version"`
	// TlsCipherSuites defines the list of ciphers to use
	TlsCipherSuites []string `json:"tlsCipherSuites" envconfig:"kuma_multizone_global_kds_tls_cipher_suites"`
	// MaxMsgSize defines a maximum size of the message that is exchanged using KDS.
	// In practice this means a limit on full list of one resource type.
	MaxMsgSize uint32 `json:"maxMsgSize" envconfig:"kuma_multizone_global_kds_max_msg_size"`
	// MsgSendTimeout defines a timeout on sending a single KDS message.
	// KDS stream between control planes is terminated if the control plane hits this timeout.
	MsgSendTimeout config_types.Duration `json:"msgSendTimeout" envconfig:"kuma_multizone_global_kds_msg_send_timeout"`
	// Backoff that is executed when the global control plane is sending the response that was previously rejected by zone control plane.
	NackBackoff config_types.Duration `json:"nackBackoff" envconfig:"kuma_multizone_global_kds_nack_backoff"`
}

var _ config.Config = &KdsServerConfig{}

func (c *KdsServerConfig) Sanitize() {
}

func (c *KdsServerConfig) Validate() error {
	var errs error
	if c.GrpcPort > 65535 {
		errs = multierr.Append(errs, errors.Errorf(".GrpcPort must be in the range [0, 65535]"))
	}
	if c.RefreshInterval.Duration <= 0 {
		errs = multierr.Append(errs, errors.New(".RefreshInterval must be positive"))
	}
	if c.ZoneInsightFlushInterval.Duration <= 0 {
		errs = multierr.Append(errs, errors.New(".ZoneInsightFlushInterval must be positive"))
	}
	if c.TlsCertFile == "" && c.TlsKeyFile != "" {
		errs = multierr.Append(errs, errors.New(".TlsCertFile cannot be empty if TlsKeyFile has been set"))
	}
	if c.TlsKeyFile == "" && c.TlsCertFile != "" {
		errs = multierr.Append(errs, errors.New(".TlsKeyFile cannot be empty if TlsCertFile has been set"))
	}
	if _, err := config_types.TLSVersion(c.TlsMinVersion); err != nil {
		errs = multierr.Append(errs, errors.New(".TlsMinVersion"+err.Error()))
	}
	if _, err := config_types.TLSVersion(c.TlsMaxVersion); err != nil {
		errs = multierr.Append(errs, errors.New(".TlsMaxVersion"+err.Error()))
	}
	if _, err := config_types.TLSCiphers(c.TlsCipherSuites); err != nil {
		errs = multierr.Append(errs, errors.New(".TlsCipherSuites"+err.Error()))
	}
	return errs
}

type KdsClientConfig struct {
	// Interval for refreshing state of the world
	RefreshInterval config_types.Duration `json:"refreshInterval" envconfig:"kuma_multizone_zone_kds_refresh_interval"`
	// RootCAFile defines a path to a file with PEM-encoded Root CA. Client will verify the server by using it.
	RootCAFile string `json:"rootCaFile" envconfig:"kuma_multizone_zone_kds_root_ca_file"`
	// MaxMsgSize defines a maximum size of the message that is exchanged using KDS.
	// In practice this means a limit on full list of one resource type.
	MaxMsgSize uint32 `json:"maxMsgSize" envconfig:"kuma_multizone_zone_kds_max_msg_size"`
	// MsgSendTimeout defines a timeout on sending a single KDS message.
	// KDS stream between control planes is terminated if the control plane hits this timeout.
	MsgSendTimeout config_types.Duration `json:"msgSendTimeout" envconfig:"kuma_multizone_zone_kds_msg_send_timeout"`
	// Backoff that is executed when the zone control plane is sending the response that was previously rejected by global control plane.
	NackBackoff config_types.Duration `json:"nackBackoff" envconfig:"kuma_multizone_zone_kds_nack_backoff"`
}

var _ config.Config = &KdsClientConfig{}

func (k KdsClientConfig) Sanitize() {
}

func (k KdsClientConfig) Validate() error {
	return nil
}
