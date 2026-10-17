package acpcheckout

//go:generate go tool oapi-codegen -config oapi-codegen.yaml spec/openapi.acpcheckout.yaml
//go:generate go tool mockgen -source=handler.go -destination=handler_mock_test.go -package=acpcheckout_test
