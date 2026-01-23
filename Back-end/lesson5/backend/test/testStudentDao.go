package testFunction

import (
	"fmt"
	"lesson5/backend/DAO"
	"lesson5/backend/models"
)

type StudentDAO struct {
	dao         *DAO.StudentDAO
	testStudent models.Student
}

func NewTestStudentDAO(dao *DAO.StudentDAO) *StudentDAO {
	return &StudentDAO{dao: dao}
}

func (u StudentDAO) Test() {
	u.TestRead()
	u.TestCreate()
	u.TestRead()
	u.TestSearch()
	u.TestDelete()
	u.TestRead()
}

func (u StudentDAO) TestRead() {
	fmt.Println("----------------------------------")
	users, err := u.dao.ListStudents()
	if err != nil {
		fmt.Println("读取Student发生错误：" + err.Error())
		return
	}
	if users == nil || len(*users) == 0 {
		fmt.Println("读取Student发生错误：空值")
		return
	}
	fmt.Println("读取Student成功：")
	fmt.Println(*users)
}

func (u StudentDAO) TestCreate() {
	fmt.Println("----------------------------------")
	u.testStudent = models.Student{
		Name: "张三",
	}
	err := u.dao.CreateStudent(&u.testStudent)
	if err != nil {
		fmt.Println("创建Student发生错误:" + err.Error())
		return
	}
	fmt.Println("创建Student成功：")

}

func (u StudentDAO) TestDelete() {
	fmt.Println("----------------------------------")
	err := u.dao.DeleteStudent(&u.testStudent)
	if err != nil {
		fmt.Println("删除Student发生错误:" + err.Error())
		return
	}
	fmt.Println("删除Student成功：")
}

func (u StudentDAO) TestSearch() {
	fmt.Println("----------------------------------")
	s, err := u.dao.GetStudentByName("张三")
	if err != nil {
		fmt.Println("按名字查找Student发生错误:" + err.Error())
		return
	}
	fmt.Println("按名字查找Student成功:", s)

	s, err = u.dao.GetStudentByID(u.testStudent.ID)
	if err != nil {
		fmt.Println("按ID查找Student发生错误:" + err.Error())
		return
	}
	fmt.Println("按ID查找Student成功:", s)
}
