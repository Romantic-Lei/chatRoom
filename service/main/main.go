package main
import (
	"fmt"
	"net"
	"time"
	"go_code/chatRoom/service/model"
)

// func readPkg(conn net.Conn) (mes message.Message, err error) {
// 	buf := make([]byte, 8096)
// 	fmt.Println("读取客户端发送的数据...")
// 	// conn.Read 在 conn没有被关闭的情况下，会阻塞
// 	// 如果客户端关闭了 conn 则不会阻塞
// 	_, err = conn.Read(buf[:4])
// 	if err != nil {
// 		// err = errors.New("read pkg header error")
// 		return 
// 	}
// 	// fmt.Println("读取到的buf =", buf[:4])

// 	// 根据 buf[:4] 转成一个 uint32 类型
// 	var pkgLen uint32
// 	pkgLen = binary.BigEndian.Uint32(buf[0:4])
// 	// 根据 pkgLen 读取消息内容
// 	n, err := conn.Read(buf[:pkgLen])
// 	if n != int(pkgLen) || err != nil {
// 		// err = errors.New("read pkg body error")
// 		return 
// 	}
// 	// 把pkgLen 反序列化成 -> message.Message
// 	err = json.Unmarshal(buf[:pkgLen], &mes)
// 	if err != nil {
// 		fmt.Println("json.Unmarshal err = ", err)
// 		return 
// 	}
// 	return 
// }

// func writePkg(conn net.Conn, data []byte) (err error) {
// 	// 先发送一个长度给对方
// 	var pkgLen uint32
// 	pkgLen = uint32(len(data))
// 	var buf [4]byte
// 	binary.BigEndian.PutUint32(buf[0:4], pkgLen)
// 	// 发送长度
// 	n, err := conn.Write(buf[:4])
// 	if n != 4 || err != nil {
// 		fmt.Println("conn.Write(buf) fail", err)
// 		return 
// 	}

// 	// 发送data 本身
// 	n, err = conn.Write(data)
// 	if n != int(pkgLen) || err != nil {
// 		fmt.Println("conn.Write(buf) fail", err)
// 		return 
// 	}
// 	return 
// }

// // 函数serverProcessLogin函数，专门处理登录请求
// func ServerProcessLogin(conn net.Conn, mes *message.Message) (err error) {
// 	// 核心代码
// 	// 1.先从mes 中取出 mes.Data, 并直接反序列化成LoginMes
// 	var loginMes message.LoginMes
// 	err = json.Unmarshal([]byte(mes.Data), &loginMes)
// 	if err != nil {
// 		fmt.Println("json.Unmarshal fail err", err)
// 		return 
// 	}

// 	// 1. 先声明一个 resMes
// 	var resMes message.Message
// 	resMes.Type = message.LoginResMesType

// 	// 2. 在声明一个LoginResMes
// 	var loginResMes message.LoginResMes

// 	// 假定用户id = 100， 密码 = 123456， 则我们认为合法
// 	if loginMes.UserId == 100 && loginMes.UserPwd == "123456" {
// 		// 合法
// 		loginResMes.Code = 200
// 	} else {
// 		// 不合法
// 		loginResMes.Code = 500 // 500我们自定为用户不存在
// 		loginResMes.Error = "该用户不存在，请确认用户名是否输入正确"
// 	}

// 	// 3. 将 loginResMes 序列化
// 	data, err := json.Marshal(loginResMes)
// 	if err != nil {
// 		fmt.Println("json.Marshal fail", err)
// 		return 
// 	}

// 	// 4. 将 data 赋值给 resMes
// 	resMes.Data = string(data)

// 	// 5. 对 resMes 进行序列化， 准备发送
// 	data, err = json.Marshal(resMes)
// 	if err != nil {
// 		fmt.Println("json.Marshal fail", err)
// 		return 
// 	}

// 	// 6. 发送 data， 将其封装到 writePkg 函数中
// 	err = writePkg(conn, data)
// 	return 
// }

// //  编写一个serverProcessMes 函数
// // 根据客户端发送消息种类的不同，确定调用哪个函数来处理
// func serverProcessMes(conn net.Conn, mes *message.Message) (err error) {
// 	switch mes.Type {
// 		case message.LoginMesType : 
// 			// 处理登录逻辑
// 			err = ServerProcessLogin(conn, mes)
// 		case message.RegisterMesType :
// 			// 处理注册逻辑
// 		default :
// 			fmt.Println("消息类型不存在，处理失败")
// 	}
// 	return 
// }

// 处理和客户端的通信
func process(conn net.Conn) {
	// 这里需要延时关闭conn
	defer conn.Close()
	
	// 这里调用总控，创建一个实例
	processor := &Processor{
		Conn : conn,
	}
	// 本类和 process2() 在同一个包里面，所以不需要大写
	err := processor.process2()
	if err != nil {
		fmt.Println("客户端和服务器通讯协程错误 = err", err)
		return 
	}
}

// 这里我们编写一个函数， 完成对 UserDao的初始化任务
func initUserDao() {
	// 这里的pool 本身就是一个全局变量，存放在本包下面的redis.go里面
	// 注意这里我们必须要先初始化initpool，不然我们的pool就是一个空的
	model.MyUserDao = model.NewUserDao(pool)
}

func main() {
	// 当服务器启动时，我们就去初始化我们的redis 的连接池
	// 不同的机器需要改成不同的ip
	initPool("192.168.0.110:6379", 16, 0, time.Second)
	initUserDao()
	// 提示信息
	fmt.Println("服务器[新的结构]在8889 端口监听...")
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