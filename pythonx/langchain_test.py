import os
import time
from pprint import pprint

from langchain_community.llms import Tongyi
from langchain_core.messages import HumanMessage, SystemMessage

from zzz.gpt import TONGYI_APIKEY

# Generate your api key from: https://platform.moonshot.cn/console/api-keys
# os.environ["MOONSHOT_API_KEY"] = KIMI_APIKEY
# model = Moonshot(model="moonshot-v1-128k")


# qwen-max（推荐） > qwen-plus(勿用)
os.environ["DASHSCOPE_API_KEY"] = TONGYI_APIKEY
model = Tongyi(model='qwen-max')

# or use a specific model
# Available models: https://platform.moonshot.cn/docs
# llm = Moonshot(model="moonshot-v1-128k")

messages = [
    SystemMessage(content="Translate the following 10x times from English into Chinese:"),
    HumanMessage(content="Welcome to LLM application development!"),
]

# 同步方式调用，会等到所有内容生成后返回，很慢
# ret = model.invoke(messages)
# pprint(ret)

# 流式方式调用，会返回一个生成器，可以迭代获取内容，速度更快
# 貌似kimi不支持流式调用
# stream = model.stream(messages)
# for response in stream:
#     pprint(response, end="|")

# 还有批处理和异步调用

time.sleep(1)

### 2. 使用提示词模板
from langchain_core.prompts import ChatPromptTemplate

prompt_template = ChatPromptTemplate.from_messages(
    [
        ("system", "Translate the following from English into Chinese:"),
        ("user", "{text}")
    ]
)

chain = prompt_template | model  # LC 使用py语法来表达了这种组件链式调用的语法，注意顺序有关
result = chain.invoke({"text": "Welcome to LLM application development!"})
pprint(f'chain.invoke:{result}')  # 某些大模型在这里返回一个对象，比如openai
time.sleep(1)

# 查看提示词拼装后的内容，和最开始的 [SystemMessage(...), HumanMessage(...)] 一样
# messages = prompt_template.invoke({"text":"Welcome to LLM application development!"})
# pprint(messages)


### 3. OutputParser
# 将输出进一步处理，String、JSON、CSV、分割符、枚举等
from langchain_core.output_parsers import StrOutputParser

parser = StrOutputParser()  # 将输出中的内容提取为字符串
chain = prompt_template | model | parser
result = chain.invoke({"text": "Welcome to LLM application development!"})
pprint(f'chain.invoke:{result}')  # 输出为字符串而不是对象
time.sleep(1)

### 3.1 OutputParser：JSON
from langchain_core.output_parsers import JsonOutputParser
from langchain_core.prompts import PromptTemplate
from pydantic import BaseModel, Field


class Work(BaseModel):
    title: str = Field(description="Title of the work")
    description: str = Field(description="Description of the work")
    time: int = Field(description="Time of the work publish")


parser = JsonOutputParser(pydantic_object=Work)
prompt = PromptTemplate(
    template="使用中文列举3部{author}的作品，将description固定到200字左右。\n{format_instructions}",
    input_variables=["author"],
    partial_variables={"format_instructions": parser.get_format_instructions()},
)
chain = prompt | model | parser
# result = chain.invoke({"author": "老舍"})
# pprint(result)

### 3.2 OutputParser：JSON (图片解析)
from langchain_community.chat_models import ChatTongyi
import base64

model = ChatTongyi(model_name='qwen-vl-max')


# from langchain_community.chat_models import ChatTongyi


class BarEvent(BaseModel):
    title: str = Field(description="Title of the bar event")
    description: str = Field(description="Description of the bar event")
    date: str = Field(description="Date of the bar event")
    hour: str = Field(description="Detail time of the bar event")


image_url = "../img/bar_event.jpg"

# 打开图片文件
with open(image_url, "rb") as image_file:
    # 将图片读取为二进制数据
    image_data = image_file.read()
    # 将二进制数据编码为Base64字符串
    img_base64 = base64.b64encode(image_data).decode('utf-8')

prompt = ChatPromptTemplate(
    messages=[
        ("system", "Extracting bar event info the image provided."),
        (
            "user",
            [
                {"image": "data:image/jpeg;base64,{image_data}"}
            ],
        ),
    ],
)
model = model.with_structured_output(BarEvent)
chain = prompt | model
result = chain.invoke({"image_data": img_base64})
pprint(result)
