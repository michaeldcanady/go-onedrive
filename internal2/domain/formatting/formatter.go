package formatting

import "io"

type Formatter interface {
    Format(w io.Writer, v any) error
}
