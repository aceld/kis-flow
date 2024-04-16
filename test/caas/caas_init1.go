package caas

import (
	"fmt"

	"github.com/aceld/kis-flow/kis"
)

// type ConnInit func(conn Connector) error

func InitConnDemo1(connector kis.Connector) error {
	fmt.Println("===> Call Connector InitDemo1")
	//config info
	connConf := connector.GetConfig()

	fmt.Println(connConf)

	// init connector

	return nil
}
