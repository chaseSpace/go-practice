import asyncio
import os

from browser_use import Agent, Browser, ChatOpenAI,Controller
from dotenv import load_dotenv

load_dotenv()

# Create browser with persistent profile to save cookies and login state
browser = Browser(
    user_data_dir='./zzz/browser_profile',  # 持久化浏览器配置文件目录
    keep_alive=True,  # 任务完成后关闭浏览器
    headless=False,  # 显示浏览器窗口
)

# 1. 创建控制器并注册人工步骤动作
controller = Controller()

@controller.action('check_if_need_login')
def check_if_need_login(reason: str):
    print(f"\n[人工干预请求] 网站需要登录: {reason}")
    print("\n请在浏览器中完成操作，完成后回到此处按回车键继续...")
    input(">>> 已手动完成操作，按回车键交还控制权...")
    return "已完成登录，请继续任务。"


agent = Agent(
    # task='访问boss直聘查看agent开发岗位的最高薪资，并输出JD（若需要可能需要登录的情况，请调用check_if_need_login）',
    task='查看腾讯新闻网站内，关于交通大学的最新十条新闻，以JSON方式写入 ./zzz/news.json，需要包含发布时间/标题/内容摘要不超过200字/标签列表',
    browser=browser,
    llm=ChatOpenAI(model=os.getenv('MODEL'), api_key=os.getenv('API_KEY'), base_url=os.getenv('API_URL')),
    tools=controller,
)


async def main():
    await agent.run()


if __name__ == "__main__":
    asyncio.run(main())
