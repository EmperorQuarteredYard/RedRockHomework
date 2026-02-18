package models

import "errors"

type Department struct {
	Department      string
	DepartmentLabel string
	Instruction     string
}

var (
	Backend  = Department{"backend", "后端", "后端开发"}
	Frontend = Department{"frontend", "前端", "前端开发"}
	Sre      = Department{"sre", "SRE", "运维工程"}
	Product  = Department{"product", "产品", "产品设计"}
	Design   = Department{"design", "视觉设计", "UI/UX 设计"}
	Android  = Department{"android", "Android", "Android 开发"}
	Ios      = Department{"ios", "iOS", "iOS 开发"}
)

var departmentMap = map[string]Department{
	"backend":  Backend,
	"frontend": Frontend,
	"sre":      Sre,
	"product":  Product,
	"design":   Design,
	"android":  Android,
	"ios":      Ios,
}

func GetDepartment(val string) (*Department, error) {
	dept, ok := departmentMap[val]
	if !ok {
		return nil, errors.New("部门不存在")
	}
	return &dept, nil
}

func GetAllDepartments() []Department {
	return []Department{Backend, Frontend, Sre, Product, Design, Android, Ios}
}

func IsValidDepartment(val string) bool {
	_, ok := departmentMap[val]
	return ok
}
