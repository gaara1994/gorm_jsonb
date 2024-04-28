# gorm jsonb

官方文档： http://www.postgres.cn/docs/12/functions-json.html



先创建表并插入数据

```go
package main

import (
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
)

var DB *gorm.DB

type Student struct {
	gorm.Model
	Name    string            `gorm:"not null;comment:'学生姓名'" json:"name"`
	Members map[string]string `gorm:"type:jsonb;comment:'家庭成员'" json:"members"`
}

func (Student) TableName() string {
	return "students"
}

func InitDB() {
	dsn := "host=localhost user=postgres password=MyNewPass4! dbname=mydatabase port=5432 sslmode=disable TimeZone=Asia/Shanghai"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		fmt.Println(err)
	}
	err = db.AutoMigrate(&Student{})
	if err != nil {
		log.Fatal(err)
	}
	DB = db
}
func main() {
	//初始化数据库
	InitDB()

	//插入多条数据
	InsertMultiple()
}

func InsertMultiple() {
	students := []Student{
		{Name: "张三", Members: map[string]string{"father": "张大", "mother": "王小花", "brother": "张四", "sister": "张春华"}},
		{Name: "李四", Members: map[string]string{"father": "李大", "mother": "赵红梅", "sister": "李林"}},
		{Name: "王五", Members: map[string]string{"father": "王大", "sister": "王艳华"}},
		{Name: "赵六", Members: map[string]string{"father": "赵大", "mother": "陈秀兰", "sister": "陈晓丽", "sister2": "陈晓美"}},
		{Name: "刘七", Members: map[string]string{"father": "刘大", "mother": "黄秋菊", "sister": "刘英"}},
		{Name: "陈八", Members: map[string]string{"mother": "郭晓婷"}},
		{Name: "杨九", Members: map[string]string{"father": "杨大", "mother": "杨晓燕", "brother": "张林"}},
	}
	DB.Create(&students)
}
```



## 查询

### 1.缺陷

对于`jsonb`的数据，可以使用`map[string]string`类型插入，但是无法读取。

```go
func SelectMultiple() {
	var students []Student
	DB.Find(&students)
	for _, record := range students {
		fmt.Println("查询多条结果=", record) //Members字段的值为空
	}
}
```



### 2.改变

重新定义结构体，接收`jsonb`字段

```go
type StudentRecord struct {
	gorm.Model
	Name    string `json:"name"`
	Members string `json:"members"`
}

func SelectMultiple2() {
	var records []StudentRecord
	DB.Table("students").Find(&records)
	for _, record := range records {
		fmt.Println("查询多条结果2=", record.Members)
	}
}
```



## 比较操作符

### =

等于

```go
// 查找所有母亲是“王小花”的学生
func select1() {
	var records []StudentRecord
	DB.Raw("SELECT * FROM students WHERE members->>'mother' = ?", "王小花").Scan(&records)
	for _, record := range records {
		fmt.Println("select1查询结果=", record)
	}
}
```



### <> 或 !=

不等于

```go
// 查找所有母亲不是是“王小花”的学生
func select2() {
	var records []StudentRecord
	DB.Raw("SELECT * FROM students WHERE members->>'mother' != ?", "王小花").Scan(&records)
	for _, record := range records {
		fmt.Println("select2查询结果=", record)
	}
}
```





## 包含操作符

### @>

左侧 JSON 值包含右侧 JSON 值

```go
// 查找所有有姐妹“张春华”的学生：
SELECT * FROM students WHERE members @> ?", `{"sister": "张春华"}
```



### <@

左侧 JSON 值被右侧 JSON 值包含

```go
SELECT * FROM students WHERE ? <@ members", `{"sister": "张春华"}
```



## 存在操作符

### ?

检查 JSON 对象包含一个指定的键的对象

```go
SELECT * FROM students WHERE members ? 'sister'
```



### ?&

检查 JSON 对象是否包含所有指定的键

```sql
SELECT * FROM students WHERE members ?& array['father','mother']
```



### ?|

检查 JSON 对象是否包含任何指定的键

```json
SELECT * FROM students WHERE members ?| array['brother','sister']
```





## 索引值访问



| 操作符 | 右操作数类型 | 返回类型          | 描述                                                    | 例子                                               | 例子结果       |
| ------ | ------------ | ----------------- | ------------------------------------------------------- | -------------------------------------------------- | -------------- |
| `->`   | `int`        | `json` or `jsonb` | 获得 JSON 数组元素（索引从 0 开始，负整数从末尾开始计） | `'[{"a":"foo"},{"b":"bar"},{"c":"baz"}]'::json->2` | `{"c":"baz"}`  |
| `->`   | `text`       | `json` or `jsonb` | 通过键获得 JSON 对象域                                  | `'{"a": {"b":"foo"}}'::json->'a'`                  | `{"b":"foo"}`  |
| `->>`  | `int`        | `text`            | 以`text`形式获得 JSON 数组元素                          | `'[1,2,3]'::json->>2`                              | `3`            |
| `->>`  | `text`       | `text`            | 以`text`形式获得 JSON 对象域                            | `'{"a":1,"b":2}'::json->>'b'`                      | `2`            |
| `#>`   | `text[]`     | `json` or `jsonb` | 获取在指定路径的 JSON 对象                              | `'{"a": {"b":{"c": "foo"}}}'::json#>'{a,b}'`       | `{"c": "foo"}` |
| `#>>`  | `text[]`     | `text`            | 以`text`形式获取在指定路径的 JSON 对象                  | `'{"a":[1,2,3],"b":[4,5,6]}'::json#>>'{a,2}'`      | `3`            |



### ->

```
SELECT members::json->'mother' FROM students WHERE members ?& array['father','mother']
```

```sql
"王小花"
"赵红梅"
"陈秀兰"
"黄秋菊"
"杨晓燕"
```



### ->>

```sql
SELECT members::json->>'mother' FROM students WHERE members ?& array['father','mother']
```

```
王小花
赵红梅
陈秀兰
黄秋菊
杨晓燕
```

