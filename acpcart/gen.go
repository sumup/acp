package acpcart

//go:generate go tool oapi-codegen -config oapi-codegen.yaml spec/openapi.acpcart.yaml
//go:generate go tool mockgen -source=handler.go -destination=handler_mock_test.go -package=acpcart_test
