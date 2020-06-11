package {{package}}

import (
	"{{service_name}}/business/{{package}}"
	"{{service_name}}/business/account"

	"github.com/gingerxman/eel"
)

type User{{plural_class_name}} struct {
	eel.RestResource
}

func (this *User{{plural_class_name}}) Resource() string {
	return "{{package}}.user_{{plural_name}}"
}

func (this *User{{plural_class_name}}) GetParameters() map[string][]string {
	return map[string][]string{
		"GET": []string{},
	}
}

func (this *User{{plural_class_name}}) Get(ctx *eel.Context) {
	req := ctx.Request
	bCtx := ctx.GetBusinessContext()

	page := req.GetPageInfo()
	filters := req.GetOrmFilters()
	repository := {{package}}.New{{class_name}}Repository(bCtx)
	user := account.GetUserFromContext(bCtx)
	{{plural_var_name}}, nextPageInfo := repository.GetAll{{plural_class_name}}ForUser(user, page, filters)

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
