from operator import itemgetter
from typing import List

import tiktoken
from langchain_chroma import Chroma
from langchain_core.chat_history import BaseChatMessageHistory, InMemoryChatMessageHistory
from langchain_core.messages import BaseMessage, HumanMessage, AIMessage, ToolMessage, SystemMessage, trim_messages
from langchain_core.prompts import ChatPromptTemplate, MessagesPlaceholder
from langchain_core.runnables import RunnablePassthrough
from langchain_core.runnables.history import RunnableWithMessageHistory
from langchain_openai.chat_models import ChatOpenAI
from langchain_community.embeddings import DashScopeEmbeddings  # tongyi
from langchain_community.chat_models import ChatTongyi
from langchain_core.output_parsers import StrOutputParser

import os
from zzz.gpt import TONGYI_APIKEY
os.environ["DASHSCOPE_API_KEY"] = TONGYI_APIKEY

vectorstore = Chroma(
    collection_name="ai_learning",
    embedding_function=DashScopeEmbeddings(model="text-embedding-v2"),
    persist_directory="vectordb"
)

retriever = vectorstore.as_retriever(search_type="similarity")


def str_token_counter(text: str) -> int:
    enc = tiktoken.get_encoding("o200k_base")
    return len(enc.encode(text))


def tiktoken_counter(messages: List[BaseMessage]) -> int:
    num_tokens = 3
    tokens_per_message = 3
    tokens_per_name = 1
    for msg in messages:
        if isinstance(msg, HumanMessage):
            role = "user"
        elif isinstance(msg, AIMessage):
            role = "assistant"
        elif isinstance(msg, ToolMessage):
            role = "tool"
        elif isinstance(msg, SystemMessage):
            role = "system"
        else:
            raise ValueError(f"Unsupported messages type {msg.__class__}")
        num_tokens += (
                tokens_per_message
                + str_token_counter(role)
                + str_token_counter(msg.content)
        )
        if msg.name:
            num_tokens += tokens_per_name + str_token_counter(msg.name)
    return num_tokens


trimmer = trim_messages(
    max_tokens=4096,
    strategy="last",
    token_counter=tiktoken_counter,
    include_system=True,
)

store = {}


def get_session_history(session_id: str) -> BaseChatMessageHistory:
    if session_id not in store:
        store[session_id] = InMemoryChatMessageHistory()
    return store[session_id]



model = ChatTongyi(model_name='qwen-vl-max')
# model = ChatOpenAI()

prompt = ChatPromptTemplate.from_messages(
    [
        (
            "system",
            """你是一个借助上下文回答问题的助手，如果上下文未提供相关信息，直接回答不知道即可，保持回答的简洁。
            Context: {context}""",
        ),
        MessagesPlaceholder(variable_name="history"),
        ("human", "{question}"),
    ]
)


# 在发送给大模型之前，需要将文本拼起来
def format_docs(docs):
    return "\n\n".join(doc.page_content for doc in docs)


# context也是链条组成
# 在 LangChain 代码里， | 运算符被用作不同组件之间的连接，其实现的关键就是大部分组件都实现了 Runnable 接口
# 在这个接口里实现了 __or__ 和 __ror__。__or__ 表示这个对象出现在| 左边时的处理，相应的 __ror__ 表示这个对象出现在右边时的处理。
# itemgetter 没有实现 __or__，但retriever实现了__ror__，也能工作
context = itemgetter("question") | retriever | format_docs
# context 最终填充到提示词模板中的同名变量
first_step = RunnablePassthrough.assign(context=context)
parser = StrOutputParser()
chain = first_step | prompt | trimmer | model | parser

with_message_history = RunnableWithMessageHistory(
    chain,
    get_session_history=get_session_history,
    input_messages_key="question",
    history_messages_key="history",
)

config = {"configurable": {"session_id": "dreamhead"}}

while True:
    user_input = input("You:> ")
    if user_input.lower() == 'exit':
        break

    if user_input.strip() == "":
        continue

    stream = with_message_history.stream(
        {"question": user_input},
        config=config
    )
    for chunk in stream:
        print(chunk, end='', flush=True)
    print()
