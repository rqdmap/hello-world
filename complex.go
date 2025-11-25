package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"time"
)

// User 用户结构体
type User struct {
	ID       int       `json:"id"`
	Name     string    `json:"name"`
	Email    string    `json:"email"`
	Age      int       `json:"age"`
	Created  time.Time `json:"created"`
	Active   bool      `json:"active"`
}

// UserManager 用户管理器
type UserManager struct {
	users map[int]User
	nextID int
	logger *log.Logger
}

// NewUserManager 创建新的用户管理器
func NewUserManager() *UserManager {
	logger := log.New(os.Stdout, "[UserManager] ", log.LstdFlags|log.Lshortfile)
	return &UserManager{
		users:  make(map[int]User),
		nextID: 1,
		logger: logger,
	}
}

// AddUser 添加用户
func (um *UserManager) AddUser(name, email string, age int) error {
	if strings.TrimSpace(name) == "" {
		err := fmt.Errorf("用户名称不能为空")
		um.logger.Printf("警告: %v", err)
		return err
	}
	if !isValidEmail(email) {
		err := fmt.Errorf("无效的邮箱地址: %s", email)
		um.logger.Printf("警告: %v", err)
		return err
	}
	if age < 0 || age > 150 {
		err := fmt.Errorf("年龄必须在0到150之间，当前为: %d", age)
		um.logger.Printf("警告: %v", err)
		return err
	}

	user := User{
		ID:      um.nextID,
		Name:    name,
		Email:   email,
		Age:     age,
		Created: time.Now(),
		Active:  true,
	}

	um.users[um.nextID] = user
	um.nextID++

	um.logger.Printf("成功添加用户: ID=%d, Name=%s, Email=%s", user.ID, user.Name, user.Email)
	return nil
}

// GetUserByID 根据ID获取用户
func (um *UserManager) GetUserByID(id int) (User, error) {
	user, exists := um.users[id]
	if !exists {
		err := fmt.Errorf("用户ID %d 不存在", id)
		um.logger.Printf("错误: %v", err)
		return User{}, err
	}
	um.logger.Printf("成功获取用户: ID=%d, Name=%s", user.ID, user.Name)
	return user, nil
}

// UpdateUser 更新用户信息
func (um *UserManager) UpdateUser(id int, name, email string, age int) error {
	user, exists := um.users[id]
	if !exists {
		err := fmt.Errorf("无法更新 - 用户ID %d 不存在", id)
		um.logger.Printf("错误: %v", err)
		return err
	}

	if strings.TrimSpace(name) != "" {
		user.Name = name
	}
	if email != "" && isValidEmail(email) {
		user.Email = email
	} else if email != "" {
		err := fmt.Errorf("无法更新用户 %d 的邮箱 - 无效邮箱地址: %s", id, email)
		um.logger.Printf("警告: %v", err)
		return err
	}
	if age > 0 {
		if age < 0 || age > 150 {
			err := fmt.Errorf("无法更新用户 %d 的年龄 - 年龄必须在0到150之间，当前为: %d", id, age)
			um.logger.Printf("警告: %v", err)
			return err
		}
		user.Age = age
	}

	um.users[id] = user
	um.logger.Printf("成功更新用户: ID=%d, Name=%s, Email=%s, Age=%d", user.ID, user.Name, user.Email, user.Age)
	return nil
}

// DeleteUser 删除用户
func (um *UserManager) DeleteUser(id int) error {
	_, exists := um.users[id]
	if !exists {
		err := fmt.Errorf("无法删除 - 用户ID %d 不存在", id)
		um.logger.Printf("错误: %v", err)
		return err
	}

	delete(um.users, id)
	um.logger.Printf("成功删除用户: ID=%d", id)
	return nil
}

// ListUsers 获取所有用户
func (um *UserManager) ListUsers() []User {
	users := make([]User, 0, len(um.users))
	for _, user := range um.users {
		users = append(users, user)
	}
	um.logger.Printf("返回 %d 个用户", len(users))
	return users
}

// SearchUsersByName 根据名称搜索用户
func (um *UserManager) SearchUsersByName(name string) []User {
	var results []User
	for _, user := range um.users {
		if strings.Contains(strings.ToLower(user.Name), strings.ToLower(name)) {
			results = append(results, user)
		}
	}
	um.logger.Printf("根据名称 '%s' 搜索到 %d 个用户", name, len(results))
	return results
}

// isValidEmail 检查邮箱格式
func isValidEmail(email string) bool {
	return strings.Contains(email, "@") && strings.Contains(email, ".")
}

// SaveToFile 保存用户数据到文件
func (um *UserManager) SaveToFile(filename string) error {
	data, err := json.MarshalIndent(um.users, "", "  ")
	if err != nil {
		um.logger.Printf("错误: 序列化用户数据失败 - %v", err)
		return err
	}

	err = ioutil.WriteFile(filename, data, 0644)
	if err != nil {
		um.logger.Printf("错误: 写入文件失败 - %v", err)
		return err
	}

	um.logger.Printf("成功保存用户数据到文件: %s", filename)
	return nil
}

// LoadFromFile 从文件加载用户数据
func (um *UserManager) LoadFromFile(filename string) error {
	_, err := os.Stat(filename)
	if os.IsNotExist(err) {
		um.logger.Printf("警告: 文件不存在，将创建新文件 - %s", filename)
		return nil
	}

	data, err := ioutil.ReadFile(filename)
	if err != nil {
		um.logger.Printf("错误: 读取文件失败 - %v", err)
		return err
	}

	err = json.Unmarshal(data, &um.users)
	if err != nil {
		um.logger.Printf("错误: 反序列化用户数据失败 - %v", err)
		return err
	}

	// 更新nextID
	maxID := 0
	for id := range um.users {
		if id > maxID {
			maxID = id
		}
	}
	um.nextID = maxID + 1

	um.logger.Printf("成功从文件加载 %d 个用户: %s", len(um.users), filename)
	return nil
}

// ActivateUser 激活用户
func (um *UserManager) ActivateUser(id int) error {
	user, exists := um.users[id]
	if !exists {
		err := fmt.Errorf("无法激活 - 用户ID %d 不存在", id)
		um.logger.Printf("错误: %v", err)
		return err
	}

	if user.Active {
		err := fmt.Errorf("用户ID %d 已经是激活状态", id)
		um.logger.Printf("警告: %v", err)
		return err
	}

	user.Active = true
	um.users[id] = user
	um.logger.Printf("成功激活用户: ID=%d, Name=%s", user.ID, user.Name)
	return nil
}

// DeactivateUser 停用用户
func (um *UserManager) DeactivateUser(id int) error {
	user, exists := um.users[id]
	if !exists {
		err := fmt.Errorf("无法停用 - 用户ID %d 不存在", id)
		um.logger.Printf("错误: %v", err)
		return err
	}

	if !user.Active {
		err := fmt.Errorf("用户ID %d 已经是停用状态", id)
		um.logger.Printf("警告: %v", err)
		return err
	}

	user.Active = false
	um.users[id] = user
	um.logger.Printf("成功停用用户: ID=%d, Name=%s", user.ID, user.Name)
	return nil
}

// GetUserCount 获取用户总数
func (um *UserManager) GetUserCount() int {
	count := len(um.users)
	um.logger.Printf("当前用户总数: %d", count)
	return count
}

// GetActiveUsers 获取活跃用户
func (um *UserManager) GetActiveUsers() []User {
	var activeUsers []User
	for _, user := range um.users {
		if user.Active {
			activeUsers = append(activeUsers, user)
		}
	}
	um.logger.Printf("当前活跃用户数: %d", len(activeUsers))
	return activeUsers
}

func main() {
	um := NewUserManager()
	
	// 从文件加载数据
	err := um.LoadFromFile("users.json")
	if err != nil {
		fmt.Printf("加载用户数据失败: %v\n", err)
	}

	// 添加一些用户
	usersToAdd := []struct {
		name  string
		email string
		age   int
	}{
		{"张三", "zhangsan@example.com", 25},
		{"李四", "lisi@example.com", 30},
		{"王五", "wangwu@example.com", 22},
		{"赵六", "zhaoliu@example.com", 28},
		{"钱七", "qianqi@example.com", 35},
	}

	for _, userData := range usersToAdd {
		err := um.AddUser(userData.name, userData.email, userData.age)
		if err != nil {
			fmt.Printf("添加用户失败: %v\n", err)
		}
	}

	// 尝试添加无效用户以演示错误处理
	fmt.Println("\n--- 测试错误处理 ---")
	um.AddUser("", "invalid@example.com", 25) // 空名称
	um.AddUser("测试用户", "invalid-email", 30) // 无效邮箱
	um.AddUser("测试用户", "test@example.com", -5) // 无效年龄

	// 列出所有用户
	fmt.Println("\n--- 所有用户 ---")
	allUsers := um.ListUsers()
	for _, user := range allUsers {
		status := "活跃"
		if !user.Active {
			status = "非活跃"
		}
		fmt.Printf("ID: %d, 姓名: %s, 邮箱: %s, 年龄: %d, 创建时间: %s, 状态: %s\n",
			user.ID, user.Name, user.Email, user.Age, user.Created.Format("2006-01-02 15:04:05"), status)
	}

	// 更新用户
	fmt.Println("\n--- 更新用户 ---")
	err = um.UpdateUser(1, "张三丰", "zhangsanfeng@example.com", 100)
	if err != nil {
		fmt.Printf("更新用户失败: %v\n", err)
	}

	// 搜索用户
	fmt.Println("\n--- 搜索用户(包含'张') ---")
	searchResults := um.SearchUsersByName("张")
	for _, user := range searchResults {
		fmt.Printf("找到用户: ID: %d, 姓名: %s\n", user.ID, user.Name)
	}

	// 停用用户
	fmt.Println("\n--- 停用用户 ---")
	err = um.DeactivateUser(2)
	if err != nil {
		fmt.Printf("停用用户失败: %v\n", err)
	}

	// 获取活跃用户
	fmt.Println("\n--- 活跃用户 ---")
	activeUsers := um.GetActiveUsers()
	for _, user := range activeUsers {
		fmt.Printf("活跃用户: ID: %d, 姓名: %s\n", user.ID, user.Name)
	}

	// 保存到文件
	fmt.Println("\n--- 保存到文件 ---")
	err = um.SaveToFile("users.json")
	if err != nil {
		fmt.Printf("保存文件失败: %v\n", err)
	}

	// 获取用户总数
	fmt.Printf("\n用户总数: %d\n", um.GetUserCount())

	// 交互式命令处理
	fmt.Println("\n--- 交互式命令 ---")
	fmt.Println("输入命令来管理用户:")
	fmt.Println("1. add <name> <email> <age> - 添加用户")
	fmt.Println("2. get <id> - 获取用户")
	fmt.Println("3. update <id> <name> <email> <age> - 更新用户")
	fmt.Println("4. delete <id> - 删除用户")
	fmt.Println("5. list - 列出所有用户")
	fmt.Println("6. search <name> - 搜索用户")
	fmt.Println("7. activate <id> - 激活用户")
	fmt.Println("8. deactivate <id> - 停用用户")
	fmt.Println("9. count - 获取用户总数")
	fmt.Println("10. active - 获取活跃用户")
	fmt.Println("11. save - 保存到文件")
	fmt.Println("12. quit - 退出程序")

	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("\n请输入命令: ")
		if !scanner.Scan() {
			break
		}
		input := strings.TrimSpace(scanner.Text())
		parts := strings.Fields(input)

		if len(parts) == 0 {
			continue
		}

		command := strings.ToLower(parts[0])
		switch command {
		case "quit", "exit":
			fmt.Println("退出程序")
			return
		case "add":
			if len(parts) != 4 {
				fmt.Println("用法: add <name> <email> <age>")
				continue
			}
			name := parts[1]
			email := parts[2]
			age := 0
			fmt.Sscanf(parts[3], "%d", &age)
			err := um.AddUser(name, email, age)
			if err != nil {
				fmt.Printf("添加用户失败: %v\n", err)
			} else {
				fmt.Println("用户添加成功")
			}
		case "get":
			if len(parts) != 2 {
				fmt.Println("用法: get <id>")
				continue
			}
			id := 0
			fmt.Sscanf(parts[1], "%d", &id)
			user, err := um.GetUserByID(id)
			if err != nil {
				fmt.Printf("获取用户失败: %v\n", err)
			} else {
				fmt.Printf("用户信息: ID=%d, Name=%s, Email=%s, Age=%d, Active=%t\n",
					user.ID, user.Name, user.Email, user.Age, user.Active)
			}
		case "update":
			if len(parts) != 5 {
				fmt.Println("用法: update <id> <name> <email> <age>")
				continue
			}
			id := 0
			fmt.Sscanf(parts[1], "%d", &id)
			name := parts[2]
			email := parts[3]
			age := 0
			fmt.Sscanf(parts[4], "%d", &age)
			err := um.UpdateUser(id, name, email, age)
			if err != nil {
				fmt.Printf("更新用户失败: %v\n", err)
			} else {
				fmt.Println("用户更新成功")
			}
		case "delete":
			if len(parts) != 2 {
				fmt.Println("用法: delete <id>")
				continue
			}
			id := 0
			fmt.Sscanf(parts[1], "%d", &id)
			err := um.DeleteUser(id)
			if err != nil {
				fmt.Printf("删除用户失败: %v\n", err)
			} else {
				fmt.Println("用户删除成功")
			}
		case "list":
			users := um.ListUsers()
			for _, user := range users {
				status := "活跃"
				if !user.Active {
					status = "非活跃"
				}
				fmt.Printf("ID: %d, 姓名: %s, 邮箱: %s, 年龄: %d, 状态: %s\n",
					user.ID, user.Name, user.Email, user.Age, status)
			}
		case "search":
			if len(parts) != 2 {
				fmt.Println("用法: search <name>")
				continue
			}
			name := parts[1]
			users := um.SearchUsersByName(name)
			for _, user := range users {
				fmt.Printf("找到用户: ID: %d, 姓名: %s\n", user.ID, user.Name)
			}
		case "activate":
			if len(parts) != 2 {
				fmt.Println("用法: activate <id>")
				continue
			}
			id := 0
			fmt.Sscanf(parts[1], "%d", &id)
			err := um.ActivateUser(id)
			if err != nil {
				fmt.Printf("激活用户失败: %v\n", err)
			} else {
				fmt.Println("用户激活成功")
			}
		case "deactivate":
			if len(parts) != 2 {
				fmt.Println("用法: deactivate <id>")
				continue
			}
			id := 0
			fmt.Sscanf(parts[1], "%d", &id)
			err := um.DeactivateUser(id)
			if err != nil {
				fmt.Printf("停用用户失败: %v\n", err)
			} else {
				fmt.Println("用户停用成功")
			}
		case "count":
			count := um.GetUserCount()
			fmt.Printf("用户总数: %d\n", count)
		case "active":
			users := um.GetActiveUsers()
			for _, user := range users {
				fmt.Printf("活跃用户: ID: %d, 姓名: %s\n", user.ID, user.Name)
			}
		case "save":
			err := um.SaveToFile("users.json")
			if err != nil {
				fmt.Printf("保存文件失败: %v\n", err)
			} else {
				fmt.Println("数据保存成功")
			}
		default:
			fmt.Println("未知命令，请重新输入")
		}
	}
}