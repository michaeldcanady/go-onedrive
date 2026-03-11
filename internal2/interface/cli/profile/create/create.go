package create

import (
	"context"
	"fmt"
	"strings"
	"time"

	logger "github.com/michaeldcanady/go-onedrive/internal2/domain/common/logger"
	"github.com/michaeldcanady/go-onedrive/internal2/domain/di"
	"github.com/michaeldcanady/go-onedrive/internal2/interface/cli/util"
	"github.com/spf13/cobra"
)

const (
	commandName = "create"
	loggerID    = "cli"

	setCurrentFlagLong  = "set-current"
	setCurrentFlagUsage = "Sets the new profile as current"

	forceFlagLong  = "force"
	forceFlagShort = "f"
	forceFlagUsage = "Overwrite an existing profile if it already exists"
)

type CreateCmd struct {
	util.BaseCommand
}

func NewCreateCmd(container di.Container) *CreateCmd {
	return &CreateCmd{
		BaseCommand: util.NewBaseCommand(container, commandName),
	}
}

func (c *CreateCmd) Run(ctx context.Context, opts Options, setCurrent, force bool) error {
	start := time.Now()

	if err := c.Initialize(loggerID); err != nil {
		return err
	}

	name := strings.ToLower(strings.TrimSpace(opts.Name))
	if name == "" {
		c.Log.Warn("profile name is empty")
		return util.NewCommandErrorWithNameWithMessage(c.Name, "name is empty")
	}

	c.Log.Info("checking if profile exists", logger.String("name", name))

	exists, err := c.Container.Profile().Exists(ctx, name)
	if err != nil {
		c.Log.Error("failed to check profile existence", logger.String("error", err.Error()))
		return util.NewCommandErrorWithNameWithError(c.Name, err)
	}

	// Handle existing profile
	if exists {
		if !force {
			c.Log.Warn("profile already exists", logger.String("name", name))
			return util.NewCommandErrorWithNameWithMessage(c.Name, "profile already exists")
		}

		c.Log.Warn("profile exists; force enabled, deleting existing profile",
			logger.String("name", name),
		)

		if err := c.Container.Profile().Delete(ctx, name); err != nil {
			c.Log.Error("failed to delete existing profile", logger.String("error", err.Error()))
			return util.NewCommandErrorWithNameWithError(c.Name, err)
		}
	}

	c.Log.Info("creating profile", logger.String("name", name))

	p, err := c.Container.Profile().Create(ctx, name)
	if err != nil {
		c.Log.Error("failed to create profile", logger.String("error", err.Error()))
		return util.NewCommandErrorWithNameWithError(c.Name, err)
	}

	if setCurrent {
		c.Log.Info("setting new profile as current", logger.String("name", name))

		if err := c.Container.State().SetCurrentProfile(p.Name); err != nil {
			c.Log.Error("failed to set current profile", logger.String("error", err.Error()))
			return util.NewCommandErrorWithNameWithError(c.Name, err)
		}
	}

	c.Log.Info("profile created successfully",
		logger.String("name", p.Name),
		logger.String("path", p.Path),
	)

	fmt.Fprintf(opts.Stdout, "Created profile %q at %s\n", p.Name, p.Path)

	c.Log.Info("profile create completed successfully",
		logger.Duration("duration", time.Since(start)),
	)
	return nil
}

func CreateCreateCmd(container di.Container) *cobra.Command {
	var (
		setCurrent bool
		force      bool
	)

	cmd := &cobra.Command{
		Use:   "create <name>",
		Short: "Create a new profile",
		Args:  cobra.ExactArgs(1),

		RunE: func(cmd *cobra.Command, args []string) error {
			opts := Options{
				Name:   args[0],
				Stdout: cmd.OutOrStdout(),
				Stderr: cmd.ErrOrStderr(),
			}

			return NewCreateCmd(container).Run(cmd.Context(), opts, setCurrent, force)
		},
	}

	cmd.Flags().BoolVar(&setCurrent, setCurrentFlagLong, false, setCurrentFlagUsage)
	cmd.Flags().BoolVarP(&force, forceFlagLong, forceFlagShort, false, forceFlagUsage)

	return cmd
}
