package faas

import (
	"context"
	"kis-flow/kis"
	"kis-flow/test/proto"
)

type AvgStuScoreIn struct {
	proto.StuScores
}

type AvgStuScoreOut struct {
	proto.StuAvgScore
}

func AvgStuScore(ctx context.Context, flow kis.Flow, rows []*AvgStuScoreIn) error {
	for _, row := range rows {
		avgScore := proto.StuAvgScore{
			StuId:    row.StuId,
			AvgScore: float64(row.Score1+row.Score2+row.Score3) / 3,
		}
		// 提交结果数据
		_ = flow.CommitRow(AvgStuScoreOut{StuAvgScore: avgScore})
	}

	return nil
}
