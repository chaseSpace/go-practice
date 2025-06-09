from typing import Literal

# from dotenv import load_dotenv
from mcp.server.session import ServerSession
from peewee import MySQLDatabase

####################################################################################
# Temporary monkeypatch which avoids crashing when a POST message is received
# before a connection has been initialized, e.g: after a deployment.
# pylint: disable-next=protected-access
old__received_request = ServerSession._received_request


async def _received_request(self, *args, **kwargs):
    try:
        return await old__received_request(self, *args, **kwargs)
    except RuntimeError:
        pass


# pylint: disable-next=protected-access
ServerSession._received_request = _received_request
####################################################################################

config = {
    'database': 'zhong_tai',
    'host': '124.70.208.39',
    'user': 'root',
    'passwd': 'zr5tJBfJ4Xxp2JNB',
    'port': 3306,
    'charset': 'utf8mb4',
}

from mcp.server.fastmcp import FastMCP

# 初始化 MCP 服务器
mcp = FastMCP("QueryMysqlDBServer")

db = MySQLDatabase(**config)
db.connect()


@mcp.tool()
def show_columns(schema, table):
    """
    Show the columns detail of a table.
    :param schema: The schema needs to query; the available arg is one item from ['zhong_tai', 'local_swoft'],
    ensure this params is right.
    :param table: Table name
    :return column details in tuples
    """
    query = """show full columns from {schema}.{table}""".format(schema=schema, table=table)
    try:
        print(f'-- show_columns: {table}')
        rows = db.execute_sql(query).fetchall()
        # for column in rows:
        #     print(column)
        return rows
    except Exception as e:
        print(f'Error: {e}')
        return f'Error: {e}'


@mcp.tool()
def execute_sql(sql):
    """
    Execute only a SELECT SQL statement. SQL must with a database name, and no multi SQL is allowed;
    Note: its best to avoid providing an SQL statement without a LIMIT or SELECT *, instead use a valid WHERE condition
    and SELECT some fields.
    :param sql: SQL statement (must with upper-case except for table name)
    :return data results in a list
    """
    try:
        print(f'-- execute_sql: {sql}')
        if not sql.startswith('SELECT'):
            return 'Error: SQL must start with SELECT'
        cursor = db.execute_sql(sql)
        columns = [desc[0] for desc in cursor.description]
        results = [dict(zip(columns, row)) for row in cursor.fetchall()]
        for column in results:
            print(column)
        return results
    except Exception as e:
        return f'Error: {e}'


@mcp.tool()
def show_tables(schema, where='1=1'):
    """
    Show all tables within the schema from params.
    :param schema: The schema needs to query; the available arg is one item from ['zhong_tai', 'local_swoft'],
    ensure this params is right.
    :param where: Conditions clause to quickly filter results, such as [table_name like '%project%']
    :return table names in a tuple
    """
    sql = f"""
        SELECT table_name
        FROM information_schema.tables
        WHERE table_schema = '{schema}' AND {where};
    """
    try:
        print(f'-- show_tables')
        rows = db.execute_sql(sql).fetchall()
        # for row in rows:
        #     print(row[0])
        return rows
    except Exception as e:
        print(f'Error: {e}')
        return f'Error: {e}'


if __name__ == '__main__':
    # show_columns('zt_provider_usage_bill_detail_base')
    # execute_sql('SELECT project_id, project_name FROM local_swoft.sf_sys_project')
    show_tables()
