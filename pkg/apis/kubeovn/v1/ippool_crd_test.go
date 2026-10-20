package v1

import (
	"net"
	"net/netip"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

// ippoolCRDPaths lists the three checked-in copies of the IPPool CRD, relative to
// this package, that must carry the same spec validation.
var ippoolCRDPaths = []string{
	"../../../../charts/kube-ovn/templates/kube-ovn-crd.yaml",
	"../../../../charts/kube-ovn-v2/crds/kube-ovn-crd.yaml",
	"../../../../dist/images/install.sh",
}

// ippoolCRDBlock returns the ippools.kubeovn.io CRD document of one of the copies.
func ippoolCRDBlock(t *testing.T, path string) string {
	t.Helper()
	raw, err := os.ReadFile(filepath.Clean(path))
	require.NoError(t, err)
	content := string(raw)
	start := strings.Index(content, "name: ippools.kubeovn.io")
	require.Positive(t, start, "%s has no ippools CRD", path)
	if i := strings.LastIndex(content[:start], "apiVersion: apiextensions.k8s.io/v1"); i >= 0 {
		start = i
	}
	end := strings.Index(content[start:], "\n---")
	require.Positive(t, end, "%s: unterminated ippools CRD", path)
	return content[start : start+end]
}

// unquote drops the double quotes the hand-maintained CRD copies put around CEL
// rules, so assertions do not depend on the quoting style.
func unquote(block string) string {
	return strings.ReplaceAll(block, `"`, "")
}

// TestIPPoolCRDSpecValidation guards the spec validations of the checked-in IPPool
// CRDs against silently dropping constraints again when the CRD is regenerated: the
// regenerated schemas once lost the subnet immutability rule, minItems and the set
// list types.
func TestIPPoolCRDSpecValidation(t *testing.T) {
	for _, path := range ippoolCRDPaths {
		block := ippoolCRDBlock(t, path)
		require.Contains(t, unquote(block), "rule: self == oldSelf", "%s: spec.subnet immutability rule missing", path)
		require.Contains(t, unquote(block), "rule: has(self.ips) != has(self.ipCount)", "%s: spec.ips/spec.ipCount exclusivity rule missing", path)
		require.Contains(t, block, "minItems: 1", "%s: spec.ips minItems missing", path)
		// keep in sync with maxIPPoolIPCount in pkg/controller, which guards pools
		// created before this bound existed; pkg/apis cannot import pkg/controller
		require.Contains(t, block, "maximum: 32", "%s: spec.ipCount upper bound missing", path)
		require.Equal(t, 2, strings.Count(block, "x-kubernetes-list-type: set"), "%s: spec.ips and spec.namespaces must be set-typed", path)
	}
}

// TestIPPoolCRDIPsItemValidation exercises the spec.ips entry validation of the
// checked-in IPPool CRDs. The schema accepts an entry when any of its items.anyOf
// arms matches: an IPv4/IPv6 address, a CIDR, or a start..end range; mixed CIDR plus
// range spellings stay rejected. The arms are format/pattern based - not CEL, whose
// cost budget a regex this size exceeds - so here they are re-evaluated with the same
// primitives (netip/net parsing like the apiserver's strfmt formats, RE2 patterns).
func TestIPPoolCRDIPsItemValidation(t *testing.T) {
	valid := []string{
		"10.0.0.1", "0.0.0.0", "255.255.255.255", "10.0.0.0/24", "0.0.0.0/0", "10.0.0.0/32",
		"10.0.0.1..10.0.0.9", "10.0.0.01..10.0.0.09", "192.168.0.0/16",
		"2001:db8::1", "::1", "::", "fe80::1/64", "2001:db8::/32", "2001:db8::/0", "2001:db8::/128",
		"2001:db8::1..2001:db8::9", "ffff:ffff:ffff:ffff:ffff:ffff:ffff:ffff",
		"1:2:3:4:5:6:7:8", "ff02::1", "2001:db8::ffff..2001:db8::1:0",
		// format: ipv6 accepts IPv4-mapped spellings; the controller treats them as IPv4
		"::ffff:10.0.0.1",
	}
	invalid := []string{
		"garbage", "10.0.0.256", "10.0.0.1/33", "10.0.0.1..10.0.0.256", "10.0.0.1.", "1.2.3",
		"2001:db8::1/129", "2001:db8::1::", "10.0.0.1..10.0.0.1..10.0.0.1", "..", "/24", "10.0.0.1/",
		"12345::",
		// a CIDR is not a range start or end
		"10.0.0.1/24..10.0.0.9", "10.0.0.1..10.0.0.9/32", "10.0.0.0/24/25",
		"2001:db8::/64..2001:db8::9", "2001:db8::1/128..2001:db8::ff",
		// ranges stay within one IP family
		"2001:db8::1..10.0.0.1", "10.0.0.1..2001:db8::1",
	}

	type arm struct {
		format  string
		pattern *regexp.Regexp
	}
	// the anyOf arms appear only under spec.ips, so scanning the whole block avoids
	// depending on the key order and indentation of the CRD copies
	parseArms := func(block string) []arm {
		var arms []arm
		for _, line := range strings.Split(block, "\n") {
			line = strings.TrimSpace(line)
			switch {
			case strings.HasPrefix(line, "- format: "):
				arms = append(arms, arm{format: strings.TrimPrefix(line, "- format: ")})
			case strings.HasPrefix(line, "- pattern: "):
				arms = append(arms, arm{pattern: regexp.MustCompile(strings.TrimPrefix(line, "- pattern: "))})
			}
		}
		return arms
	}

	formatOK := func(format, s string) bool {
		switch format {
		case "ipv4":
			addr, err := netip.ParseAddr(s)
			return err == nil && addr.Is4()
		case "ipv6":
			_, err := netip.ParseAddr(s)
			return err == nil && strings.Contains(s, ":")
		case "cidr":
			_, _, err := net.ParseCIDR(s)
			return err == nil
		}
		return false
	}

	for _, path := range ippoolCRDPaths {
		arms := parseArms(ippoolCRDBlock(t, path))
		require.Len(t, arms, 5, "%s: spec.ips must validate entries as one of ipv4/ipv6/cidr or an address range", path)
		accepted := func(s string) bool {
			for _, a := range arms {
				if a.pattern != nil && a.pattern.MatchString(s) {
					return true
				}
				if a.format != "" && formatOK(a.format, s) {
					return true
				}
			}
			return false
		}
		for _, s := range valid {
			require.True(t, accepted(s), "%s: %q must be accepted", path, s)
		}
		for _, s := range invalid {
			require.False(t, accepted(s), "%s: %q must be rejected", path, s)
		}
	}
}
