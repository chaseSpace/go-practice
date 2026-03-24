## Agent 学习


### SKILL

- name/description
- instructions，定义输入输出规则，可以根据需求指示LLM使用下面的额外资料或脚本
    - reference
    - script

------

- 渐进式披露机制
    - 元数据层（始终加载）：name/description
    - 指令层（按需加载）：SKILL.md 中除了name/description以外的内容
    - 资源层（按需加载）：reference/script

#### 与MCP的区别

- MCP为大模型提供数据
- SKILL为大模型提供处理数据的方法
