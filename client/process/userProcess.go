package process
import (
	"fmt"
	"net"
	"go_code/chatRoom/common/message"
	"go_code/chatRoom/client/utils"
	"encoding/json"
	"encoding/binary"
	_"errors"
)

type UserProcess struct {
	// 暂时不需要字段

}

// 给关联一个用户登录的方法
// 写一个函数，完成登陆
func (this *UserProcess) Login(userId int, userPwd string) (err error) {
	// 下一步就开始定协议...
	// fmt.Printf("userId = %d userPwd = %s", userId, userPwd)
	// return nil
	
	// 1. 链接到服务器
	conn, err := net.Dial("tcp", "127.0.0.1:8889")
	if err != nil {
		fmt.Println("net.Dial err =", err)
		return 
	}
	// 延时关闭
	defer conn.Close()

	// 2.准备通过conn发送消息给服务器
	var mes message.Message
	mes.Type = message.LoginMesType
	// 3. 创建一个LoginMes 结构体
	var loginMes message.LoginMes
	loginMes.UserId = userId
	loginMes.UserPwd = userPwd

	// 4.将loginMes 序列化
	data, err := json.Marshal(loginMes)
	if err != nil {
		fmt.Println("json.Marshal err =", err)
		return 
	}
	// 5. 把data赋给 mes.Data字段
	mes.Data = string(data)

	// 6.将 mes进行序列化
	data, err = json.Marshal(mes)
	if err != nil {
		fmt.Println("json.Marshal err =", err)
		return 
	}

	// 7. 到这个时候 data就是我们要发送的消息
	// 7.1 先把 data的长度发送给服务器
	// 先获取到 data 的长度，然后转成一个表示长度的切片
	var pkgLen uint32
	pkgLen = uint32(len(data))
	var buf [4]byte
	binary.BigEndian.PutUint32(buf[0:4], pkgLen)
	// 发送长度
	n, err := conn.Write(buf[:4])
	if n != 4 || err != nil {
		fmt.Println("conn.Write(buf) fail", err)
		return 
	}

	fmt.Printf("客服端发送消息长度 =%d, 内容 =%s", len(data), string(data))

	// 发送消息本身
	_, err = conn.Write(data)
	if err != nil {
		fmt.Println("conn.Write(data) fail", err)
		return 
	}

	// 这里还需要处理服务器端返回的消息
	// 创建一个Transfer 实例
	tf := &utils.Transfer {
		Conn : conn,
	}
	mes, err = tf.ReadPkg() // mes 就是
	if err != nil {
		fmt.Println("readPkg fail", err)
		return
	}
	// 将 mes 的data部分反序列化成LoginResMes
	var loginResMes message.LoginResMes
	err = json.Unmarshal([]byte(mes.Data), &loginResMes)
	if loginResMes.Code == 200 {
		// fmt.Println("登录成功!")

		// 这里我们还需要在客户端启动一个协程
		// 该协程保持和服务器端的通讯，如果服务器有数据推送给客户端
		// 则接受并显示在客户端的终端
		go serverProcessMes(conn)
		// 1. 显示我们登录成功的菜单[我们需要循环显示]
		for {
			ShowMenu()
		}

	} else {
		fmt.Println(loginResMes.Error)
	}
	return 
}