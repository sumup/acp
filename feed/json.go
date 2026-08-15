package feed

import (
	"compress/gzip"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"iter"
)

// Result contains either a decoded product or an error encountered while reading.
type Result struct {
	Product *Product
	Err     error
}

// JSONLGzWriter streams products as gzip-compressed JSON Lines.
type JSONLGzWriter struct {
	w      *gzip.Writer
	enc    *json.Encoder
	err    error
	closed bool
}

// NewJSONLGzWriter creates a writer for gzip-compressed JSON Lines output.
func NewJSONLGzWriter(w io.Writer) *JSONLGzWriter {
	gz := gzip.NewWriter(w)
	enc := json.NewEncoder(gz)
	enc.SetEscapeHTML(false)
	return &JSONLGzWriter{
		w:   gz,
		enc: enc,
	}
}

// Write appends a single product record as one JSON Lines entry.
func (w *JSONLGzWriter) Write(p *Product) error {
	if w.closed {
		return errors.New("jsonl writer is closed")
	}
	if w.err != nil {
		return w.err
	}
	if p == nil {
		w.err = errors.New("product is nil")
		return w.err
	}
	if err := w.enc.Encode(p); err != nil {
		w.err = fmt.Errorf("encode jsonl record: %w", err)
		return w.err
	}
	return nil
}

// Close flushes and finalizes the gzip stream.
func (w *JSONLGzWriter) Close() error {
	if w.closed {
		return nil
	}
	w.closed = true

	closeErr := w.w.Close()
	if closeErr != nil {
		closeErr = fmt.Errorf("close gzip writer: %w", closeErr)
	}
	return errors.Join(w.err, closeErr)
}

// JSONLReadSeq reads products from a legacy JSON Lines feed.
func JSONLReadSeq(r io.Reader) iter.Seq[Result] {
	return func(yield func(Result) bool) {
		dec := json.NewDecoder(r)
		for {
			var product Product
			if err := dec.Decode(&product); err != nil {
				if err == io.EOF {
					return
				}
				_ = yield(Result{Err: fmt.Errorf("decode jsonl record: %w", err)})
				return
			}
			if !yield(Result{Product: &product}) {
				return
			}
		}
	}
}

// JSONLGzReadSeq reads products from a gzip-compressed legacy JSON Lines feed.
func JSONLGzReadSeq(r io.Reader) iter.Seq[Result] {
	return func(yield func(Result) bool) {
		gz, err := gzip.NewReader(r)
		if err != nil {
			_ = yield(Result{Err: fmt.Errorf("open gzip reader: %w", err)})
			return
		}

		interrupted := false
		for result := range JSONLReadSeq(gz) {
			if !yield(result) {
				interrupted = true
				break
			}
		}

		if err := gz.Close(); err != nil && !interrupted {
			_ = yield(Result{Err: fmt.Errorf("close gzip reader: %w", err)})
		}
	}
}
