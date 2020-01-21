package process
import (
	"fmt"
	"go_code/chatRoom/common/message"
	"go_code/chatRoom/client/model"
)

// 客户端要维护的 map 
var onlineUsers map[int]*message.User = make(map[int]*message.User, 10)
var CurUser model.CurUser // 用户在登录成功后，完成对 curUser进行初始化

// 在客户端显示当前在线的用户
func outputOnlineUser() {
	// 遍历一把 onlineUsers
	fmt.Println("当前在线用户列表")
	for id, _ := range onlineUsers {
		fmt.Println("用户id: \t", id)
	}
}

// 编写一个方法，处理返回的NotifyUserStatusMes
func updateUserStatus(notifyUserStatusMes *message.NotifyUserStatusMes) {

	// 适当的优化一下
	user, ok := onlineUsers[notifyUserStatusMes.UserId]
	if !ok {
		// 原来没有的情况下
		user = &message.User{
			UserId : notifyUserStatusMes.UserId,
		}
	}
	user.UserStatus = notifyUserStatusMes.Status
	onlineUsers[notifyUserStatusMes.UserId] = user
	outputOnlineUser()
}