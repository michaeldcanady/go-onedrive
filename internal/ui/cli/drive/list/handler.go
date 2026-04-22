package list

import (
	"fmt"
	"reflect"
	"slices"

	"github.com/michaeldcanady/go-onedrive/internal/drive"
	"github.com/michaeldcanady/go-onedrive/internal/features/logger"
	formatting "github.com/michaeldcanady/go-onedrive/pkg/format"
)

var supportedFormats = []formatting.Format{
	formatting.FormatJSON,
	formatting.FormatYAML,
	formatting.FormatTable,
}

// Command executes the drive list operation.
type Command struct {
	drive            drive.Service
	formatterFactory *formatting.FormatterFactory
	log              logger.Logger
}

// NewCommand initializes a new instance of the drive list Command.
func NewCommand(d drive.Service, ff *formatting.FormatterFactory, l logger.Logger) *Command {
	return &Command{
		drive:            d,
		formatterFactory: ff,
		log:              l,
	}
}

// Validate prepares and validates the options for the drive list operation.
func (c *Command) Validate(ctx *CommandContext) error {
	if ctx.Format = formatting.NewFormat(ctx.Options.Format); ctx.Format == formatting.FormatUnknown {
		return fmt.Errorf("unknown format: %s; expected %s", ctx.Options.Format, supportedFormats)
	}

	if !slices.Contains(supportedFormats, ctx.Format) {
		return fmt.Errorf("unsupported format: %s; expected %s", ctx.Format, supportedFormats)
	}

	return nil
}

func init() {
	formatting.GlobalRegistry.RegisterTable(reflect.TypeOf(drive.Drive{}), []formatting.Column{
		formatting.NewColumn("Name", func(i any) string {
			d := i.(drive.Drive)
			return "  " + d.Name
		}),
		formatting.NewColumn("ID", func(i any) string {
			return i.(drive.Drive).ID
		}),
		formatting.NewColumn("Type", func(i any) string {
			return i.(drive.Drive).Type
		}),
	})
}

// Execute retrieves and displays all available OneDrive drives.
func (c *Command) Execute(ctx *CommandContext) error {
	log := c.log.WithContext(ctx.Ctx)

	log.Debug("fetching all drives", logger.String("identity", ctx.Options.IdentityID))
	drives, err := c.drive.ListDrives(ctx.Ctx, ctx.Options.IdentityID)
	if err != nil {
		log.Error("failed to list drives", logger.Error(err))
		return fmt.Errorf("failed to list drives: %w", err)
	}

	// Prepare data for formatting
	var items []any
	for _, d := range drives {
		items = append(items, d)
	}

	// Resolve formatter
	formatter, err := c.formatterFactory.Create(ctx.Format)
	if err != nil {
		return err
	}

	log.Info("drives retrieved successfully", logger.Int("count", len(drives)))
	return formatter.Format(ctx.Options.Stdout, items)
}

// Finalize performs any necessary cleanup after the drive list operation.
func (c *Command) Finalize(ctx *CommandContext) error {
	return nil
}
