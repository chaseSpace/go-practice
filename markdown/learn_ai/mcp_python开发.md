## Python开发MCP Server

### Init Project

```shell
uv init mcp-server-demo
cd mcp-server-demo

uv add "mcp[cli]" # 添加mcp cli工具

# uv run mcp [OPTIONS] COMMAND [ARGS]...

uv run mcp version # 查看mcp cli版本


# 初始化venv
uv venv
source .venv/bin/activate
```

### Coding

。。。

### Debug

#### 1. 使用 MCP Inspector 测试

```shell
uv run mcp dev main.py
```

#### 2. 使用其他MCP Client测试

配置参数：

```shell
# 1. 使用uv
uv --directory /project_path run main.py

# 2. 使用python
python -m /path/to/main.py
```

### Publish

#### 1. 再次初始化

之前的方式是将项目初始化为一个App目录结构，但如果要发布为一个Lib（二进制工具），就必须重新初始化目录结构。

```shell
cd mcp-server-demo

# 先将 pyproject.toml 重命名为 pyproject.toml1，否则不能执行下面的init
uv init . --lib

# 现在，会创建src/mcp-server-demo/__init__.py 结构

# 然后，手动将 pyproject.toml1 中的dependencies部分复制到新的pyproject.toml中。
```

然后打开 pyproject.toml，将其中所有的 `mcp-server-demo` 替换为 `mcp-server-demo-lei`，这是因为它会作为发布到pypi的包名，必须是唯一的，否则发布失败。
注意，还要重命名`src/mcp-server-demo` 为 `src/mcp-server-demo-lei`，否则也无法发布。

参考下面的详细说明：

```toml
[project]
name = "mcp-server-demo-lei"  # <------------- 修改
version = "0.1.0"
description = "Add your description here"
readme = "README.md"
authors = [
    { name = "chasespace", email = "random2035@qq.com" }
]
requires-python = ">=3.13"
dependencies = []

[project.scripts]
mcp-server-demo-lei = "mcp_server_demo_lei:main"  # <------------- 修改

[build-system]
requires = ["hatchling"]
build-backend = "hatchling.build"
```

然后将 `main.py` 移动到 `src/mcp-server-demo-lei/` 。然后使用新的命令运行项目：

```shell
uv run mcp-server-demo-lei
```

#### 2. 在 init 文件中添加启动入口

在 `src/mcp-server-demo/__init__.py` 中添加启动入口：

```python
from .main import main

__all__ = ["main"]
```

这是对应 `pyproject.toml` 中的 `[project.scripts]`部分。

#### 3. 准备pypi账号

python写的MCP Server则发布到[Pypi仓库](https://pypi.org)，需要提前注册账号。

为了等会免密上传包，需要在系统的`HOME`目录下提前创建 `.pypirc` 文件，内容如下：

```
[distutils]
index-servers =
    pypi

[pypi]
repository = https://upload.pypi.org/legacy/
username = __token__
password = ****   # token 从 https://pypi.org/manage/account/token/ 获取
```

#### 4. 最后一次测试

前面使用 MCP Inspector 测试是不够的，那只能作为单测来测试MCP接口。我们还需要模拟真实的MCP Client使用uvx来测试。

```shell
uv build
uv pip install -e . # 安装当前项目（作为命令工具）

# 命令测试（注意提前激活venv）
(mcp-server-demo) PS D:\Users\Desktop\Go\go-practice\pythonx\mcp-server-demo> mcp-server-demo-lei
以标准 I/O 方式运行 MCP 服务器...
```

如果这里能运行，说明MCP Server可以成功启动，你可以在代码中添加打印来验证，然后通过 Ctrl+C 退出完成测试。

#### 5. 发布

发布简要流程：

```
# 在项目根目录下执行

uv build # 构建，完成后会生成目录 dist，里面放着压缩包
uv publish dist/* # 需要输入密码
twine upload dist/* # 或者使用这个命令，配置 ～/.pypirc后不需要手动输入密码
```

最后，在MCP Client中使用时，使用uvx命令下载即可。

> 在MCP Client中使用钱，你还可以在命令行中执行`uvx mcp-server-demo-lei`

> `pyproject.toml`中的name作为你上传的包名，如果被pypi拒绝，则需要修改，然后删除dist目录下内容重新打包。