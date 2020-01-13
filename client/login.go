package main
import (
	"fmt"
	"encoding/json"
	"encoding/binary"
	"net"
	"go_code/chatRoom/common/message"
)

// 写一个函数，完成登陆
func login(userId int, userPwd string) (err error) {
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

	// fmt.Printf("客服端发送消息长度 =%d, 内容 =%s", len(data), string(data))

	// 发送消息本身
	_, err = conn.Write(data)
	if err != nil {
		fmt.Println("conn.Write(data) fail", err)
		return 
	}

	// 这里还需要处理服务器端返回的消息


	return 
}