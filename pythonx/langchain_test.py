import os
import time
from pprint import pprint

from langchain_community.llms import Tongyi
from langchain_community.llms.moonshot import Moonshot
from langchain_core.messages import HumanMessage, SystemMessage

from zzz.gpt import KIMI_APIKEY,TONGYI_APIKEY

# Generate your api key from: https://platform.moonshot.cn/console/api-keys
os.environ["MOONSHOT_API_KEY"] = KIMI_APIKEY
model = Moonshot(model="moonshot-v1-128k")


# os.environ["DASHSCOPE_API_KEY"] = TONGYI_APIKEY
# model = Tongyi(model='qwq-plus')

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
    template="使用中文列举3部{author}的作品。\n{format_instructions}",
    input_variables=["author"],
    partial_variables={"format_instructions": parser.get_format_instructions()},
)
chain = prompt | model | parser
result = chain.invoke({"author": "老舍"})
pprint(result)
