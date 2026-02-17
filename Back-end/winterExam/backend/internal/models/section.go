package models

import "errors"

type Section struct {
	Department      string `json:"department"`
	DepartmentLabel string `json:"department_label"`
	Instruction     string `json:"instruction"`
}

var (
	Backend = Section{
		Department:      "backend",
		DepartmentLabel: "后端",
		Instruction:     "后端开发",
	}

	Frontend = Section{
		Department:      "frontend",
		DepartmentLabel: "前端",
		Instruction:     "前端开发",
	}

	Sre = Section{
		Department:      "sre",
		DepartmentLabel: "SRE",
		Instruction:     "运维工程",
	}

	Product = Section{
		Department:      "product",
		DepartmentLabel: "产品",
		Instruction:     "产品设计",
	}

	Design = Section{
		Department:      "design",
		DepartmentLabel: "视觉设计",
		Instruction:     "UI/UX 设计",
	}

	Android = Section{
		Department:      "android",
		DepartmentLabel: "Android",
		Instruction:     "Android 开发",
	}

	Ios = Section{
		Department:      "ios",
		DepartmentLabel: "iOS",
		Instruction:     "iOS 开发",
	}
)

var sections = map[string]Section{
	"frontend": Frontend,
	"product":  Product,
	"design":   Design,
	"android":  Android,
	"ios":      Ios,
	"sre":      Sre,
	"backend":  Backend,
}

func GetSection(name string) (*Section, error) {
	section, ok := sections[name]
	if !ok {
		return nil, errors.New("section not found")
	}
	return &section, nil
}
