// Package feed provides types and helpers for exporting [ACP product feeds].
// Product feed is used to index and display your products with up-to-date price
// and availability in AI platforms such as OpenAI.
//
// Before you begin, you will have to sign up at [chatgpt.com/merchants](https://chatgpt.com/merchants).
//
// # How it works
//
//   - Prepare and export your feed.
//   - Share the feed using the preferred delivery method and file format described in the [integration section](https://developers.openai.com/commerce/specs/feed#integration-overview).
//   - OpenAI ingests the feed, validates records, and indexes product metadata for retrieval and ranking in ChatGPT.
//   - Update the feed whenever products, pricing, or availability change to ensure users see accurate information.
//
// # Example
//
//	  feed := New([]Product{
//	  	{
//	  		EnableSearch:    true,
//	  		EnableCheckout:  true,
//	  		ID:              "sku_123",
//	  		Title:           "Everyday T-Shirt",
//	  		Description:     "Soft cotton tee in classic fit.",
//	  		Link:            "https://store.example.com/products/sku_123",
//	  		ProductCategory: "Apparel & Accessories > Clothing",
//	  		ImageLink:       "https://store.example.com/images/sku_123.jpg",
//	  		Price:           "19.00 EUR",
//	  		Availability:    "in_stock",
//	  		SellerName:      "Example Store",
//	  		SellerURL:       "https://store.example.com",
//	  	},
//	  	{
//	  		EnableSearch:    true,
//	  		EnableCheckout:  true,
//	  		ID:              "sku_456",
//	  		Title:           "Canvas Tote",
//	  		Description:     "Reusable tote bag for daily errands.",
//	  		Link:            "https://store.example.com/products/sku_456",
//	  		ProductCategory: "Apparel & Accessories > Bags",
//	  		ImageLink:       "https://store.example.com/images/sku_456.jpg",
//	  		Price:           "12.00 EUR",
//	  		Availability:    "in_stock",
//	  		SellerName:      "Example Store",
//	  		SellerURL:       "https://store.example.com",
//	  	},
//	  })
//
//		var buf bytes.Buffer
//		writer := NewJSONLGzWriter(&buf)
//		for i := range feed {
//			_ = writer.Write(feed[i])
//		}
//		_ = writer.Close()
//
// [ACP product feeds]: https://developers.openai.com/commerce/specs/feed
package feed
