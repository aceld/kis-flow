# kis-flow


#### KisFlow(Keep It Simple Flowing) 

基于Golang的流式计算框架. 为保持简单的流动，强调在进行各种活动或工作时保持简洁、清晰、流畅的过程。



---

## 开发及教程文档(持续更新中...)


KisFlow 采用开发教程同步执行 

教程文档地址：https://www.yuque.com/aceld/hsa94o



---





## KisFlow系统定位

KisFlow为业务上游计算层，上层接数仓/其他业务方ODS层、下游接本业务存储数据中心。<br />
![yuque_diagram (2)](https://github.com/aceld/kis-flow/assets/7778936/b9e1957a-2d11-45d9-84c1-e92c9ac833cc)


<a name="elhiR"></a>
## KisFlow整体架构图

| 层级 | 层级说明 | 包括子模块 |
| --- | --- | --- |
| 流式计算层 | 为KisFlow上游计算层，直接对接业务存储及数仓ODS层，如上游可以为Mysql Binlog、日志、接口数据等，为被动消费模式，提供KisFlow实时计算能力。 | **KisFlow**：分布式批量消费者,一个KisFlow是由多个KisFunction组合。<br /><br />**KisConnectors**：计算数据流流中间状态持久存储及连接器。<br /><br />**KisFunctions**：支持算子表达式拼接，Connectors集成、策略配置、Stateful Function模式、Slink流式拼接等。<br /><br />**KisConfig：**KisFunction的绑定的流处理策略，可以绑定ReSource让Function具有固定的独立流处理能力。<br /><br />**KisSource：**对接ODS的数据源 |
| 任务调度层 | 定时任务调度及执行器业务逻辑，包括任务调度平台、执行器管理、调度日志及用户管理等。提供KisFlow的定时任务、统计、聚合运算等调度计算能力。 | **任务调度平台可视化**：包括任务的运行报表、调度报表、成功比例、任务管理、配置管理、GLUE IDE等可视化管理平台。<br /><br />执行器管理**KisJobs**：Golang SDK及计算自定义业务逻辑、执行器的自动注册、任务触发、终止及摘除等。<br /><br />**执行器场景KisScenes：**根据业务划分的逻辑任务集合。<br /><br />**调度日志及用户管理**：任务调度日志收集、调度详细、调度流程痕迹等。 |

![KisFlow架构图drawio](https://github.com/aceld/kis-flow/assets/7778936/3b829bdb-600d-4ab9-9e62-e14f90737cc3)

| 流 | 组成 |
| --- | --- |
| KisFlow(1) | KisFunction(V) + KisFunction(S) + KisFunction(C) + KisFunction(E) |
| KisFlow(2) | KisFunction(V) + KisFunction(L) + KisFunction(S) + KisFunction(C) + KisFunction(E) |
| KisFlow(3) | KisFunction(V) + KisFunction(L) + KisFunction(C) + KisFunction(E) |


通过 KisFunction(S) 和 KisFunction(L)的并流组合关系，各个KisFlow有如下关系：
```yaml
KisFlow(2) = KisFlow(1) + KisFlow(2)
KisFlow(3) = KisFlow(1) + KisFlow(2) + KisFlow(3)
```


#### (1) KisFunction配置
```yaml
kistype: func
fid: 测试KisFunction_S1
fname: 测试KisFunction_S1
fmode: Save
source:
  name: 被校验的测试数据源1-用户订单维度
  must:
    - userid
    - orderid
    
option:
  cid: 测试KisConnector_1
  retry_times: 3
  retry_duration: 500
  default_params:
    default1: default1_param
    default2: default2_param
```

#### (2) KisFlow配置
```yaml
kistype: flow
flow_id: MyFlow1
status: 1
flow_name: MyFlow1
flows:
  - fid: 测试PrintInput
    params:
      args1: value1
      args2: value2
  - fid: 测试KisFunction_S1
  - fid: 测试PrintInput
    params:
      args1: value11
      args2: value22
      default2: newDefault
  - fid: 测试PrintInput
  - fid: 测试KisFunction_S1
    params:
      my_user_param1: ffffffxxxxxx
  - fid: 测试PrintInput
```


![KisFlow架构设计-KisFlow整体结构 drawio](https://github.com/aceld/kis-flow/assets/7778936/efc1b29d-9dd4-4945-a35a-fb9a618002d7)


KisFlow是一种流式概念形态，具体表现的特征如下：<br />

1、一个KisFlow可以由任意KisFunction组成，且KisFlow可以动态的调整长度。<br />

2、一个KisFunction可以随时动态的加入到某个KisFlow中，且KisFlow和KisFlow之间的关系可以通过KisFunction的Load和Save节点的加入，进行动态的并流和分流动作。<br />

3、KisFlow在编程行为上，从面向流进行数据业务编程，变成了面向KisFunction的函数单计算逻辑的开发，接近FaaS(Function as a service)体系。

#### (3) KisConnector配置

```yaml
kistype: conn
cid: 测试KisConnector_1
cname: 测试KisConnector_1
addrs: '0.0.0.0:9988,0.0.0.0:9999,0.0.0.0:9990'
type: redis
key: userid_orderid_option
params:
  args1: value1
  args2: value2
load: null
save:
  - 测试KisFunction_S1
```

#### (4) KisFlow全局配置

```yaml
#kistype Global为kisflow的全局配置
kistype: global
#是否启动prometheus监控
prometheus_enable: true
#是否需要kisflow单独启动端口监听
prometheus_listen: true
#prometheus取点监听地址
prometheus_serve: 0.0.0.0:20004
```


