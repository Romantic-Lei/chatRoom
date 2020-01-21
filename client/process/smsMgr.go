package process
import (
	"fmt"
	"encoding/json"
	"go_code/chatRoom/common/message"
)

func outputGroupMes(mes *message.Message) {
	// 显示消息
	// 1. 反序列化
	var smsMes message.SmsMes
	err := json.Unmarshal([]byte(mes.Data), &smsMes)
	if err != nil {
		fmt.Println("json.Unmarshal err = ", err.Error())
		return 
	}
	// 显示
	info := fmt.Sprintf("用户id : \t %d 对大家说：\t %s", smsMes.UserId, smsMes.Content)
	fmt.Println(info)
	fmt.Println()
}