package json_util

import "github.com/mailru/easyjson"

/*
go get github.com/mailru/easyjson
go install  github.com/mailru/easyjson/easyjson

@优势：性能是标准库json的5倍
@缺点：每定义一个结构体以及修改tag时需要重新生成easyjson代码


-- 如何生成easyjson需要的marshal/unmarshal代码
	进入当前目录json_util/
	# windows PowerShell
		执行 build_win.bat
*/

// User 支持json tag
// - 修改tag后 要再次运行脚本生成easyjson代码
type User struct {
	Name  string `json:"name1,omitempty,nocopy"`
	Desc  string `json:"desc2,intern"`
	other string
}

func EasyJsonExample() {
	user := &User{
		Name:  "xxx",
		Desc:  "xxx",
		other: "1",
	}
	b, _ := easyjson.Marshal(user)
	println(string(b))
}
