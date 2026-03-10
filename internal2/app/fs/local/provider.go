package local

import (
	"context"
	"errors"
	"io"
	"os"
	"path/filepath"

	"github.com/michaeldcanady/go-onedrive/internal2/domain/common/logger"
	domainfs "github.com/michaeldcanady/go-onedrive/internal2/domain/fs"
	"github.com/michaeldcanady/go-onedrive/internal2/interface/cli/util"
)

var _ domainfs.Service = (*Provider)(nil)

type Provider struct {
	log logger.Logger
}

func NewProvider(l logger.Logger) *Provider {
	return &Provider{
		log: l,
	}
}

func (p *Provider) buildLogger(ctx context.Context) logger.Logger {
	correlationID := util.CorrelationIDFromContext(ctx)
	return p.log.WithContext(ctx).With(
		logger.String("correlation_id", correlationID),
	)
}

func (p *Provider) Get(ctx context.Context, path string) (domainfs.Item, error) {
	log := p.buildLogger(ctx).With(logger.String("path", path))
	log.Debug("local.Get")

	info, err := os.Stat(path)
	if err != nil {
		return domainfs.Item{}, err
	}
	return p.mapInfoToItem(path, info), nil
}

func (p *Provider) List(ctx context.Context, path string, opts domainfs.ListOptions) ([]domainfs.Item, error) {
	log := p.buildLogger(ctx).With(logger.String("path", path))
	log.Debug("local.List", logger.Bool("recursive", opts.Recursive))

	entries, err := os.ReadDir(path)
	if err != nil {
		return nil, err
	}

	var items []domainfs.Item
	for _, entry := range entries {
		info, err := entry.Info()
		if err != nil {
			continue
		}
		item := p.mapInfoToItem(filepath.Join(path, entry.Name()), info)
		items = append(items, item)

		if opts.Recursive && entry.IsDir() {
			children, err := p.List(ctx, filepath.Join(path, entry.Name()), opts)
			if err == nil {
				items = append(items, children...)
			}
		}
	}
	return items, nil
}

func (p *Provider) Stat(ctx context.Context, path string, opts domainfs.StatOptions) (domainfs.Item, error) {
	return p.Get(ctx, path)
}

func (p *Provider) ReadFile(ctx context.Context, path string, opts domainfs.ReadOptions) (io.ReadCloser, error) {
	log := p.buildLogger(ctx).With(logger.String("path", path))
	log.Debug("local.ReadFile")

	return os.Open(path)
}

func (p *Provider) WriteFile(ctx context.Context, path string, r io.Reader, opts domainfs.WriteOptions) (domainfs.Item, error) {
	log := p.buildLogger(ctx).With(logger.String("path", path))
	log.Debug("local.WriteFile", logger.Bool("overwrite", opts.Overwrite))

	flag := os.O_WRONLY | os.O_CREATE
	if opts.Overwrite {
		flag |= os.O_TRUNC
	} else {
		flag |= os.O_EXCL
	}

	f, err := os.OpenFile(path, flag, 0644)
	if err != nil {
		return domainfs.Item{}, err
	}
	defer f.Close()

	_, err = io.Copy(f, r)
	if err != nil {
		return domainfs.Item{}, err
	}

	return p.Get(ctx, path)
}

func (p *Provider) Mkdir(ctx context.Context, path string, opts domainfs.MKDirOptions) error {
	log := p.buildLogger(ctx).With(logger.String("path", path))
	log.Debug("local.Mkdir", logger.Bool("parents", opts.Parents))

	if opts.Parents {
		return os.MkdirAll(path, 0755)
	}
	return os.Mkdir(path, 0755)
}

func (p *Provider) Remove(ctx context.Context, path string, opts domainfs.RemoveOptions) error {
	log := p.buildLogger(ctx).With(logger.String("path", path))
	log.Debug("local.Remove")

	return os.RemoveAll(path)
}

func (p *Provider) Copy(ctx context.Context, src, dst string, opts domainfs.CopyOptions) error {
	log := p.buildLogger(ctx).With(logger.String("src", src), logger.String("dst", dst))
	log.Debug("local.Copy")

	srcItem, err := p.Stat(ctx, src, domainfs.StatOptions{})
	if err != nil {
		return err
	}

	if srcItem.Type == domainfs.ItemTypeFolder {
		if !opts.Recursive {
			return errors.New("source is a directory, use recursive flag")
		}

		if err := p.Mkdir(ctx, dst, domainfs.MKDirOptions{Parents: true}); err != nil {
			return err
		}

		children, err := p.List(ctx, src, domainfs.ListOptions{Recursive: false})
		if err != nil {
			return err
		}

		for _, child := range children {
			childSrc := child.Path
			childDst := filepath.Join(dst, child.Name)
			if err := p.Copy(ctx, childSrc, childDst, opts); err != nil {
				return err
			}
		}
		return nil
	}

	r, err := p.ReadFile(ctx, src, domainfs.ReadOptions{})
	if err != nil {
		return err
	}
	defer r.Close()

	_, err = p.WriteFile(ctx, dst, r, domainfs.WriteOptions{
		Overwrite: opts.Overwrite,
	})
	return err
}

func (p *Provider) Move(ctx context.Context, src, dst string, opts domainfs.MoveOptions) error {
	log := p.buildLogger(ctx).With(logger.String("src", src), logger.String("dst", dst))
	log.Debug("local.Move")

	return os.Rename(src, dst)
}

func (p *Provider) Upload(ctx context.Context, src, dst string, opts domainfs.UploadOptions) (domainfs.Item, error) {
	log := p.buildLogger(ctx).With(logger.String("src", src), logger.String("dst", dst))
	log.Debug("local.Upload")

	f, err := os.Open(src)
	if err != nil {
		return domainfs.Item{}, err
	}
	defer f.Close()

	return p.WriteFile(ctx, dst, f, domainfs.WriteOptions{Overwrite: opts.Overwrite})
}

func (p *Provider) Touch(ctx context.Context, path string, opts domainfs.TouchOptions) (domainfs.Item, error) {
	log := p.buildLogger(ctx).With(logger.String("path", path))
	log.Debug("local.Touch")

	f, err := os.OpenFile(path, os.O_RDONLY|os.O_CREATE, 0644)
	if err != nil {
		return domainfs.Item{}, err
	}
	f.Close()
	return p.Get(ctx, path)
}

func (p *Provider) mapInfoToItem(path string, info os.FileInfo) domainfs.Item {
	itemType := domainfs.ItemTypeFile
	if info.IsDir() {
		itemType = domainfs.ItemTypeFolder
	}

	return domainfs.Item{
		Path:     path,
		Name:     info.Name(),
		Size:     info.Size(),
		Type:     itemType,
		Modified: info.ModTime(),
	}
}
