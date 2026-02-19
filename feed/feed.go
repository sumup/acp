package feed

import (
	"strconv"
	"strings"
)

// Feed is a collection of products n the ACP product feed.
type Feed []*Product

// New creates new feed from a list of products.
func New(products []*Product) Feed {
	return Feed(products)
}

// ProductCondition is the condition of the product.
type ProductCondition string

const (
	ProductConditionNew         ProductCondition = "new"
	ProductConditionRefurbished ProductCondition = "refurbished"
	ProductConditionUsed        ProductCondition = "used"
)

// ProductAgeGroup is the target demographic of the product.
type ProductAgeGroup string

const (
	ProductAgeGroupNewborn ProductAgeGroup = "newborn"
	ProductAgeGroupInfant  ProductAgeGroup = "infant"
	ProductAgeGroupToddler ProductAgeGroup = "toddler"
	ProductAgeGroupKids    ProductAgeGroup = "kids"
	ProductAgeGroupAdult   ProductAgeGroup = "adult"
)

// Product describes a single product entry in an ACP product feed.
// Field names follow the feed specification for JSONL and CSV export.
type Product struct {
	// EnableSearch controls whether the product can be surfaced in ChatGPT search results.
	EnableSearch bool `json:"enable_search" csv:"enable_search"`
	// EnableCheckout allows direct purchase inside ChatGPT when EnableSearch is true.
	EnableCheckout bool `json:"enable_checkout" csv:"enable_checkout"`

	// Basic Product Data
	//
	// Provide the core identifiers and descriptive text needed to uniquely reference each product.
	// These fields establish the canonical record that ChatGPT Search uses to display and link to your product.

	// ID is the merchant product identifier and must remain stable over time.
	//
	// Example: SKU12345
	ID string `json:"id" csv:"id"`
	// GTIN is a universal product identifier such as GTIN, UPC, or ISBN.
	//
	// Example: 123456789543
	GTIN string `json:"gtin,omitempty" csv:"gtin"`
	// MPN is the manufacturer part number, required if GTIN is absent.
	//
	// Example: GPT5
	MPN string `json:"mpn,omitempty" csv:"mpn"`
	// Title is the product title shown to shoppers.
	//
	// Example: Men's Trail Running Shoes Black
	Title string `json:"title" csv:"title"`
	// Description is the full product description, plain text.
	//
	// Example: Waterproof trail shoe with cushioned sole…
	Description string `json:"description" csv:"description"`
	// Link is the product detail page URL.
	//
	// Example: https://example.com/product/SKU12345
	Link string `json:"link" csv:"link"`

	// Item Information
	//
	// Capture the physical characteristics and classification details of the product.
	// This data helps ensure accurate categorization, filtering, and search relevance.

	// Condition is the product condition such as new, refurbished, or used.
	//
	// Example: new
	Condition ProductCondition `json:"condition,omitempty" csv:"condition"`
	// ProductCategory is the taxonomy path using a ">" separator.
	ProductCategory string `json:"product_category" csv:"product_category"`
	// Brand is the product brand name.
	Brand string `json:"brand,omitempty" csv:"brand"`
	// Material describes the primary material.
	Material string `json:"material,omitempty" csv:"material"`
	// Dimensions is the overall size formatted as LxWxH with a unit.
	Dimensions string `json:"dimensions,omitempty" csv:"dimensions"`
	// Length is the individual length dimension with a unit.
	Length string `json:"length,omitempty" csv:"length"`
	// Width is the individual width dimension with a unit.
	Width string `json:"width,omitempty" csv:"width"`
	// Height is the individual height dimension with a unit.
	Height string `json:"height,omitempty" csv:"height"`
	// Weight is the product weight with a unit.
	Weight string `json:"weight,omitempty" csv:"weight"`
	// AgeGroup is the target demographic such as newborn, infant, toddler, kids, or adult.
	AgeGroup ProductAgeGroup `json:"age_group,omitempty" csv:"age_group"`

	// Media
	//
	// Supply visual and rich media assets that represent the product.
	// High-quality images and optional videos or 3D models improve user trust and engagement.

	// ImageLink is the main product image URL.
	ImageLink string `json:"image_link" csv:"image_link"`
	// AdditionalImageLink lists extra image URLs (comma-separated in CSV).
	AdditionalImageLink []string `json:"additional_image_link,omitempty" csv:"additional_image_link"`
	// VideoLink is a publicly accessible product video URL.
	VideoLink string `json:"video_link,omitempty" csv:"video_link"`
	// Model3DLink is a 3D model URL (GLB/GLTF preferred).
	Model3DLink string `json:"model_3d_link,omitempty" csv:"model_3d_link"`

	// Price & Promotions
	//
	// Define standard and promotional pricing information.
	// These attributes power price display, discount messaging, and offer comparisons.

	// Price is the regular price with ISO 4217 currency code.
	Price string `json:"price" csv:"price"`
	// SalePrice is the discounted price with currency code.
	SalePrice string `json:"sale_price,omitempty" csv:"sale_price"`
	// SalePriceEffectiveDate is the ISO 8601 sale window date range.
	SalePriceEffectiveDate string `json:"sale_price_effective_date,omitempty" csv:"sale_price_effective_date"`
	// UnitPricingMeasure and BaseMeasure describe unit pricing (both required together).
	UnitPricingMeasure string `json:"unit_pricing_measure,omitempty" csv:"unit_pricing_measure"`
	// BaseMeasure is the base unit used for unit pricing.
	BaseMeasure string `json:"base_measure,omitempty" csv:"base_measure"`
	// PricingTrend is a short string like "Lowest price in N months".
	PricingTrend string `json:"pricing_trend,omitempty" csv:"pricing_trend"`

	// Availability & Inventory
	//
	// Describe current stock levels and key timing signals for product availability.
	// Accurate inventory data ensures users only see items they can actually purchase.

	// Availability is the stock status: in_stock, out_of_stock, or preorder.
	Availability string `json:"availability" csv:"availability"`
	// AvailabilityDate is the ISO 8601 availability date for preorder items.
	AvailabilityDate string `json:"availability_date,omitempty" csv:"availability_date"`
	// InventoryQuantity is the non-negative stock count.
	InventoryQuantity *int `json:"inventory_quantity,omitempty" csv:"inventory_quantity"`
	// ExpirationDate is the ISO 8601 date to remove the product after.
	ExpirationDate string `json:"expiration_date,omitempty" csv:"expiration_date"`
	// PickupMethod specifies pickup options: in_store, reserve, or not_supported.
	PickupMethod string `json:"pickup_method,omitempty" csv:"pickup_method"`
	// PickupSLA is the pickup service-level agreement, like "1 day".
	PickupSLA string `json:"pickup_sla,omitempty" csv:"pickup_sla"`

	// Variants
	//
	// Specify variant relationships and distinguishing attributes such as color or size.
	// These fields allow ChatGPT to group related SKUs and surface variant-specific details.
	//
	// The item_group_id value should represent how the product is presented on the merchant’s website
	// (the canonical product page or parent listing shown to customers).
	// If you are submitting variant rows (e.g., by color or size),
	// you must include the same item_group_id for every variant.
	// Do not submit individual variant SKUs without a group id.

	// ItemGroupID groups variants under a canonical product listing.
	ItemGroupID string `json:"item_group_id,omitempty" csv:"item_group_id"`
	// ItemGroupTitle is the title for the variant group.
	ItemGroupTitle string `json:"item_group_title,omitempty" csv:"item_group_title"`
	// Color is the variant color.
	Color string `json:"color,omitempty" csv:"color"`
	// Size is the variant size.
	Size string `json:"size,omitempty" csv:"size"`
	// SizeSystem is the ISO 3166 size system code, such as US.
	SizeSystem string `json:"size_system,omitempty" csv:"size_system"`
	// Gender is the target gender: male, female, or unisex.
	Gender string `json:"gender,omitempty" csv:"gender"`
	// OfferID identifies a specific offer (SKU + seller + price), unique within the feed.
	OfferID string `json:"offer_id,omitempty" csv:"offer_id"`
	// CustomVariant1Category names the first custom variant dimension.
	CustomVariant1Category string `json:"custom_variant1_category,omitempty" csv:"custom_variant1_category"`
	// CustomVariant1Option provides the option value for custom variant 1.
	CustomVariant1Option string `json:"custom_variant1_option,omitempty" csv:"custom_variant1_option"`
	// CustomVariant2Category names the second custom variant dimension.
	CustomVariant2Category string `json:"custom_variant2_category,omitempty" csv:"custom_variant2_category"`
	// CustomVariant2Option provides the option value for custom variant 2.
	CustomVariant2Option string `json:"custom_variant2_option,omitempty" csv:"custom_variant2_option"`
	// CustomVariant3Category names the third custom variant dimension.
	CustomVariant3Category string `json:"custom_variant3_category,omitempty" csv:"custom_variant3_category"`
	// CustomVariant3Option provides the option value for custom variant 3.
	CustomVariant3Option string `json:"custom_variant3_option,omitempty" csv:"custom_variant3_option"`

	// Fulfillment
	//
	// Outline shipping methods, costs, and estimated delivery times.
	// Providing detailed shipping information helps users understand fulfillment options upfront.

	// Shipping lists shipping entries in country:region:service_class:price format.
	Shipping []string `json:"shipping,omitempty" csv:"shipping"`
	// DeliveryEstimate is the ISO 8601 estimated arrival date.
	DeliveryEstimate string `json:"delivery_estimate,omitempty" csv:"delivery_estimate"`

	// Merchant Info
	//
	// Identify the seller and link to any relevant merchant policies or storefront pages.
	// This ensures proper attribution and enables users to review seller credentials.

	// SellerName is the merchant display name.
	SellerName string `json:"seller_name" csv:"seller_name"`
	// SellerURL is the merchant storefront URL.
	SellerURL string `json:"seller_url" csv:"seller_url"`
	// SellerPrivacyPolicy is the seller-specific privacy policy URL.
	SellerPrivacyPolicy string `json:"seller_privacy_policy,omitempty" csv:"seller_privacy_policy"`
	// SellerTOS is the seller-specific terms of service URL.
	SellerTOS string `json:"seller_tos,omitempty" csv:"seller_tos"`

	// Returns
	//
	// Provide return policies and time windows to set clear expectations for buyers.
	// Transparent return data builds trust and reduces post-purchase confusion.

	// ReturnPolicy is the return policy URL.
	ReturnPolicy string `json:"return_policy,omitempty" csv:"return_policy"`
	// ReturnWindow is the number of days allowed for returns.
	ReturnWindow *int `json:"return_window,omitempty" csv:"return_window"`

	// Performance Signals
	//
	// Share popularity and return-rate metrics where available.
	// These signals can be used to enhance ranking and highlight high-performing products.

	// PopularityScore is a popularity indicator (for example, 0-5 scale).
	PopularityScore *float64 `json:"popularity_score,omitempty" csv:"popularity_score"`
	// ReturnRate is the percentage of returns, 0-100%.
	ReturnRate string `json:"return_rate,omitempty" csv:"return_rate"`

	// Compliance
	//
	// Include regulatory warnings, disclaimers, or age restrictions.
	// Compliance fields help meet legal obligations and protect consumers.

	// Warning is a product disclaimer or regulatory warning.
	Warning string `json:"warning,omitempty" csv:"warning"`
	// WarningURL links to warning details and must resolve.
	WarningURL string `json:"warning_url,omitempty" csv:"warning_url"`
	// AgeRestriction is the minimum purchase age.
	AgeRestriction *int `json:"age_restriction,omitempty" csv:"age_restriction"`

	// Reviews and Q&A
	//
	// Supply aggregated review statistics and frequently asked questions.
	// User-generated insights strengthen credibility and help shoppers make informed decisions.

	// ProductReviewCount is the number of product reviews.
	ProductReviewCount *int `json:"product_review_count,omitempty" csv:"product_review_count"`
	// ProductReviewRating is the average product review score.
	ProductReviewRating *float64 `json:"product_review_rating,omitempty" csv:"product_review_rating"`
	// StoreReviewCount is the number of brand or store reviews.
	StoreReviewCount *int `json:"store_review_count,omitempty" csv:"store_review_count"`
	// StoreReviewRating is the average brand or store rating.
	StoreReviewRating *float64 `json:"store_review_rating,omitempty" csv:"store_review_rating"`
	// QAndA is FAQ content in plain text.
	QAndA string `json:"q_and_a,omitempty" csv:"q_and_a"`
	// RawReviewData contains raw review payloads and may include JSON blobs.
	RawReviewData string `json:"raw_review_data,omitempty" csv:"raw_review_data"`

	// Related Products
	//
	// List products that are commonly bought together or act as substitutes.
	// This enables basket-building recommendations and cross-sell opportunities.

	// RelatedProductID lists associated product IDs (comma-separated in CSV).
	RelatedProductID []string `json:"related_product_id,omitempty" csv:"related_product_id"`
	// RelationshipType describes how related products connect (for example, part_of_set).
	RelationshipType string `json:"relationship_type,omitempty" csv:"relationship_type"`

	// Geo Tagging
	//
	// Indicate any region-specific pricing or availability overrides
	//  Geo data allows ChatGPT to present accurate offers and stock status by location.

	// GeoPrice lists country-specific prices using ISO 3166-1 country codes.
	//
	// Example: 79.99 EUR (California)
	GeoPrice []string `json:"geo_price,omitempty" csv:"geo_price"`
	// GeoAvailability lists country-specific availability using ISO 3166-1 country codes.
	//
	// Example: in_stock (Texas), out_of_stock (New York)
	GeoAvailability []string `json:"geo_availability,omitempty" csv:"geo_availability"`
}

// joinStrings combines list values using a comma, returning empty for nil slices.
func joinStrings(values []string) string {
	if len(values) == 0 {
		return ""
	}
	return strings.Join(values, ",")
}

// formatInt converts an optional int to a string, returning empty when nil.
func formatInt(value *int) string {
	if value == nil {
		return ""
	}
	return strconv.Itoa(*value)
}

// formatFloat converts an optional float to a string, returning empty when nil.
func formatFloat(value *float64) string {
	if value == nil {
		return ""
	}
	return strconv.FormatFloat(*value, 'f', -1, 64)
}
