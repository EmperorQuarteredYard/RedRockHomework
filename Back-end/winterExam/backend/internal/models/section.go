package models

import "errors"

type Department struct {
	Department      string `json:"department"`
	DepartmentLabel string `json:"department_label"`
	Instruction     string `json:"instruction"`
}

var (
	Backend = Department{
		Department:      "backend",
		DepartmentLabel: "后端",
		Instruction:     "后端开发",
	}

	Frontend = Department{
		Department:      "frontend",
		DepartmentLabel: "前端",
		Instruction:     "前端开发",
	}

	Sre = Department{
		Department:      "sre",
		DepartmentLabel: "SRE",
		Instruction:     "运维工程",
	}

	Product = Department{
		Department:      "product",
		DepartmentLabel: "产品",
		Instruction:     "产品设计",
	}

	Design = Department{
		Department:      "design",
		DepartmentLabel: "视觉设计",
		Instruction:     "UI/UX 设计",
	}

	Android = Department{
		Department:      "android",
		DepartmentLabel: "Android",
		Instruction:     "Android 开发",
	}

	Ios = Department{
		Department:      "ios",
		DepartmentLabel: "iOS",
		Instruction:     "iOS 开发",
	}
)

var departments = map[string]Department{
	"frontend": Frontend,
	"product":  Product,
	"design":   Design,
	"android":  Android,
	"ios":      Ios,
	"sre":      Sre,
	"backend":  Backend,
}

func GetDepartment(name string) (*Department, error) {
	section, ok := departments[name]
	if !ok {
		return nil, errors.New("section not found")
	}
	return &section, nil
}
