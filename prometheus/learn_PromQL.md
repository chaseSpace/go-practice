# 学习 PromQL

PromQL是Prometheus内置的数据查询语言，其提供对时间序列数据丰富的查询，聚合以及逻辑运算能力的支持。并且被广泛应用在Prometheus的日常应用当中，包括对数据查询、可视化、告警处理当中。

可以这么说，PromQL是Prometheus所有应用场景的基础，理解和掌握PromQL是Prometheus入门的第一课。

## 1. 初识PromQL
### 1.1 标签匹配
```shell
# 1.直接使用指标名作为表达式查询，返回所有数据
http_requests_total  # 等同于 http_requests_total{}

# 2.标签匹配
http_requests_total{instance="localhost:9090"} 
# - 或
http_requests_total{instance!="localhost:9090"}
# 2.1 标签使用正则匹配： label=~regx 是包含匹配   label!~regx是排除匹配
http_requests_total{environment=~"staging|testing|development",method!="GET"}
```

#### 瞬时向量
上面的是非范围查询，返回的结果叫做**瞬时向量**，这样的表达式叫**瞬时向量表达式**。   
**注意**：当指标没有label时，瞬时向量只返回一条时间序列数据，如下：
```shell
# 查询表达式：tcp_long_conn_req_total
# - 返回的是最新的一条数据，即现在的tcp长连接请求数总量
tcp_long_conn_req_total{}     1912
```
当指标存在label时且上报的label存在多个值时，瞬时向量会返回多条时间序列数据，如下：
```shell
# 查询表达式：rpc_total
# - 返回也是最新数据，但却不仅是时间最新，还按每个不同label值最新，即有多少label值，就会有多少条数据返回
# - 这里存在func/instance/job三个label条目，且各自又包含不同的label值，那么返回的数据条数等于它们的组合总数。
rpc_total{func="Attend", instance="asldkqo2d", job="goim.comet"}  6
rpc_total{func="BaseDecorateConf", instance="asldkqo2d", job="goim.comet"}  272
rpc_total{func="BillDetail", instance="asldkqo2d", job="goim.comet"}  25
```
#### 区间向量
与瞬时向量相对，区间向量是一段时间内的时间序列数据，使用的是**区间向量表达式**，它使用`[]`进行时间范围过滤。

### 1.2 范围查询
多数时候我们需要获取一段时间内的时间序列数据，即区间向量。

```shell
# 5m内的范围数据，时间单位支持s m h d w y
http_requests_total{}[5m]
```

### 1.3 时间位移操作
```shell
http_request_total offset 5m  # 5m前的瞬时值
http_request_total offset 1d  # 1天前的瞬时值
```

### 1.4 使用聚合操作
在返回多条数据时，我们可以使用聚合函数聚合时间，形成新的数据，并且聚合时可以按label名进行分组。

```shell
# 查询系统所有http请求的总量(不分组)
sum(http_request_total)

# 按instance进行分组求和
sum(http_request_total) by (instance)

# 按照mode计算主机CPU的平均使用时间
avg(node_cpu) by (mode)
```

### 1.5 标量和字符串
- 标量(Scalar)：一个浮点型的数字值，类似瞬时向量，但不是，比如`count(http_requests_total)`是瞬时向量，需要通过`scalar()`转化为标量；
- 字符串：一个简单的字符串值，使用单/双/反引号包含的内容，也可作为表达式；

### 1.6 合法的PromQL表达式
所有的PromQL表达式都必须至少包含一个指标名称(例如http_request_total)，或者一个不会匹配到空字符串的标签过滤器(例如{code="200"})。

因此以下两种方式，均为合法的表达式：
```shell
http_request_total # 合法
http_request_total{} # 合法
{method="get"} # 合法
```
而如下表达式，则不合法：
```shell
{job=~".*"} # 不合法
```

同时，除了使用<metric name>{label=value}的形式以外，我们还可以使用内置的`__name__`标签来指定监控指标名称：
```shell
# 这允许我们明确指定查询多个指标
{__name__=~"http_request_total"} # 合法
{__name__=~"node_disk_bytes_read|node_disk_bytes_written"} # 合法
```

## 2. PromQL操作符
### 2.1 数学运算
```shell
# 字节换算为MB
node_memory_free_bytes_total / (1024 * 1024)
```
注意，这个示例演示的是一个**瞬时向量**与一个标量之间的数学运算，数学运算符会依次作用域瞬时向量中的每一个样本值，从而得到一组新的时间序列。

#### 当两个瞬时向量之间进行数学运算
过程会相对复杂一点，如，如果我们想根据node_disk_bytes_written和node_disk_bytes_read获取主机磁盘IO的总量，可以使用如下表达式：
```shell
node_disk_bytes_written + node_disk_bytes_read
```
**查询原理**：依次找到与左边向量元素匹配（标签完全一致）的右边向量元素进行运算，如果没找到匹配元素，则直接丢弃。同时新的时间序列将不会包含指标名称。 该表达式返回结果的示例如下所示：
```shell
{device="sda",instance="localhost:9100",job="node_exporter"}=>1634967552@1518146427.807 + 864551424@1518146427.807
{device="sdb",instance="localhost:9100",job="node_exporter"}=>0@1518146427.807 + 1744384@1518146427.807
```
PromQL支持的数学运算符：`+ - * / %(取模) ^(幂运算)`

### 2.2 使用布尔运算过滤时间序列
```shell
# 查询当前所有主机节点的内存使用率
(node_memory_bytes_total - node_memory_free_bytes_total) / node_memory_bytes_total
# 进一步筛选出内存使用率超过95%的主机
(node_memory_bytes_total - node_memory_free_bytes_total) / node_memory_bytes_total > 0.95
```
瞬时向量与**标量**进行布尔运算时，PromQL依次比较向量中的所有时间序列样本的值，如果比较结果为true则保留，反之丢弃。

瞬时向量与**瞬时向量**直接进行布尔运算时，同样遵循默认的匹配模式：依次找到与左边向量元素匹配（标签完全一致）的右边向量元素进行相应的操作，如果没找到匹配元素，则直接丢弃。

目前，Prometheus支持以下布尔运算符如下：`== != > < >= <=`
### 2.3 使用bool修饰符改变布尔运算符的行为
布尔运算符的默认行为是对时序数据进行过滤。而在其它的情况下我们可能需要的是真正的布尔结果。例如，只需要知道当前模块的HTTP请求量是否>=1000，如果大于等于1000则返回1（true）否则返回0（false）。这时可以使用bool修饰符改变布尔运算的默认行为。 例如：
```shell
http_requests_total > bool 1000
```
使用bool修改符后，布尔运算**不会对时间序列进行过滤**，而是直接依次瞬时向量中的各个样本数据与标量的比较结果0或者1。从而形成一条新的时间序列。
```shell
http_requests_total{code="200",handler="query",instance="localhost:9090",job="prometheus",method="get"}  1
http_requests_total{code="200",handler="query_range",instance="localhost:9090",job="prometheus",method="get"}  0
```
同时需要注意的是，如果是在两个标量之间使用布尔运算，则必须使用bool修饰符:
```shell
2 == bool 2 # 结果为1
```

### 2.3 使用集合运算符
使用瞬时向量表达式能够获取到一个包含多个时间序列的集合，我们称为瞬时向量。 通过集合运算，可以在两个瞬时向量与瞬时向量之间进行相应的集合操作。目前，Prometheus支持以下集合运算符：
- and (并且)
- or (或者)
- unless (排除)

#### 1. `vector1 and vector2` 会产生一个**由vector1的元素组成**的新的向量。该向量包含vector1中完全匹配vector2中的元素组成。
   假设有两个瞬时向量：
```shell
A = {foo="bar", baz="qux"} 2
B = {foo="bar", baz="quxx"} 3
```
那么查询表达式`A and B`查询结果为空，因为标签不完全匹配，如果修改B的标签值为`baz=qux`，那么可查询到结果如下：
```shell
{foo="bar", baz="qux"} 2
```
请注意，这里返回结果的值是A的数据，而不是B。

#### 2. `vector1 or vector2` 会产生一个新的向量，该向量包含vector1中所有的样本数据，以及vector2中没有与vector1匹配到的样本数据。
>简言之，返回结果是两者的并集。
假设有两个瞬时向量：
```shell
A = {foo="bar", baz="qux"} 2
B = {foo="bar", baz="quxx"} 3
```
那么`A or B`可查询到结果如下：
```shell
{foo="bar", baz="qux"} 2
{foo="bar", baz="quxx"} 3
```

#### 3. `vector1 unless vector2` 会产生一个新的向量，新向量中的元素由vector1中没有与vector2匹配的元素组成。
假设有两个瞬时向量：
```shell
A ：{foo="bar", baz="louis"} 2
    {foo="bar", baz="quxx"}  3
    
B = {foo="bar", baz="quxx"} 3
```
那么`A unless B`可查询到结果如下：
```shell
{foo="bar", baz="louis"} 2 # A的数据
```

### 2.4 操作符优先级
从高到底：
```shell
^
*, /, %
+, -
==, !=, <=, <, >=, >
and, unless
or
```

### 2.5 匹配模式详解
向量与向量之间进行运算操作时会基于默认的匹配规则：依次找到与左边向量元素匹配（标签完全一致）的右边向量元素进行运算，如果没找到匹配元素，则直接丢弃。
接下来将介绍在PromQL中有两种典型的匹配模式：一对一（one-to-one）,多对一（many-to-one）或一对多（one-to-many）。

#### 1. 一对一匹配
一对一匹配模式会从操作符两边表达式获取的瞬时向量依次比较并找到唯一匹配(标签完全一致)的样本值。默认情况下，使用表达式：
```shell
vector1 <operator> vector2
```
在操作符两边表达式标签不一致的情况下，可以使用on(label list)或者ignoring(label list）来修改便签的匹配行为。使用ignoreing可以在匹配时忽略某些便签。而on则用于将匹配行为限定在某些便签之内。
```shell
<vector expr> <bin-op> ignoring(<label list>) <vector expr>
<vector expr> <bin-op> on(<label list>) <vector expr>
```
例如当存在样本：
```shell
method_code:http_errors:rate5m{method="get", code="500"}  24
method_code:http_errors:rate5m{method="get", code="404"}  30
method_code:http_errors:rate5m{method="put", code="501"}  3
method_code:http_errors:rate5m{method="post", code="500"} 6
method_code:http_errors:rate5m{method="post", code="404"} 21

method:http_requests:rate5m{method="get"}  600
method:http_requests:rate5m{method="del"}  34
method:http_requests:rate5m{method="post"} 120
```
使用PromQL表达式：
```shell
method_code:http_errors:rate5m{code="500"} / ignoring(code) method:http_requests:rate5m
```
该表达式会返回在过去5分钟内，HTTP请求状态码为500的在所有请求中的比例。如果没有使用ignoring(code)，操作符两边表达式返回的瞬时向量中将找不到任何一个标签完全相同的匹配项。
因此结果如下：
```shell
{method="get"}  0.04            //  24 / 600
{method="post"} 0.05            //   6 / 120
```
同时由于method为put和del的样本找不到匹配项，因此不会出现在结果当中。

#### 2. 多对一和一对多
在一对一模式中，总是以左边向量为基准，在右边向量寻找标签完全匹配（除非除了ignore或on）的数据，并且完全是一对一产生结果。
>一对一指的是，当找到两边完全匹配的数据时，将从两边同时抽出这条数据加入到结果列表中。（不能再用这条数据继续匹配）
