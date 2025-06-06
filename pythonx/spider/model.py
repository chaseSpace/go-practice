import datetime

import peewee

# peewee 教程
# https://benpaodewoniu.github.io/2020/03/14/python77/

db = peewee.SqliteDatabase('jump_jump_tiger.db')

db.connect()


# 定义数据库模型
class WorldTopBars(peewee.Model):
    rank_no = peewee.IntegerField(null=True)  # 为0表示 individual award bar
    name = peewee.CharField(null=True)
    city = peewee.CharField(null=True)
    year = peewee.IntegerField(null=True)
    area = peewee.CharField(null=True)
    src = peewee.CharField(null=True)
    honor_desc = peewee.CharField(null=True)
    location = peewee.CharField(null=True)
    site_url = peewee.CharField(null=True)
    facebook_url = peewee.CharField(null=True)
    instagram_url = peewee.CharField(null=True)
    youtube_url = peewee.CharField(null=True)
    phone = peewee.CharField(null=True)
    created = peewee.DateTimeField(default=datetime.datetime.now)
    intro = peewee.CharField(null=True)
    img_cover = peewee.BlobField(null=True)
    img_cover2 = peewee.BlobField(null=True)
    img_cover3 = peewee.BlobField(null=True)

    class Meta:
        database = db
        db_table = 'world_top_bars'

# WorldTopBars.create_table()
