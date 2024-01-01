package kis

import (
	"context"
	"kis-flow/config"
)

type Flow interface {
	// Run 调度Flow，依次调度Flow中的Function并且执行
	Run(ctx context.Context) error
	// Link 将Flow中的Function按照配置文件中的配置进行连接
	Link(fConf *config.KisFuncConfig, fParams config.FParam) error
}
