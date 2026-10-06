package controller

import (
	"testing"

	"github.com/spf13/pflag"
	"github.com/stretchr/testify/require"
)

func TestConfigurationValidateModeFlags(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name          string
		enableBgpLb   bool
		enableLbSvc   bool
		expectErr     bool
		errorContains string
	}{
		{
			name:        "both disabled is valid",
			enableBgpLb: false,
			enableLbSvc: false,
			expectErr:   false,
		},
		{
			name:        "only bgp lb eip enabled is valid",
			enableBgpLb: true,
			enableLbSvc: false,
			expectErr:   false,
		},
		{
			name:        "only lb svc enabled is valid",
			enableBgpLb: false,
			enableLbSvc: true,
			expectErr:   false,
		},
		{
			name:          "both enabled is invalid",
			enableBgpLb:   true,
			enableLbSvc:   true,
			expectErr:     true,
			errorContains: "mutually exclusive",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			cfg := &Configuration{
				EnableBgpLbVip: tt.enableBgpLb,
				EnablePodLbSvc: tt.enableLbSvc,
			}

			err := cfg.validateModeFlags()
			if tt.expectErr {
				require.Error(t, err)
				require.ErrorContains(t, err, tt.errorContains)
				return
			}
			require.NoError(t, err)
		})
	}
}

func TestRegisterBgpLbVipFlags(t *testing.T) {
	t.Parallel()

	t.Run("enable via vip flag", func(t *testing.T) {
		t.Parallel()

		fs := pflag.NewFlagSet("test", pflag.ContinueOnError)
		enabled := fs.Bool("enable-bgp-lb-vip", false, "test flag")

		require.NoError(t, fs.Parse([]string{"--enable-bgp-lb-vip=true"}))
		require.True(t, *enabled)
	})

	t.Run("unknown bgp lb flag is rejected", func(t *testing.T) {
		t.Parallel()

		fs := pflag.NewFlagSet("test", pflag.ContinueOnError)
		enabled := fs.Bool("enable-bgp-lb-vip", false, "test flag")

		err := fs.Parse([]string{"--enable-bgp-lb-unknown=true"})
		require.Error(t, err)
		require.False(t, *enabled)
	})
}

func TestConfigurationValidateServiceFeatureGates(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		config  Configuration
		wantErr string
	}{
		{name: "all disabled"},
		{name: "ovn and gateway load balancer coexist", config: Configuration{EnableOvnLB: true, EnableGwNftableLbSvc: true}},
		{name: "both gateway identities coexist", config: Configuration{EnableGwNftableLbSvc: true, EnableGwNftableSvcClusterIP: true}},
		{
			name:    "pod and gateway load balancer conflict",
			config:  Configuration{EnableOvnLB: true, EnablePodLbSvc: true, EnableGwNftableLbSvc: true},
			wantErr: "--enable-lb-svc and --enable-gw-nftable-lb-svc are mutually exclusive",
		},
		{
			name:    "pod load balancer and gateway cluster ip conflict",
			config:  Configuration{EnableOvnLB: true, EnablePodLbSvc: true, EnableGwNftableSvcClusterIP: true},
			wantErr: "--enable-lb-svc and --enable-gw-nftable-svc-cluster-ip are mutually exclusive",
		},
		{
			name:    "ovn and gateway cluster ip conflict",
			config:  Configuration{EnableOvnLB: true, EnableGwNftableSvcClusterIP: true},
			wantErr: "--enable-gw-nftable-svc-cluster-ip and --enable-lb are mutually exclusive",
		},
		{
			// the lanIP vip switch is a pure add-on: it must combine with every mode,
			// including on its own
			name:   "lanip vip combines with every mode",
			config: Configuration{EnableOvnLB: true, EnableGwNftableLbSvc: true, EnableGwNftableLanipVip: true},
		},
		{name: "lanip vip alone", config: Configuration{EnableGwNftableLanipVip: true}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.config.validateServiceFeatureGates()
			if tt.wantErr == "" {
				require.NoError(t, err)
			} else {
				require.EqualError(t, err, tt.wantErr)
			}
		})
	}
}

// resolveServiceFeatureGateDefaults must compose the default-on gateway nftables modes with the
// historical switches: explicit operator intent always wins, and a bare default startup must pass
// validation with the OVN switch LB replaced by the gateway ClusterIP mode.
func TestResolveServiceFeatureGateDefaults(t *testing.T) {
	t.Parallel()

	changedNone := func(string) bool { return false }
	changedOf := func(flags ...string) func(string) bool {
		return func(name string) bool {
			for _, f := range flags {
				if f == name {
					return true
				}
			}
			return false
		}
	}
	defaults := func() Configuration {
		// the flag defaults as registered in ParseFlags
		return Configuration{
			EnableOvnLB:                 true, // --enable-lb historical default
			EnableGwNftableLbSvc:        true,
			EnableGwNftableSvcClusterIP: true,
			EnableGwNftableLanipVip:     true,
		}
	}

	tests := []struct {
		name         string
		changed      func(string) bool
		mutate       func(*Configuration)
		want         Configuration
		wantValidErr bool
	}{
		{
			name:    "bare defaults: gateway modes on, OVN LB auto-off",
			changed: changedNone,
			want:    Configuration{EnableGwNftableLbSvc: true, EnableGwNftableSvcClusterIP: true, EnableGwNftableLanipVip: true},
		},
		{
			name:    "explicit enable-lb wins over the default-on gateway cluster-ip mode",
			changed: changedOf("enable-lb"),
			want:    Configuration{EnableOvnLB: true, EnableGwNftableLbSvc: true, EnableGwNftableLanipVip: true},
		},
		{
			name:         "both sides pinned: resolver leaves them for the validator to reject",
			changed:      changedOf("enable-lb", "enable-gw-nftable-svc-cluster-ip"),
			want:         Configuration{EnableOvnLB: true, EnableGwNftableLbSvc: true, EnableGwNftableSvcClusterIP: true, EnableGwNftableLanipVip: true},
			wantValidErr: true,
		},
		{
			name:    "opt-out of gateway cluster-ip restores the OVN LB default",
			changed: changedOf("enable-gw-nftable-svc-cluster-ip"),
			mutate:  func(c *Configuration) { c.EnableGwNftableSvcClusterIP = false },
			want:    Configuration{EnableOvnLB: true, EnableGwNftableLbSvc: true, EnableGwNftableLanipVip: true},
		},
		{
			name:    "explicit enable-lb-svc makes the default-on gateway modes yield",
			changed: changedOf("enable-lb-svc"),
			mutate:  func(c *Configuration) { c.EnablePodLbSvc = true },
			want:    Configuration{EnableOvnLB: true, EnablePodLbSvc: true, EnableGwNftableLanipVip: true},
		},
		{
			name:    "enable-lb-svc pinned together with the gateway lb mode stays a validator error",
			changed: changedOf("enable-lb-svc", "enable-gw-nftable-lb-svc"),
			mutate:  func(c *Configuration) { c.EnablePodLbSvc = true },
			// gw-lb-svc stays on (it was pinned, so the validator must reject the pair); the
			// unpinned cluster-ip mode yields to the pinned enable-lb-svc.
			want:         Configuration{EnableOvnLB: true, EnablePodLbSvc: true, EnableGwNftableLbSvc: true, EnableGwNftableLanipVip: true},
			wantValidErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			cfg := defaults()
			if tt.mutate != nil {
				tt.mutate(&cfg)
			}
			resolveServiceFeatureGateDefaults(&cfg, tt.changed)
			require.Equal(t, tt.want, cfg)
			err := cfg.validateServiceFeatureGates()
			if tt.wantValidErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
