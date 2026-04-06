package feed

import (
	"bytes"
	"os"
	"strings"
	"testing"
)

func TestNewJSONLGzWriter(t *testing.T) {
	t.Parallel()

	t.Run("writes expected jsonl gzip output", func(t *testing.T) {
		t.Parallel()

		feed := sampleFeed()

		var buf bytes.Buffer
		writer := NewJSONLGzWriter(&buf)
		for i := range feed {
			if err := writer.Write(feed[i]); err != nil {
				t.Fatalf("write jsonl record: %v", err)
			}
		}
		if err := writer.Close(); err != nil {
			t.Fatalf("close jsonl writer: %v", err)
		}

		got, err := gunzipBytes(buf.Bytes())
		if err != nil {
			t.Fatalf("read jsonl gzip: %v", err)
		}

		want, err := os.ReadFile("testdata/feed.jsonl")
		if err != nil {
			t.Fatalf("read expected jsonl: %v", err)
		}

		if !bytes.Equal(got, want) {
			t.Fatalf("jsonl output mismatch\nwant:\n%s\n\n got:\n%s", want, got)
		}
	})
}

func TestJSONLGzWriter_Write(t *testing.T) {
	t.Parallel()

	t.Run("returns error for nil product", func(t *testing.T) {
		t.Parallel()

		var buf bytes.Buffer
		writer := NewJSONLGzWriter(&buf)

		if err := writer.Write(nil); err == nil {
			t.Fatalf("expected write error for nil product")
		}
	})

	t.Run("returns error after close", func(t *testing.T) {
		t.Parallel()

		var buf bytes.Buffer
		writer := NewJSONLGzWriter(&buf)
		if err := writer.Close(); err != nil {
			t.Fatalf("close jsonl writer: %v", err)
		}
		if err := writer.Write(&Product{ID: "sku"}); err == nil {
			t.Fatalf("expected write error after close")
		}
	})
}

func TestJSONLGzWriter_Close(t *testing.T) {
	t.Parallel()

	t.Run("is idempotent", func(t *testing.T) {
		t.Parallel()

		var buf bytes.Buffer
		writer := NewJSONLGzWriter(&buf)

		if err := writer.Close(); err != nil {
			t.Fatalf("first close failed: %v", err)
		}
		if err := writer.Close(); err != nil {
			t.Fatalf("second close failed: %v", err)
		}
	})
}

func TestJSONLReadSeq(t *testing.T) {
	t.Parallel()

	t.Run("reads all products", func(t *testing.T) {
		t.Parallel()

		file, err := os.Open("testdata/feed.jsonl")
		if err != nil {
			t.Fatalf("open jsonl: %v", err)
		}
		defer func() {
			_ = file.Close()
		}()

		got := make(Feed, 0)
		for result := range JSONLReadSeq(file) {
			if result.Err != nil {
				t.Fatalf("read jsonl: %v", result.Err)
			}
			got = append(got, result.Product)
		}

		want := sampleFeed()
		if len(got) != len(want) {
			t.Fatalf("product count mismatch: got %d want %d", len(got), len(want))
		}
		if got[0].ID != want[0].ID || got[1].ID != want[1].ID {
			t.Fatalf("unexpected ids: got [%s, %s], want [%s, %s]", got[0].ID, got[1].ID, want[0].ID, want[1].ID)
		}
	})

	t.Run("returns decode error", func(t *testing.T) {
		t.Parallel()

		var gotErr error
		for result := range JSONLReadSeq(strings.NewReader("{")) {
			if result.Err != nil {
				gotErr = result.Err
				break
			}
		}
		if gotErr == nil {
			t.Fatalf("expected jsonl decode error")
		}
	})
}

func TestJSONLGzReadSeq(t *testing.T) {
	t.Parallel()

	t.Run("reads all products from gzip input", func(t *testing.T) {
		t.Parallel()

		feed := sampleFeed()

		var buf bytes.Buffer
		writer := NewJSONLGzWriter(&buf)
		for i := range feed {
			if err := writer.Write(feed[i]); err != nil {
				t.Fatalf("write jsonl record: %v", err)
			}
		}
		if err := writer.Close(); err != nil {
			t.Fatalf("close jsonl writer: %v", err)
		}

		got := make(Feed, 0)
		for result := range JSONLGzReadSeq(bytes.NewReader(buf.Bytes())) {
			if result.Err != nil {
				t.Fatalf("read jsonl gzip: %v", result.Err)
			}
			got = append(got, result.Product)
		}

		want := sampleFeed()
		if len(got) != len(want) {
			t.Fatalf("product count mismatch: got %d want %d", len(got), len(want))
		}
		if got[0].ID != want[0].ID || got[1].ID != want[1].ID {
			t.Fatalf("unexpected ids: got [%s, %s], want [%s, %s]", got[0].ID, got[1].ID, want[0].ID, want[1].ID)
		}
	})

	t.Run("returns gzip open error", func(t *testing.T) {
		t.Parallel()

		var gotErr error
		for result := range JSONLGzReadSeq(strings.NewReader("not gzip")) {
			if result.Err != nil {
				gotErr = result.Err
				break
			}
		}
		if gotErr == nil {
			t.Fatalf("expected gzip error")
		}
	})
}
