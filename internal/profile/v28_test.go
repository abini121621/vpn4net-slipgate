package profile

import (
	"encoding/base64"
	"strings"
	"testing"

	"github.com/anonvector/slipgate/internal/config"
)

func TestDraftFromTunnelV28Layout(t *testing.T) {
	draft, err := DraftFromTunnel(config.TunnelConfig{
		Tag:       "edge",
		Transport: config.TransportSSH,
		Port:      443,
	}, DraftInput{
		Name:     "edge|primary",
		Domain:   "vpn.example.test",
		Username: "alice",
		Password: "secret",
	})
	if err != nil {
		t.Fatal(err)
	}
	if !strings.HasPrefix(draft, Scheme) {
		t.Fatalf("scheme = %q", draft)
	}
	decoded, err := base64.RawStdEncoding.DecodeString(strings.TrimPrefix(draft, Scheme))
	if err != nil {
		t.Fatal(err)
	}
	fields := strings.Split(string(decoded), "|")
	if len(fields) != 79 {
		t.Fatalf("field count = %d, want 79", len(fields))
	}
	if fields[0] != Version || fields[1] != "ssh" {
		t.Fatalf("header = %q/%q", fields[0], fields[1])
	}
	if fields[2] != "edgeprimary" {
		t.Fatalf("name was not sanitized: %q", fields[2])
	}
	if fields[15] != "alice" || fields[16] != "secret" {
		t.Fatalf("SSH credentials not mapped")
	}
}

func TestDraftRejectsMissingDomain(t *testing.T) {
	_, err := DraftFromTunnel(config.TunnelConfig{Tag: "dns", Transport: config.TransportDNSTT}, DraftInput{})
	if err == nil {
		t.Fatal("expected missing-domain error")
	}
}
