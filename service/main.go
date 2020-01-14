package main
import (
	"fmt"
	"net"
	"go_code/chatroom/common/message"
	"encoding/binary"
	"encoding/json"
	_"errors"
	"io"
)

func readPkg(conn net.Conn) (mes message.Message, err error) {
	buf := make([]byte, 8096)
	fmt.Println("读取客户端发送的数据...")
	// conn.Read 在 conn没有被关闭的情况下，会阻塞
	// 如果客户端关闭了 conn 则不会阻塞
	_, err = conn.Read(buf[:4])
	if err != nil {
		// err = errors.New("read pkg header error")
		return 
	}
	// fmt.Println("读取到的buf =", buf[:4])

	// 根据 buf[:4] 转成一个 uint32 类型
	var pkgLen uint32
	pkgLen = binary.BigEndian.Uint32(buf[0:4])
	// 根据 pkgLen 读取消息内容
	n, err := conn.Read(buf[:pkgLen])
	if n != int(pkgLen) || err != nil {
		// err = errors.New("read pkg body error")
		return 
	}
	// 把pkgLen 反序列化成 -> message.Message
	err = json.Unmarshal(buf[:pkgLen], &mes)
	if err != nil {
		fmt.Println("json.Unmarshal err = ", err)
		return 
	}
	return 
}

// 函数serverProcessLogin函数，专门处理登录请求
func ServerProcessLogin(conn net.Conn, mes *message.Message) (err error) {
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
	resMes.Type = message.LoginResType

	// 2. 在声明一个LoginResMes
	var loginResMes message.LoginResMesType

	// 假定用户id = 100， 密码 = 123456， 则我们认为合法
	if loginMes.UserId == 100 && loginMes.UserPwd == "123456" {
		// 合法
		loginResMes.Code = 200
	} else {
		// 不合法
		loginResMes.Code = 500 // 500我们自定为用户不存在
		loginResMes.Error = "该用户不存在，请确认用户名是否思路正确"
	}

	// 3. 将 loginResMes 序列化
	data, err := json.Marshal(loginResMes)
	if err != nil {
		fmt.Println("json.Marshal fail", err)
		return 
	}

	// 4. 将 data 赋值给 resMes
	resMes.Data = string(data)
}

//  编写一个ServerProcessMes 函数
// 根据客户端发送消息种类的不同，确定调用哪个函数来处理
func ServerProcessMes(conn net.Conn, mes *message.Message) (err error) {
	switch mes.Type {
		case message.LoginMesType : 
			// 处理登录逻辑
			err = ServerProcessLogin(conn, mes)
		case message.RegisterMesType :
			// 处理注册逻辑
		default :
			fmt.Println("消息类型不存在，处理失败")
	}
}

// 处理和客户端的通信
func process(conn net.Conn) {
	// 这里需要延时关闭conn
	defer conn.Close()
	// 循环读取客户端发送的消息
	for {
		// 这里我们将读取数据包，直接封装成一个函数readPkg()
		mes, err := readPkg(conn)
		if err != nil {
			if err == io.EOF {
				fmt.Println("客户端已退出，服务器端也退出")
				return 
			} else {
				fmt.Println("readPkg err = ", err)
				return 
			}
		}

		fmt.Println("mes =", mes)
		
	}
}

func main() {
	// 提示信息
	fmt.Println("服务器在8889 端口监听...")
	listen, err := net.Listen("tcp", "127.0.0.1:8889")
	defer listen.Close()

	if err != nil {
		fmt.Println("net.Listen err =", err)
		return 
	}
	// 监听成功，等待客户端连接服务器
	for {
		fmt.Println("等待客户端来链接服务器...")
		conn, err := listen.Accept()
		if err != nil {
			fmt.Println("listen.Accept err =", err)

		}
		// 一旦链接成功，则启动一个协程和客户端保存通讯
		go process(conn)
	}
}