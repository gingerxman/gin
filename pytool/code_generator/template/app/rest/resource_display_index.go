package {{package}}

import (
	"fmt"
	"{{service_name}}/business/{{package}}"

	"github.com/gingerxman/eel"
)

type {{class_name}}DisplayIndex struct {
	eel.RestResource
}

func (this *{{class_name}}DisplayIndex) Resource() string {
	return "{{package}}.{{name}}_display_index"
}

func (this *{{class_name}}DisplayIndex) GetParameters() map[string][]string {
	return map[string][]string{
		"POST": []string{
			"id:int",
			"action",
		},
	}
}

func (this *{{class_name}}DisplayIndex) Post(ctx *eel.Context) {
	req := ctx.Request
	bCtx := ctx.GetBusinessContext()

	id, _ := req.GetInt("id")
	{{var_name}} := {{package}}.New{{class_name}}Repository(bCtx).Get{{class_name}}(id)

	if {{var_name}} == nil {
		ctx.Response.Error( "{{name}}_display_index:invalid_{{var_name}}", fmt.Sprintf("id(%d)", id))
		return
	}

	action := req.GetString("action")
	{{var_name}}.UpdateDisplayIndex(action)

	ctx.Response.JSON(eel.Map{})
}
