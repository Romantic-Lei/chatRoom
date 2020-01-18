package model
import (
	"fmt"
	"github.com/garyburd/redigo/redis"
	"encoding/json"
)

// 我们在服务器启动后，就初始化一个 userDao实例，
// 把他做成全局的变量，在需要和redis 操作时，就直接使用即可
var (
	MyUserDao *UserDao
)

// 定义一个 userDao 结构体
// 定义对 User 结构体的各种操作
type UserDao struct {
	pool *redis.Pool
}

// 使用工厂模式，创建一个UserDao 实例
func NewUserDao(pool *redis.Pool) (userDao *UserDao) {
	userDao = &UserDao{
		pool : pool,
	}
	return userDao
}

// 在userDao 中，应该有哪些提供给我们
// 1. 根据用户id 返回 一个 User实例和 err
func (this *UserDao) getUserById(conn redis.Conn, id int) (user *User, err error) {
	// 通过给定 id去redis查询这个用户
	res, err := redis.String(conn.Do("HGet", "users", id))
	if err != nil {
		// 错误
		if err == redis.ErrNil {
			// 此错误表示在 users哈希中找不到这样的id
			err = ERROR_USER_NOTEXISTS
		}
		return 
	}

	user = &User{}
	// 这里我们需要把 res 反序列化成一个User实例
	err = json.Unmarshal([]byte(res), user)
	if err != nil {
		fmt.Println("json.Unmarshal err = ", err)
		return 
	}
	return 
}

// 完成登录的校验
// 1. Login 完成对用户的验证
// 2. 如果用户的id 和pwd 都正确，则返回一个user实例
// 3. 如果用户的 id 或者 pwd 有错误，则返回对应的错误信息
func (this *UserDao) Login(userId int, userPwd string) (user *User, err error) {
	// 先从UserDao的连接池中取出一个连接
	conn := this.pool.Get()
	defer conn.Close()
	user, err = this.getUserById(conn, userId)
	if err != nil {
		fmt.Println("this.getUserById err = ", err)
		return 
	}

	// 这是可以证明用户是可以获取到的。我们接下来需要校验密码是正确的
	fmt.Println("user.UserPwd = ", user.UserPwd, "userPwd = ", userPwd)
	if user.UserPwd != userPwd {
		err = ERROR_USER_PWD
		return 
	}
	return 
}
