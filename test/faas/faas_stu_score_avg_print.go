package faas

import (
	"context"
	"fmt"
	"kis-flow/kis"
	"kis-flow/serialize"
	"kis-flow/test/proto"
)

type PrintStuAvgScoreIn struct {
	serialize.DefaultSerialize
	proto.StuAvgScore
}

type PrintStuAvgScoreOut struct {
	serialize.DefaultSerialize
}

func PrintStuAvgScore(ctx context.Context, flow kis.Flow, rows []*PrintStuAvgScoreIn) error {

	for _, row := range rows {
		fmt.Printf("stuid: [%+v], avg score: [%+v]\n", row.StuId, row.AvgScore)
	}

	return nil
}
