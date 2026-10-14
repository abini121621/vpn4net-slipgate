package cmd

import (
	"fmt"
	"strings"

	v28 "github.com/anonvector/slipgate/internal/profile"
	"github.com/spf13/cobra"
)

func init() {
	var tag, name, domain, username, password, sshHost string
	var sshPort int

	profileCmd := &cobra.Command{
		Use:   "profile",
		Short: "Create v4.0.0-compatible client profile drafts",
	}
	draftCmd := &cobra.Command{
		Use:   "draft",
		Short: "Generate an unlocked vpnfornet://secure/ v28 draft",
		Long:  "Generate an unlocked v4.0.0 profile draft. The Admin app must lock/encrypt it before sharing as a .v4n file.",
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := loadOrDefaultConfig()
			if err != nil {
				return err
			}
			if strings.TrimSpace(tag) == "" {
				if len(cfg.Tunnels) != 1 {
					return fmt.Errorf("--tag is required when more than one tunnel is configured")
				}
				tag = cfg.Tunnels[0].Tag
			}
			tunnel := cfg.GetTunnel(tag)
			if tunnel == nil {
				return fmt.Errorf("tunnel %q not found", tag)
			}
			uri, err := v28.DraftFromTunnel(*tunnel, v28.DraftInput{
				Name: name, Domain: domain, Username: username, Password: password,
				SSHHost: sshHost, SSHPort: sshPort,
			})
			if err != nil {
				return err
			}
			fmt.Fprintln(cmd.OutOrStdout(), uri)
			fmt.Fprintln(cmd.ErrOrStderr(), "warning: this is an unlocked draft; do not distribute it as a production .v4n profile")
			return nil
		},
	}
	draftCmd.Flags().StringVar(&tag, "tag", "", "managed tunnel tag")
	draftCmd.Flags().StringVar(&name, "name", "", "profile name (defaults to tunnel tag)")
	draftCmd.Flags().StringVar(&domain, "domain", "", "server domain")
	draftCmd.Flags().StringVar(&username, "username", "", "SSH/SOCKS username")
	draftCmd.Flags().StringVar(&password, "password", "", "SSH/SOCKS password")
	draftCmd.Flags().StringVar(&sshHost, "ssh-host", "", "SSH host (defaults to profile domain)")
	draftCmd.Flags().IntVar(&sshPort, "ssh-port", 22, "SSH port")

	profileCmd.AddCommand(draftCmd)
	rootCmd.AddCommand(profileCmd)
}
