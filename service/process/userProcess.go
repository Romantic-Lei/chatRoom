package process2
import (
	"fmt"
	"net"
	"go_code/chatRoom/common/message"
	"go_code/chatRoom/service/utils"
	"go_code/chatRoom/service/model"
	"encoding/json"
)

type UserProcess struct {
	Conn net.Conn
	// 表示该 Conn是哪个用户的
	UserId int
}

// 通知所有在线用户的方法
// userId 要通知其他的在线用户
func (this *UserProcess) NotifyOthersOnlineUser(userId int) {
	// 遍历 onlineUsers,然后一个一个的发送 NotifyUserStatusMes
	for id, up := range userMgr.onlineUsers {
		// 过滤掉自己
		if id == userId {
			continue
		}
		// 开始通知(单独的写一个方法)
		up.NotifyMeOnlie(userId)
	}
}

func (this *UserProcess) NotifyMeOnlie(userId int) {
	// 组装我们的 NotifyUserStatusMes 消息
	var mes message.Message
	mes.Type = message.NotifyUserStatusMesType

	var NotifyUserStatusMes message.NotifyUserStatusMes
	NotifyUserStatusMes.UserId = userId
	NotifyUserStatusMes.Status = message.UserOnline

	// 将 notifyUserStatusMes序列化
	data, err := json.Marshal(NotifyUserStatusMes)
	if err != nil {
		fmt.Println("json.Marshal err", err)
		return
	}
	// 将序列化后的 notifyUserStatusMes 赋值给 mes.Data
	mes.Data = string(data)

	// 对mes 再次序列化，准备发送
	data, err = json.Marshal(mes)
	if err != nil {
		fmt.Println("json.Marshal err", err)
		return
	}

	// 发送，创建一个 Transfer 实例，发送
	tf := &utils.Transfer {
		Conn : this.Conn,
	}

	err = tf.WritePkg(data)
	if err != nil {
		fmt.Println("NotifyMeOnlie err = ", err)
		return 
	}

}

// 注册用户的方法
func (this *UserProcess) ServerProcessRegister(mes *message.Message) (err error) {
	// 1.先从mes 中取出 mes.Data, 并直接反序列化成registerMes
	var registerMes message.RegisterMes
	err = json.Unmarshal([]byte(mes.Data), &registerMes)
	if err != nil {
		fmt.Println("json.Unmarshal fail err", err)
		return 
	}

	// 1. 先声明一个 resMes
	var resMes message.Message
	resMes.Type = message.RegisterMesType

	// 2. 在声明一个RegisterResMes
	var registerResMes message.RegisterResMes

	// 我们需要到 redis数据库去完成注册
	// 1.使用 model.MyUserDao 到redis验证
	err = model.MyUserDao.Register(&registerMes.User)

	if err != nil {
		if err == model.ERROR_USER_EXISTS {
			registerResMes.Code = 505
			registerResMes.Error = model.ERROR_USER_EXISTS.Error()
		} else {
			registerResMes.Code = 506
			registerResMes.Error = "注册发生未知错误..."
		}
	} else {
		registerResMes.Code = 200
	}
	// 3. 将 registerResMes 序列化
	data, err := json.Marshal(registerResMes)
	if err != nil {
		fmt.Println("json.Marshal fail", err)
		return 
	}
	// 4. 将 data 赋值给 resMes
	resMes.Data = string(data)

	// 5. 对 resMes 进行序列化， 准备发送
	data, err = json.Marshal(resMes)
	if err != nil {
		fmt.Println("json.Marshal fail", err)
		return 
	}

	// 6. 发送 data， 将其封装到 writePkg 函数中
	// 因为使用了分层模式（mvc），我们先创建一个 Transfer 实例，然后读取
	tf := &utils.Transfer{
		Conn: this.Conn,

	}
	err = tf.WritePkg(data)
	return 
}

// 函数serverProcessLogin函数，专门处理登录请求
func (this *UserProcess) ServerProcessLogin(mes *message.Message) (err error) {
	// 核心代码
	// 1.先从mes 中取出 mes.Data, 并直接反序列化成LoginMes
	var loginMes message.LoginMes
	err = json.Unmarshal([]byte(mes.Data), &loginMes)
	if err != nil {
		fmt.Println("json.Unmarshal fail err", err)
		return 
	}

	// 1. 先声明一个 resMes
	var resMes message.Message
	resMes.Type = message.LoginResMesType

	// 2. 在声明一个LoginResMes
	var loginResMes message.LoginResMes
	// 我们需要到 redis数据库去完成校验
	// 1.使用 model.MyUserDao 到redis验证
	user, err := model.MyUserDao.Login(loginMes.UserId, loginMes.UserPwd)
	
	if err != nil {
		if err == model.ERROR_USER_NOTEXISTS {
			loginResMes.Code = 500
			loginResMes.Error = err.Error()
		} else if err == model.ERROR_USER_PWD {
			loginResMes.Code = 403
			loginResMes.Error = err.Error()
		} else {
			loginResMes.Code = 505
			loginResMes.Error = "服务器内部错误"
		}
		// 这里我们先测试成功然后我们根据不听
		} else {
			loginResMes.Code = 200
			// 这里，因为用户登录成功，我们就把该登录成功的用户放到 userMgr中
			// 将登录成功的用户userId 赋值给 this
			this.UserId = loginMes.UserId
			userMgr.AddOnlineUser(this)
			// 通知其他在线的用户，当前用户上线了
			this.NotifyOthersOnlineUser(loginMes.UserId)
			// 将当前在线用户的id， 放入到 loginResMes.UsersId
			// 遍历 UserMgr.onlineUsers 
			for id, _ := range userMgr.onlineUsers {
				loginResMes.UsersId = append(loginResMes.UsersId, id)
			}

			fmt.Println(user, "登录成功")
	}

	// // 假定用户id = 100， 密码 = 123456， 则我们认为合法
	// if loginMes.UserId == 100 && loginMes.UserPwd == "123456" {
	// 	// 合法
	// 	loginResMes.Code = 200
	// } else {
	// 	// 不合法
	// 	loginResMes.Code = 500 // 500我们自定为用户不存在
	// 	loginResMes.Error = "该用户不存在，请确认用户名是否输入正确"
	// }

	// 3. 将 loginResMes 序列化
	data, err := json.Marshal(loginResMes)
	if err != nil {
		fmt.Println("json.Marshal fail", err)
		return 
	}

	// 4. 将 data 赋值给 resMes
	resMes.Data = string(data)

	// 5. 对 resMes 进行序列化， 准备发送
	data, err = json.Marshal(resMes)
	if err != nil {
		fmt.Println("json.Marshal fail", err)
		return 
	}

	// 6. 发送 data， 将其封装到 writePkg 函数中
	// 因为使用了分层模式（mvc），我们先创建一个 Transfer 实例，然后读取
	tf := &utils.Transfer{
		Conn: this.Conn,

	}
	err = tf.WritePkg(data)
	return 
}