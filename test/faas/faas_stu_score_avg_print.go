package faas

import (
	"context"
	"fmt"
	"kis-flow/kis"
	"kis-flow/test/proto"
)

type PrintStuAvgScoreIn struct {
	proto.StuAvgScore
}

type PrintStuAvgScoreOut struct {
}

func PrintStuAvgScore(ctx context.Context, flow kis.Flow, rows []*PrintStuAvgScoreIn) error {

	for _, row := range rows {
		fmt.Printf("stuid: [%+v], avg score: [%+v]\n", row.StuId, row.AvgScore)
	}

	return nil
}
