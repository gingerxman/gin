package {{package}}

import (
	"fmt"
	{%- if belong_to_corp or belong_to_user -%}"{{service_name}}/business/account"{%- endif %}
	b_{{package}} "{{service_name}}/business/{{package}}"

	"github.com/gingerxman/eel"
)

type {{class_name}} struct {
	eel.RestResource
}

func (this *{{class_name}}) Resource() string {
	return "{{package}}.{{name}}"
}

func (this *{{class_name}}) GetParameters() map[string][]string {
	return map[string][]string{
		"GET": []string{"id:int"},
		"PUT": []string{
			{%- for field in creatable_fields %}
			"{{ field.snake_name }}:{{ field.rest_type }}",{{""-}}
			{% endfor %}
			
			{%- for refer in refers %}
			{%- if refer.create_nto1_1 %}
			"{{refer.resource.name}}_id:int",{{""-}}
			{%- endif %}
			{%- endfor %}

			{%- for refer in refers %}
			{%- if refer.create_nto1_n or refer.create_nton %}
			"{{refer.resource.name}}_ids:json-array",{{""-}}
			{%- endif %}
			{%- endfor %}
		},
		"POST": []string{
			"id:int",
			{%- for field in updatable_fields %}
			"{{ field.snake_name }}:{{ field.rest_type }}",{{""-}}
			{% endfor %}

			{%- for refer in refers %}
			{%- if refer.update_nto1_1 %}
			"{{refer.resource.name}}_id:int",{{""-}}
			{%- endif %}
			{%- endfor %}

			{%- for refer in refers %}
			{%- if refer.update_nto1_n or refer.update_nton %}
			"{{refer.resource.name}}_ids:json-array",{{""-}}
			{%- endif %}
			{%- endfor %}
		},
		"DELETE": []string{"id:int"},
	}
}

func (this *{{class_name}}) Get(ctx *eel.Context) {
	req := ctx.Request
	bCtx := ctx.GetBusinessContext()
	id, _ := req.GetInt("id", 0)
	repository := b_{{package}}.New{{class_name}}Repository(bCtx)
	{% if belong_to_user -%}
	user := account.GetUserFromContext(bCtx)
	{{var_name}} := repository.Get{{class_name}}ForUser(user, id)
	{% endif -%}
	{% if belong_to_corp -%}
	corp := account.GetCorpFromContext(bCtx)
	{{var_name}} := repository.Get{{class_name}}InCorp(corp, id)
	{% endif -%}
	{% if belong_to_platform -%}
	{{var_name}} := repository.Get{{class_name}}ById(id)
	{% endif -%}

	if {{var_name}} == nil {
		ctx.Response.Error( "{{var_name}}:invalid_{{var_name}}", fmt.Sprintf("id(%d)", id))
		return
	}

	fillService := b_{{package}}.NewFill{{class_name}}Service(bCtx)
	fillService.Fill([]*b_{{package}}.{{class_name}}{ {{var_name}} }, eel.FillOption{
		{%- for refer in refers %}
		{%- if refer.enable_fill_nto1_1 or refer.enable_fill_nto1_n or refer.enable_fill_nton %}
		"with_{{ refer.resource.name }}": true,
		{%- endif %}
		{%- endfor %}
	})

	encodeService := b_{{package}}.NewEncode{{class_name}}Service(bCtx)
	respData := encodeService.Encode({{var_name}})

	ctx.Response.JSON(respData)
}

func (this *{{class_name}}) Put(ctx *eel.Context) {
	req := ctx.Request
	bCtx := ctx.GetBusinessContext()

	{%- for field in creatable_fields %}
	{% if field.rest_type == "int" -%}
	{{ field.var_name}}, _ := req.GetInt("{{ field.snake_name }}")
	{%- endif -%}
	{% if field.rest_type == "bool" -%}
	{{ field.var_name}}, _ := req.GetBool("{{ field.snake_name }}")
	{%- endif -%}
	{% if field.rest_type == "float" -%}
	{{ field.var_name}}, _ := req.GetFloat("{{ field.snake_name }}")
	{%- endif -%}
	{% if field.rest_type == "string" -%}
	{{ field.var_name }} := req.GetString("{{ field.snake_name }}")
	{%- endif -%}
	{% endfor %}
	
	{% for refer in refers %}
	{% if refer.create_nto1_1 -%}
	{{refer.resource.var_name}}Id, _ := req.GetInt("{{ refer.resource.name }}_id")
	{{refer.resource.var_name}} := b_{{package}}.New{{refer.resource.class_name}}Repository(bCtx).Get{{refer.resource.class_name}}ById({{refer.resource.var_name}}Id)
	if {{refer.resource.var_name}} == nil {
		ctx.Response.Error( "{{var_name}}:invalid_{{refer.resource.var_name}}", fmt.Sprintf("id(%d)", {{refer.resource.var_name}}Id))
		return
	}
	{%- endif %}
	{%- endfor %}

	{%- for refer in refers %}
	{%- if refer.create_nto1_n or refer.create_nton %}
	{{refer.resource.var_name}}Ids := req.GetIntArray("{{ refer.resource.name }}_ids")
	{%- endif %}
	{%- endfor %}

	{% if belong_to_user -%}
	user := account.GetUserFromContext(bCtx)
	{%- endif %}
	{% if belong_to_corp -%}
	corp := account.GetCorpFromContext(bCtx)
	{%- endif %}
	{{var_name}} := b_{{package}}.New{{class_name}}(
		bCtx,
		{% if belong_to_user -%}
		user,
		{%- endif %}
		{%- if belong_to_corp -%}
		corp,
		{%- endif -%}
		
		{%- for field in creatable_fields %}
		{{ field.var_name }},{{""-}}
		{%- endfor %}

		{%- for refer in refers %}
		{%- if refer.create_nto1_1 %}
		{{refer.resource.var_name}},{{""-}}
		{%- endif %}
		{%- endfor %}

		{%- for refer in refers %}
		{%- if refer.create_nto1_n or refer.create_nton %}
		{{refer.resource.var_name}}Ids,{{""-}}
		{%- endif %}
		{%- endfor %}
	)

	ctx.Response.JSON(eel.Map{
		"id": {{var_name}}.Id,
	})
}

func (this *{{class_name}}) Post(ctx *eel.Context) {
	req := ctx.Request
	bCtx := ctx.GetBusinessContext()

	id, _ := req.GetInt("id")
	{%- for field in updatable_fields %}
	{% if field.rest_type == "int" -%}
	{{ field.var_name}}, _ := req.GetInt("{{ field.snake_name }}")
	{%- endif -%}
	{% if field.rest_type == "bool" -%}
	{{ field.var_name}}, _ := req.GetBool("{{ field.snake_name }}")
	{%- endif -%}
	{% if field.rest_type == "float" -%}
	{{ field.var_name}}, _ := req.GetFloat("{{ field.snake_name }}")
	{%- endif -%}
	{% if field.rest_type == "string" -%}
	{{ field.var_name }} := req.GetString("{{ field.snake_name }}")
	{%- endif -%}
	{% endfor %}

	{%- for refer in refers %}
	{%- if refer.update_nto1_1 %}
	{{refer.resource.var_name}}Id, _ := req.GetInt("{{ refer.resource.name }}_id")
	{%- endif %}
	{%- endfor %}

	{%- for refer in refers %}
	{%- if refer.update_nto1_n or refer.update_nton %}
	{{refer.resource.var_name}}Ids := req.GetIntArray("{{ refer.resource.name }}_ids")
	{%- endif %}
	{%- endfor %}

	repository := b_{{package}}.New{{class_name}}Repository(bCtx)
	{% if belong_to_user -%}
	user := account.GetUserFromContext(bCtx)
	{{var_name}} := repository.Get{{class_name}}ForUser(user, id)
	{% endif -%}
	{% if belong_to_corp -%}
	corp := account.GetCorpFromContext(bCtx)
	{{var_name}} := repository.Get{{class_name}}InCorp(corp, id)
	{% endif -%}
	{% if belong_to_platform -%}
	{{var_name}} := repository.Get{{class_name}}ById(id)
	{% endif -%}
	if {{var_name}} == nil {
		ctx.Response.Error( "{{var_name}}:invalid_{{var_name}}", fmt.Sprintf("id(%d)", id))
		return
	}

	_ = {{var_name}}.Update(
		{%- for field in updatable_fields %}
		{{ field.var_name }},{{""-}}
		{%- endfor %}

		{%- for refer in refers %}
		{%- if refer.update_nto1_1 %}
		{{refer.resource.var_name}}Id,{{""-}}
		{%- endif %}
		{%- endfor %}

		{%- for refer in refers %}
		{%- if refer.update_nto1_n or refer.update_nton %}
		{{refer.resource.var_name}}Ids,{{""-}}
		{%- endif %}
		{%- endfor %}
	)

	ctx.Response.JSON(eel.Map{})
}

func (this *{{class_name}}) Delete(ctx *eel.Context) {
	req := ctx.Request
	bCtx := ctx.GetBusinessContext()
	
	id, _ := req.GetInt("id")
	repository := b_{{package}}.New{{class_name}}Repository(bCtx)

	{% if belong_to_user -%}
	user := account.GetUserFromContext(bCtx)
	{{var_name}} := repository.Get{{class_name}}ForUser(user, id)
	{% endif -%}
	{% if belong_to_corp -%}
	corp := account.GetCorpFromContext(bCtx)
	{{var_name}} := repository.Get{{class_name}}InCorp(corp, id)
	{% endif -%}
	{% if belong_to_platform -%}
	{{var_name}} := repository.Get{{class_name}}ById(id)
	{% endif -%}
	if {{var_name}} == nil {
		ctx.Response.Error( "{{var_name}}:invalid_{{var_name}}", fmt.Sprintf("id(%d)", id))
		return
	}

	err := {{var_name}}.Delete()

	if err != nil {
		eel.Logger.Error(err)
		response := eel.MakeErrorResponse(500, "{{var_name}}:delete_fail", err.Error())
		ctx.Response.JSON(response)
	} else {
		ctx.Response.JSON(eel.Map{})
	}
}