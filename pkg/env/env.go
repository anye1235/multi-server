package env

import (
	"flag"
	"fmt"
	"strings"
)

var (
	active Environment
	dev    Environment = &environment{value: "dev"}
	test   Environment = &environment{value: "test"}
	fat    Environment = &environment{value: "fat"}
	uat    Environment = &environment{value: "uat"}
	prod   Environment = &environment{value: "prod"}
)

var _ Environment = (*environment)(nil)

// Environment 环境配置
type Environment interface {
	Value() string
	IsDev() bool
	IsTest() bool
	IsFat() bool
	IsUat() bool
	IsProd() bool
	t()
}

type environment struct {
	value string
}

func (e *environment) Value() string {
	return e.value
}

func (e *environment) IsDev() bool {
	return e.value == "dev"
}

func (e *environment) IsTest() bool {
	return e.value == "test"
}

func (e *environment) IsFat() bool {
	return e.value == "fat"
}

func (e *environment) IsUat() bool {
	return e.value == "uat"
}

func (e *environment) IsProd() bool {
	return e.value == "prod"
}

func (e *environment) t() {}

func init() {
	env := flag.String("env", "", "请输入运行环境:\n dev:开发环境\n fat:测试环境\n uat:预上线环境\n pro:正式环境\n")
	flag.Parse()

	switch strings.ToLower(strings.TrimSpace(*env)) {
	case "dev":
		active = dev
	case "test":
		active = test
	case "fat":
		active = fat
	case "uat":
		active = uat
	case "prod":
		active = prod
	default:
		active = dev
		fmt.Println("Warning: '-env' cannot be found, or it is illegal. The default 'dev' will be used.")
	}
}

// Active 当前配置的env
func Active() Environment {
	return active
}
