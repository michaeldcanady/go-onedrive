package list

import (
	"context"
	"fmt"
	"sort"
	"time"

	"github.com/michaeldcanady/go-onedrive/internal/cli/util"
	domainlogger "github.com/michaeldcanady/go-onedrive/internal/core/logger/domain"
	didomain "github.com/michaeldcanady/go-onedrive/internal/di/domain"
)

type ListCmd struct {
	util.BaseCommand
}

func NewListCmd(container didomain.Container) *ListCmd {
	return &ListCmd{
		BaseCommand: util.NewBaseCommand(container, commandName),
	}
}

func (c *ListCmd) Run(ctx context.Context, opts Options) error {
	start := time.Now()

	if err := c.Initialize(loggerID); err != nil {
		return err
	}

	c.Log.Info("listing drive aliases")

	aliases, err := c.Container.State().ListDriveAliases()
	if err != nil {
		c.Log.Error("failed to list drive aliases",
			domainlogger.Error(err),
		)
		c.RenderError(opts.Stderr, err)
		return util.NewCommandError(c.Name, "failed to list drive aliases", err)
	}

	if len(aliases) == 0 {
		fmt.Fprintln(opts.Stdout, "No drive aliases found.")
		return nil
	}

	// Sort aliases for consistent output
	keys := make([]string, 0, len(aliases))
	for k := range aliases {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	fmt.Fprintf(opts.Stdout, "%-20s %-20s\n", "ALIAS", "DRIVE ID")
	fmt.Fprintf(opts.Stdout, "%-20s %-20s\n", "-----", "--------")
	for _, k := range keys {
		fmt.Fprintf(opts.Stdout, "%-20s %-20s\n", k, aliases[k])
	}

	c.Log.Info("drive alias list completed successfully",
		domainlogger.Duration("duration", time.Since(start)),
	)
	return nil
}
