package caas

import (
	"fmt"
	"kis-flow/kis"
)

// type ConnInit func(conn Connector) error

func InitConnDemo1(connector kis.Connector) error {
	fmt.Println("===> Call Connector InitDemo1")
	//config info
	connConf := connector.GetConfig()

	fmt.Println(connConf)

	// init connector , 如 初始化数据库链接等

	return nil
}
