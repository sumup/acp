package acp

import (
	"encoding/json"
	"testing"
)

func TestAuthenticationMetadataUnmarshalWithoutChannel(t *testing.T) {
	t.Parallel()

	payload := []byte(`{
		"acquirer_details": {
			"acquirer_bin": "123456",
			"acquirer_country": "US",
			"acquirer_merchant_id": "merchant_123",
			"merchant_name": "ACME Store"
		},
		"directory_server": "visa"
	}`)

	var got AuthenticationMetadata
	if err := json.Unmarshal(payload, &got); err != nil {
		t.Fatalf("unmarshal authentication metadata: %v", err)
	}
	if got.AcquirerDetails.AcquirerBin != "123456" {
		t.Fatalf("unexpected acquirer_bin %q", got.AcquirerDetails.AcquirerBin)
	}
	if got.DirectoryServer != AuthenticationDirectoryServerVisa {
		t.Fatalf("unexpected directory_server %q", got.DirectoryServer)
	}
}

func TestMessageErrorResolutionAndInterventionRequired(t *testing.T) {
	t.Parallel()

	resolution := MessageResolutionRequiresBuyerReview
	msg := MessageError{
		Type:        "error",
		Code:        InterventionRequired,
		ContentType: MessageErrorContentTypePlain,
		Content:     "manual review required",
		Resolution:  &resolution,
	}

	data, err := json.Marshal(msg)
	if err != nil {
		t.Fatalf("marshal message error: %v", err)
	}

	var got MessageError
	if err := json.Unmarshal(data, &got); err != nil {
		t.Fatalf("unmarshal message error: %v", err)
	}
	if got.Code != InterventionRequired {
		t.Fatalf("unexpected code %q", got.Code)
	}
	if got.Resolution == nil || *got.Resolution != MessageResolutionRequiresBuyerReview {
		t.Fatalf("unexpected resolution %v", got.Resolution)
	}
}

func TestErrorSupportedVersions(t *testing.T) {
	t.Parallel()

	versions := []string{"2026-01-30", "2026-01-16"}
	errPayload := NewHTTPError(
		400,
		RequestNotIdempotentType,
		RequestNotIdempotent,
		"idempotency violation",
		WithSupportedVersions(versions),
	)

	versions[0] = "mutated"
	if errPayload.SupportedVersions[0] != "2026-01-30" {
		t.Fatalf("supported_versions should be copied, got %q", errPayload.SupportedVersions[0])
	}

	data, err := json.Marshal(errPayload)
	if err != nil {
		t.Fatalf("marshal error payload: %v", err)
	}
	var body map[string]any
	if err := json.Unmarshal(data, &body); err != nil {
		t.Fatalf("unmarshal error payload: %v", err)
	}
	gotType, _ := body["type"].(string)
	if gotType != string(RequestNotIdempotentType) {
		t.Fatalf("unexpected type %q", gotType)
	}
	gotVersions, ok := body["supported_versions"].([]any)
	if !ok || len(gotVersions) != 2 {
		t.Fatalf("unexpected supported_versions %v", body["supported_versions"])
	}
}

func TestWebhookEventTypeMatchesSpec(t *testing.T) {
	t.Parallel()

	if WebhookEventTypeOrderCreated != "order_create" {
		t.Fatalf("unexpected order_created event value: %q", WebhookEventTypeOrderCreated)
	}
	if WebhookEventTypeOrderUpdated != "order_update" {
		t.Fatalf("unexpected order_updated event value: %q", WebhookEventTypeOrderUpdated)
	}
}

func TestTotalAmountRefundedMarshals(t *testing.T) {
	t.Parallel()

	total := Total{
		Type:        TotalTypeAmountRefunded,
		DisplayText: "Refunded",
		Amount:      500,
	}
	data, err := json.Marshal(total)
	if err != nil {
		t.Fatalf("marshal total: %v", err)
	}
	var body map[string]any
	if err := json.Unmarshal(data, &body); err != nil {
		t.Fatalf("unmarshal total: %v", err)
	}
	if got, _ := body["type"].(string); got != "amount_refunded" {
		t.Fatalf("unexpected total type %q", got)
	}
}

func TestPaymentDataHandlerIDAndInstrumentMarshals(t *testing.T) {
	t.Parallel()

	p := PaymentData{
		HandlerID: ptr("handler_sumup_card"),
		Instrument: &PaymentInstrument{
			Type: "card",
			Credential: PaymentCredential{
				Type:  "spt",
				Token: "spt_123",
			},
		},
	}
	data, err := json.Marshal(p)
	if err != nil {
		t.Fatalf("marshal payment_data: %v", err)
	}
	var body map[string]any
	if err := json.Unmarshal(data, &body); err != nil {
		t.Fatalf("unmarshal payment_data: %v", err)
	}
	if got, _ := body["handler_id"].(string); got != "handler_sumup_card" {
		t.Fatalf("unexpected handler_id %q", got)
	}
	instrument, ok := body["instrument"].(map[string]any)
	if !ok {
		t.Fatalf("expected instrument object, got %T", body["instrument"])
	}
	credential, ok := instrument["credential"].(map[string]any)
	if !ok {
		t.Fatalf("expected credential object, got %T", instrument["credential"])
	}
	if got, _ := credential["token"].(string); got != "spt_123" {
		t.Fatalf("unexpected credential token %q", got)
	}
}
