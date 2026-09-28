# AgentGo 后端

Go + Eino 的数据库规范检查服务。任务、Worker 进度、分片检查点、Finding 和报告快照保存在独立 SQLite 任务库中；目标数据库密码只保留在运行内存。

## 运行

```bash
go run ./cmd/server
```

默认配置位于 [config/config.yaml](config/config.yaml)，其中统一建模：

- `server`：监听地址与 HTTP 超时
- `security`：CORS 来源和请求体上限
- `task_store`：SQLite 路径与任务过期时间
- `scanner`：数据库连接超时和 Worker 并发配置
- `execution`：受控执行开关与目标主机白名单
- `llm`：必填的 Eino ChatModel 端点、模型、密钥环境变量名与索引建议阈值

可通过环境变量指定配置文件位置：

```bash
AGENTGO_CONFIG=./config/config.yaml go run ./cmd/server
```

前端开发服务器默认使用 `http://localhost:8080/api/v1`。

## 当前能力

- MySQL URL 或结构化连接创建只读检查任务。
- Eino Chain 编排 Schema 采集、并行规则检查与报告生成。
- SQLite（WAL）实时保存任务、Worker 和表分片进度及每个 Worker 的输入/输出/总 Token，用于任务总用量汇总，支持暂停、继续、取消。
- 任务/报告访问令牌只在任务创建响应中返回一次；后续请求使用 `X-Task-Access-Token`。
- 支持报告 JSON、HTML、Markdown 预览/下载，以及 DML 建议下载。

服务重启后，任务进度可恢复；但出于安全原因，继续任务时需重新提交同一目标数据库的连接信息。

## Worker 规则、Prompt 与必填 LLM

服务启动时会完整自检 YAML：网络与安全参数、任务库、扫描器、受控执行白名单、`rule_and_prompt.yaml` 以及 LLM 配置。LLM 是必填依赖：`llm.enabled`、`index_advice.enabled`、`report_summary.enabled` 必须为 `true`，且 `endpoint`、`model`、`api_key_env` 对应的环境变量均须有效；随后会在 `llm.startup_check_timeout_seconds` 时限内通过 Eino ChatModel 发送最小健康检查请求。任一配置、模型初始化或可达性检查失败，服务均拒绝启动，不会降级为确定性报告。

运行时若大表索引 Worker 或报告 Worker 的模型调用失败，对应 Worker 会记录安全的失败摘要，任务不会伪装为已完成；请恢复模型服务后重新执行任务。

先在 [config/config.yaml](config/config.yaml) 填入实际模型端点与模型名，并设置密钥：

```bash
export AGENTGO_LLM_API_KEY=your-key
go run ./cmd/server
```

[config/rule_and_prompt.yaml](config/rule_and_prompt.yaml) 是与 `config.yaml` 同目录自动加载的唯一 Worker 规则集和 AI Rule Prompt 配置源，无需也不应在主配置中重复声明它。Worker 只声明调度与 `rules`；每条 AI 指令必须写在 `ai_rules.<rule_id>.prompt`：

- 每个 Worker 都显式声明规则 ID；启动时会校验内置规则是否完整、重复或错配。
- `索引与性能检查(AI参与)` 的 Prompt 只接收表名、行数估算、列类型、现有索引与外键数。
- `语义与治理检查(AI参与)` 仅审查静态元数据中的枚举注释、审计字段、业务字段注释与明显语义类型冲突，所有结论均要求人工复核。
- `DDL 建议生成(AI参与)` 审查索引、主键、外键等确定性 DDL 变更候选的风险、顺序与回滚关注点；模型不会生成或改写可执行 SQL。
- `报告生成(AI参与)` 根据已持久化的 Worker 进度、Finding、分级统计和未检查对象生成阅读友好的最终 Markdown 文章。

规则的等级、DML 模板和安全边界仍由后端确定性代码固定，YAML 不允许把未知规则或任意 SQL 注入工作流。

## P1 / P2 可选能力

PostgreSQL 连接可使用 `postgres://用户名:密码@主机:5432/数据库`，或在前端逐项选择 PostgreSQL。

`索引与性能检查(AI参与)` 会按 `min_estimated_rows` 阈值，向模型发送**表名、估算行数、列类型、现有索引和外键数量**，生成静态索引审查建议。它不发送业务数据、不读取慢查询，也不假设查询需求；所有建议均标记为人工复核项。

受控 DML 执行默认关闭。启用时必须同时设置允许的目标主机；前端还需要提交连接 URL 并输入 `EXECUTE` 二次确认：

```bash
AGENTGO_EXECUTION_ENABLED=true
AGENTGO_EXECUTION_ALLOWED_HOSTS=db.example,127.0.0.1
```

该能力只允许单条 `UPDATE`、`INSERT` 或 `DELETE`，并记录任务审计事件；不要将生产数据库加入白名单，除非已有变更审批、备份和演练流程。
