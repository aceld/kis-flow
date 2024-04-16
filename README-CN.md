# <img width="250px" src="https://github.com/aceld/kis-flow/assets/7778936/8729d750-897c-4ba3-98b4-c346188d034e" />
[English](README.md) | 简体中文

[![License](https://img.shields.io/badge/License-MIT-black.svg)](LICENSE)
[![Discord](https://img.shields.io/badge/KisFlow-Discord-blue.svg)](https://discord.gg/xQ8Xxfyfcz)
[![KisFlow-tutorial](https://img.shields.io/badge/KisFlowTutorial-YuQue-red.svg)](https://www.yuque.com/aceld/kis-flow) 
[![KisFlow-Doc](https://img.shields.io/badge/KisFlow-Doc-green.svg)](https://www.yuque.com/aceld/kis-flow-doc)


#### KisFlow(Keep It Simple Flowing)

基于Golang的流式计算框架. 为保持简单的流动，强调在进行各种活动或工作时保持简洁、清晰、流畅的过程。



## KisFlow源代码

Github
Git: https://github.com/aceld/kis-flow

Gitee(China)
Git: https://gitee.com/Aceld/kis-flow

## 《KisFlow开发者文档》

[ < KisFlow Wiki : English > ](https://github.com/aceld/kis-flow/wiki)

[ < KisFlow 文档 : 简体中文> ](https://www.yuque.com/aceld/kis-flow-doc)


## 在线开发教程


| platform | Entry                                                                                                                                              | 
| ---- |----------------------------------------------------------------------------------------------------------------------------------------------------| 
| <img src="https://user-images.githubusercontent.com/7778936/236784004-b6d99e26-b1ab-4bc3-988e-7a46108b85fe.png" width = "100" height = "100" alt="" align=center />| [Practical Tutorial for a Streaming Computation Framework Based on Golang](https://dev.to/aceld/1building-basic-services-with-zinx-framework-296e) | 
|<img src="https://user-images.githubusercontent.com/7778936/236784168-6528a9b8-d37b-4b02-a37c-b9988d7508d8.jpeg" width = "100" height = "100" alt="" align=center />| [《基于Golang的流式计算框架实战教程》](https://www.yuque.com/aceld/hsa94o)                                                                                              |


## KisFlow系统定位

KisFlow为业务上游计算层，上层接数仓/其他业务方ODS层、下游接本业务存储数据中心。<br />

<img width="700px" src="https://github.com/aceld/kis-flow/assets/7778936/b9e1957a-2d11-45d9-84c1-e92c9ac833cc" />


## KisFlow整体架构图

| 层级    | 层级说明                                                                               | 包括子模块                                                                                                                                                                                                                                                                                                           |
|-------|------------------------------------------------------------------------------------|-----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| 流式计算层 | 为KisFlow上游计算层，直接对接业务存储及数仓ODS层，如上游可以为Mysql Binlog、日志、接口数据等，为被动消费模式，提供KisFlow实时计算能力。 | **KisFlow**：分布式批量消费者,一个KisFlow是由多个KisFunction组合。<br /><br />**KisConnectors**：计算数据流流中间状态持久存储及连接器。<br /><br />**KisFunctions**：支持算子表达式拼接，Connectors集成、策略配置、Stateful Function模式、Slink流式拼接等。<br /><br />**KisConfig：** KisFunction的绑定的流处理策略，可以绑定ReSource让Function具有固定的独立流处理能力。<br /><br />**KisSource：** 对接ODS的数据源 |
| 任务调度层 | 定时任务调度及执行器业务逻辑，包括任务调度平台、执行器管理、调度日志及用户管理等。提供KisFlow的定时任务、统计、聚合运算等调度计算能力。            | **任务调度平台可视化**：包括任务的运行报表、调度报表、成功比例、任务管理、配置管理、GLUE IDE等可视化管理平台。<br /><br />执行器管理**KisJobs**：Golang SDK及计算自定义业务逻辑、执行器的自动注册、任务触发、终止及摘除等。<br /><br />**执行器场景KisScenes：** 根据业务划分的逻辑任务集合。<br /><br />**调度日志及用户管理**：任务调度日志收集、调度详细、调度流程痕迹等。                                                                              |

![KisFlow架构图drawio](https://github.com/aceld/kis-flow/assets/7778936/3b829bdb-600d-4ab9-9e62-e14f90737cc3)

![KisFlow架构设计-KisFlow整体结构 drawio](https://github.com/aceld/kis-flow/assets/7778936/efc1b29d-9dd4-4945-a35a-fb9a618002d7)

KisFlow是一种流式概念形态，具体表现的特征如下：<br />

1、一个KisFlow可以由任意KisFunction组成，且KisFlow可以动态的调整长度。<br />

2、一个KisFunction可以随时动态的加入到某个KisFlow中，且KisFlow和KisFlow之间的关系可以通过KisFunction的Load和Save节点的加入，进行动态的并流和分流动作。<br />

3、KisFlow在编程行为上，从面向流进行数据业务编程，变成了面向KisFunction的函数单计算逻辑的开发，接近FaaS(Function as a
service)体系。

## Example

下面是简单的应用场景案例，具体应用单元用例请 参考

https://github.com/aceld/kis-flow-usage

#### 《KisFlow开发者文档》

https://www.yuque.com/aceld/kis-flow-doc

#### 安装KisFlow

```bash
$go get github.com/aceld/kis-flow
```

<details>
<summary>1. Quick Start（快速开始）</summary>

### 案例源代码

https://github.com/aceld/kis-flow-usage/tree/main/1-quick_start

### 项目目录

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
<summary>2. Quick Start With Config（快速开始）</summary>

### 案例源代码

https://github.com/aceld/kis-flow-usage/tree/main/2-quick_start_with_config

项目目录

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
  name: 学生学分
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
  name: 学生学分
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
Add FlowRouter FlowName=CalStuAvgScore
stuid: [101], avg score: [90]
stuid: [102], avg score: [76.66666666666667]
```

</details>


---

### 开发者

* 刘丹冰([@aceld](https://github.com/aceld))
* 胡辰豪([@ChenHaoHu](https://github.com/ChenHaoHu))

Thanks to all the developers who contributed to KisFlow!

<a href="https://github.com/aceld/kis-flow/graphs/contributors">
  <img src="https://contrib.rocks/image?repo=aceld/kis-flow" />
</a>    

### 加入KisFlow 社区

| platform                                                                                                                                                             | Entry                                                                                                                                                                                                    | 
|----------------------------------------------------------------------------------------------------------------------------------------------------------------------|----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------| 
| <img src="https://user-images.githubusercontent.com/7778936/236775008-6bd488e3-249a-4d43-8885-7e3889e11e2d.png" width = "100" height = "100" alt="" align=center />  | https://discord.gg/xQ8Xxfyfcz                                                                                                                                                                            | 
| <img src="https://user-images.githubusercontent.com/7778936/236775137-5381f8a6-f534-49c4-8628-e52bf245c3bc.jpeg" width = "100" height = "100" alt="" align=center /> | 加微信: `ace_ld`  或扫二维码，备注`flow`即可。</br><img src="https://user-images.githubusercontent.com/7778936/236781258-2f0371bd-5797-49e8-a74c-680e9f15843d.png" width = "150" height = "150" alt="" align=center /> |
| <img src="https://user-images.githubusercontent.com/7778936/236778547-9cdadfb6-0f62-48ac-851a-b940389038d0.jpeg" width = "100" height = "100" alt="" align=center /> | <img src="https://s1.ax1x.com/2020/07/07/UFyUdx.th.jpg" height = "150"  alt="" align=center /> **WeChat Public Account**                                                                                 |
| <img src="https://user-images.githubusercontent.com/7778936/236779000-70f16c8f-0eec-4b5f-9faa-e1d5229a43e0.png" width = "100" height = "100" alt="" align=center />  | <img src="https://github.com/aceld/zinx/assets/7778936/461b409f-6337-48a8-826b-a7a746aaee31" width = "150" height = "150" alt="" align=center /> **QQ Group**                                                                                 |

