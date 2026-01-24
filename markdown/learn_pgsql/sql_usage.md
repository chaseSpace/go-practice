## PgSQL 建表

### 建库

```sql
CREATE
DATABASE mydatabase ENCODING = 'gbk' LC_COLLATE = 'en_US.utf8' LC_CTYPE = 'en_US.utf8';
```

PgSQL 不同的是，它支持字符排序规则、分类规则拆开设置，而MySQL则是合并到 `COLLATE` 中。

### 建立 schema

schema 并非物理上的概念，而是在库的维度上对表做逻辑上的分组。

```sql
CREATE DATABASE mydatabase;
use mydatabase; -- 此时已进入数据库 mydatabase 的 public schema
    
-- 也可以创建一个 schema
CREATE SCHEMA myschema;
```


### 建表

```sql
CREATE TABLE employees
(
    id          SERIAL PRIMARY KEY,
    name        VARCHAR(100) NOT NULL,
    email       VARCHAR(100) UNIQUE,
    age         INT,
    salary      NUMERIC(10, 2),
    description TEXT,
    hire_date   DATE,
    CONSTRAINT chk_age CHECK (age >= 18 AND age <= 100)
);
```

- **自增列(SERIAL)**：在 PostgreSQL 中，SERIAL 数据类型用来定义自增列，实际上它是 INTEGER 类型并自动创建一个序列（sequence）。
- **没有存储引擎选项**：PostgreSQL 使用单一的存储引擎，它的存储引擎与表类型的关系不需要显式声明。PostgreSQL 自动处理所有存储和事务管理。
- **浮点数**：PostgreSQL 使用 `NUMERIC` 数据类型来存储高精度浮点数。
- **TEXT**：PostgreSQL 使用 `TEXT` 数据类型来存储任意长度的字符串；不同于MySQL会区分`TINYTEXT`/`text`/`MEDIUMTEXT`/
  `LONGTEXT`。
- **约束**：pgsql和MySQL（v8）都支持使用 CONSTRAINT 来约束字段。

### 常用命令

```sql

-- 查所有db， CLI使用 \l
-- 数据库版本
SELECT version();
show
server_version;

-- 当前数据库名
SELECT current_database();

-- 当前用户
SELECT current_user;

-- 连接信息
SELECT inet_client_addr(), inet_client_port(), inet_server_addr(), inet_server_port();

-- 数据库大小
SELECT pg_size_pretty(pg_database_size(current_database()));


-- 查时区
show
timezone;
SELECT now();
SELECT current_setting('timezone');


-- 查看客户端编码
SHOW
client_encoding;
-- 或者
SELECT current_setting('client_encoding');

-- 查看所有schema，或者CLI中使用 \dn+
SELECT schema_name
FROM information_schema.schemata;

-- 查看当前连接使用的 schema
SHOW
search_path;

-- 查询schema下的表，或者CLI中使用 \dt+ schema_name.*
SELECT table_name, table_schema
FROM information_schema.tables
WHERE table_schema = 'public';

-- 进程模型的connect不支持切换db
-- use xx

-- 在指定的 Schema 中创建表
CREATE TABLE myschema.employees
(
    id   SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL
);
```