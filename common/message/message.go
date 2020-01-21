package message

// 定义一些常量
const (
	LoginMesType        	 = "LoginMes"
	LoginResMesType   	 	 = "LoginResMes"
	RegisterMesType   		 = "RegisterMes"
	RegisterResMesType		 = "RegisterResMes"
	NotifyUserStatusMesType  = "NotifyUserStatusMes"
	SmsMesType				 = "SmsMes"
)

// 定义几个用户状态的常量
const (
	UserOnline = iota // 在线 iota -> 自动向下递增
	UserOffline // 离线
	UserBusyStatus // 繁忙
)

type Message struct {
	Type string `json:"type"` // 消息的类型
	Data string `json:"data"`// 消息的内容
}

// 登录用 ↓
// 定义两个消息..后面需要在增加
type LoginMes struct {
	UserId int `json:"userid"` // 用户id
	UserPwd string `json:"userPwd"` // 用户密码
	UserName string `json:"userName"` // 用户名
}

type LoginResMes struct {
	Code int `json:"code"` // 返回的状态码，我们自定义500 表示该用户未注册 200表示登录成功
	UsersId []int // 增加字段，保存用户 id 的切片
	Error string `json:"error"` // 返回错误信息
}

// 注册用 ↓
type RegisterMes struct {
	User User `json:"user"`// 类型就是User结构体
}

type RegisterResMes struct{
	Code int `json:"code"` // 返回的状态码，我们自定义400 表示该用户已占用 200表示登录成功
	Error string `json:"error"` // 返回错误信息
}

// 为了配合服务器端推送用户状态变化的消息
type NotifyUserStatusMes struct {
	UserId int `json:"userId"` // 用户id
	Status int `json:"status"` // 用户状态
}

// 增加一个SmsMes  发送的消息
type SmsMes struct {
	Content string `json:"contene"`
	User // 匿名结构体， 好似继承的关系
}