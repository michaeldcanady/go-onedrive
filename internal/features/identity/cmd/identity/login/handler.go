package login

import (
	"fmt"
)

// Validate performs initial validation of the command options.
func (c *Command) Validate(ctx *CommandContext) error {
	return nil
}

// Resolve performs argument resolution.
func (c *Command) Resolve(ctx *CommandContext) error {
	return c.BaseResolve(ctx)
}

// Execute performs the core business logic of the command.
func (c *Command) Execute(ctx *CommandContext) error {
	provider := ctx.Options.Provider
	if provider == "" {
		val, err := c.config.Get("auth.provider")
		if err != nil {
			provider = "azure" // ultimate fallback
		} else {
			provider = fmt.Sprintf("%v", val)
		}
	}

	options := make(map[string]string)

	// Helper to get from flags or config
	getOption := func(flag string, configKey string) string {
		if flag != "" {
			return flag
		}
		val, _ := c.config.Get(configKey)
		if val == nil {
			return ""
		}
		return fmt.Sprintf("%v", val)
	}

	options["method"] = getOption(ctx.Options.Method, "identity."+provider+".method")
	options["client_id"] = getOption(ctx.Options.ClientId, "identity."+provider+".client_id")
	options["tenant_id"] = getOption(ctx.Options.TenantId, "identity."+provider+".tenant_id")
	options["client_secret"] = getOption(ctx.Options.ClientSecret, "identity."+provider+".client_secret")
	options["scopes"] = getOption(ctx.Options.Scopes, "identity."+provider+".scopes")
	options["redirect_uri"] = getOption("", "identity."+provider+".redirect_uri")

	_, err := c.identity.Login(ctx.Ctx, provider, options)
	return err
}

// Finalize performs any cleanup or final output formatting.
func (c *Command) Finalize(ctx *CommandContext) error {
	fmt.Println("Login successful")
	return nil
}
