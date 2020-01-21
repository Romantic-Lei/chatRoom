package process2
import (
	"fmt"
	"net"
	"encoding/json"
	"go_code/chatRoom/common/message"
	"go_code/chatRoom/service/utils"
)

type SmsProcess struct {

}

// 写一个方法转发消息
func (this *SmsProcess) SendGroupMes(mes *message.Message) {
	// 遍历服务器端的 onlineUsers map[int]*UserProcess
	// 将消息转发出去
	// 取出 mes 的内容 SmsMes
	var smsMes message.SmsMes
	err := json.Unmarshal([]byte(mes.Data), &smsMes)
	if err != nil {
		fmt.Println("json.Unmarshal err =", err)
		return 
	}

	data, err := json.Marshal(mes)
	if err != nil {
		fmt.Println("json.Marshal err = ", err)
		return 
	}
	for id, up := range userMgr.onlineUsers {
		// 这里我们还需要过滤掉自己，即不要发送给自己
		if id == smsMes.UserId {
			continue
		}
		this.sendMesToEachOnlineUser(data, up.Conn)
	}
}

func (this *SmsProcess) sendMesToEachOnlineUser(data []byte, conn net.Conn) {

	// 创建一个 Transfer 实例， 发送data
	tf := &utils.Transfer{
		Conn : conn,
	}
	err := tf.WritePkg(data)
	if err != nil {
		fmt.Println("转发消息失败 err = ", err)
	}
} 