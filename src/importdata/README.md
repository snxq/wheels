# import data

目前仅实现了从 JSON 文件/文件夹 到 MySQL 数据库的数据单向导入。

## Usage

```
importdata -dsn='root:@(localhost:3306)/database' -path=filepath

-dsn string
    data source name
-path string
    data file path
```
