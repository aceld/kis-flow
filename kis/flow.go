package kis

import (
	"context"
	"kis-flow/common"
	"kis-flow/config"
)

type Flow interface {
	// Run 调度Flow，依次调度Flow中的Function并且执行
	Run(ctx context.Context) error
	// Link 将Flow中的Function按照配置文件中的配置进行连接
	Link(fConf *config.KisFuncConfig, fParams config.FParam) error
	// CommitRow 提交Flow数据到即将执行的Function层
	CommitRow(row interface{}) error
	// Input 得到flow当前执行Function的输入源数据
	Input() common.KisRowArr
	// GetName 得到Flow的名称
	GetName() string
	// GetThisFunction 得到当前正在执行的Function
	GetThisFunction() Function
	// GetThisFuncConf 得到当前正在执行的Function的配置
	GetThisFuncConf() *config.KisFuncConfig
}
