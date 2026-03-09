package list

import (
	"context"
	"fmt"
	"sort"

	"github.com/michaeldcanady/go-onedrive/internal2/domain/di"
	infralogging "github.com/michaeldcanady/go-onedrive/internal2/infra/common/logging"
	"github.com/michaeldcanady/go-onedrive/internal2/interface/cli/util"
)

type Command struct {
	container di.Container
	logger    infralogging.Logger
}

func NewCmd(container di.Container) *Command {
	return &Command{
		container: container,
	}
}

func (c *Command) Run(ctx context.Context, opts Options) error {
	if ctx == nil {
		ctx = context.Background()
	}

	if c.logger == nil {
		logger, err := util.EnsureLogger(c.container, loggerID)
		if err != nil {
			return util.NewCommandErrorWithNameWithError(commandName, err)
		}
		c.logger = logger
	}

	c.logger.Info("listing drive aliases")

	aliases, err := c.container.State().ListDriveAliases()
	if err != nil {
		c.logger.Error("failed to list drive aliases",
			infralogging.Error(err),
		)
		return util.NewCommandError(commandName, "failed to list drive aliases", err)
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

	return nil
}
