package {{package}}

import (
	"fmt"
	"{{service_name}}/business/account"
	{% if name == package -%}
	b_{{package}} "{{service_name}}/business/{{package}}"
	{% else -%}
	"{{service_name}}/business/{{package}}"
	{% endif -%}

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
			{%- for field in fields %}
			"{{ field.snake_name }}:{{ field.rest_type }}",{{""-}}
			{% endfor %}
			
			{%- for refer in refers %}
			{%- if refer.enable_fill_object %}
			"{{refer.resource.name}}_id:int",{{""-}}
			{%- endif %}
			{%- endfor %}

			{%- for refer in refers %}
			{%- if refer.enable_fill_objects %}
			"{{refer.resource.name}}_ids:json-array",{{""-}}
			{%- endif %}
			{%- endfor %}
		},
		"POST": []string{
			"id:int",
			{%- for field in fields %}
			"{{ field.snake_name }}:{{ field.rest_type }}",{{""-}}
			{% endfor %}

			{%- for refer in refers %}
			{%- if refer.enable_fill_object %}
			"{{refer.resource.name}}_id:int",{{""-}}
			{%- endif %}
			{%- endfor %}

			{%- for refer in refers %}
			{%- if refer.enable_fill_objects %}
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
	repository := {% if name == package -%}b_{% endif %}{{package}}.New{{class_name}}Repository(bCtx)
	{{var_name}} := repository.Get{{class_name}}(id)

	if {{var_name}} == nil {
		ctx.Response.Error( "{{var_name}}:invalid_{{var_name}}", fmt.Sprintf("id(%d)", id))
		return
	}

	fillService := {% if name == package -%}b_{% endif %}{{package}}.NewFill{{class_name}}Service(bCtx)
	fillService.Fill([]*{% if name == package -%}b_{% endif %}{{package}}.{{class_name}}{ {{var_name}} }, eel.FillOption{
		{%- for refer in refers %}
		{%- if refer.enable_fill_object or refer.enable_fill_objects %}
		"with_{{ refer.resource.name }}": true,
		{%- endif %}
		{%- endfor %}
	})

	encodeService := {% if name == package -%}b_{% endif %}{{package}}.NewEncode{{class_name}}Service(bCtx)
	respData := encodeService.Encode({{var_name}})

	ctx.Response.JSON(respData)
}

func (this *{{class_name}}) Put(ctx *eel.Context) {
	req := ctx.Request
	bCtx := ctx.GetBusinessContext()

	{%- for field in fields %}
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
	{%- if refer.enable_fill_object %}
	{{refer.resource.var_name}}Id, _ := req.GetInt("{{ refer.resource.name }}_id")
	{%- endif %}
	{%- endfor %}

	{%- for refer in refers %}
	{%- if refer.enable_fill_objects %}
	{{refer.resource.var_name}}Ids := req.GetIntArray("{{ refer.resource.name }}_ids")
	{%- endif %}
	{%- endfor %}

	corp := account.GetCorpFromContext(bCtx)
	{{var_name}} := {% if name == package -%}b_{% endif %}{{package}}.New{{class_name}}(
		bCtx, 
		corp, 
		{%- for field in fields %}
		{{ field.var_name }},{{""-}}
		{%- endfor %}

		{%- for refer in refers %}
		{%- if refer.enable_fill_object %}
		{{refer.resource.var_name}}Id,{{""-}}
		{%- endif %}
		{%- endfor %}

		{%- for refer in refers %}
		{%- if refer.enable_fill_objects %}
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
	{%- for field in fields %}
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
	{%- if refer.enable_fill_object %}
	{{refer.resource.var_name}}Id, _ := req.GetInt("{{ refer.resource.name }}_id")
	{%- endif %}
	{%- endfor %}

	{%- for refer in refers %}
	{%- if refer.enable_fill_objects %}
	{{refer.resource.var_name}}Ids := req.GetIntArray("{{ refer.resource.name }}_ids")
	{%- endif %}
	{%- endfor %}

	repository := {% if name == package -%}b_{% endif %}{{package}}.New{{class_name}}Repository(bCtx)
	{{var_name}} := repository.Get{{class_name}}(id)

	if {{var_name}} == nil {
		ctx.Response.Error( "{{var_name}}:invalid_{{var_name}}", fmt.Sprintf("id(%d)", id))
		return
	}

	_ = {{var_name}}.Update(
		{%- for field in fields %}
		{{ field.var_name }},{{""-}}
		{%- endfor %}

		{%- for refer in refers %}
		{%- if refer.enable_fill_object %}
		{{refer.resource.var_name}}Id,{{""-}}
		{%- endif %}
		{%- endfor %}

		{%- for refer in refers %}
		{%- if refer.enable_fill_objects %}
		{{refer.resource.var_name}}Ids,{{""-}}
		{%- endif %}
		{%- endfor %}
	)

	ctx.Response.JSON(eel.Map{})
}

func (this *{{class_name}}) Delete(ctx *eel.Context) {
	req := ctx.Request
	bCtx := ctx.GetBusinessContext()

	corp := account.GetCorpFromContext(bCtx)
	id, _ := req.GetInt("id")
	repository := {% if name == package -%}b_{% endif %}{{package}}.New{{class_name}}Repository(bCtx)
	{{var_name}} := repository.Get{{class_name}}InCorp(corp, id)
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