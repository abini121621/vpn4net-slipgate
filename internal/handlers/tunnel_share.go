package handlers

import (
	"github.com/anonvector/slipgate/internal/actions"
	"github.com/anonvector/slipgate/internal/config"
	"github.com/anonvector/slipgate/internal/profile"
	"github.com/anonvector/slipgate/internal/prompt"
)

func handleTunnelShare(ctx *actions.Context) error {
	cfg := ctx.Config.(*config.Config)
	tag := ctx.GetArg("tag")

	if tag == "" {
		return actions.NewError(actions.TunnelShare, "tunnel tag is required", nil)
	}

	tunnel := cfg.GetTunnel(tag)
	if tunnel == nil {
		return actions.NewErrorWithHint(actions.TunnelShare, "tunnel not found",
			"Run 'xopapi tunnel status' to see available tunnels", nil)
	}

	var username, password string

	// Ask which user's credentials to embed
	if len(cfg.Users) > 0 {
		userOpts := make([]actions.SelectOption, 0, len(cfg.Users)+1)
		userOpts = append(userOpts, actions.SelectOption{Value: "", Label: "No credentials"})
		for _, u := range cfg.Users {
			userOpts = append(userOpts, actions.SelectOption{Value: u.Username, Label: u.Username})
		}
		username, err := prompt.Select("User", userOpts)
		if err != nil {
			return err
		}
		if user := cfg.GetUser(username); user != nil {
			username = user.Username
			password = user.Password
		}
	}

	uri, err := profile.DraftFromTunnel(*tunnel, profile.DraftInput{
		Name: tunnel.Tag, Domain: tunnel.Domain,
		Username: username, Password: password,
		SSHHost: tunnel.Domain, SSHPort: 22,
	})
	if err != nil {
		return actions.NewError(actions.TunnelShare, "failed to generate v28 profile draft", err)
	}

	ctx.Output.Print(uri)
	ctx.Output.Warning("This is an unlocked draft; lock it in the Admin app before sharing as .v4n")
	return nil
}
