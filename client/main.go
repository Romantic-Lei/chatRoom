package main
import (
	"fmt"
	"os"
)

// 定义两个变量，一个表示用户id，一个表示用户密码
var userId int
var	userPwd string

func main() {
	// 接收用户的选择
	var key int
	// 判断是否还继续显示菜单
	var loop = true
	for loop {
		fmt.Println("-----------------------欢迎登陆本聊天系统-----------------------")
		fmt.Println("\t\t\t 1 登陆聊天室")
		fmt.Println("\t\t\t 2 用户注册")
		fmt.Println("\t\t\t 3 退出系统")
		fmt.Println("\t\t\t 请选择(1-3)")

		fmt.Scanf("%d\n", &key)
		switch key {
			case 1 : 
				fmt.Println("登陆聊天室")
				loop = false
			case 2 : 
				fmt.Println("用户注册")
				loop = false
			case 3 : 
				fmt.Println("退出系统")
				os.Exit(0)
				loop = false
			default :
				fmt.Println("你输入有误，请重新输入")
		}
	}
	// 根据用户的输入，显示新的提示信息
	if key == 1 {
		// 说明用户要登陆
		fmt.Println("请输入用户的id")
		fmt.Scanf("%d\n", &userId)
		fmt.Println("请输入用户的密码")
		fmt.Scanf("%s\n", &userPwd)
		// 先把登陆的函数，写到另外一个文件login.go,
		// 因为他们在同一个包下，所以login小写我们也可以访问到
		err := login(userId, userPwd)
		if err != nil {
			fmt.Println("登陆失败")
		} else {
			fmt.Println("登陆成功")
		}
	}else if key == 2 {
		fmt.Println("用户注册")
	}
}