import sqlite3

# 1. 连接数据库（如果文件不存在会自动创建）
conn = sqlite3.connect('example.db')  # 也可使用 ':memory:' 创建内存数据库

# 2. 创建游标
cursor = conn.cursor()

# 3. 建表（使用 IF NOT EXISTS 更安全）
create_table_sql = '''
CREATE TABLE IF NOT EXISTS users (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL,
    age INTEGER,
    email TEXT UNIQUE
);
'''

cursor.execute(create_table_sql)

# 4. 提交更改
conn.commit()

# 5. 关闭连接
conn.close()

print("数据库和表创建成功！")
