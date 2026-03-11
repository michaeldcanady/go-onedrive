package cp

import (
	"context"
	"fmt"
	"os"
	"path"
	"time"

	"github.com/michaeldcanady/go-onedrive/internal/cli/util"
	domainlogger "github.com/michaeldcanady/go-onedrive/internal/core/logger/domain"
	didomain "github.com/michaeldcanady/go-onedrive/internal/di/domain"
	domainfs "github.com/michaeldcanady/go-onedrive/internal/fs/domain"
)

type CpCmd struct {
	util.BaseCommand
}

func NewCpCmd(container didomain.Container) *CpCmd {
	return &CpCmd{
		BaseCommand: util.NewBaseCommand(container, commandName),
	}
}

func (c *CpCmd) Run(ctx context.Context, opts Options) error {
	start := time.Now()

	if err := c.Initialize(loggerID); err != nil {
		return err
	}

	c.Log.Info("starting cp command",
		domainlogger.String("src", opts.Source),
		domainlogger.String("dst", opts.Dest),
		domainlogger.Bool("overwrite", opts.Overwrite),
		domainlogger.String("ignoreFile", opts.IgnoreFile),
	)

	fsSvc := c.Container.FS()
	if fsSvc == nil {
		return util.NewCommandErrorWithNameWithMessage(c.Name, "filesystem service is nil")
	}

	matcher, err := c.loadIgnoreMatcher(ctx, opts.IgnoreFile)
	if err != nil {
		c.Log.Warn("failed to load ignore file", domainlogger.String("path", opts.IgnoreFile), domainlogger.Error(err))
	}

	if matcher != nil {
		err = c.copyRecursive(ctx, fsSvc, opts.Source, opts.Dest, opts.Overwrite, matcher)
	} else {
		err = fsSvc.Copy(ctx, opts.Source, opts.Dest, domainfs.CopyOptions{Overwrite: opts.Overwrite})
	}

	if err != nil {
		c.RenderError(opts.Stderr, err)
		return util.NewCommandError(c.Name, "failed to copy item", err)
	}

	c.Log.Info("cp completed successfully",
		domainlogger.Duration("duration", time.Since(start)),
	)

	fmt.Fprintf(opts.Stdout, "Successfully copied \"%s\" to \"%s\"\n", opts.Source, opts.Dest)

	return nil
}

func (c *CpCmd) loadIgnoreMatcher(ctx context.Context, ignorePath string) (domainfs.IgnoreMatcher, error) {
	if ignorePath == "" {
		return nil, nil
	}

	f, err := os.Open(ignorePath)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	factory := c.Container.IgnoreMatcherFactory()
	if factory == nil {
		return nil, nil
	}

	return factory.CreateMatcher(ctx, f)
}

func (c *CpCmd) copyRecursive(ctx context.Context, fsSvc domainfs.Service, src, dst string, overwrite bool, matcher domainfs.IgnoreMatcher) error {
	item, err := fsSvc.Get(ctx, src)
	if err != nil {
		return err
	}

	if matcher != nil && matcher.ShouldIgnore(item.Path, item.Type == domainfs.ItemTypeFolder) {
		c.Log.Debug("ignoring item", domainlogger.String("path", item.Path))
		return nil
	}

	if item.Type == domainfs.ItemTypeFile {
		return fsSvc.Copy(ctx, src, dst, domainfs.CopyOptions{Overwrite: overwrite})
	}

	// It's a folder, ensure destination exists
	if err := fsSvc.Mkdir(ctx, dst, domainfs.MKDirOptions{Parents: true}); err != nil {
		// Ignore error if it already exists? FS implementation should handle it or we check here.
		c.Log.Debug("mkdir destination", domainlogger.String("path", dst), domainlogger.Error(err))
	}

	children, err := fsSvc.List(ctx, src, domainfs.ListOptions{})
	if err != nil {
		return err
	}

	for _, child := range children {
		// Better way: use item names
		// childDst := path.Join(dst, child.Name)
		// How to get provider-aware child path?
		// For now, assume simple join if it's the same provider

		// If src is "local:./foo", and child.Name is "bar.txt"
		// we want "local:./foo/bar.txt"

		provider, subPath := util.ParsePath(src)
		newSrc := fmt.Sprintf("%s:%s", provider, path.Join(subPath, child.Name))

		dstProvider, dstSubPath := util.ParsePath(dst)
		newDst := fmt.Sprintf("%s:%s", dstProvider, path.Join(dstSubPath, child.Name))

		if err := c.copyRecursive(ctx, fsSvc, newSrc, newDst, overwrite, matcher); err != nil {
			return err
		}
	}

	return nil
}
