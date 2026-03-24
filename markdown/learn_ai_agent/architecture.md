## Agent 架构

这是 AI Agent 工程的真正核心。

### 1. ReAct Agent

最经典架构：Thought -> Action -> Observation -> Thought -> ...

流程：

```plain
用户问题
   ↓
LLM 推理 Thought
   ↓
调用 Tool / API
   ↓
返回 Observation
   ↓
LLM 再推理
   ↓
循环直到得到答案
```

优点：

- 实现简单
- 灵活
- 支持工具调用

缺点：

- 长任务能力弱
- 容易陷入循环
- 规划能力弱

### 2. Plan-and-Execute 架构

流程：

```plain
User Query
   ↓
Planner (LLM)
   ↓
Task List
   ↓
Executor
   ↓
Tools / APIs
   ↓
Result
```

优点：

- 适合复杂任务
- 可控性强
- 成本可控

缺点：

- Planner 质量决定一切
- 动态调整能力弱

### 3. Multi-Agent 架构

**核心思想**：把复杂任务拆成多个专业 Agent

架构：

```plain
             Orchestrator
                   │
    ┌──────────────┼───────────────┐
 Researcher     Coder           Writer
    │              │               │
  Tools         Tools           Tools
```

## 优点

- 可扩展
- 专业分工
- 复杂任务能力强

## 缺点

- 系统复杂
- 调度难
- 成本高

### 4. Memory Agent

普通 Agent 是 无状态 的。对于Memory Agent：

- 短期记忆  (context window)
- 长期记忆  (vector DB)

架构：

```plain
User Query
   ↓
Memory Retrieve
   ↓
LLM Reasoning
   ↓
Action
   ↓
Memory Update
```

实现：

```plain
Vector memory
Key-value memory
Graph memory
```