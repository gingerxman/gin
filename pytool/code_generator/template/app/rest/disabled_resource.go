package {{package}}

import (
	"fmt"
	"{{service_name}}/business/account"

	"github.com/gingerxman/eel"
	b_{{package}} "{{service_name}}/business/{{package}}"
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
	
	id, _ := req.GetInt("id")
	{% if belong_to_user -%}
	user := account.GetUserFromContext(bCtx)
	{{var_name}} := b_{{package}}.New{{class_name}}Repository(bCtx).Get{{class_name}}ForUser(user, id)
	{% endif -%}
	{% if belong_to_corp -%}
	corp := account.GetCorpFromContext(bCtx)
	{{var_name}} := b_{{package}}.New{{class_name}}Repository(bCtx).Get{{class_name}}InCorp(corp, id)
	{% endif -%}

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
	
	id, _ := req.GetInt("id")
	{% if belong_to_user -%}
	user := account.GetUserFromContext(bCtx)
	{{var_name}} := b_{{package}}.New{{class_name}}Repository(bCtx).Get{{class_name}}ForUser(user, id)
	{% endif -%}
	{% if belong_to_corp -%}
	corp := account.GetCorpFromContext(bCtx)
	{{var_name}} := b_{{package}}.New{{class_name}}Repository(bCtx).Get{{class_name}}InCorp(corp, id)
	{% endif -%}

	if {{var_name}} == nil {
		ctx.Response.Error( "disabled_{{name}}:invalid_{{var_name}}", fmt.Sprintf("id(%d)", id))
		return
	}
	
	{{var_name}}.Enable()
	
	ctx.Response.JSON(eel.Map{})
}
