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