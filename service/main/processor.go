package main
import (
	"fmt"
	"net"
	"io"
	"go_code/chatRoom/common/message"
	"go_code/chatRoom/service/utils"
	"go_code/chatRoom/service/process"
)

// 先创建一个Processor 的结构体
type Processor struct {
	Conn net.Conn
}

//  编写一个serverProcessMes 函数
// 根据客户端发送消息种类的不同，确定调用哪个函数来处理
func (this *Processor) serverProcessMes(mes *message.Message) (err error) {
	switch mes.Type {
		case message.LoginMesType : 
			// 处理登录逻辑
			// 创建一个 UserProcess实例
			up := &process2.UserProcess{
				Conn : this.Conn,
			}
			err = up.ServerProcessLogin(mes)
		case message.RegisterMesType :
			// 处理注册逻辑
		default :
			fmt.Println("消息类型不存在，处理失败")
	}
	return 
}

func (this *Processor) process2() (err error) {
	// 循环读取客户端发送的消息
	
	for {
		// 这里我们将读取数据包，直接封装成一个函数readPkg(),返回Message， err
		// 创建一个Transfer 实例，完成读包的任务
		tf := &utils.Transfer {
			Conn : this.Conn,
		}
		mes, err := tf.ReadPkg()
		if err != nil {
			if err == io.EOF {
				fmt.Println("客户端已退出，服务器端也退出")
				return err
			} else {
				fmt.Println("readPkg err = ", err)
				return err
			}
		}

		err = this.serverProcessMes(&mes)
		if err != nil {
			return err
		}
	}
}