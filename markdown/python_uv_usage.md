## uv使用

常用命令：

```shell
uv python list # 列出所有python版本，含本地。加上 --all-versions 显示所有可下载版本，--only-installed 显示已安装版本
uv python install # 安装最新python（若已安装，则需要指定版本才会继续）

uv python install 3.10  # 安装指定版本（修订号使用最新）
uv python uninstall 3.10 3.11 # 卸载指定版本（所有修订号）
```

初始化：

```shell
uv init my_project # 也可以进入已存在的项目下执行，此时不加参数

uv venv --python 3.10 # 创建虚拟环境，依赖库默认安装在虚拟环境下，而不是全局缓存目录

# 激活虚拟环境
source .venv/bin/activate  # Linux/macOS
.venv\Scripts\activate     # Windows

# 退出虚拟环境
source .venv/bin/deactivate # Linux/macOS
.venv\Scripts\deactivate # Windows

# 不常见的操作
uv venv p313 --python 3.13 # 指定虚拟环境目录名称，默认为 .venv; 删除目录即可删除环境
```

启动项目：

```shell
uv run main.py # --python 3.10 指定版本（首次执行时，会自动创建venv）
uv python pin 3.10 # 修改为 3.10，谨慎！会更新 .python-version 文件，并在运行代码时重建venv

# 不常见
uv init --script cow3.py --python 3.13  # 为单个文件创建临时环境
uv add --script cow3.py cowsay rich # 为这个文件下载依赖库
```

按照以下顺序查找可以使用的py版本：

- 目前資料夾下的 `.python-version` 檔內設定的版本。
- 目前啟用的虛擬環境。
- 目前資料夾下的 `.venv` 資料夾內設定的虛擬環境。
- uv 自己安裝的 Python。
- 系統環境變數設定的 Python 環境。

操作本地依赖库：

```shell
uv add cowsay
uv remove cowsay

uv cache dir #  存放依赖库的目录
uv cache clean # 清空缓存，谨慎！

uv tree # 列出依赖库关系
uv lock --upgrade-package cowsay # 更新依赖库版本
uv lock --upgrade # 更新所有依赖库版本，谨慎！

# 如果手动修改了 [pyproject.toml]，需要更新lock文件
uv lock
uv pip list # 查看目前安装的依赖库
uv sync # 同步依赖库（让环境与lock文件一致）
```

使用依赖库提供的指令（类似npx）：
```shell
# 会自动下载依赖库
uvx cowsay -t 'hello, uv'

# 如果工具名称不同于依赖库名称，则需要指定 --from 参数
uvx --from httpie http -p=b GET https://flagtech.github.io/flag.txt
```

將依赖库提供的指令安裝到系統上：
```shell
#如果某個套件中的指令很使用，你也可以使用 uv tool install 把它安裝到系統上
# 例如剛剛 httpie 套件的 http 指令，就很適合安裝起來替代 curl 使用
uv tool install httpie

# uv 會建立一個獨立的虛擬環境來安裝套件
http -p=b GET https://flagtech.github.io/flag.txt

# 如果想知道實際安裝的路徑，可以透過 uv tool dir 指令：
uv tool dir

# 更新這些指令的版本
uv tool upgrade httpie

# 卸載工具
uv tool uninstall httpie
```