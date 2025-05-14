import logging
import time
from contextlib import asynccontextmanager

import uvicorn
from fastapi import FastAPI, Request, status, Depends, File, UploadFile
from fastapi import HTTPException
from fastapi import Header, Cookie
from fastapi.exceptions import RequestValidationError
from fastapi.middleware.cors import CORSMiddleware
from fastapi.middleware.gzip import GZipMiddleware
from fastapi.responses import JSONResponse

from fastapi_demo.middlewares.mw import ProcessTimeMiddleware
from fastapi_demo.routers import users, items
from fastapi_demo.utils.logging import configure_logging

configure_logging()


@asynccontextmanager
async def lifespan(app: FastAPI):
    """生命周期管理"""
    # 启动时执行
    logging.info("Application starting up...")
    # 应用启动逻辑...
    yield
    # 关闭时执行
    logging.info("Application shutting down...")
    # 应用关闭逻辑...


app = FastAPI(
    title="FastAPI Demo",
    description="A demo FastAPI project with best practices",
    version="0.1.0",
    lifespan=lifespan
)

# 添加中间件
app.add_middleware(ProcessTimeMiddleware, header_namespace="middleware")
app.add_middleware(GZipMiddleware, minimum_size=1000, compresslevel=5)  # 返回内容小于N字节时不压缩
app.add_middleware(
    CORSMiddleware,
    allow_origins=["*"],
    allow_credentials=True,
    allow_methods=["*"],
    allow_headers=["*"]
)

# 包含路由
app.include_router(users.router)
app.include_router(items.router, prefix="/items", tags=["items"])


# 自定义异常处理
@app.exception_handler(RequestValidationError)
async def validation_exception_handler(request: Request, exc: RequestValidationError):
    """处理请求验证错误"""
    logging.error(f"Validation error: {exc.errors()}")
    return JSONResponse(
        status_code=status.HTTP_422_UNPROCESSABLE_ENTITY,
        content={"detail": exc.errors(), "body": exc.body},
    )


@app.get("/", tags=["root"])
async def root():
    """根路由"""
    return {"message": "Welcome to FastAPI Demo"}


@app.get("/health", tags=["health"])
async def health_check():
    """健康检查端点"""
    return {"status": "ok"}


@app.get("/get_header_cookie/{code}")  # 通过类型注解的方式获取header和cookie字段
def read_item(code: int, user_agent: str = Header(None), session_token: str = Cookie(None)):
    time.sleep(1)
    if code == 404:
        raise HTTPException(status_code=404, detail="Not found*")
    elif code == 500:
        raise HTTPException(status_code=500, detail="Internal Server Error*")
    return {"User-Agent": user_agent, "Session-Token": session_token}


@app.get("/customize_resp")  # e.g. /customize_resp?item_id=1
def customize_resp(item_id: int):
    content = {"item_id": item_id}
    headers = {"X-Custom-Header": "custom-header-value"}
    return JSONResponse(status_code=200, content=content, headers=headers)


# 依赖项类（对比函数所拥有的的优势是让编辑器可以自动补全）
class CommonQueryParams:
    def __init__(self, q: str | None = None, skip: int = 0, limit: int = 100):
        self.q = q
        self.skip = skip
        self.limit = limit


# 路由操作函数
@app.get("/depend_example")  # Depends可以在函数执行前调用依赖项函数处理路由参数
async def read_items(commons: CommonQueryParams = Depends()):
    response = {}
    if commons.q:
        response.update({"q": commons.q})
    return response


# 示例：文件上传
@app.post("/files/")
async def create_file(file: UploadFile = File(...)):
    return {"filename": file.filename}


# 打印环境信息
# python -m uvicorn --version

if __name__ == "__main__":
    uvicorn.run(
        "main:app",
        port=8000,
        reload=True,
        reload_delay=3,
        reload_excludes=[".logs/*", "temp/*", ".git/*"],
        workers=1,
        log_config=None,  # 使用我们自己的日志配置
        timeout_keep_alive=60,
        limit_concurrency=100,  # 并发连接数的上限,默认为None不限制
        limit_max_requests=100,  # 每个工作进程的最大请求数,默认为None不限制
        timeout_graceful_shutdown=10  # 优雅重启的超时时间
    )
