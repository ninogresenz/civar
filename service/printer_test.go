package service_test

import (
	"testing"

	"github.com/bradleyjkemp/cupaloy"

	"github.com/ninogresenz/civar/service"
)

func TestPrinter(t *testing.T) {
	tests := map[string]service.CiPrinter{
		"TestDotenvPrinter": service.PrinterProvider("dotenv"),
		"TestPrettyPrinter": service.PrinterProvider("pretty"),
		"TestJsonPrinter":   service.PrinterProvider("json"),
	}
	for testName, printer := range tests {
		t.Run(testName, func(t *testing.T) {
			cupaloy.SnapshotT(t, printer.Print(getVars()))
		})
	}
}
