package message

// 定义一些常量
const (
	LoginMesType       = "LoginMes"
	LoginResMesType    = "LoginResMes"
	RegisterMesType    = "RegisterMes"
	RegisterResMesType = "RegisterResMes"
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