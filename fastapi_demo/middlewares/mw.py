import time
from fastapi_demo.utils.time.time_format import seconds_to_readable
from fastapi import Request
from starlette.middleware.base import BaseHTTPMiddleware


class ProcessTimeMiddleware(BaseHTTPMiddleware):
    def __init__(self, app, header_namespace: str):
        super().__init__(app)
        self.header_namespace = header_namespace

    async def dispatch(self, request: Request, call_next):
        # 获取headers
        # he = dict(request.scope["headers"])
        start_time = time.perf_counter()

        response = await call_next(request)  # 继续请求

        process_time = time.perf_counter() - start_time
        response.headers["X-Process-Time"] = seconds_to_readable(process_time)
        return response
