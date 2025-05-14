import base64
import io

from PIL import Image

"""
*** 使用langchain接入硅基流动的文生文和视觉模型 ***
"""


def convert_image_to_webp_base64(input_image_path):
    try:
        with Image.open(input_image_path) as img:
            byte_arr = io.BytesIO()
            img.save(byte_arr, format='webp')
            byte_arr = byte_arr.getvalue()
            base64_str = base64.b64encode(byte_arr).decode('utf-8')
            return base64_str
    except IOError:
        print(f"Error: Unable to open or convert the image {input_image_path}")
        return None


###

import os

# 例子1：langchain-硅基流动模型-对话
from langchain_openai import ChatOpenAI
from langchain_core.output_parsers import JsonOutputParser
from langchain_core.prompts import PromptTemplate, ChatPromptTemplate
from pydantic import BaseModel, Field
from pprint import pprint
from zzz.gpt import GUIJI_APIKEY

os.environ["OPENAI_API_KEY"] = GUIJI_APIKEY

llm = ChatOpenAI(
    # model="Qwen/Qwen2.5-VL-32B-Instruct", # ￥1.89/ M Tokens
    model='Qwen/Qwen3-14B',  # ￥1/M Tokens
    temperature=0.5,
    base_url="https://api.siliconflow.cn/v1"
)


class Work(BaseModel):
    title: str = Field(description="Title of the work")
    description: str = Field(description="Description of the work")
    time: int = Field(description="Time of the work publish")


parser = JsonOutputParser(pydantic_object=Work)

prompt = PromptTemplate(
    template="使用中文列举3部{author}的作品，将description固定到100字左右。\n{format_instructions}",
    input_variables=["author"],
    partial_variables={"format_instructions": parser.get_format_instructions()},
)
chain = prompt | llm | parser

# pprint(chain.invoke({"author": "老舍"}))
# result = chain.stream({"author": "老舍"})
# for chunk in result:
#     print(chunk)

# 例子2：langchain-硅基流动模型-视觉
llm = ChatOpenAI(
    model="Qwen/Qwen2.5-VL-32B-Instruct",  # ￥1.89/ M Tokens
    temperature=0.5,
    base_url="https://api.siliconflow.cn/v1"
)


class BarEvent(BaseModel):
    is_event: bool = Field(description="Whether the image is an event or not")
    title: str = Field(description="Title of the bar event")
    description: str = Field(description="Description of the bar event")
    date: str = Field(description="Date of the bar event, use `yyyy-mm-dd` format")
    location: str = Field(description="Location of the bar event")
    hour_start: str = Field(
        description="Start hour of the bar event, use `h:i` format in 24H pattern. Keep empty if not provided.")
    hour_end: str = Field(
        description="End hour time of the bar event, use `h:i` format in 24H pattern. "
                    "Set to `late` or empty if not provided.")


parser = JsonOutputParser(pydantic_object=BarEvent)

prompt = ChatPromptTemplate(
    messages=[
        (
            "user", [
                {
                    "type": "image_url",
                    "image_url": {
                        "url": "data:image/jpeg;base64,{base64_image}",
                        "detail": "low"
                    }
                },
                {
                    "type": "text",
                    "text": "Extracting bar event info the image provided. \n{format_instructions}"
                }]
        ),
    ],
    partial_variables={"format_instructions": parser.get_format_instructions()},
)

chain = prompt | llm | parser
base64_image = convert_image_to_webp_base64("./assets/chinese_event2.jpg")
if not base64_image:
    raise None

result = chain.invoke({"base64_image": base64_image})
pprint(result)
