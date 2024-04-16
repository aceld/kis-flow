package id

import (
	"strings"

	"github.com/aceld/kis-flow/common"
	"github.com/google/uuid"
)

// KisID generates a random instance ID.
// The format is "prefix1-[prefix2-][prefix3-]ID"
// Example: flow-1234567890
// Example: func-1234567890
// Example: conn-1234567890
// Example: func-1-1234567890
func KisID(prefix ...string) (kisId string) {

	idStr := strings.Replace(uuid.New().String(), "-", "", -1)
	kisId = formatKisID(idStr, prefix...)

	return
}

func formatKisID(idStr string, prefix ...string) string {
	var kisId string

	for _, fix := range prefix {
		kisId += fix
		kisId += common.KisIDJoinChar
	}

	kisId += idStr

	return kisId
}
