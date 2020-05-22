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
		"GET": []string{},
	}
}

func (this *{{plural_class_name}}) Get(ctx *eel.Context) {
	req := ctx.Request
	bCtx := ctx.GetBusinessContext()

	page := req.GetPageInfo()
	filters := req.GetOrmFilters()
	repository := {{package}}.New{{class_name}}Repository(bCtx)
	corp := account.GetCorpFromContext(bCtx)
	{{plural_var_name}}, nextPageInfo := repository.GetEnabled{{plural_class_name}}ForCorp(corp, page, filters)

	fillService := {{package}}.NewFill{{class_name}}Service(bCtx)
	fillService.Fill({{plural_var_name}}, eel.FillOption{
		{%- for refer in refers %}
		{%- if refer.enable_fill_object or refer.enable_fill_objects %}
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
