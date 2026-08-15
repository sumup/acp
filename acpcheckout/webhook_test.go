package acpcheckout

import "testing"

func TestSignWebhookPayload(t *testing.T) {
	t.Parallel()

	got := signWebhookPayload([]byte("test-secret"), 1_709_123_456, []byte(`{"x":1}`))
	want := "t=1709123456,v1=bcabec3508154efb9c54bb592b8634986ad1dc4970874a1735ffa36c47d2742c"
	if got != want {
		t.Fatalf("signWebhookPayload() = %q, want %q", got, want)
	}
}
