# 常见py项目布局

```
my_project/
├── app/
│   ├── __init__.py
│   ├── main.py                # 入口文件，启动 FastAPI 应用
│   ├── config.py              # 配置文件，存储数据库配置、环境变量等
│   ├── models/                # 数据库模型定义
│   │   ├── __init__.py
│   │   ├── base.py            # 基础模型类
│   │   ├── user.py            # 用户模型
│   │   └── post.py            # 帖子模型
│   ├── schemas/               # Pydantic 数据模型定义
│   │   ├── __init__.py
│   │   ├── user.py            # 用户数据模型
│   │   └── post.py            # 帖子数据模型
│   ├── crud/                  # 数据访问层
│   │   ├── __init__.py
│   │   ├── base.py            # 基础 CRUD 操作
│   │   ├── user.py            # 用户 CRUD 操作
│   │   └── post.py            # 帖子 CRUD 操作
│   ├── routes/                # API 路由定义
│   │   ├── __init__.py
│   │   ├── user.py            # 用户相关路由
│   │   └── post.py            # 帖子相关路由
│   ├── services/              # 业务逻辑层
│   │   ├── __init__.py
│   │   ├── user.py            # 用户相关业务逻辑
│   │   └── post.py            # 帖子相关业务逻辑
│   ├── utils/                 # 工具函数和模块
│   │   ├── __init__.py
│   │   ├── auth.py            # 认证相关工具
│   │   └── logger.py          # 日志工具
│   ├── tests/                 # 测试代码
│   │   ├── __init__.py
│   │   ├── test_user.py        # 用户相关测试
│   │   └── test_post.py        # 帖子相关测试
│   └── static/                # 静态文件（如图片、CSS 等）
│       └── favicon.ico
├── .env                       # 环境变量文件
├── Dockerfile                 # Docker 部署文件
├── docker-compose.yml         # Docker Compose 配置文件
├── README.md                  # 项目说明文档
