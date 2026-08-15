package feed

import (
	"compress/gzip"
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"iter"
	"strconv"
	"strings"
)

// csvHeader defines the ordered CSV column list for legacy feeds.
var csvHeader = []string{
	"enable_search",
	"enable_checkout",
	"id",
	"gtin",
	"mpn",
	"title",
	"description",
	"link",
	"condition",
	"product_category",
	"brand",
	"material",
	"dimensions",
	"length",
	"width",
	"height",
	"weight",
	"age_group",
	"image_link",
	"additional_image_link",
	"video_link",
	"model_3d_link",
	"price",
	"sale_price",
	"sale_price_effective_date",
	"unit_pricing_measure",
	"base_measure",
	"pricing_trend",
	"availability",
	"availability_date",
	"inventory_quantity",
	"expiration_date",
	"pickup_method",
	"pickup_sla",
	"item_group_id",
	"item_group_title",
	"color",
	"size",
	"size_system",
	"gender",
	"offer_id",
	"custom_variant1_category",
	"custom_variant1_option",
	"custom_variant2_category",
	"custom_variant2_option",
	"custom_variant3_category",
	"custom_variant3_option",
	"shipping",
	"delivery_estimate",
	"seller_name",
	"seller_url",
	"seller_privacy_policy",
	"seller_tos",
	"return_policy",
	"return_window",
	"popularity_score",
	"return_rate",
	"warning",
	"warning_url",
	"age_restriction",
	"product_review_count",
	"product_review_rating",
	"store_review_count",
	"store_review_rating",
	"q_and_a",
	"raw_review_data",
	"related_product_id",
	"relationship_type",
	"geo_price",
	"geo_availability",
}

// csvRecord formats a product into a CSV record matching csvHeader order.
func (p Product) csvRecord() []string {
	return []string{
		strconv.FormatBool(p.EnableSearch),
		strconv.FormatBool(p.EnableCheckout),
		p.ID,
		p.GTIN,
		p.MPN,
		p.Title,
		p.Description,
		p.Link,
		string(p.Condition),
		p.ProductCategory,
		p.Brand,
		p.Material,
		p.Dimensions,
		p.Length,
		p.Width,
		p.Height,
		p.Weight,
		string(p.AgeGroup),
		p.ImageLink,
		joinStrings(p.AdditionalImageLink),
		p.VideoLink,
		p.Model3DLink,
		p.Price,
		p.SalePrice,
		p.SalePriceEffectiveDate,
		p.UnitPricingMeasure,
		p.BaseMeasure,
		p.PricingTrend,
		p.Availability,
		p.AvailabilityDate,
		formatInt(p.InventoryQuantity),
		p.ExpirationDate,
		p.PickupMethod,
		p.PickupSLA,
		p.ItemGroupID,
		p.ItemGroupTitle,
		p.Color,
		p.Size,
		p.SizeSystem,
		p.Gender,
		p.OfferID,
		p.CustomVariant1Category,
		p.CustomVariant1Option,
		p.CustomVariant2Category,
		p.CustomVariant2Option,
		p.CustomVariant3Category,
		p.CustomVariant3Option,
		joinStrings(p.Shipping),
		p.DeliveryEstimate,
		p.SellerName,
		p.SellerURL,
		p.SellerPrivacyPolicy,
		p.SellerTOS,
		p.ReturnPolicy,
		formatInt(p.ReturnWindow),
		formatFloat(p.PopularityScore),
		p.ReturnRate,
		p.Warning,
		p.WarningURL,
		formatInt(p.AgeRestriction),
		formatInt(p.ProductReviewCount),
		formatFloat(p.ProductReviewRating),
		formatInt(p.StoreReviewCount),
		formatFloat(p.StoreReviewRating),
		p.QAndA,
		p.RawReviewData,
		joinStrings(p.RelatedProductID),
		p.RelationshipType,
		joinStrings(p.GeoPrice),
		joinStrings(p.GeoAvailability),
	}
}

// CSVGzWriter streams products as gzip-compressed CSV.
type CSVGzWriter struct {
	gz     *gzip.Writer
	cw     *csv.Writer
	err    error
	closed bool
}

// NewCSVGzWriter creates a writer for gzip-compressed CSV output.
func NewCSVGzWriter(w io.Writer) *CSVGzWriter {
	gz := gzip.NewWriter(w)
	cw := csv.NewWriter(gz)
	err := cw.Write(csvHeader)
	if err != nil {
		err = fmt.Errorf("write csv header: %w", err)
	}
	return &CSVGzWriter{
		gz:  gz,
		cw:  cw,
		err: err,
	}
}

// Write appends a single product record to the CSV feed.
func (w *CSVGzWriter) Write(p *Product) error {
	if w.closed {
		return errors.New("csv writer is closed")
	}
	if w.err != nil {
		return w.err
	}
	if p == nil {
		w.err = errors.New("product is nil")
		return w.err
	}
	if err := w.cw.Write(p.csvRecord()); err != nil {
		w.err = fmt.Errorf("write csv row: %w", err)
		return w.err
	}
	return nil
}

// Close flushes CSV data and finalizes the gzip stream.
func (w *CSVGzWriter) Close() error {
	if w.closed {
		return nil
	}
	w.closed = true

	if w.err == nil {
		w.cw.Flush()
		if err := w.cw.Error(); err != nil {
			w.err = fmt.Errorf("flush csv: %w", err)
		}
	}

	closeErr := w.gz.Close()
	if closeErr != nil {
		closeErr = fmt.Errorf("close gzip writer: %w", closeErr)
	}
	return errors.Join(w.err, closeErr)
}

// CSVReadSeq reads products from a legacy CSV feed.
func CSVReadSeq(r io.Reader) iter.Seq[Result] {
	return func(yield func(Result) bool) {
		reader := csv.NewReader(r)
		header, err := reader.Read()
		if err != nil {
			if err == io.EOF {
				return
			}
			_ = yield(Result{Err: fmt.Errorf("read csv header: %w", err)})
			return
		}

		headerIndex := make(map[string]int, len(header))
		for i := range header {
			headerIndex[header[i]] = i
		}
		if err := csvValidateHeader(headerIndex); err != nil {
			_ = yield(Result{Err: err})
			return
		}

		rowNumber := 2
		for {
			record, err := reader.Read()
			if err != nil {
				if err == io.EOF {
					return
				}
				_ = yield(Result{Err: fmt.Errorf("read csv row %d: %w", rowNumber, err)})
				return
			}
			product, err := csvRecordToProduct(record, headerIndex)
			if err != nil {
				_ = yield(Result{Err: fmt.Errorf("parse csv row %d: %w", rowNumber, err)})
				return
			}
			if !yield(Result{Product: product}) {
				return
			}
			rowNumber++
		}
	}
}

// CSVGzReadSeq reads products from a gzip-compressed legacy CSV feed.
func CSVGzReadSeq(r io.Reader) iter.Seq[Result] {
	return func(yield func(Result) bool) {
		gz, err := gzip.NewReader(r)
		if err != nil {
			_ = yield(Result{Err: fmt.Errorf("open gzip reader: %w", err)})
			return
		}

		interrupted := false
		for result := range CSVReadSeq(gz) {
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

// csvRecordToProduct maps one CSV row to a product.
func csvRecordToProduct(record []string, headerIndex map[string]int) (*Product, error) {
	enableSearch, err := csvParseBool(csvField(record, headerIndex, "enable_search"), "enable_search")
	if err != nil {
		return nil, err
	}
	enableCheckout, err := csvParseBool(csvField(record, headerIndex, "enable_checkout"), "enable_checkout")
	if err != nil {
		return nil, err
	}
	inventoryQuantity, err := csvParseIntPointer(csvField(record, headerIndex, "inventory_quantity"), "inventory_quantity")
	if err != nil {
		return nil, err
	}
	returnWindow, err := csvParseIntPointer(csvField(record, headerIndex, "return_window"), "return_window")
	if err != nil {
		return nil, err
	}
	popularityScore, err := csvParseFloatPointer(csvField(record, headerIndex, "popularity_score"), "popularity_score")
	if err != nil {
		return nil, err
	}
	ageRestriction, err := csvParseIntPointer(csvField(record, headerIndex, "age_restriction"), "age_restriction")
	if err != nil {
		return nil, err
	}
	productReviewCount, err := csvParseIntPointer(csvField(record, headerIndex, "product_review_count"), "product_review_count")
	if err != nil {
		return nil, err
	}
	productReviewRating, err := csvParseFloatPointer(csvField(record, headerIndex, "product_review_rating"), "product_review_rating")
	if err != nil {
		return nil, err
	}
	storeReviewCount, err := csvParseIntPointer(csvField(record, headerIndex, "store_review_count"), "store_review_count")
	if err != nil {
		return nil, err
	}
	storeReviewRating, err := csvParseFloatPointer(csvField(record, headerIndex, "store_review_rating"), "store_review_rating")
	if err != nil {
		return nil, err
	}

	product := &Product{
		EnableSearch:           enableSearch,
		EnableCheckout:         enableCheckout,
		ID:                     csvField(record, headerIndex, "id"),
		GTIN:                   csvField(record, headerIndex, "gtin"),
		MPN:                    csvField(record, headerIndex, "mpn"),
		Title:                  csvField(record, headerIndex, "title"),
		Description:            csvField(record, headerIndex, "description"),
		Link:                   csvField(record, headerIndex, "link"),
		Condition:              ProductCondition(csvField(record, headerIndex, "condition")),
		ProductCategory:        csvField(record, headerIndex, "product_category"),
		Brand:                  csvField(record, headerIndex, "brand"),
		Material:               csvField(record, headerIndex, "material"),
		Dimensions:             csvField(record, headerIndex, "dimensions"),
		Length:                 csvField(record, headerIndex, "length"),
		Width:                  csvField(record, headerIndex, "width"),
		Height:                 csvField(record, headerIndex, "height"),
		Weight:                 csvField(record, headerIndex, "weight"),
		AgeGroup:               ProductAgeGroup(csvField(record, headerIndex, "age_group")),
		ImageLink:              csvField(record, headerIndex, "image_link"),
		AdditionalImageLink:    csvSplit(csvField(record, headerIndex, "additional_image_link")),
		VideoLink:              csvField(record, headerIndex, "video_link"),
		Model3DLink:            csvField(record, headerIndex, "model_3d_link"),
		Price:                  csvField(record, headerIndex, "price"),
		SalePrice:              csvField(record, headerIndex, "sale_price"),
		SalePriceEffectiveDate: csvField(record, headerIndex, "sale_price_effective_date"),
		UnitPricingMeasure:     csvField(record, headerIndex, "unit_pricing_measure"),
		BaseMeasure:            csvField(record, headerIndex, "base_measure"),
		PricingTrend:           csvField(record, headerIndex, "pricing_trend"),
		Availability:           csvField(record, headerIndex, "availability"),
		AvailabilityDate:       csvField(record, headerIndex, "availability_date"),
		InventoryQuantity:      inventoryQuantity,
		ExpirationDate:         csvField(record, headerIndex, "expiration_date"),
		PickupMethod:           csvField(record, headerIndex, "pickup_method"),
		PickupSLA:              csvField(record, headerIndex, "pickup_sla"),
		ItemGroupID:            csvField(record, headerIndex, "item_group_id"),
		ItemGroupTitle:         csvField(record, headerIndex, "item_group_title"),
		Color:                  csvField(record, headerIndex, "color"),
		Size:                   csvField(record, headerIndex, "size"),
		SizeSystem:             csvField(record, headerIndex, "size_system"),
		Gender:                 csvField(record, headerIndex, "gender"),
		OfferID:                csvField(record, headerIndex, "offer_id"),
		CustomVariant1Category: csvField(record, headerIndex, "custom_variant1_category"),
		CustomVariant1Option:   csvField(record, headerIndex, "custom_variant1_option"),
		CustomVariant2Category: csvField(record, headerIndex, "custom_variant2_category"),
		CustomVariant2Option:   csvField(record, headerIndex, "custom_variant2_option"),
		CustomVariant3Category: csvField(record, headerIndex, "custom_variant3_category"),
		CustomVariant3Option:   csvField(record, headerIndex, "custom_variant3_option"),
		Shipping:               csvSplit(csvField(record, headerIndex, "shipping")),
		DeliveryEstimate:       csvField(record, headerIndex, "delivery_estimate"),
		SellerName:             csvField(record, headerIndex, "seller_name"),
		SellerURL:              csvField(record, headerIndex, "seller_url"),
		SellerPrivacyPolicy:    csvField(record, headerIndex, "seller_privacy_policy"),
		SellerTOS:              csvField(record, headerIndex, "seller_tos"),
		ReturnPolicy:           csvField(record, headerIndex, "return_policy"),
		ReturnWindow:           returnWindow,
		PopularityScore:        popularityScore,
		ReturnRate:             csvField(record, headerIndex, "return_rate"),
		Warning:                csvField(record, headerIndex, "warning"),
		WarningURL:             csvField(record, headerIndex, "warning_url"),
		AgeRestriction:         ageRestriction,
		ProductReviewCount:     productReviewCount,
		ProductReviewRating:    productReviewRating,
		StoreReviewCount:       storeReviewCount,
		StoreReviewRating:      storeReviewRating,
		QAndA:                  csvField(record, headerIndex, "q_and_a"),
		RawReviewData:          csvField(record, headerIndex, "raw_review_data"),
		RelatedProductID:       csvSplit(csvField(record, headerIndex, "related_product_id")),
		RelationshipType:       csvField(record, headerIndex, "relationship_type"),
		GeoPrice:               csvSplit(csvField(record, headerIndex, "geo_price")),
		GeoAvailability:        csvSplit(csvField(record, headerIndex, "geo_availability")),
	}
	return product, nil
}

// csvField returns a CSV cell value by column name.
func csvField(record []string, headerIndex map[string]int, column string) string {
	index, ok := headerIndex[column]
	if !ok || index >= len(record) {
		return ""
	}
	return record[index]
}

// csvValidateHeader ensures all expected columns are present.
func csvValidateHeader(headerIndex map[string]int) error {
	for i := range csvHeader {
		if _, ok := headerIndex[csvHeader[i]]; !ok {
			return fmt.Errorf("missing csv column %q", csvHeader[i])
		}
	}
	return nil
}

// csvSplit parses comma-separated values from a single CSV cell.
func csvSplit(value string) []string {
	if value == "" {
		return nil
	}
	return strings.Split(value, ",")
}

// csvParseBool parses a boolean value from a CSV cell.
func csvParseBool(value, column string) (bool, error) {
	if value == "" {
		return false, nil
	}
	parsed, err := strconv.ParseBool(value)
	if err != nil {
		return false, fmt.Errorf("invalid bool in %q: %q", column, value)
	}
	return parsed, nil
}

// csvParseIntPointer parses an optional int pointer from a CSV cell.
func csvParseIntPointer(value, column string) (*int, error) {
	if value == "" {
		return nil, nil
	}
	parsed, err := strconv.Atoi(value)
	if err != nil {
		return nil, fmt.Errorf("invalid int in %q: %q", column, value)
	}
	return &parsed, nil
}

// csvParseFloatPointer parses an optional float pointer from a CSV cell.
func csvParseFloatPointer(value, column string) (*float64, error) {
	if value == "" {
		return nil, nil
	}
	parsed, err := strconv.ParseFloat(value, 64)
	if err != nil {
		return nil, fmt.Errorf("invalid float in %q: %q", column, value)
	}
	return &parsed, nil
}
