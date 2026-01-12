package feed

import (
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
)

// WriteJSONLGz writes the feed in JSON Lines format, compressed with gzip.
func (f Feed) WriteJSONLGz(w io.Writer) error {
	gz := gzip.NewWriter(w)
	enc := json.NewEncoder(gz)
	enc.SetEscapeHTML(false)
	for i := range f {
		if err := enc.Encode(f[i]); err != nil {
			_ = gz.Close()
			return fmt.Errorf("encode jsonl: %w", err)
		}
	}
	if err := gz.Close(); err != nil {
		return fmt.Errorf("close gzip: %w", err)
	}
	return nil
}
