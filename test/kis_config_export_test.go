package test

import (
	"fmt"
	"testing"

	"github.com/aceld/kis-flow/file"
	"github.com/aceld/kis-flow/kis"
)

func TestConfigExportYaml(t *testing.T) {

	// 1. Load the configuration file and build the Flow
	if err := file.ConfigImportYaml("load_conf/"); err != nil {
		fmt.Println("Wrong Config Yaml Path!")
		panic(err)
	}

	// 2. Export the built memory KisFlow structure configuration to files
	flows := kis.Pool().GetFlows()
	for _, flow := range flows {
		if err := file.ConfigExportYaml(flow, "/Users/Aceld/go/src/kis-flow/test/export_conf/"); err != nil {
			panic(err)
		}
	}
}
