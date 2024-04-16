package test

import (
	"context"
	"testing"

	"github.com/aceld/kis-flow/common"
	"github.com/aceld/kis-flow/config"
	"github.com/aceld/kis-flow/file"
	"github.com/aceld/kis-flow/flow"
	"github.com/aceld/kis-flow/kis"
	"github.com/aceld/kis-flow/test/faas"
	"github.com/aceld/kis-flow/test/proto"
)

func TestAutoInjectParamWithConfig(t *testing.T) {
	ctx := context.Background()

	kis.Pool().FaaS("AvgStuScore", faas.AvgStuScore)
	kis.Pool().FaaS("PrintStuAvgScore", faas.PrintStuAvgScore)

	// 1. Load the configuration file and build the Flow
	if err := file.ConfigImportYaml("load_conf/"); err != nil {
		panic(err)
	}

	// 2. Get the Flow
	flow1 := kis.Pool().GetFlow("StuAvg")
	if flow1 == nil {
		panic("flow1 is nil")
	}

	// 3. Commit original data
	_ = flow1.CommitRow(&faas.AvgStuScoreIn{
		StuScores: proto.StuScores{
			StuId:  100,
			Score1: 1,
			Score2: 2,
			Score3: 3,
		},
	})
	_ = flow1.CommitRow(faas.AvgStuScoreIn{
		StuScores: proto.StuScores{
			StuId:  100,
			Score1: 1,
			Score2: 2,
			Score3: 3,
		},
	})

	// Commit original data (as JSON string)
	_ = flow1.CommitRow(`{"stu_id":101}`)

	// 4. Execute flow1
	if err := flow1.Run(ctx); err != nil {
		panic(err)
	}
}

func TestAutoInjectParam(t *testing.T) {
	ctx := context.Background()

	kis.Pool().FaaS("AvgStuScore", faas.AvgStuScore)
	kis.Pool().FaaS("PrintStuAvgScore", faas.PrintStuAvgScore)

	source1 := config.KisSource{
		Name: "Test",
		Must: []string{},
	}

	avgStuScoreConfig := config.NewFuncConfig("AvgStuScore", common.C, &source1, nil)
	if avgStuScoreConfig == nil {
		panic("AvgStuScore is nil")
	}

	printStuAvgScoreConfig := config.NewFuncConfig("PrintStuAvgScore", common.C, &source1, nil)
	if printStuAvgScoreConfig == nil {
		panic("printStuAvgScoreConfig is nil")
	}

	myFlowConfig1 := config.NewFlowConfig("cal_stu_avg_score", common.FlowEnable)

	flow1 := flow.NewKisFlow(myFlowConfig1)

	// 4. Link Functions to Flow
	if err := flow1.Link(avgStuScoreConfig, nil); err != nil {
		panic(err)
	}
	if err := flow1.Link(printStuAvgScoreConfig, nil); err != nil {
		panic(err)
	}

	// 3. Commit original data
	_ = flow1.CommitRow(&faas.AvgStuScoreIn{
		StuScores: proto.StuScores{
			StuId:  100,
			Score1: 1,
			Score2: 2,
			Score3: 3,
		},
	})
	_ = flow1.CommitRow(`{"stu_id":101}`)
	_ = flow1.CommitRow(faas.AvgStuScoreIn{
		StuScores: proto.StuScores{
			StuId:  100,
			Score1: 1,
			Score2: 2,
			Score3: 3,
		},
	})

	// 4. Execute flow1
	if err := flow1.Run(ctx); err != nil {
		panic(err)
	}
}
