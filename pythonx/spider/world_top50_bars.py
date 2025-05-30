import asyncio, aiometer
from urllib.parse import urljoin
from functools import partial
from lxml import html, etree

from infra import send_request, remove_path_simple, save_html, init_session, close_session
from model import WorldTopBars

count = 0


# 当前任务的运行入口
async def run():
    url = 'https://www.theworlds50best.com/bars/list/1-50'
    tasks = [asyncio.create_task(start_scrape(url))]

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

    # 并发执行：抓取1-50个酒吧（含详情页）
    tasks = [partial(parse_detail_page, base_url, href) for href in hrefs]
    await aiometer.run_all(tasks, max_at_once=3)  # 并发数限制为3

    # 并发执行：抓取51-100个酒吧（不含详情页）
    tasks = [partial(parse_51_100, base_url, i, div) for i, div in
             enumerate(tree.xpath('//div[@data-list="51-100"]/div'))]
    await aiometer.run_all(tasks, max_at_once=3)  # 并发数限制为3


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
    WorldTopBars.insert(
        rank_no=rank_no,
        name=bar_name[0],
        city=city,
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
    ).on_conflict(action='IGNORE').execute()

    global count
    count = count + 1
    print(f'数据入库 - 数量：{count} - bar_name：{bar_name[0]} - rank_no：{rank_no} - city：{city}')


async def parse_51_100(baseurl, i, div: etree._ElementTree):
    rank_no = int(div.xpath('.//p[@class="position "]/text()')[0])
    img_src = div.xpath('.//img/@src')[0]
    bar_name = div.xpath('.//h2/text()')[0]
    city = div.xpath('.//p[2]/text()')[0]
    # print(rank_no, img_src, bar_name, city)
    WorldTopBars.insert(
        rank_no=rank_no,
        name=bar_name,
        city=city,
        img_cover=await download_img_from_url(urljoin(baseurl, img_src)),
    ).on_conflict(action='IGNORE').execute()
    print(f'(51-100)数据入库 - 数量：{i + 1} - bar_name：{bar_name} - rank_no：{rank_no} - city：{city}')


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


async def download_img_from_url(url):
    r = await send_request('download_img', url, 'GET')
    if r.status != 200:
        raise Exception('Failed to fetch img url: ' + url)
    return r.content

# 运行
# asyncio.run(run())
