package main

import (
	"fmt"
	"io/ioutil"
	"os"
)

// User 用户结构体
type User struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// MinimalManager 最小化管理器
type MinimalManager struct {
	users  map[int]User
	nextID int
}

// NewMinimalManager 创建管理器
func NewMinimalManager() *MinimalManager {
	return &MinimalManager{
		users:  make(map[int]User),
		nextID: 1,
	}
}

// AddUser 添加用户
func (m *MinimalManager) AddUser(name string) {
	user := User{ID: m.nextID, Name: name}
	m.users[m.nextID] = user
	m.nextID++
}

// ShowUsers 显示所有用户
func (m *MinimalManager) ShowUsers() {
	fmt.Println("用户列表:")
	for _, user := range m.users {
		fmt.Printf("ID: %d, 姓名: %s\n", user.ID, user.Name)
	}
}

// SaveToFile 保存到文件
func (m *MinimalManager) SaveToFile() {
	data := ""
	for _, user := range m.users {
		data += fmt.Sprintf("%d,%s\n", user.ID, user.Name)
	}
	ioutil.WriteFile("users.txt", []byte(data), 0644)
}

func main() {
	manager := NewMinimalManager()

	// 添加用户
	manager.AddUser("张三")
	manager.AddUser("李四")
	manager.AddUser("王五")

	// 显示用户
	manager.ShowUsers()

	// 保存文件
	manager.SaveToFile()
	fmt.Println("数据已保存到 users.txt")
}
