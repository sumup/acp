package feed

import (
	"bytes"
	"compress/gzip"
	"io"
	"os"
	"testing"
)

func TestWriteJSONLGz(t *testing.T) {
	feed := sampleFeed()

	var buf bytes.Buffer
	if err := feed.WriteJSONLGz(&buf); err != nil {
		t.Fatalf("write jsonl gzip: %v", err)
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
}

func TestWriteCSVGz(t *testing.T) {
	feed := sampleFeed()

	var buf bytes.Buffer
	if err := feed.WriteCSVGz(&buf); err != nil {
		t.Fatalf("write csv gzip: %v", err)
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
}

func gunzipBytes(data []byte) ([]byte, error) {
	reader, err := gzip.NewReader(bytes.NewReader(data))
	if err != nil {
		return nil, err
	}
	defer func() { _ = reader.Close() }()

	return io.ReadAll(reader)
}

func sampleFeed() Feed {
	inventory := 42
	returnWindow := 30
	ageRestriction := 18
	productReviewCount := 124
	storeReviewCount := 2120
	popularityScore := 4.7
	productReviewRating := 4.8
	storeReviewRating := 4.6

	return Feed{
		{
			EnableSearch:    true,
			EnableCheckout:  true,
			ID:              "sku-1",
			GTIN:            "000123",
			MPN:             "MPN-1",
			Title:           "Horizon Mug",
			Description:     "Ceramic mug 12oz",
			Link:            "https://example.com/p/sku-1",
			Condition:       "new",
			ProductCategory: "Home & Garden > Kitchen & Dining",
			Brand:           "SumUp",
			Material:        "Ceramic",
			Dimensions:      "12x8x8 cm",
			Length:          "12 cm",
			Width:           "8 cm",
			Height:          "8 cm",
			Weight:          "0.4 kg",
			AgeGroup:        "adult",
			ImageLink:       "https://example.com/img/sku-1.jpg",
			AdditionalImageLink: []string{
				"https://example.com/img/sku-1-1.jpg",
				"https://example.com/img/sku-1-2.jpg",
			},
			VideoLink:              "https://example.com/video/sku-1.mp4",
			Model3DLink:            "https://example.com/model/sku-1.glb",
			Price:                  "12.99 USD",
			SalePrice:              "9.99 USD",
			SalePriceEffectiveDate: "2026-01-01T00:00:00Z/2026-01-31T23:59:59Z",
			UnitPricingMeasure:     "1 lb",
			BaseMeasure:            "1 lb",
			PricingTrend:           "Lowest price in 30 days",
			Availability:           "in_stock",
			InventoryQuantity:      &inventory,
			ExpirationDate:         "2026-12-31",
			PickupMethod:           "in_store",
			PickupSLA:              "1 day",
			ItemGroupID:            "mug-1",
			ItemGroupTitle:         "Horizon Mugs",
			Color:                  "blue",
			Size:                   "12oz",
			SizeSystem:             "US",
			Gender:                 "unisex",
			OfferID:                "offer-1",
			CustomVariant1Category: "pattern",
			CustomVariant1Option:   "speckled",
			CustomVariant2Category: "finish",
			CustomVariant2Option:   "matte",
			CustomVariant3Category: "limited",
			CustomVariant3Option:   "winter",
			Shipping: []string{
				"US:::5.00 USD",
				"CA:::7.00 CAD",
			},
			DeliveryEstimate:    "2026-01-20",
			SellerName:          "SumUp Shop",
			SellerURL:           "https://example.com",
			SellerPrivacyPolicy: "https://example.com/privacy",
			SellerTOS:           "https://example.com/terms",
			ReturnPolicy:        "https://example.com/returns",
			ReturnWindow:        &returnWindow,
			PopularityScore:     &popularityScore,
			ReturnRate:          "2%",
			Warning:             "Dishwasher safe",
			WarningURL:          "https://example.com/warnings",
			AgeRestriction:      &ageRestriction,
			ProductReviewCount:  &productReviewCount,
			ProductReviewRating: &productReviewRating,
			StoreReviewCount:    &storeReviewCount,
			StoreReviewRating:   &storeReviewRating,
			QAndA:               "Is it microwave safe? Yes.",
			RawReviewData:       `{"source":"internal"}`,
			RelatedProductID: []string{
				"sku-2",
				"sku-3",
			},
			RelationshipType: "accessory",
			GeoPrice: []string{
				"US:12.99 USD",
				"CA:16.99 CAD",
			},
			GeoAvailability: []string{
				"US:in_stock",
				"CA:out_of_stock",
			},
		},
		{
			EnableSearch:     false,
			EnableCheckout:   false,
			ID:               "sku-2",
			Title:            "Linen Apron",
			Description:      "Lightweight apron",
			Link:             "https://example.com/p/sku-2",
			ProductCategory:  "Home & Garden > Kitchen & Dining",
			ImageLink:        "https://example.com/img/sku-2.jpg",
			Price:            "24.00 USD",
			Availability:     "preorder",
			AvailabilityDate: "2026-02-01",
			SellerName:       "SumUp Shop",
			SellerURL:        "https://example.com",
		},
	}
}
