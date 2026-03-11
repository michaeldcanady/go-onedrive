package get

import (
	"context"
	"strings"
	"time"

	"github.com/michaeldcanady/go-onedrive/internal/cli/util"
	"github.com/michaeldcanady/go-onedrive/internal/common/formatting"
	domainlogger "github.com/michaeldcanady/go-onedrive/internal/core/logger/domain"
	didomain "github.com/michaeldcanady/go-onedrive/internal/di/domain"
	domaindrive "github.com/michaeldcanady/go-onedrive/internal/drive/domain"
)

type GetCmd struct {
	util.BaseCommand
}

func NewGetCmd(container didomain.Container) *GetCmd {
	return &GetCmd{
		BaseCommand: util.NewBaseCommand(container, commandName),
	}
}

func (c *GetCmd) Run(ctx context.Context, opts Options) error {
	start := time.Now()

	if err := c.Initialize(loggerID); err != nil {
		return util.NewCommandError(c.Name, "failed to initialize command", err)
	}

	id := strings.ToLower(strings.TrimSpace(opts.DriveIDOrAlias))
	if id == "" {
		c.Log.Warn("id is empty", domainlogger.String("command", c.Name))
		return util.NewCommandErrorWithNameWithMessage(c.Name, "id is empty")
	}

	c.Log.Info("retrieving drive details", domainlogger.String("target", id))

	drive, err := c.Container.Drive().ResolveDrive(ctx, id)
	if err != nil {
		c.Log.Warn("failed to retrieve drive", domainlogger.Error(err), domainlogger.String("target", id))
		return util.NewCommandErrorWithNameWithError(c.Name, err)
	}

	columns := []formatting.Column[*domaindrive.Drive]{
		formatting.NewColumn("ID", func(item *domaindrive.Drive) string { return item.ID }),
		formatting.NewColumn("Name", func(item *domaindrive.Drive) string { return item.Name }),
		formatting.NewColumn("Owner", func(item *domaindrive.Drive) string { return string(item.Owner) }),
		formatting.NewColumn("ReadOnly", func(item *domaindrive.Drive) string {
			if item.ReadOnly {
				return "true"
			}
			return "false"
		}),
		formatting.NewColumn("Type", func(item *domaindrive.Drive) string { return item.Type.String() }),
	}

	formatter := formatting.NewTableFormatter(columns...)
	if err := formatter.Format(opts.Stdout, []*domaindrive.Drive{drive}); err != nil {
		c.Log.Warn("failed to format output", domainlogger.Error(err))
		return util.NewCommandErrorWithNameWithError(c.Name, err)
	}

	c.Log.Info("drive get completed successfully",
		domainlogger.Duration("duration", time.Since(start)),
	)
	return nil
}
