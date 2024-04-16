# <img width="250px" src="https://github.com/aceld/kis-flow/assets/7778936/8729d750-897c-4ba3-98b4-c346188d034e" />
English | [简体中文](README-CN.md)


[![License](https://img.shields.io/badge/License-MIT-black.svg)](LICENSE)
[![Discord](https://img.shields.io/badge/KisFlow-Discord-blue.svg)](https://discord.gg/xQ8Xxfyfcz)
[![KisFlow-tutorial](https://img.shields.io/badge/KisFlowTutorial-YuQue-red.svg)](https://www.yuque.com/aceld/kis-flow)
[![KisFlow-Doc](https://img.shields.io/badge/KisFlow-Doc-green.svg)](https://www.yuque.com/aceld/kis-flow-doc)


#### KisFlow(Keep It Simple Flowing)

A Streaming Computation Framework Based on Golang. Emphasizes maintaining a simple, clear, and smooth process while performing various activities or tasks.

## Source of KisFlow

Github
Git: https://github.com/aceld/kis-flow

Gitee(China)
Git: https://gitee.com/Aceld/kis-flow

## Documentation

[ < KisFlow Wiki : English > ](https://github.com/aceld/kis-flow/wiki)

[ < KisFlow 文档 : 简体中文> ](https://www.yuque.com/aceld/kis-flow-doc)



## Online Tutorial



| platform | Entry                                                                                                                                              | 
| ---- |----------------------------------------------------------------------------------------------------------------------------------------------------| 
| <img src="https://user-images.githubusercontent.com/7778936/236784004-b6d99e26-b1ab-4bc3-988e-7a46108b85fe.png" width = "100" height = "100" alt="" align=center />| [Practical Tutorial for a Streaming Computation Framework Based on Golang](https://dev.to/aceld/1building-basic-services-with-zinx-framework-296e) | 
|<img src="https://user-images.githubusercontent.com/7778936/236784168-6528a9b8-d37b-4b02-a37c-b9988d7508d8.jpeg" width = "100" height = "100" alt="" align=center />| [《基于Golang的流式计算框架实战教程》](https://www.yuque.com/aceld/hsa94o)                                                                                              |


## Positioning of the KisFlow System

KisFlow serves as the upstream computing layer for business, connecting to the ODS layer of data warehouses or other business methods upstream, and connecting to the data center of this business's storage downstream. <br />

<img width="700px" src="https://github.com/aceld/kis-flow/assets/7778936/c7db6df0-fcec-4b4d-94b9-4aba71cd55c3" />


## KisFlow Overall Architecture Diagram


| Levels    | Level Explanation                                                                               | Sub-modules                                                                                                                                                                                                                                                                                                           |
|-------|------------------------------------------------------------------------------------|-----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| Flowing Computation Layer | The upstream computing layer for KisFlow, which directly connects to business storage and the ODS (Operational Data Store) layer of data warehouses. The upstream can be MySQL Binlog, logs, interface data, etc., and it supports a passive consumption mode, providing KisFlow with real-time computing capabilities. | **KisFlow**: Distributed batch consumer; a KisFlow is composed of multiple KisFunctions. <br /><br />**KisConnectors**: Computing data stream intermediate state persistence and connectors. <br /><br />**KisFunctions**: Supports operator expression splicing, connector integration, strategy configuration, Stateful Function mode, and Slink stream splicing. <br /><br />**KisConfig**: Binding of flow processing policies for KisFunctions, allowing Functions to have fixed independent processing capabilities. <br /><br />**KisSource**: Interface for connecting to ODS data sources. |
| Task Scheduling Layer | Timed task scheduling and execution business logic, including task scheduling platform, executor management, scheduling logs, and user management. Provides KisFlow's timed task, statistics, and aggregation calculation capabilities.  | **The task scheduling platform has a visual interface.**：ncludes running reports, scheduling reports, success rate, task management, configuration management, and GLUE IDE as visual management platforms. <br /><br /> **Executor management KisJobs**: Golang SDK, custom business logic, executor automatic registration, task triggering, termination, and removal.<br /><br /> **Executor scenarios KisScenes**: Logical task sets divided according to business needs.<br /><br /> **Scheduling logs and user management**: Collection of task scheduling logs, detailed scheduling, and scheduling process traces.                                                                              |

![KisFlow](https://github.com/aceld/kis-flow/assets/7778936/59dceaf3-16e7-41fc-979f-98297638416f)


<img width="700px" src="https://github.com/aceld/kis-flow/assets/7778936/efc1b29d-9dd4-4945-a35a-fb9a618002d7" />


KisFlow is a flow-based conceptual form, and its specific characteristics are as follows: <br />

1. A KisFlow can be composed of any KisFunction(s), and the length of a KisFlow can be dynamically adjusted. <br />

2. A KisFunction can be dynamically added to a specific KisFlow at any time, and the relationship between KisFlows can be dynamically adjusted through the addition of KisFunction's Load and Save nodes for parallel and branching actions.<br />

3. In programming behavior, KisFlow has shifted from data business programming to function-based single computing logic development, approaching the FaaS (Function as a Service) system.

## Example

Below is a simple application scenario case, please refer to the specific application unit cases.

https://github.com/aceld/kis-flow-usage

#### [KisFlow Developer Documentation](https://github.com/aceld/kis-flow/wiki)



#### Install KisFlow

```bash
$go get github.com/aceld/kis-flow
```

<details>
<summary>1. Quick Start </summary>

### Source Code

https://github.com/aceld/kis-flow-usage/tree/main/1-quick_start

### Project Directory

```bash
├── faas_stu_score_avg.go
├── faas_stu_score_avg_print.go
└── main.go
```

### Flow

<img width="770" alt="image" src="https://github.com/aceld/kis-flow/assets/7778936/3747ed10-aba1-417e-a3c1-c6205a02444b">

### Main

> main.go

```go
package main

import (
	"context"
	"fmt"
	"github.com/aceld/kis-flow/common"
	"github.com/aceld/kis-flow/config"
	"github.com/aceld/kis-flow/flow"
	"github.com/aceld/kis-flow/kis"
)

func main() {
	ctx := context.Background()

	// Create a new flow configuration
	myFlowConfig1 := config.NewFlowConfig("CalStuAvgScore", common.FlowEnable)

	// Create new function configuration
	avgStuScoreConfig := config.NewFuncConfig("AvgStuScore", common.C, nil, nil)
	printStuScoreConfig := config.NewFuncConfig("PrintStuAvgScore", common.E, nil, nil)

	// Create a new flow
	flow1 := flow.NewKisFlow(myFlowConfig1)

	// Link functions to the flow
	_ = flow1.Link(avgStuScoreConfig, nil)
	_ = flow1.Link(printStuScoreConfig, nil)

	// Submit a string
	_ = flow1.CommitRow(`{"stu_id":101, "score_1":100, "score_2":90, "score_3":80}`)
	// Submit a string
	_ = flow1.CommitRow(`{"stu_id":102, "score_1":100, "score_2":70, "score_3":60}`)

	// Run the flow
	if err := flow1.Run(ctx); err != nil {
		fmt.Println("err: ", err)
	}

	return
}

func init() {
	// Register functions
	kis.Pool().FaaS("AvgStuScore", AvgStuScore)
	kis.Pool().FaaS("PrintStuAvgScore", PrintStuAvgScore)
}
```

### Function1

> faas_stu_score_avg.go

```go
package main

import (
	"context"
	"github.com/aceld/kis-flow/kis"
	"github.com/aceld/kis-flow/serialize"
)

type AvgStuScoreIn struct {
	serialize.DefaultSerialize
	StuId  int `json:"stu_id"`
	Score1 int `json:"score_1"`
	Score2 int `json:"score_2"`
	Score3 int `json:"score_3"`
}

type AvgStuScoreOut struct {
	serialize.DefaultSerialize
	StuId    int     `json:"stu_id"`
	AvgScore float64 `json:"avg_score"`
}

// AvgStuScore(FaaS) 计算学生平均分
func AvgStuScore(ctx context.Context, flow kis.Flow, rows []*AvgStuScoreIn) error {
	for _, row := range rows {

		out := AvgStuScoreOut{
			StuId:    row.StuId,
			AvgScore: float64(row.Score1+row.Score2+row.Score3) / 3,
		}

		// 提交结果数据
		_ = flow.CommitRow(out)
	}

	return nil
}
```

### Function2

> faas_stu_score_avg_print.go

```go
package main

import (
	"context"
	"fmt"
	"github.com/aceld/kis-flow/kis"
	"github.com/aceld/kis-flow/serialize"
)

type PrintStuAvgScoreIn struct {
	serialize.DefaultSerialize
	StuId    int     `json:"stu_id"`
	AvgScore float64 `json:"avg_score"`
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
```

### OutPut

```bash
Add KisPool FuncName=AvgStuScore
Add KisPool FuncName=PrintStuAvgScore
funcName NewConfig source is nil, funcName = AvgStuScore, use default unNamed Source.
funcName NewConfig source is nil, funcName = PrintStuAvgScore, use default unNamed Source.
stuid: [101], avg score: [90]
stuid: [102], avg score: [76.66666666666667]
```

</details>


<details>
<summary>2. Quick Start With Config </summary>

### Source Code

https://github.com/aceld/kis-flow-usage/tree/main/2-quick_start_with_config

### Project Directory

```bash
├── Makefile
├── conf
│   ├── flow-CalStuAvgScore.yml
│   ├── func-AvgStuScore.yml
│   └── func-PrintStuAvgScore.yml
├── faas_stu_score_avg.go
├── faas_stu_score_avg_print.go
└── main.go
```

### Flow

<img width="770" alt="image" src="https://github.com/aceld/kis-flow/assets/7778936/3747ed10-aba1-417e-a3c1-c6205a02444b">

### Config

#### (1) Flow Config

> conf/flow-CalStuAvgScore.yml

```yaml
kistype: flow
status: 1
flow_name: CalStuAvgScore
flows:
  - fname: AvgStuScore
  - fname: PrintStuAvgScore
```

#### (2) Function1 Config

> conf/func-AvgStuScore.yml

```yaml
kistype: func
fname: AvgStuScore
fmode: Calculate
source:
  name: StudentScore
  must:
    - stu_id
```

#### (3) Function2(Slink) Config

> conf/func-PrintStuAvgScore.yml

```yaml
kistype: func
fname: PrintStuAvgScore
fmode: Expand
source:
  name: StudentScore
  must:
    - stu_id
```

### Main

> main.go

```go
package main

import (
	"context"
	"fmt"
	"github.com/aceld/kis-flow/file"
	"github.com/aceld/kis-flow/kis"
)

func main() {
	ctx := context.Background()

	// Load Configuration from file
	if err := file.ConfigImportYaml("conf/"); err != nil {
		panic(err)
	}

	// Get the flow
	flow1 := kis.Pool().GetFlow("CalStuAvgScore")
	if flow1 == nil {
		panic("flow1 is nil")
	}

	// Submit a string
	_ = flow1.CommitRow(`{"stu_id":101, "score_1":100, "score_2":90, "score_3":80}`)
	// Submit a string
	_ = flow1.CommitRow(`{"stu_id":102, "score_1":100, "score_2":70, "score_3":60}`)

	// Run the flow
	if err := flow1.Run(ctx); err != nil {
		fmt.Println("err: ", err)
	}

	return
}

func init() {
	// Register functions
	kis.Pool().FaaS("AvgStuScore", AvgStuScore)
	kis.Pool().FaaS("PrintStuAvgScore", PrintStuAvgScore)
}
```

### Function1

> faas_stu_score_avg.go

```go
package main

import (
	"context"
	"github.com/aceld/kis-flow/kis"
	"github.com/aceld/kis-flow/serialize"
)

type AvgStuScoreIn struct {
	serialize.DefaultSerialize
	StuId  int `json:"stu_id"`
	Score1 int `json:"score_1"`
	Score2 int `json:"score_2"`
	Score3 int `json:"score_3"`
}

type AvgStuScoreOut struct {
	serialize.DefaultSerialize
	StuId    int     `json:"stu_id"`
	AvgScore float64 `json:"avg_score"`
}

// AvgStuScore(FaaS) 计
func AvgStuScore(ctx context.Context, flow kis.Flow, rows []*AvgStuScoreIn) error {
	for _, row := range rows {

		out := AvgStuScoreOut{
			StuId:    row.StuId,
			AvgScore: float64(row.Score1+row.Score2+row.Score3) / 3,
		}

		// Commit Result Data
		_ = flow.CommitRow(out)
	}

	return nil
}
```

### Function2

> faas_stu_score_avg_print.go

```go
package main

import (
	"context"
	"fmt"
	"github.com/aceld/kis-flow/kis"
	"github.com/aceld/kis-flow/serialize"
)

type PrintStuAvgScoreIn struct {
	serialize.DefaultSerialize
	StuId    int     `json:"stu_id"`
	AvgScore float64 `json:"avg_score"`
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
```

### OutPut

```bash
Add KisPool FuncName=AvgStuScore
Add KisPool FuncName=PrintStuAvgScore
Add FlowRouter FlowName=CalStuAvgScore
stuid: [101], avg score: [90]
stuid: [102], avg score: [76.66666666666667]
```

</details>


---

### KisFlow Contributors

* 刘丹冰([@aceld](https://github.com/aceld))
* 胡辰豪([@ChenHaoHu](https://github.com/ChenHaoHu))

Thanks to all the developers who contributed to KisFlow!

<a href="https://github.com/aceld/kis-flow/graphs/contributors">
  <img src="https://contrib.rocks/image?repo=aceld/kis-flow" />
</a>    

### Join the KisFlow Community

| platform                                                                                                                                                             | Entry                                                                                                                                                                                                                                      | 
|----------------------------------------------------------------------------------------------------------------------------------------------------------------------|--------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------| 
| <img src="https://user-images.githubusercontent.com/7778936/236775008-6bd488e3-249a-4d43-8885-7e3889e11e2d.png" width = "100" height = "100" alt="" align=center />  | https://discord.gg/xQ8Xxfyfcz                                                                                                                                                                                                              | 
| <img src="https://user-images.githubusercontent.com/7778936/236775137-5381f8a6-f534-49c4-8628-e52bf245c3bc.jpeg" width = "100" height = "100" alt="" align=center /> | Add WeChat: ace_ld or scan the QR code, and note flow to proceed. </br><img src="https://user-images.githubusercontent.com/7778936/236781258-2f0371bd-5797-49e8-a74c-680e9f15843d.png" width = "150" height = "150" alt="" align=center /> |
| <img src="https://user-images.githubusercontent.com/7778936/236778547-9cdadfb6-0f62-48ac-851a-b940389038d0.jpeg" width = "100" height = "100" alt="" align=center /> | <img src="https://s1.ax1x.com/2020/07/07/UFyUdx.th.jpg" height = "150"  alt="" align=center /> **WeChat Public Account**                                                                                                                   |
| <img src="https://user-images.githubusercontent.com/7778936/236779000-70f16c8f-0eec-4b5f-9faa-e1d5229a43e0.png" width = "100" height = "100" alt="" align=center />  | <img src="https://github.com/aceld/zinx/assets/7778936/461b409f-6337-48a8-826b-a7a746aaee31" width = "150" height = "150" alt="" align=center /> **QQ Group**                                                                              |

