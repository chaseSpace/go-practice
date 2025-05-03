# 为了支持聊天历史，LangChain 引入了一个抽象叫 ChatMessageHistory
from langchain_core.chat_history import BaseChatMessageHistory, InMemoryChatMessageHistory
from langchain_core.runnables import RunnableWithMessageHistory
from langchain_openai import ChatOpenAI
from langchain_core.messages import HumanMessage

#
# 这里的 Runnable 是一个接口，它表示一个工作单元。我们前面说过，组成链是由一个一个的组件组成的。
#
# 严格地说，这些组件都实现了 Runnable 接口，甚至链本身也实现了 Runnable 接口，我们之前讨论的 invoke、stream 等接口都是定义在 Runnable 里，
# 可以说，Runnable 是真正的基础类型，LCEL 之所以能够以声明式的方式起作用，Runnable 接口是关键。
#
# 不过，在真实的编码过程中，我们很少会直接面对 Runnable，大多数时候我们看见的都是各种具体类型。只是你会在很多具体类的名字中见到 Runnable，
# 这里的 RunnableWithMessageHistory 就是其中一个。


chat_model = ChatOpenAI(model="gpt-4o-mini")

store = {}

def get_session_history(session_id: str) -> BaseChatMessageHistory:
    if session_id not in store:
        store[session_id] = InMemoryChatMessageHistory()
    return store[session_id]

# 是一个把聊天历史和链封装到一起的一个类
with_message_history = RunnableWithMessageHistory(chat_model, get_session_history)

config = {"configurable": {"session_id": "dreamhead"}}

while True:
    user_input = input("You:> ")
    if user_input.lower() == 'exit':
        break
    stream = with_message_history.stream(
        [HumanMessage(content=user_input)],
        config=config
    )
    for chunk in stream:
        print(chunk.content, end='', flush=True)
    print()