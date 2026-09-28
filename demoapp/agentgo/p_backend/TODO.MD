# 后端 TODO · DB 开发规范检查服务

## 目标与边界

在 `p_backend` 使用 **Go + Eino** 提供 Web API、任务编排、MySQL 只读元数据采集、规则检测、报告与 SQL 建议生成能力，服务于免登录前端。Eino
负责 Agent、Tool 和工作流编排；确定性规则引擎仍负责最终问题判定与等级。

- 不实现 CLI；所有能力经 HTTP API 提供。
- 首期只支持 MySQL；连接必须使用只读账号，且不读取业务数据。
- 无登录不等于无保护：服务应只部署在本地或受控内网，使用 HTTPS、访问来源限制、限流、短生命周期任务令牌和不可预测任务 ID。
- 连接信息可由数据库 URL 或结构化字段提交，但请求中只能使用一种方式。解析后的密码、完整 URL
  和规范化后的完整连接串仅用于创建任务和连接数据库：不得写入任务存储、日志、错误、报告、SQL 文件或任何队列消息；任务结束、取消或失败后立即释放。
- 默认禁止执行变更 SQL；只生成经过风险分级的建议。

## Go + Eino 技术路线

- 使用 Go 实现 API、任务状态机、数据库方言 Adapter、规则引擎、报告存储和安全控制。
- 使用 Eino 将各能力封装为可测试的 Tool / Component，并由 Workflow 编排：`SchemaCollector` → `TableSplitter` →
  `RuleWorkers` → `FindingNormalizer` → `Evaluator` → `ReportGenerator` / `DMLGenerator`。
- `SchemaCollector` Tool 仅使用只读数据库连接；Tool 输入和输出均使用脱敏后的结构化模型，禁止向 Agent 上下文传递密码或业务数据。
- 规则 Worker 可以按命名、索引性能、字段类型与约束等维度并行执行。Eino 负责任务流转、并行聚合、超时、重试和取消；等级计算必须由确定性规则输出，不能由
  LLM 自主决定。
- 如使用 Eino 的 ChatModel 生成解释或修复说明，必须是显式可关闭的辅助能力；传入内容仅限脱敏 Schema 与 Finding，最终 DML/SQL
  必须经过模板化和风险校验。

```text
HTTP Handler
    |
    v
Eino Workflow
    ├── SchemaCollector Tool (MySQL 只读元数据)
    ├── TableSplitter Component
    ├── RuleWorkers (命名 / 索引 / 类型约束，可并行)
    ├── FindingNormalizer Component
    ├── Evaluator Component (确定性分级)
    └── ReportGenerator + DMLGenerator Tools
```

## 核心模型

| 模型                 | 关键字段                                                                                        |
|--------------------|---------------------------------------------------------------------------------------------|
| `Task`             | `id`、`status`、`progress`、`stage`、`createdAt`、`expiresAt`、暂停请求/暂停时间、访问令牌哈希、脱敏后的目标信息、失败摘要、版本号 |
| `ConnectionConfig` | 方言、主机、端口、数据库、用户名、密码、连接参数；只在请求处理和任务运行内存中存在                                                   |
| `TaskProgress`     | `stage`、整体处理总量/完成量、百分比、当前对象、`workerFlow`（顺序/依赖/并行组）、`workers`、最近更新时间                        |
| `WorkerProgress`   | `workerId`、名称、状态、处理总量/完成量、当前对象、发现问题数、输入/输出/总 Token、检查点、失败摘要、开始/更新时间/完成时间                    |
| `WorkUnit`         | `taskId`、`workerId`、分片/表标识、状态、检查点、租约、尝试次数、完成时间；用于精确恢复未完成工作                                  |
| `SchemaSnapshot`   | 方言、库、表、列、索引、主外键、注释、字符集/排序规则、行数估算                                                            |
| `Rule`             | `id`、方言、等级、检测条件、证据模板、建议模板、SQL 策略、规则版本                                                       |
| `Finding`          | `id`、`ruleId`、`severity`、对象定位、证据、影响、建议、置信度、人工复核标记、DML 修复方案、执行前校验 SQL、回滚说明、关联 SQL            |
| `Report`           | 任务信息、规则版本、分级统计、Finding 列表、未检查对象、失败原因、生成时间、脱敏后的 Markdown/HTML 报告快照                           |

### 任务状态

`pending` → `running` → `completed`

手动暂停：`running` → `pausing` → `paused`；继续：`paused` → `running`。运行中若部分表/规则失败：`running` →`partial_success`
；不可恢复错误：`running` → `failed`；用户终止：`pending/running/pausing/paused` → `cancelled`。

### 任务、阶段与 Worker 进度

任务详情接口必须让前端同时获得三个层次的可视状态：

| 层次     | 字段/枚举                | 说明                                                                                                                                                  |
|--------|----------------------|-----------------------------------------------------------------------------------------------------------------------------------------------------|
| 任务     | `status`             | 整个任务的生命周期状态：`pending`、`running`、`pausing`、`paused`、`partial_success`、`failed`、`completed`、`cancelled`                                               |
| 阶段     | `progress.stage`     | 当前主要阶段：`validating`、`collecting_schema`、`splitting`、`analyzing`、`generating_dml`、`generating_report`、`completed`                                    |
| Worker | `progress.workers[]` | 各 Eino Tool / Rule Worker 的独立状态与处理进度；`status` 为 `pending`、`running`、`pausing`、`paused`、`completed`、`partial_success`、`failed`、`skipped`、`cancelled` |

`WorkerProgress` 最小响应示例：

```json
{
  "workerId": "index-performance",
  "name": "索引与性能检查",
  "status": "running",
  "totalUnits": 50,
  "completedUnits": 31,
  "currentObject": "orders",
  "findingCount": 6,
  "errorSummary": null,
  "startedAt": "2026-09-02T10:00:00Z",
  "updatedAt": "2026-09-02T10:00:21Z",
  "finishedAt": null
}
```

首期至少上报 `schema-collector`、`naming-convention`、`index-performance`、`type-constraint`、`finding-normalizer`、
`dml-generator`、`report-generator` 七类 Worker。Worker 失败不能阻塞其他 Worker；任务进入 `partial_success` 时，必须保留失败
Worker 的 `errorSummary` 和已完成的结果。对象名遵从报告脱敏选项，错误摘要不得包含密码、完整连接串或驱动原始错误。

任务详情还应返回固定的 `workerFlow[]`：每个节点包含 `workerId`、名称、阶段、`order`、`dependsOn`、`parallelGroup` 与
`aiInvolved`。首期执行图为 `Schema 采集` → 并行的 `命名规范 / 索引与性能(AI参与) / 类型与约束` → `问题聚合` →
`DML 建议生成` → `报告生成(AI参与)`；前端以它作为流程图布局依据，实时状态仍以 `workers[]` 为准。

### 实时持久化、暂停与继续语义

首期使用独立于被检查数据库的本地 SQLite 任务库（WAL 模式）保存任务运行态；通过 `TaskRepository`
隔离持久化实现，后续可替换为服务数据库。严禁将目标数据库的 URL、密码或业务数据写入任务库。

- 任务状态、Worker 状态、`WorkUnit` 检查点、Finding、报告快照片段和 DML 结果在每次状态迁移及每个原子单元完成后立即写入；完成一个单元时，
  `WorkUnit=completed`、`WorkerProgress` 与 Finding 必须在同一事务提交。
- `POST /pause` 只接受 `running` 任务；服务端先持久化暂停请求并标记为 `pausing`，不强行中断正在处理的数据库操作。
- 每个 Eino Workflow 节点在处理单表、单个规则分片或单个报告阶段前后读取持久化暂停标记；完成当前原子单元后写入
  `WorkerProgress.checkpoint`，更新为 `paused`，释放数据库连接和 Worker 执行资源。
- 所有活跃 Worker 均抵达安全检查点后，任务才以事务方式转为 `paused`。Worker 已完成的 Finding、报告快照片段与 DML 结果必须保留。
- `POST /resume` 只接受 `paused` 任务，先持久化状态变化，再按 `WorkUnit` 检查点将状态改回 `running`
  ；仅领取未完成或因重启回收的分片，不重新采集或分析已确认完成的表。若运行内存中不存在连接信息，则响应
  `connectionRequired=true`，调用方须在同一次继续请求中重新提交 URL 或结构化连接。
- 每个执行中的 `WorkUnit` 需要短期租约。服务启动时回收过期租约，将未提交完成的单元置回待处理，并把原 `running`/`pausing`
  任务转为 `paused`、标记“服务重启待恢复”。目标库密码和完整 URL 不落库，故重启后必须由用户重新提供连接信息才可恢复执行。

## 报告分级与 SQL 策略

后端固定返回下列枚举，前端映射为中文标签。规则必须显式声明等级，禁止由模型自由推断。

| 枚举           | 页面名称 | 判定标准                | SQL 策略                                                            |
|--------------|------|---------------------|-------------------------------------------------------------------|
| `urgent`     | 紧急   | 数据正确性、可用性或重大性能存在高风险 | **必须提供 DML 修复方案**；包含执行前校验、影响行数预估、分批/事务策略和回滚说明，仅供人工审核、复制或下载，禁止自动执行 |
| `normal`     | 一般   | 影响维护性、一致性或中期性能      | 可生成标注风险、前置条件和回滚说明的 SQL 建议                                         |
| `suggestion` | 建议   | 优化或规范改进，不影响当前核心正确性  | 优先给出说明；仅低风险场景生成幂等 SQL 建议                                          |

每条 Finding 必须带 `severityReason`、`evidence`、`recommendation`、`requiresManualReview`、`sqlRisk`、`dmlSuggestion`、
`precheckSql` 和 `rollbackPlan`。所有紧急规则都必须提供 `dmlSuggestion`：对可直接修复数据的问题生成 `UPDATE` / `INSERT` /
`DELETE` 的受控脚本；对只能以 DDL 修复的问题生成数据核验、迁移准备或回填 DML，并将 DDL 明确标为独立的人工操作。DML
生成器需识别破坏性操作、影响行数、锁表风险、字段缩短和索引重建，命中时禁止标记为低风险。

## API 草案

| 方法     | 路径                                              | 用途                                                                         |
|--------|-------------------------------------------------|----------------------------------------------------------------------------|
| `POST` | `/api/v1/tasks`                                 | 提交连接信息与检查配置，创建任务；返回 `taskId`、`taskAccessToken`、`expiresAt`；密码只在本请求和运行内存中使用 |
| `GET`  | `/api/v1/tasks/{taskId}`                        | 查询任务状态、当前阶段、整体进度、各 Worker 进度/失败摘要与报告是否就绪；需要访问令牌                            |
| `POST` | `/api/v1/tasks/{taskId}/pause`                  | 请求在当前安全检查点暂停任务；需要访问令牌                                                      |
| `POST` | `/api/v1/tasks/{taskId}/resume`                 | 从保存的检查点继续已暂停任务；需要访问令牌                                                      |
| `POST` | `/api/v1/tasks/{taskId}/cancel`                 | 取消尚未完成的任务；需要访问令牌                                                           |
| `GET`  | `/api/v1/reports/{taskId}`                      | 分页获取报告与 Finding，支持等级、规则、对象、关键字筛选；需要访问令牌                                    |
| `GET`  | `/api/v1/reports/{taskId}/findings/{findingId}` | 获取单条 Finding 详情、紧急 DML 修复方案与 SQL 建议；需要访问令牌                                 |
| `GET`  | `/api/v1/reports/{taskId}/preview`              | 获取固定、脱敏后的完整报告；`format` 参数可取 `html` 或 `markdown`，用于浏览器在线预览；需要访问令牌           |
| `GET`  | `/api/v1/reports/{taskId}/sql`                  | 下载经过风险标记的 DML/SQL 建议文件；需要访问令牌                                              |
| `GET`  | `/api/v1/reports/{taskId}/export`               | 下载脱敏报告快照；`format` 参数可取 `json`、`html` 或 `markdown`；需要访问令牌                   |

报告预览与下载均从同一份已完成的报告快照读取，保证内容一致；Markdown 由结构化 `Report` 模板生成，HTML 必须由受控模板或净化后的
Markdown 渲染生成，禁止接收或回显未净化 HTML。预览接口使用 `Content-Disposition: inline`，并配置限制脚本、外部资源和内联事件的安全响应头。

除 `POST /api/v1/tasks` 外，所有任务和报告接口均要求请求头 `X-Task-Access-Token`
。创建任务时生成不可预测的短期令牌，仅返回一次原文，服务端只保存令牌哈希；任务过期、令牌缺失或校验失败时返回统一的安全响应，避免透露任务是否存在。默认过期时间通过配置控制，并由创建响应的
`expiresAt` 明确告知前端。

所有接口的错误响应使用统一错误码和用户安全提示；禁止回传完整连接串、密码、数据库驱动原始错误或内部堆栈。任务 ID 使用不可预测
UUID/随机值，并设置过期时间和访问频率限制。

`POST /api/v1/tasks` 的请求体支持以下两种互斥方式，服务端先解析为内存中的 `ConnectionConfig` 再建立连接：

```json
{
  "connectionUrl": "mysql://readonly:password@db.example:3306/app_db"
}
```

```json
{
  "connection": {
    "dialect": "mysql",
    "host": "db.example",
    "port": 3306,
    "database": "app_db",
    "username": "readonly",
    "password": "password"
  }
}
```

校验规则：请求必须且只能提供 `connectionUrl` 或 `connection`；URL 模式仅接受已支持方言的 URL
格式，并校验协议、主机、端口、数据库名和允许的连接参数。解析失败或两种方式同时提交时返回安全的参数错误，不回显原始 URL。

## 后端开发 TODO

### P0：检查服务闭环

- [x] 初始化 Go 后端工程、依赖管理和 Eino 运行配置；约定 Eino Component 的输入/输出模型、上下文取消和错误包装规范。
- [x] 初始化 Web 服务、配置管理、路由、中间件、统一响应与错误码；实现健康检查。
- [x] 实现 `TaskRepository` 与 SQLite（WAL）迁移：持久化 `Task`、`WorkerProgress`、`WorkUnit`、Finding、报告快照和 DML
  结果；数据库路径、任务过期和清理策略可配置，且数据库中不保存任何目标库凭据。
- [x] 实现匿名任务访问控制：创建任务时生成一次性返回的 `taskAccessToken`、服务端仅保存哈希、按配置设置 `expiresAt`
  ；任务/报告接口校验 `X-Task-Access-Token`，过期或无效时返回统一安全响应。
- [x] 实现任务暂停/继续状态机与 `/pause`、`/resume` 接口：暂停请求、状态变化与检查点均实时持久化；只允许合法状态迁移，重复请求返回幂等结果或明确状态错误，暂停/继续/取消均需访问令牌。
- [x] 实现 `POST /tasks` 请求校验与连接信息处理：支持 `connectionUrl` 和结构化 `connection` 两种互斥输入，统一解析为内存中的
  `ConnectionConfig`；实现最小权限提示、密码内存隔离、完整 URL 脱敏、SSRF/内网地址策略和连接超时。
- [x] 定义并实现 `Task`、`TaskProgress`、`WorkerProgress`、`SchemaSnapshot`、`Rule`、`Finding`、`Report` 模型及持久化/缓存边界；敏感凭据不落库。
- [x] 将 MySQL 只读采集器实现为 Eino `SchemaCollector` Tool：采集表、列、索引、主外键、注释、字符集/排序规则、行数估算；处理空库、权限不足、超时和单表失败。
- [x] 建立规则引擎：按方言和启用状态加载规则，输出结构化 Finding，记录规则版本和检测证据；索引 Worker 在大表场景加入 LLM
  静态建议，并新增 `语义与治理检查(AI参与)` 对敏感字段、关联、状态、审计和注释候选给出静态人工复核建议。
- [x] 实现首批基础规则，并为每条规则明确三级等级与理由：无主键、缺注释、命名不规范、重复索引、外键列无索引、金额浮点、时间字段、NULL
  默认值、字符集不一致、超宽索引。
- [x] 实现 Eino Workflow：由 `TableSplitter` 按表创建持久化 `WorkUnit`，调度多个 `RuleWorker` 并行检查，在
  `FindingNormalizer` 聚合；每个节点持续持久化 `WorkerProgress`，并返回有序 `workerFlow`（依赖、并行组、AI
  标记）供前端流程图使用；在单表/规则分片/报告阶段边界以事务提交检查点、Finding 和完成状态，支持租约回收、并发上限、超时、重试、暂停、继续、取消、整体/Worker
  进度聚合，单表失败转换为 `partial_success`。
- [x] 实现报告聚合：三级统计、规则/对象/关键字筛选、分页、未检查对象和失败原因。
- [x] 实现报告快照生成器：`报告生成(AI参与)` 使用 Eino ChatModel 基于持久化的 Worker 结果、Finding 与分级统计生成阅读友好的脱敏
  Markdown；将快照与规则版本一同保存，保证在线预览和下载内容一致。
- [x] 实现 DML/SQL 建议生成器：紧急 Finding 必须输出 DML 修复方案；统一生成方言模板、执行前校验、影响行数预估、分批/事务策略、风险级别、回滚说明和语法预检；严格禁止自动执行。
- [x] 实现任务状态、阶段、整体/Worker 进度、报告、报告在线预览、Finding 详情、暂停、继续、取消、SQL 下载与导出 API；覆盖`pausing`、
  `paused`、`cancelled` 状态。API 文档按要求不生成。
- [x] 为预览接口补充快照一致性、Markdown/HTML 安全渲染、过期任务和脱敏回归测试。
- [x] 为 Eino Tool、Workflow 分支、上下文取消和 Agent 辅助开关补充单元测试；对同一 Schema 的确定性规则结果进行稳定性回归测试。
- [x] 覆盖实时进度持久化与暂停/继续测试：每个 Worker/分片进度实时落库、暂停请求落在单表执行中、多个 Worker
  到达检查点、连接释放、继续不重复处理、暂停后取消、非法状态迁移、服务重启后的租约回收和从检查点恢复。
- [x] 编写单元和 MySQL 集成测试，重点覆盖凭据脱敏、采集映射、分级、规则证据、部分失败和 SQL 风险拦截。

**P0 验收标准**：免登录前端能够通过 API 完成一次 MySQL 只读检查；创建响应返回短期访问令牌，而其余接口缺少/使用无效令牌时均无法读取任务；任务、Worker
和分片进度实时持久化，服务重启后重新提交连接信息即可从最后一个已提交检查点继续；运行中任务可在安全检查点变为 `paused`
，数据库连接被释放，继续后只处理未完成分片；报告中准确汇总紧急/一般/建议；每条紧急 Finding 都返回 DML
修复方案、执行前校验与回滚说明；任务取消后状态为 `cancelled`；任一表失败仍可输出部分成功报告；密码、完整连接串与令牌原文不出现在持久层、日志、响应、报告和导出文件中。

### P1：稳定性与集成

- [ ] 增加按 IP/任务的限流、请求大小限制、CORS 白名单、安全响应头和审计事件。
- [ ] 支持规则配置、项目规则覆盖、对象白名单/忽略列表、规则版本追踪与已知问题标记。
- [x] 增加任务日志的分级、指标、trace ID、运行耗时与规则命中统计；日志只保留脱敏对象。
- [ ] 增加 JSON/HTML/SARIF 导出、报告基线比较和“仅新增问题”查询。
- [ ] 为大 Schema 增加背压、资源配额、连接池、任务取消回收和压测基准。
- [x] 实现 PostgreSQL adapter，并验证同一规则模型下的方言差异。

### P2：治理与受控执行

- [x] 使用 YAML 统一管理 Worker 规则集与 Prompt，并在严格脱敏前提下接入必填的 Eino LLM：索引建议和报告文章由 LLM
  辅助生成；规则判定和等级仍以确定性规则引擎为准。
- [ ] 实现修复方案合并、SQL 冲突检测、执行顺序、回滚 SQL 与人工复核理由。
- [x] 若业务允许，再设计受控执行：环境白名单、审批/变更单、备份检查、双重确认、完整审计；默认关闭。
- [ ] 建立规则误报反馈、规则效果趋势与插件化方言/规则包。

### P3

- [ ] 优化 AI-SEMANTIC-001
  - 
  
```1. 敏感字段标记：`tg_id`（Telegram用户ID）、`out_trade_no`（资金流水号）及`reward_amount`/`reward_currency`（财务金额）可能涉及隐私与资金安全，
需人工复核是否按敏感字段进行脱敏或访问控制治理。 2. 疑似关联字段：`player_id`、`rule_id`、`activity_code`、`source_id`均无外键约束，
需人工复核其与用户表、规则表、活动配置表及外部系统的逻辑关联是否应在Schema层面约束。 3. 状态/枚举字段：`claim_status`使用varchar(16)存储状态枚举，
需人工复核是否应改为ENUM/TINYINT或增加CHECK约束以提升数据一致性。 4. 注释质量：`id`注释仅为“ID”，过于简略，需人工复核是否补充为“活动奖励领取记录ID”；
`reward_currency`注释列出了枚举值但类型为varchar，需人工复核是否需增加约束。 5. 索引冗余：存在多组重复索引（如针对`activity_code`、`claim_status`、
`tg_id`的单列索引各有两组，且联合唯一索引包含`activity_code`导致前缀冗余），需人工复核清理。 （以上所有结论均为人工复核候选）```