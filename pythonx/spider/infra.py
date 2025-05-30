import time
from urllib.parse import urlparse

from aiohttp import ClientSession


# 定义一个静态类来封装aiohttp返回的异步response，便于传递
class HttpResponse:
    content: bytes
    status: int

    def __init__(self, content, status):
        self.content = content
        self.status = status

    def text(self):
        return self.content.decode('utf-8')


# 统一请求
async def send_request(hint, url, method, body=None):
    time.sleep(0.1)  # 限制请求频率

    st = time.time()
    ct, status = await core_req(url, method, body)
    et = time.time()
    print(f'{hint} - send_request - {et - st:.2f}s {method}: {url}')

    return HttpResponse(ct, status)


session: ClientSession


def init_session():
    global session
    session = ClientSession()


async def close_session():
    global session
    await session.close()


# 异步请求
async def core_req(url, method, body):
    if method == 'GET':
        async with session.get(url, params=body) as response:
            return await response.read(), response.status
    elif method == 'POST':
        async with session.post(url, params=body) as response:
            return await response.read(), response.status
    else:
        raise ValueError('Invalid method')


def remove_path_simple(url):
    parsed = urlparse(url)
    base_url = f"{parsed.scheme}://{parsed.netloc}"
    return base_url


# 保存不符合预期的html以供分析
def save_html(url, text):
    with open('exception.html', 'w', encoding='utf-8') as f:
        f.write(url + "\n\n" + text)
