package feed

import (
	"compress/gzip"
	"encoding/csv"
	"fmt"
	"io"
	"strconv"
)

// csvHeader defines the ordered CSV column list for ACP feeds.
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
		p.AgeGroup,
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

// WriteCSVGz writes the feed as CSV, compressed with gzip.
func (f Feed) WriteCSVGz(w io.Writer) error {
	gz := gzip.NewWriter(w)
	writer := csv.NewWriter(gz)
	if err := writer.Write(csvHeader); err != nil {
		_ = gz.Close()
		return fmt.Errorf("write csv header: %w", err)
	}
	for i := range f {
		if err := writer.Write(f[i].csvRecord()); err != nil {
			_ = gz.Close()
			return fmt.Errorf("write csv row: %w", err)
		}
	}
	writer.Flush()
	if err := writer.Error(); err != nil {
		_ = gz.Close()
		return fmt.Errorf("flush csv: %w", err)
	}
	if err := gz.Close(); err != nil {
		return fmt.Errorf("close gzip: %w", err)
	}
	return nil
}
