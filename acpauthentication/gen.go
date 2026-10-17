package acpauthentication

//go:generate go tool oapi-codegen -config oapi-codegen.yaml spec/openapi.acpauthentication.yaml
//go:generate go tool mockgen -source=handler.go -destination=handler_mock_test.go -package=acpauthentication_test
