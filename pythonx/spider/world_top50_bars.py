import asyncio
from urllib.parse import urljoin

from lxml import html

from infra import send_request, remove_path_simple, save_html, init_session, close_session
from model import WorldTopBars

count = 0


# 当前任务的运行入口
async def run():
    entry_urls = [
        # 'https://www.theworlds50best.com/bars/list/1-50',
        # 'https://www.theworlds50best.com/bars/northamerica/list/1-50',
        'https://www.theworlds50best.com/bars/asia/list/1-50'
    ]

    tasks = []
    for entry_url in entry_urls:
        tasks.append(asyncio.create_task(start_scrape(entry_url)))

    init_session()
    await asyncio.gather(*tasks)
    await close_session()


async def start_scrape(entry_url):
    r = await send_request('start_scrape', entry_url, 'GET')
    if r.status != 200:
        raise Exception('Failed to fetch entry url: ' + entry_url)

    tree = html.fromstring(r.text())
    hrefs = tree.xpath('//div[@data-list="1-50"]/a/@href')  # 查出所有二级页面地址
    if len(hrefs) == 0:
        print('No hrefs found, check pls!')
        return

    # 从入口地址中得到基础url
    base_url = remove_path_simple(entry_url)
    # 开启并发 max=3
    semaphore = asyncio.Semaphore(3)

    # 定义一个协程来包装parse_detail_page函数，增加Semaphore的 acquire 和 release 调用
    async def wrapped_parse_detail_page(href):
        async with semaphore:
            await parse_detail_page(base_url, href)

    # 创建并启动并发任务
    tasks = [wrapped_parse_detail_page(href) for href in hrefs]
    await asyncio.gather(*tasks)


# 解析详情页
async def parse_detail_page(base_url, url):
    url = urljoin(base_url, url)
    r = await send_request('parse_detail_page', url, 'GET')
    if r.status != 200:
        raise Exception('Failed to fetch entry url: ' + url)

    text = r.text()
    tree = html.fromstring(text)

    # bar_name = tree.xpath('//div[@class="content profile"]')
    bar_name = tree.xpath('//div[@class="content profile"]/h1/text()')
    if not bar_name:
        print(f'No bar name found, check pls!')
        save_html(url, text)
        exit(1)
    content_arr = tree.xpath('//div/p/text()')

    # print(content_arr, len(content_arr))
    rank_no = int(content_arr[1])
    city = content_arr[2]

    honor_desc = '\n'.join(tree.xpath('//ul[@compact="Accolades"]/li/text()'))
    intro = '\n'.join(tree.xpath('//div[@class="content profile"]/p[3]/text()'))
    location = ''.join(tree.xpath('//div[@class="details"]/p[1]/text()'))
    site_url = ''.join(tree.xpath('//div[@class="details"]/a[contains(@class, "website")]/@href'))
    phone = ''.join(tree.xpath('//div[@class="details"]/a[contains(@class, "telephone")]/@href'))
    fb_url = ''.join(tree.xpath('//div[@class="details"]/a[contains(@class, "facebook")]/@href'))
    ins_url = ''.join(tree.xpath('//div[@class="details"]/a[contains(@class, "instagram")]/@href'))
    # youtube_url = ''.join(tree.xpath('//div[@class="iframe"]/iframe/@src'))
    # print(youtube_url)
    img_blobs = await download_img(base_url, tree)

    # 入库
    WorldTopBars(
        rank_no=rank_no,
        name=bar_name[0],
        city=city,
        year=2024,
        area='Asia',
        src='theworlds50best.com',
        honor_desc=honor_desc,
        location=location,
        site_url=site_url,
        facebook_url=fb_url,
        instagram_url=ins_url,
        youtube_url='',
        phone=phone,
        img_cover=img_blobs[0],
        img_cover2=img_blobs[1],
        img_cover3=img_blobs[2],
        intro=intro,
    ).save()

    global count
    count = count + 1
    print(f'数据入库 - 数量：{count} - bar_name：{bar_name[0]} - rank_no：{rank_no} - city：{city}')


# 下载图片
async def download_img(base_url, tree):
    urls = tree.xpath('//div[@class="swiper-wrapper"]//img/@src')

    img_blobs = []
    for url in urls:
        url = urljoin(base_url, url)
        r = await send_request('download_img', url, 'GET')
        if r.status != 200:
            raise Exception('Failed to fetch img url: ' + url)
        img_blobs.append(r.content)
    while len(img_blobs) < 3:
        img_blobs.append(None)
    return img_blobs


# 运行
# asyncio.run(run())