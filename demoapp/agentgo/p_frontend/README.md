# AgentGo 前端

Vue 3 + Vite + TypeScript + Element Plus 的单页 Web 界面。

## 运行

```bash
pnpm install
pnpm dev
```

默认请求 `http://localhost:8080/api/v1`。若后端地址不同，创建 `.env.local`：

```bash
VITE_API_BASE_URL=http://localhost:18080/api/v1
```

Vite 开发服务器默认开启热更新（HMR）：修改 `src/` 下的 Vue、TypeScript 或样式文件会立即更新页面，无需手动刷新。局域网或反向代理开发时，如 WebSocket 主机与当前访问地址不同，可额外设置：

在 WSL2 的 Windows 挂载目录（例如 `/mnt/c/...`）中，开发服务器会自动启用文件轮询，避免文件系统事件丢失。也可显式指定：

```bash
VITE_USE_POLLING=true pnpm dev
```

```bash
VITE_HMR_HOST=你的开发机地址 VITE_HMR_CLIENT_PORT=5173 pnpm dev
```

## 当前页面

- 数据库 URL / 逐项连接输入
- 任务状态、阶段和 Worker 实时进度
- 暂停、继续、取消及服务重启后的连接信息重新提交
- 报告筛选、问题详情、DML 复制与下载
- HTML / Markdown 在线预览和报告下载
