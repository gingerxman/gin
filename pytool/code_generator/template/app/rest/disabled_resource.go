package {{package}}

import (
	"fmt"
	"{{service_name}}/business/account"

	"github.com/gingerxman/eel"
	"{{service_name}}/business/{{package}}"
)

type Disabled{{class_name}} struct {
	eel.RestResource
}

func (this *Disabled{{class_name}}) Resource() string {
	return "{{package}}.disabled_{{name}}"
}

func (this *Disabled{{class_name}}) GetParameters() map[string][]string {
	return map[string][]string{
		"PUT": []string{"id:int"},
		"DELETE": []string{"id:int"},
	}
}

func (this *Disabled{{class_name}}) Put(ctx *eel.Context) {
	req := ctx.Request
	bCtx := ctx.GetBusinessContext()
	
	corp := account.GetCorpFromContext(bCtx)
	id, _ := req.GetInt("id")
	repository := {{package}}.New{{class_name}}Repository(bCtx)
	{{var_name}} := repository.Get{{class_name}}InCorp(corp, id)

	if {{var_name}} == nil {
		ctx.Response.Error( "disabled_{{name}}:invalid_{{var_name}}", fmt.Sprintf("id(%d)", id))
		return
	}
	
	{{var_name}}.Disable()
	
	ctx.Response.JSON(eel.Map{})
}

func (this *Disabled{{class_name}}) Delete(ctx *eel.Context) {
	req := ctx.Request
	bCtx := ctx.GetBusinessContext()
	
	corp := account.GetCorpFromContext(bCtx)
	id, _ := req.GetInt("id")
	repository := {{package}}.New{{class_name}}Repository(bCtx)
	{{var_name}} := repository.Get{{class_name}}InCorp(corp, id)

	if {{var_name}} == nil {
		ctx.Response.Error( "disabled_{{name}}:invalid_{{var_name}}", fmt.Sprintf("id(%d)", id))
		return
	}
	
	{{var_name}}.Enable()
	
	ctx.Response.JSON(eel.Map{})
}
