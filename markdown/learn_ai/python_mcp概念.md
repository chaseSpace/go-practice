# MCP 核心概念完全指南

## 1. Server (服务器)

**作用**：所有功能的运行容器
代码示例：

```python
from mcp.server.fastmcp import FastMCP

mcp = FastMCP("我的服务")  # 创建实例
mcp.run()  # 启动服务
```

**场景**：

- 对接Claude时作为后台服务
- 开发AI插件时的主程序

## 2. Resource (资源)

**作用**：提供只读数据
代码示例：

```python
@mcp.resource("weather://{city}")
def get_weather(city: str) -> str:
    return f"{city}天气：晴，25℃"
```

**场景**：

- 提供数据库查询结果
- 返回静态配置文件
- 实时数据接口（股票/天气）

## 3. Tool (工具)

**作用**：执行具体操作
代码示例：

```python
@mcp.tool()
def translate(text: str, to_lang: str) -> str:
    return f"Translated({to_lang}): {text}"
```

**场景**：

- 调用外部API
- 执行计算任务
- 修改数据库记录

## 4. Prompt (提示)

**作用**：定义对话模板
代码示例：

```python
@mcp.prompt()
def customer_service(question: str) -> str:
    return f"""你是一名客服，请专业地回答：
    用户问题：{question}
    回答要求：用中文，不超过100字"""
```

**触发方式**：

```
@customer_service question="如何退款？"
→ AI回复："尊敬的客户，退款流程如下..."
```

**场景**：

- 标准化问答流程
- 多轮对话预设
- 内容生成模板

## 5. Image (图像)

**作用**：处理图片数据
代码示例：

```python
@mcp.tool()
def resize_image(img: Image) -> Image:
    img.data = compress(img.data)
    return img
```

**特点**：

- 自动处理二进制数据
- 支持常见图片格式
- 专为视觉任务设计

## 6. Context (上下文)

**作用**：管理会话状态
代码示例：

```python
@mcp.tool()
async def long_task(ctx: Context):
    ctx.info("任务开始")  # 记录日志
    await ctx.report_progress(1, 5)  # 更新进度
    data = await ctx.read_resource("data.csv")  # 读取资源
```

**包含功能**：

- 进度报告
- 临时存储
- 跨工具通信

## 7. 对比总结：

| 概念       | 数据流方向     | 是否持久化 | 典型用例      |
|----------|-----------|-------|-----------|
| Resource | 服务端 → 客户端 | 是     | 提供数据看板    |
| Tool     | 客户端 → 服务端 | 否     | 执行订单取消操作  |
| Prompt   | 双向        | 是     | 标准化客服对话流程 |
| Image    | 双向        | 否     | 图片缩略图生成   |
| Context  | 会话内       | 否     | 跟踪多步骤任务进度 |

## 8. 完整示例代码

```python
from PIL import Image as PILImage
from mcp.server.fastmcp import FastMCP, Image, Context

mcp = FastMCP("电商助手")


@mcp.resource("product://{id}")
def get_product(id: str):
    return f"商品{id}详情：..."


@mcp.tool()
def create_order(item: str, qty: int):
    return f"已创建{item}x{qty}的订单"


@mcp.prompt()
def sales_advice(style: str):
    return f"你是一名{style}风格的销售，请推荐商品..."


@mcp.tool()
def process_photo(img: Image) -> Image:
    pil_img = PILImage.open(img.data)
    return Image(data=pil_img.resize((800, 600)))


@mcp.tool()
async def export_data(ctx: Context):
    ctx.info("导出开始")
    await ctx.report_progress(0, 3)

```