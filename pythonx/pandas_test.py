import pandas as pd

# 进行多列数据的合并

df = pd.read_excel('./assets/ai_table.xlsx')

# 获取列名列表
columns = df.columns.tolist()


def process(start, end):
    columns_to_merge = columns[start: end]
    df.iloc[:, [start]] = (df[columns_to_merge].
                           apply(lambda row: '; '.join([str(x) for x in row if pd.notnull(x)]), axis=1))


list = ((13, 21), (25, 30), (31, 36), (37, 42), (43, 48), (49, 54), (55, 60))

for s, e in list:
    process(s, e)

for (s, e) in list:
    df.drop(columns=columns[s + 1:e + 1], inplace=True)

# 将合并后的数据写回到一个新的Excel文件
df.to_excel('merged_example.xlsx', index=False)

print("数据合并完成，并已保存到 'merged_example.xlsx'")
