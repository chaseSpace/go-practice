from pathlib import Path

from openai import OpenAI

from zzz.gpt import KIMI_APIKEY

client = OpenAI(
    api_key=KIMI_APIKEY,
    base_url="https://api.moonshot.cn/v1",
)

file_object = client.files.create(file=Path("haibao.jpg"), purpose="file-extract")
file_content = client.files.content(file_id=file_object.id).text

# 把它放进请求中
messages = [
    {
        "role": "system",
        "content": "提取出图片中的文字",
    },
    {
        "role": "system",
        "content": file_content,
    },
    {"role": "user",
     "content": "请提取这张图中的信息，并分析是否一个酒吧活动海报，若是则返回JSON，其中time_range字段表示活动时间范围，topic字段表示活动主题，若不是活动则返回空的{}"},
]

# 然后调用 chat-completion, 获取 Kimi 的回答
completion = client.chat.completions.create(
    model="moonshot-v1-8k",
    messages=messages,
    temperature=0.3,
)

print(completion.choices[0].message.content)
