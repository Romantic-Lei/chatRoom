package process

import (
	"fmt"
	"os"
	"net"
	"encoding/json"
	"go_code/chatRoom/client/utils"
	"go_code/chatRoom/common/message"
)

// 先试试登录成功之后的界面
func ShowMenu() {
	fmt.Println("---------恭喜您，登录成功!---------")
	fmt.Println("---------1. 显示在线用户列表---------")
	fmt.Println("---------2. 发送消息---------")
	fmt.Println("---------3. 信息列表---------")
	fmt.Println("---------4. 退出系统---------")
	fmt.Println("请选择(1 - 4): ")
	var key int 
	fmt.Scanf("%d\n", &key)
	switch key {
		case 1:
			// fmt.Println("显示在线用户列表")
			outputOnlineUser()
		case 2:
			fmt.Println("发送消息")
		case 3:
			fmt.Println("信息列表")
		case 4:
			fmt.Println("您已选择退出系统")
			os.Exit(0) 
		default :
			fmt.Println("您输入的选项不正确")
	}
}

// 和服务器端保持通讯
func serverProcessMes(conn net.Conn) {
	// 创建一个transfer实例，不停读取服务器发送的消息
	tf := &utils.Transfer{
		Conn : conn,
	}
	for {
		fmt.Println("客户端正在等待读取服务器发送的消息" )
		mes, err := tf.ReadPkg()
		if err != nil {
			fmt.Println("tf.ReadPkg() err = ", err)
			return 
		}
		// 如果读到消息，则到下一步处理逻辑
		switch mes.Type {
			case message.NotifyUserStatusMesType : // 有用户上线
				// 处理逻辑
				// 1. 取出 NotifyUserStatusMesType 
				var notifyUserStatusMes message.NotifyUserStatusMes
				json.Unmarshal([]byte(mes.Data), &notifyUserStatusMes)
				// 2. 把用户的信息，状态保存到客户 map 中
				updateUserStatus(&notifyUserStatusMes)
			default :
				fmt.Println("服务器端返回了一个未知的消息类型")
		}
		// fmt.Printf("mes = %v\n", mes)
	}
}