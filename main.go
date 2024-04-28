package main

import (
	"fmt"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
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

func InsertMultiple() {
	students := []Student{
		{Name: "张三", Members: map[string]string{"father": "张大", "mother": "王小花", "brother": "张四", "sister": "张春华"}},
		{Name: "李四", Members: map[string]string{"father": "李大", "mother": "赵红梅", "sister": "李林"}},
		{Name: "王五", Members: map[string]string{"father": "王大", "sister": "王艳华"}},
		{Name: "赵六", Members: map[string]string{"father": "赵大", "mother": "陈秀兰", "sister": "陈晓丽", "sister2": "陈晓美"}},
		{Name: "刘七", Members: map[string]string{"father": "刘大", "mother": "黄秋菊", "sister": "刘英"}},
		{Name: "陈八", Members: map[string]string{"mother": "郭晓婷"}},
		{Name: "杨九", Members: map[string]string{"father": "杨大", "mother": "杨晓燕", "brother": "张林"}},
		{Name: "韩信", Members: map[string]string{}},
		{Name: "朱元璋"},
	}
	DB.Create(&students)
}

func SelectMultiple() {
	var students []Student
	DB.Find(&students)
	for _, record := range students {
		fmt.Println("查询多条结果=", record) //Members字段的值为空
	}
}

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

// 查找所有母亲是“王小花”的学生
func select1() {
	var records []StudentRecord
	DB.Raw("SELECT * FROM students WHERE members->>'mother' = ?", "王小花").Scan(&records)
	for _, record := range records {
		fmt.Println("select1查询结果=", record)
	}
}

// 查找所有母亲不是是“王小花”的学生
func select2() {
	var records []StudentRecord
	DB.Raw("SELECT * FROM students WHERE members->>'mother' != ?", "王小花").Scan(&records)
	for _, record := range records {
		fmt.Println("select2查询结果=", record)
	}
}

// 查找所有有姐妹“张春华”的学生：
func select3() {
	var records []StudentRecord
	DB.Raw("SELECT * FROM students WHERE members @> ?", `{"sister": "张春华"}`).Scan(&records)
	for _, record := range records {
		fmt.Println("select3查询结果=", record)
	}
}

func select4() {
	var records []StudentRecord
	DB.Raw("SELECT * FROM students WHERE ? <@ members", `{"sister": "张春华"}`).Scan(&records)
	for _, record := range records {
		fmt.Println("select4查询结果=", record)
	}
}

// 查找所有有姐妹的学生：
func select5() {
	var records []StudentRecord
	DB.Raw("SELECT * FROM students WHERE members ? 'sister'").Scan(&records)
	for _, record := range records {
		fmt.Println("select5查询结果=", record)
	}
}

// 查找所有有父母的学生：
func select6() {
	var records []StudentRecord
	DB.Raw("SELECT * FROM students WHERE members ?& array['father','mother']").Scan(&records)
	for _, record := range records {
		fmt.Println("select6查询结果=", record)
	}
}

// 查找所有有兄弟或者姐妹的学生：
func select7() {
	var records []StudentRecord
	DB.Raw("SELECT * FROM students WHERE members ?| array['brother','sister']").Scan(&records)
	for _, record := range records {
		fmt.Println("select7查询结果=", record)
	}
}

// 查找所有姐妹列表中包含“李林”的学生（注意这里 sister 是一个数组）
func select8() {
	var records []StudentRecord
	DB.Raw("SELECT * FROM students WHERE members ?| array['杨大'] ").Scan(&records)
	for _, record := range records {
		fmt.Println("select8查询结果=", record)
	}
}

// 索引值访问
// 获取 JSON 对象中指定键的值（返回 jsonb）
func select9() {
	var records []StudentRecord
	DB.Raw("SELECT * FROM students WHERE members ?| array['杨大'] ").Scan(&records)
	for _, record := range records {
		fmt.Println("select8查询结果=", record)
	}
}
func main() {
	//初始化数据库
	InitDB()
	//1.插入多条数据
	// InsertMultiple()

	select6()
}
