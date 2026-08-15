// Package acpcart serves the seller-hosted ACP Cart API.
//
// Carts provide a lightweight pre-checkout phase for item collection without
// payment configuration or a status lifecycle. Totals are estimates until
// checkout. Updates replace the complete mutable cart state.
//
// Implement [Provider] and pass it to [NewHandler].
package acpcart
