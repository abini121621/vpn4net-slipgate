package profile

import (
	"encoding/base64"
	"fmt"
	"strconv"
	"strings"

	"github.com/anonvector/slipgate/internal/config"
)

const (
	Scheme  = "vpnfornet://secure/"
	Version = "28"
)

// DraftInput describes the public server details needed to create an unlocked
// v28 profile draft. Locked/encrypted profiles must still be produced by the
// trusted Admin app; this package deliberately contains no Android key.
type DraftInput struct {
	Name     string
	Domain   string
	Username string
	Password string
	SSHHost  string
	SSHPort  int
}

// DraftFromTunnel creates a v28-compatible, base64-wrapped profile draft for
// a managed tunnel. It is intentionally not a locked .v4n payload.
func DraftFromTunnel(t config.TunnelConfig, in DraftInput) (string, error) {
	if strings.TrimSpace(in.Name) == "" {
		in.Name = t.Tag
	}
	if strings.TrimSpace(in.Domain) == "" && t.Transport != config.TransportSSH && t.Transport != config.TransportSOCKS {
		return "", fmt.Errorf("profile domain is required")
	}

	mode, err := modeForTransport(t.Transport)
	if err != nil {
		return "", err
	}
	sshPort := in.SSHPort
	if sshPort == 0 {
		sshPort = 22
	}
	sshHost := in.SSHHost
	if sshHost == "" {
		sshHost = in.Domain
	}

	// v28 has 79 pipe-delimited fields (indexes 0..78). Empty values are
	// retained because field positions are part of the Android contract.
	f := make([]string, 79)
	f[0] = Version
	f[1] = mode
	f[2] = sanitize(in.Name)
	f[3] = sanitize(in.Domain)
	f[6] = "30"
	f[8] = strconv.Itoa(t.Port)
	f[14] = boolString(usesSSH(t.Transport))
	f[15] = sanitize(in.Username)
	f[16] = sanitize(in.Password)
	f[17] = strconv.Itoa(sshPort)
	f[19] = sanitize(sshHost)
	f[29] = "443"
	f[47] = "0"
	f[51] = "0"
	f[57] = "0"
	f[58] = "0"
	f[59] = "0"
	f[60] = "0"
	f[61] = "0"
	f[62] = "0"
	f[63] = "0"
	f[64] = "0"
	f[65] = "0"
	f[66] = "0"
	f[67] = "0"
	f[68] = "0"
	f[69] = "0"
	f[70] = "0"
	f[71] = ""
	f[72] = "0"
	f[73] = "0"
	f[74] = "0"
	f[75] = "0"
	f[76] = ""
	f[77] = "0"
	f[78] = ""

	data := strings.Join(f, "|")
	return Scheme + base64.RawStdEncoding.EncodeToString([]byte(data)), nil
}

func modeForTransport(transport string) (string, error) {
	switch transport {
	case config.TransportSSH, config.TransportStunTLS:
		return "ssh", nil
	case config.TransportSOCKS:
		return "socks5", nil
	case config.TransportDNSTT:
		return "dnstt", nil
	case config.TransportSlipstream:
		return "ss", nil
	case config.TransportVayDNS:
		return "vaydns", nil
	case config.TransportNaive:
		return "naive", nil
	default:
		return "", fmt.Errorf("transport %q has no v28 mapping", transport)
	}
}

func usesSSH(transport string) bool {
	return transport == config.TransportSSH || transport == config.TransportStunTLS
}

func boolString(v bool) string {
	if v {
		return "1"
	}
	return "0"
}

func sanitize(v string) string {
	return strings.ReplaceAll(v, "|", "")
}
