package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/sumup/acp/feed"
)

func main() {
	products := sampleProducts()
	productFeed := feed.New(products)

	outDir := filepath.Join("examples", "feed", "output")
	if err := os.MkdirAll(outDir, 0o755); err != nil {
		log.Fatalf("create output directory: %v", err)
	}

	jsonlPath := filepath.Join(outDir, "product_feed.jsonl.gz")
	if err := writeJSONL(productFeed, jsonlPath); err != nil {
		log.Fatalf("write jsonl feed: %v", err)
	}

	csvPath := filepath.Join(outDir, "product_feed.csv.gz")
	if err := writeCSV(productFeed, csvPath); err != nil {
		log.Fatalf("write csv feed: %v", err)
	}

	log.Printf("wrote %s", jsonlPath)
	log.Printf("wrote %s", csvPath)
}

func writeJSONL(productFeed feed.Feed, path string) error {
	file, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("create file: %w", err)
	}
	defer func() {
		_ = file.Close()
	}()

	return productFeed.WriteJSONLGz(file)
}

func writeCSV(productFeed feed.Feed, path string) error {
	file, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("create file: %w", err)
	}
	defer func() {
		_ = file.Close()
	}()

	return productFeed.WriteCSVGz(file)
}

func sampleProducts() []feed.Product {
	return []feed.Product{
		{
			EnableSearch:      true,
			EnableCheckout:    true,
			ID:                "sku_latte_12oz",
			Title:             "Oat Milk Latte (12oz)",
			Description:       "Smooth espresso with steamed oat milk.",
			Link:              "https://store.example.com/products/sku_latte_12oz",
			ProductCategory:   "Food, Beverages & Tobacco > Beverages > Coffee",
			ImageLink:         "https://store.example.com/images/sku_latte_12oz.jpg",
			Price:             "6.50 USD",
			Availability:      "in_stock",
			SellerName:        "Example Roasters",
			SellerURL:         "https://store.example.com",
			InventoryQuantity: ptr(120),
			Shipping: []string{
				"US:CA:Ground:5.00 USD",
				"US:NY:Ground:6.00 USD",
			},
			PickupMethod: "in_store",
			PickupSLA:    "1 day",
		},
		{
			EnableSearch:    true,
			EnableCheckout:  true,
			ID:              "sku_mug_slate",
			Title:           "Stoneware Mug (Slate)",
			Description:     "12oz stoneware mug with matte finish.",
			Link:            "https://store.example.com/products/sku_mug_slate",
			ProductCategory: "Home & Garden > Kitchen & Dining > Drinkware",
			ImageLink:       "https://store.example.com/images/sku_mug_slate.jpg",
			AdditionalImageLink: []string{
				"https://store.example.com/images/sku_mug_slate_alt1.jpg",
				"https://store.example.com/images/sku_mug_slate_alt2.jpg",
			},
			Price:        "15.00 USD",
			SalePrice:    "12.00 USD",
			Availability: "in_stock",
			Condition:    feed.ProductConditionNew,
			Brand:        "Example Roasters",
			Material:     "Stoneware",
			Dimensions:   "4x4x4 in",
			Shipping: []string{
				"US:::7.00 USD",
			},
			SellerName: "Example Roasters",
			SellerURL:  "https://store.example.com",
		},
		{
			EnableSearch:    true,
			EnableCheckout:  true,
			ID:              "sku_shirt_black_s",
			Title:           "Logo T-Shirt (Black, S)",
			Description:     "Soft cotton tee with embroidered logo.",
			Link:            "https://store.example.com/products/sku_shirt_black_s",
			ProductCategory: "Apparel & Accessories > Clothing",
			ImageLink:       "https://store.example.com/images/sku_shirt_black.jpg",
			Price:           "24.00 USD",
			Availability:    "in_stock",
			ItemGroupID:     "sku_shirt_black",
			ItemGroupTitle:  "Logo T-Shirt (Black)",
			Color:           "Black",
			Size:            "S",
			SizeSystem:      "US",
			Gender:          "unisex",
			SellerName:      "Example Roasters",
			SellerURL:       "https://store.example.com",
		},
	}
}

func ptr[V any](v V) *V {
	return &v
}
