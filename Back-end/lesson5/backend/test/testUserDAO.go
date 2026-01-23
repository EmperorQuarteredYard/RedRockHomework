package testFunction

import (
	"fmt"
	"lesson5/backend/DAO"
)

type UserDAO struct {
	dao *DAO.UserDAO
}

func NewTestUserDAO(dao *DAO.UserDAO) *UserDAO {
	return &UserDAO{dao: dao}
}
func (u UserDAO) Test() {
	u.TestRead()
	u.TestCreate()
	u.TestRead()
	u.TestSearch()
	u.TestDelete()
	u.TestRead()
}
func (u UserDAO) TestRead() {
	fmt.Println("----------------------------------")
	users, err := u.dao.ListUsers()
	if err != nil {
		fmt.Println("读取User发生错误：" + err.Error())
		return
	}
	if users == nil || len(*users) == 0 {
		fmt.Println("读取User发生错误：空值")
		return
	}
	fmt.Println("读取Users成功：")
	fmt.Println(*users)
}

func (u UserDAO) TestCreate() {
	fmt.Println("----------------------------------")
	err := u.dao.CreateUser("zhangsan1", "lisi666", "administrator")
	if err != nil {
		fmt.Println("创建User发生错误:" + err.Error())
		return
	}
	fmt.Println("创建Users成功：")

}
func (u UserDAO) TestDelete() {
	fmt.Println("----------------------------------")
	err := u.dao.DeleteUserByName("zhangsan1")
	if err != nil {
		fmt.Println("删除User发生错误:" + err.Error())
		return
	}
	fmt.Println("删除Users成功：")
}
func (u UserDAO) TestSearch() {
	fmt.Println("----------------------------------")
	p, err := u.dao.VerifyUser("zhangsan1", "lisi666")
	if err != nil {
		fmt.Println("验证User发生错误:" + err.Error())
		return
	}
	fmt.Println("验证Users成功，权限等级为：" + string(p))
}
