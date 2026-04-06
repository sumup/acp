package feed

import (
	"bytes"
	"os"
	"strings"
	"testing"
)

func TestNewCSVGzWriter(t *testing.T) {
	t.Parallel()

	t.Run("writes expected csv gzip output", func(t *testing.T) {
		t.Parallel()

		feed := sampleFeed()

		var buf bytes.Buffer
		writer := NewCSVGzWriter(&buf)
		for i := range feed {
			if err := writer.Write(feed[i]); err != nil {
				t.Fatalf("write csv row: %v", err)
			}
		}
		if err := writer.Close(); err != nil {
			t.Fatalf("close csv writer: %v", err)
		}

		got, err := gunzipBytes(buf.Bytes())
		if err != nil {
			t.Fatalf("read csv gzip: %v", err)
		}

		want, err := os.ReadFile("testdata/feed.csv")
		if err != nil {
			t.Fatalf("read expected csv: %v", err)
		}

		if !bytes.Equal(got, want) {
			t.Fatalf("csv output mismatch\nwant:\n%s\n\n got:\n%s", want, got)
		}
	})
}

func TestCSVGzWriter_Write(t *testing.T) {
	t.Parallel()

	t.Run("returns error for nil product", func(t *testing.T) {
		t.Parallel()

		var buf bytes.Buffer
		writer := NewCSVGzWriter(&buf)

		if err := writer.Write(nil); err == nil {
			t.Fatalf("expected write error for nil product")
		}
	})

	t.Run("returns error after close", func(t *testing.T) {
		t.Parallel()

		var buf bytes.Buffer
		writer := NewCSVGzWriter(&buf)
		if err := writer.Close(); err != nil {
			t.Fatalf("close csv writer: %v", err)
		}
		if err := writer.Write(&Product{ID: "sku"}); err == nil {
			t.Fatalf("expected write error after close")
		}
	})
}

func TestCSVGzWriter_Close(t *testing.T) {
	t.Parallel()

	t.Run("is idempotent", func(t *testing.T) {
		t.Parallel()

		var buf bytes.Buffer
		writer := NewCSVGzWriter(&buf)

		if err := writer.Close(); err != nil {
			t.Fatalf("first close failed: %v", err)
		}
		if err := writer.Close(); err != nil {
			t.Fatalf("second close failed: %v", err)
		}
	})
}

func TestCSVReadSeq(t *testing.T) {
	t.Parallel()

	t.Run("reads all products", func(t *testing.T) {
		t.Parallel()

		file, err := os.Open("testdata/feed.csv")
		if err != nil {
			t.Fatalf("open csv: %v", err)
		}
		defer func() {
			_ = file.Close()
		}()

		got := make(Feed, 0)
		for result := range CSVReadSeq(file) {
			if result.Err != nil {
				t.Fatalf("read csv: %v", result.Err)
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
		if got[0].Price != want[0].Price || got[1].Price != want[1].Price {
			t.Fatalf("unexpected prices: got [%s, %s], want [%s, %s]", got[0].Price, got[1].Price, want[0].Price, want[1].Price)
		}
	})

	t.Run("returns row parse error", func(t *testing.T) {
		t.Parallel()

		var gotErr error
		for result := range CSVReadSeq(strings.NewReader("id\n\"unterminated")) {
			if result.Err != nil {
				gotErr = result.Err
				break
			}
		}
		if gotErr == nil {
			t.Fatalf("expected csv read error")
		}
	})
}

func TestCSVGzReadSeq(t *testing.T) {
	t.Parallel()

	t.Run("reads all products from gzip input", func(t *testing.T) {
		t.Parallel()

		feed := sampleFeed()

		var buf bytes.Buffer
		writer := NewCSVGzWriter(&buf)
		for i := range feed {
			if err := writer.Write(feed[i]); err != nil {
				t.Fatalf("write csv row: %v", err)
			}
		}
		if err := writer.Close(); err != nil {
			t.Fatalf("close csv writer: %v", err)
		}

		got := make(Feed, 0)
		for result := range CSVGzReadSeq(bytes.NewReader(buf.Bytes())) {
			if result.Err != nil {
				t.Fatalf("read csv gzip: %v", result.Err)
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
		for result := range CSVGzReadSeq(strings.NewReader("not gzip")) {
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
