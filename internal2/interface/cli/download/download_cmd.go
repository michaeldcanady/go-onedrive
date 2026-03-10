package download

import (
	"context"
	"strings"

	"github.com/michaeldcanady/go-onedrive/internal2/domain/di"
	"github.com/michaeldcanady/go-onedrive/internal2/interface/cli/cp"
	"github.com/michaeldcanady/go-onedrive/internal2/interface/cli/util"
)

type DownloadCmd struct {
	util.BaseCommand
}

func NewDownloadCmd(container di.Container) *DownloadCmd {
	return &DownloadCmd{
		BaseCommand: util.NewBaseCommand(container, commandName),
	}
}

func (c *DownloadCmd) Run(ctx context.Context, opts Options) error {
	src := opts.Source
	if !strings.HasPrefix(src, "onedrive:") && !strings.HasPrefix(src, "local:") {
		src = "onedrive:" + src
	}

	dst := opts.Destination
	if !strings.HasPrefix(dst, "onedrive:") && !strings.HasPrefix(dst, "local:") {
		dst = "local:" + dst
	}

	cpOpts := cp.Options{
		Source:    src,
		Dest:      dst,
		Overwrite: opts.Overwrite,
		Recursive: opts.Recursive,
		Stdin:     opts.Stdin,
		Stdout:    opts.Stdout,
		Stderr:    opts.Stderr,
	}

	cpCmd := cp.NewCpCmd(c.Container)
	return cpCmd.Run(ctx, cpOpts)
}
