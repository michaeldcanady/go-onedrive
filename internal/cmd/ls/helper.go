package ls

import (
	"cmp"
	"errors"
	"fmt"
	"os"
	"slices"
	"strings"

	driveservice "github.com/michaeldcanady/go-onedrive/internal/app/file_service"
	"golang.org/x/term"
)

func filterHiddenDomain(items []*driveservice.DriveItem) []*driveservice.DriveItem {
	out := items[:0]
	for _, it := range items {
		if !strings.HasPrefix(it.Name, ".") {
			out = append(out, it)
		}
	}
	return out
}

func sortDomainItems(items []*driveservice.DriveItem) {
	slices.SortFunc(items, func(a, b *driveservice.DriveItem) int {
		return cmp.Compare(a.Name, b.Name)
	})
}

func detectTerminalWidth() int {
	w, _, err := term.GetSize(int(os.Stdout.Fd()))
	if err != nil {
		return 80
	}
	return w
}

func handleDomainError(op, path string, err error) error {
	var derr *driveservice.DomainError

	if errors.As(err, &derr) {
		switch derr.Kind {
		case driveservice.ErrNotFound:
			return fmt.Errorf("%s: '%s' not found", op, path)

		case driveservice.ErrNotFolder:
			return fmt.Errorf("%s: '%s' is not a folder", op, path)

		case driveservice.ErrUnauthorized:
			return fmt.Errorf("%s: unauthorized (are you logged in?)", op)

		case driveservice.ErrForbidden:
			return fmt.Errorf("%s: access denied to '%s'", op, path)

		case driveservice.ErrConflict:
			return fmt.Errorf("%s: conflict accessing '%s'", op, path)

		case driveservice.ErrPrecondition:
			return fmt.Errorf("%s: precondition failed for '%s'", op, path)

		case driveservice.ErrTransient:
			return fmt.Errorf("%s: temporary service issue, try again later", op)
		}
	}

	// Fallback for unexpected errors
	return fmt.Errorf("%s: unexpected error: %w", op, err)
}
