package faas

import (
	"context"

	"github.com/aceld/kis-flow/kis"
	"github.com/aceld/kis-flow/serialize"
	"github.com/aceld/kis-flow/test/proto"
)

type AvgStuScoreIn struct {
	serialize.DefaultSerialize
	proto.StuScores
}

type AvgStuScoreOut struct {
	serialize.DefaultSerialize
	proto.StuAvgScore
}

// AvgStuScore(FaaS) calculates the average score of students
func AvgStuScore(ctx context.Context, flow kis.Flow, rows []*AvgStuScoreIn) error {
	for _, row := range rows {
		avgScore := proto.StuAvgScore{
			StuId:    row.StuId,
			AvgScore: float64(row.Score1+row.Score2+row.Score3) / 3,
		}

		// Submit result data
		_ = flow.CommitRow(AvgStuScoreOut{StuAvgScore: avgScore})
	}

	return nil
}
