package flags

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strings"

	"blog/internal/global"
	"blog/internal/model"
	"blog/internal/model/enum"
	"blog/internal/utils/pwd"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type FlagUser struct{}

// Create 创建用户
func (FlagUser) Create() {
	if global.DB == nil {
		logrus.Error("数据库连接未初始化，无法创建用户")
		return
	}

	reader := bufio.NewReader(os.Stdin)

	role, ok := selectRole(reader)
	if !ok {
		return
	}

	username, ok := promptInput(reader, "请输入用户名:")
	if !ok {
		return
	}
	username = strings.TrimSpace(username)
	if username == "" {
		logrus.Error("用户名不能为空")
		return
	}

	// 检查是否已存在
	var exist model.UserModel
	err := global.DB.Take(&exist, "username = ?", username).Error
	if err == nil {
		logrus.Errorf("用户名 %s 已存在", username)
		return
	}
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		logrus.Errorf("查询用户失败: %v", err)
		return
	}

	password, ok := promptInput(reader, "请输入密码:")
	if !ok {
		return
	}
	confirm, ok := promptInput(reader, "请再次输入密码:")
	if !ok {
		return
	}
	if password != confirm {
		logrus.Error("两次密码不一致")
		return
	}

	if password == "" {
		logrus.Error("密码不能为空")
		return
	}

	hashed, err := pwd.GenerateFromPassword(password)
	if err != nil {
		logrus.Errorf("密码加密失败: %v", err)
		return
	}

	user := model.UserModel{
		Username:  username,
		Password:  hashed,
		Avatar:    "",
		RegSource: enum.RegTerminalSourceType,
		Role:      role,
	}

	if err := global.DB.Create(&user).Error; err != nil {
		logrus.Errorf("创建用户失败: %v", err)
		return
	}
	logrus.Infof("创建用户成功，ID=%d 用户名=%s 角色=%d", user.ID, user.Username, user.Role)
}

func promptInput(reader *bufio.Reader, msg string) (string, bool) {
	fmt.Println(msg)
	text, err := reader.ReadString('\n')
	if err != nil {
		logrus.Errorf("读取输入失败: %v", err)
		return "", false
	}
	return strings.TrimSpace(text), true
}

func selectRole(reader *bufio.Reader) (enum.RoleType, bool) {
	fmt.Println("请选择角色：1 普通用户  2 管理员")
	text, err := reader.ReadString('\n')
	if err != nil {
		logrus.Errorf("读取角色失败: %v", err)
		return 0, false
	}
	switch strings.TrimSpace(text) {
	case "1":
		return enum.RoleUserType, true
	case "2":
		return enum.RoleAdminType, true
	default:
		logrus.Error("角色输入无效，只能为 1 或 2")
		return 0, false
	}
}
