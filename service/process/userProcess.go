package process2
import (
	"fmt"
	"net"
	"go_code/chatRoom/common/message"
	"go_code/chatRoom/service/utils"
	"encoding/json"
)

type UserProcess struct {
	Conn net.Conn
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

	// 假定用户id = 100， 密码 = 123456， 则我们认为合法
	if loginMes.UserId == 100 && loginMes.UserPwd == "123456" {
		// 合法
		loginResMes.Code = 200
	} else {
		// 不合法
		loginResMes.Code = 500 // 500我们自定为用户不存在
		loginResMes.Error = "该用户不存在，请确认用户名是否输入正确"
	}

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