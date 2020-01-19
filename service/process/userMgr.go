package process2
import (
	"fmt"
)

// 因为 UserMgr 实例在服务器端有且只有一个
// 因为在很多的地方都会使用到，因此我们将其定义为全局变量
var (
	userMgr *UserMgr
)

type UserMgr struct {
	onlineUsers map[int]*UserProcess
}

// 完成对 userMgr 初始化工作
func init() {
	userMgr = &UserMgr{
		onlineUsers : make(map[int]*UserProcess, 1024),
	}
}

// 完成对onlineUsers添加
func (this *UserMgr) AddOnlineUser(up *UserProcess) {
	this.onlineUsers[up.UserId] = up
}

// 删除
func (this *UserMgr) DelOnlineUser(userID int) {
	delete(this.onlineUsers, userID)
}

// 返回当前所有在线用户
func (this *UserMgr) GetAllOnlineUser() map[int]*UserProcess {
	return this.onlineUsers
}

// 根据id 查询对应的值
func (this *UserMgr) DelOnlineUsers(userID int) (up *UserProcess, err error) {
	// 如何从map中取出一个值，带检测方式
	up, ok := this.onlineUsers[userID]
	if !ok {
		// 说明需要查找的用户当前不在线
		err = fmt.Errorf("用户 %v 不在线", userID)
		return 
	}
	return 
}