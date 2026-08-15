package feed_test

import (
	"bytes"
	"fmt"

	"github.com/sumup/acp/feed"
)

func ExampleNewJSONLGzWriter() {
	var compressed bytes.Buffer
	writer := feed.NewJSONLGzWriter(&compressed)
	if err := writer.Write(&feed.Product{ID: "sku_123", Title: "Classic Tee"}); err != nil {
		panic(err)
	}
	if err := writer.Close(); err != nil {
		panic(err)
	}

	for result := range feed.JSONLGzReadSeq(&compressed) {
		if result.Err != nil {
			panic(result.Err)
		}
		fmt.Println(result.Product.ID, result.Product.Title)
	}
	// Output: sku_123 Classic Tee
}
