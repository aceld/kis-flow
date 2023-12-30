# kis-flow

基于Golang的流式计算框架.

KisFlow(Keep It Simple Flow).为保持简单的流动，强调在进行各种活动或工作时保持简洁、清晰、流畅的过程。


## 开发及教程文档

KisFlow 采用开发教程同步执行，教程文档地址：https://www.yuque.com/aceld/hsa94o





<a name="KPh6H"></a>
## 1.1 为什么需要KisFlow
一些大型toB企业级的项目，需要大量的业务数据，多数的数据需要流式实时计算的能力，但是很多公司还不足以承担一个数仓类似，Flink + Hadoop/HBase 等等。 但是业务数据的实时计算需求依然存在，所以大多数的企业依然会让业务工程师来消化这些业务数据计算的工作。

而这样只能直接查询业务数据库，这样会间接影响的业务能力，或定时任务/脚本来做定时计算，这些都不是好的办法。本人亲身经历过一个大规模的系统，多达上千个需要计算的业务数据字段，而早期因为没有规划好，最后导致存在1000+的脚本在定时跑，最后导致了脚本之间对数据的影响，数据始终无法准确，导致业务数据经常性的报数据问题错误。<br />       如下面一个场景：某个业务计算字段的值，正确为100，错误为99， 但是由于历史代码的臃肿，会有多个计算脚本对其值做修复补丁计算，会有各个脚本相互冲突，在一定的时间间隔内会存在数据值抖动，可能最终一个补丁修复正确，但是这种情况就会存在一定时间范围内业务数据不正确，最终却奇迹正确的情况，很让用户苦恼。<br />
![](https://cdn.nlark.com/yuque/0/2023/jpeg/26269664/1702626843206-995cb619-e329-4f5b-83eb-e47780dbe277.jpeg)

KisFlow就是为了解决当企业不具备数仓平台的计算能力，又依然存在大量数据实时计算的场景，让业务工程师可以投入到数据流式计算的业务中来，并且可以复用常用和通用的计算逻辑。


<a name="aSEWt"></a>
## 1.2 KisFlow实要支持的能力
<a name="d4Nt0"></a>
### 流式计算
1、分布式批量消费能力（基于上游ODS消费配置：如Binlog、Kafka等）<br />2、Stateful Function能力，基于有状态的流式计算节点拼接，流式计算横纵向扩展。<br />3、数据流监控及修复能力，消费服务监控。<br />4、多流拼接及第三方中间件存储插件化。
<a name="oV4gp"></a>
### 分布式任务调度
5、分布式定时任务调度、日志监控、任务调度状态。<br />6、可视化调度平台。

<a name="TQFqe"></a>
## 1.3 KisFlow系统定位
KisFlow为业务上游计算层，上层接数仓/其他业务方ODS层、下游接本业务存储数据中心。<br />
![](https://cdn.nlark.com/yuque/0/2023/jpeg/26269664/1702626531446-964eaeee-bf3c-4ef8-a1db-04cb99f1b1cc.jpeg)

<a name="elhiR"></a>
## 1.4  KisFlow整体架构图

| 层级 | 层级说明 | 包括子模块 |
| --- | --- | --- |
| 流式计算层 | 为KisFlow上游计算层，直接对接业务存储及数仓ODS层，如上游可以为Mysql Binlog、日志、接口数据等，为被动消费模式，提供KisFlow实时计算能力。 | **KisFlow**：分布式批量消费者,一个KisFlow是由多个KisFunction组合。<br /><br />**KisConnectors**：计算数据流流中间状态持久存储及连接器。<br /><br />**KisFunctions**：支持算子表达式拼接，Connectors集成、策略配置、Stateful Function模式、Slink流式拼接等。<br /><br />**KisConfig：**KisFunction的绑定的流处理策略，可以绑定ReSource让Function具有固定的独立流处理能力。<br /><br />**KisSource：**对接ODS的数据源 |
| 任务调度层 | 定时任务调度及执行器业务逻辑，包括任务调度平台、执行器管理、调度日志及用户管理等。提供KisFlow的定时任务、统计、聚合运算等调度计算能力。 | **任务调度平台可视化**：包括任务的运行报表、调度报表、成功比例、任务管理、配置管理、GLUE IDE等可视化管理平台。<br /><br />执行器管理**KisJobs**：Golang SDK及计算自定义业务逻辑、执行器的自动注册、任务触发、终止及摘除等。<br /><br />**执行器场景KisScenes：**根据业务划分的逻辑任务集合。<br /><br />**调度日志及用户管理**：任务调度日志收集、调度详细、调度流程痕迹等。 |

![KisFlow架构图drawio.png](https://cdn.nlark.com/yuque/0/2023/png/26269664/1703834438819-88fc68ca-c078-475e-8733-98729d0ec3da.png#averageHue=%23b5cc5b&clientId=ua6b75298-2e7b-4&from=drop&id=u7a4debd6&originHeight=3023&originWidth=2932&originalType=binary&ratio=2&rotation=0&showTitle=false&size=1819694&status=done&style=none&taskId=u46c44eb8-1224-4d34-86aa-4d0f0afd935&title=)

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

![KisFlow架构设计-KisFlow整体结构.drawio.png](https://cdn.nlark.com/yuque/0/2023/png/26269664/1703140767204-477b49a3-f5a9-4515-a171-ed649c0ca8b8.png#averageHue=%23cccd6f&clientId=ue4e14b5e-7e82-4&from=drop&id=u6f81a11a&originHeight=1658&originWidth=1962&originalType=binary&ratio=2&rotation=0&showTitle=false&size=344404&status=done&style=none&taskId=u29d95873-83b8-4c16-827a-88f44eb5f28&title=)
KisFlow是一种流式概念形态，具体表现的特征如下：<br />1、一个KisFlow可以由任意KisFunction组成，且KisFlow可以动态的调整长度。<br />2、一个KisFunction可以随时动态的加入到某个KisFlow中，且KisFlow和KisFlow之间的关系可以通过KisFunction的Load和Save节点的加入，进行动态的并流和分流动作。<br />3、KisFlow在编程行为上，从面向流进行数据业务编程，变成了面向KisFunction的函数单计算逻辑的开发，接近FaaS(Function as a service)体系。

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
flow:
  - MyFlow1
```

#### (4) KisFlow全局配置

```yaml
#kistype Global为kisflow的全局配置
kistype: global
#是否启动prometheus监控
prometheus_enable: true
#是否需要nsflow单独启动端口监听
prometheus_listen: true
#prometheus取点监听地址
prometheus_serve: 0.0.0.0:20004
```


