package list

import (
	"context"
	"fmt"
	"sort"
	"time"

	"github.com/michaeldcanady/go-onedrive/internal2/domain/di"
	logger "github.com/michaeldcanady/go-onedrive/internal2/domain/common/logger"
	"github.com/michaeldcanady/go-onedrive/internal2/interface/cli/util"
)

type ListCmd struct {
	util.BaseCommand
}

func NewListCmd(container di.Container) *ListCmd {
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
			logger.Error(err),
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
		logger.Duration("duration", time.Since(start)),
	)
	return nil
}
