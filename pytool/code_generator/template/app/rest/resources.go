package {{package}}

import (
	"{{service_name}}/business/{{package}}"
	"{{service_name}}/business/account"

	"github.com/gingerxman/eel"
)

type {{plural_class_name}} struct {
	eel.RestResource
}

func (this *{{plural_class_name}}) Resource() string {
	return "{{package}}.{{plural_name}}"
}

func (this *{{plural_class_name}}) GetParameters() map[string][]string {
	return map[string][]string{
		"GET": []string{"?filters:json", "?fill_options:json"},
	}
}

func (this *{{plural_class_name}}) Get(ctx *eel.Context) {
	req := ctx.Request
	bCtx := ctx.GetBusinessContext()

	page := req.GetPageInfo()
	filters := req.GetOrmFilters()
	repository := {{package}}.New{{class_name}}Repository(bCtx)
	{% if belong_to_user -%}
	user := account.GetUserFromContext(bCtx)
	{{plural_var_name}}, nextPageInfo := repository.GetEnabled{{plural_class_name}}ForUser(user, filters, page)
	{%- endif %}
	{% if belong_to_corp -%}
	corp := account.GetCorpFromContext(bCtx)
	{{plural_var_name}}, nextPageInfo := repository.GetEnabled{{plural_class_name}}ForCorp(corp, filters, page)
	{%- endif %}
	fillService := {{package}}.NewFill{{class_name}}Service(bCtx)
	fillService.Fill({{plural_var_name}}, eel.FillOption{
		{%- for refer in refers %}
		{%- if refer.enable_fill_nto1_1 or refer.enable_fill_nto1_n or enable_fill_nton %}
		"with_{{ refer.resource.name }}": true,
		{%- endif %}
		{%- endfor %}
	})

	encodeService := {{package}}.NewEncode{{class_name}}Service(bCtx)
	rows := encodeService.EncodeMany({{plural_var_name}})
	
	ctx.Response.JSON(eel.Map{
		"{{plural_name}}": rows,
		"pageinfo": nextPageInfo.ToMap(),
	})
}
