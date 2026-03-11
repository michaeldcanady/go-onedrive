package list

import (
	"context"
	"sort"
	"time"

	"github.com/michaeldcanady/go-onedrive/internal/cli/util"
	"github.com/michaeldcanady/go-onedrive/internal/common/formatting"
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

type aliasItem struct {
	Alias   string
	DriveID string
}

func (c *ListCmd) Run(ctx context.Context, opts Options) error {
	start := time.Now()

	if err := c.Initialize(loggerID); err != nil {
		return util.NewCommandError(c.Name, "failed to initialize command", err)
	}

	c.Log.Info("listing drive aliases")

	aliases, err := c.Container.State().ListDriveAliases()
	if err != nil {
		c.Log.Error("failed to list drive aliases",
			domainlogger.Error(err),
		)
		return util.NewCommandError(c.Name, "failed to list drive aliases", err)
	}

	if len(aliases) == 0 {
		c.RenderInfo(opts.Stdout, "no drive aliases found")
		return nil
	}

	items := make([]aliasItem, 0, len(aliases))
	for k, v := range aliases {
		items = append(items, aliasItem{Alias: k, DriveID: v})
	}

	sort.Slice(items, func(i, j int) bool {
		return items[i].Alias < items[j].Alias
	})

	formatter := formatting.NewTableFormatter(
		formatting.NewColumn("ALIAS", func(i aliasItem) string { return i.Alias }),
		formatting.NewColumn("DRIVE ID", func(i aliasItem) string { return i.DriveID }),
	)

	if err := formatter.Format(opts.Stdout, items); err != nil {
		return util.NewCommandError(c.Name, "failed to format output", err)
	}

	c.Log.Info("drive alias list completed successfully",
		domainlogger.Duration("duration", time.Since(start)),
	)
	return nil
}
