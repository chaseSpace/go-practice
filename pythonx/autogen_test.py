from autogen import ConversableAgent

from zzz.gpt import KIMI_APIKEY

config_list = [{"model": "moonshot-v1-8k", "api_key": KIMI_APIKEY, "base_url": "https://api.moonshot.cn/v1"}]

student_agent = ConversableAgent(
    name="学生",
    system_message="你是一个学生，现在正在学习数学",
    llm_config={"config_list": config_list},
)
teacher_agent = ConversableAgent(
    name="老师",
    system_message="你是一名老师，现在正在帮助学生学习数学",
    llm_config={"config_list": config_list},
)

# 将英文描述转换为中文
ConversableAgent.DEFAULT_SUMMARY_PROMPT = "从对话中总结收获，不要添加任何介绍性短语"

chat_result = student_agent.initiate_chat(
    teacher_agent,
    message="什么是线性代数？",
    summary_method="reflection_with_llm",
    max_turns=2,
)

print("总结--\n",chat_result.summary)  # 打印【学生-老师】智能体之间的对话总结
