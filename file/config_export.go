package file

import (
	"errors"
	"fmt"
	"kis-flow/common"
	"kis-flow/kis"
	"os"

	"gopkg.in/yaml.v3"
)

// ConfigExportYaml 将flow配置输出，且存储本地
func ConfigExportYaml(flow kis.Flow, savaPath string) error {

	if data, err := yaml.Marshal(flow.GetConfig()); err != nil {
		return err
	} else {
		// flow
		err := os.WriteFile(savaPath+common.KisIdTypeFlow+"-"+flow.GetName()+".yaml", data, 0644)
		if err != nil {
			return err
		}

		// function
		for _, fp := range flow.GetConfig().Flows {
			fConf := flow.GetFuncConfigByName(fp.FuncName)
			if fConf == nil {
				return errors.New(fmt.Sprintf("function name = %s config is nil ", fp.FuncName))
			}

			if fdata, err := yaml.Marshal(fConf); err != nil {
				return err
			} else {
				if err := os.WriteFile(savaPath+common.KisIdTypeFunction+"-"+fp.FuncName+".yaml", fdata, 0644); err != nil {
					return err
				}
			}

			// Connector
			if fConf.Option.CName != "" {
				cConf, err := fConf.GetConnConfig()
				if err != nil {
					return err
				}
				if cdata, err := yaml.Marshal(cConf); err != nil {
					return err
				} else {
					if err := os.WriteFile(savaPath+common.KisIdTypeConnector+"-"+cConf.CName+".yaml", cdata, 0644); err != nil {
						return err
					}
				}
			}
		}
	}

	return nil
}
