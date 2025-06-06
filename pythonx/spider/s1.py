import re
from model import WorldTopBars

import pandas as pd

df = pd.read_excel("theworlds50best-2025-06-05.xlsx")

ct = 0
# 打印每一行数据
for index, row in enumerate(df.values):
    # if not pd.isna(row[4]):
    #     pass
    # rank_no = re.search(r'\d+',row[0]).group()
    #
    # # print(111, row[2])
    # city = re.search(r'\((.*?)/',row[2]).group(1)
    # site_url = row[4] if row[4].startswith('http') else None
    # WorldTopBars(
    #     src='top500bars.com',
    #     year=2023,
    #     rank_no=rank_no,
    #     name=row[1],
    #     city=city,
    #     intro=row[3],
    #     site_url=site_url,
    # ).save(True)
    v = re.search(r'^(\d+)$', str(row[0]))
    if v and v.group() == str(row[0]):
        print(row)
        ct += 1
        site_url = row[1] if not pd.isna(row[1]) and row[1].startswith('http') else None
        name = re.search(r'^(.*?)\(', row[2]).group(1)
        city = re.search(r'\((.*?)/', row[2]).group(1)
        WorldTopBars(
            src='top500bars.com',
            year=2023,
            rank_no=int(row[0]),
            name=name,
            city=city,
            site_url=site_url,
        ).save(True)
        print(222, int(row[0]), name, city, site_url)

print(ct)
