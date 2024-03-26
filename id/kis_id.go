package id

import (
	"github.com/aceld/kis-flow/common"
	"github.com/google/uuid"
	"strings"
)

// KisID 获取一个中随机实例ID
// 格式为  "prefix1-[prefix2-][prefix3-]ID"
// 如：flow-1234567890
// 如：func-1234567890
// 如: conn-1234567890
// 如: func-1-1234567890
func KisID(prefix ...string) (kisId string) {

	idStr := strings.Replace(uuid.New().String(), "-", "", -1)
	kisId = formatKisID(idStr, prefix...)

	return
}

func formatKisID(idStr string, prefix ...string) string {
	var kisId string

	for _, fix := range prefix {
		kisId += fix
		kisId += common.KisIdJoinChar
	}

	kisId += idStr

	return kisId
}
