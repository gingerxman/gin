package {{package}}

import (
	"context"

	"github.com/gingerxman/eel"
)

type Encode{{class_name}}Service struct {
	eel.ServiceBase
}

func NewEncode{{class_name}}Service(ctx context.Context) *Encode{{class_name}}Service {
	service := new(Encode{{class_name}}Service)
	service.Ctx = ctx
	return service
}

//Encode 对单个实体对象进行编码
func (this *Encode{{class_name}}Service) Encode({{var_name}} *{{class_name}}) *R{{class_name}} {
	if {{var_name}} == nil {
		return nil
	}

	{%- for refer in refers %}
	{%- if refer.enable_fill_nto1_1 %}
	r{{refer.resource.class_name}} := NewEncode{{refer.resource.class_name}}Service(this.Ctx).Encode({{var_name}}.{{refer.resource.class_name}})
	{%- endif %}
	{%- if refer.enable_fill_nto1_n or refer.enable_fill_nton %}
	r{{refer.resource.plural_class_name}} := NewEncode{{refer.resource.class_name}}Service(this.Ctx).EncodeMany({{var_name}}.{{refer.resource.plural_class_name}})
	{%- endif %}
	{%- endfor %}

	return &R{{class_name}}{
		Id: {{var_name}}.Id,
		{%- for field in fields %}
		{%- if field.type == "time.Time" %}
		{{ field.name }}: {{var_name}}.{{ field.name }}.Format("2006-01-02 15:04:05"),{{""-}}
		{%- else %}
		{{ field.name }}: {{var_name}}.{{ field.name }},{{""-}}
		{%- endif %}
		{%- endfor %}

		{%- for refer in refers %}
		{%- if refer.enable_fill_nto1_1 %}
		{{refer.resource.class_name}}: r{{refer.resource.class_name}},{{""-}}
		{%- endif %}
		{%- if refer.enable_fill_nto1_n or refer.enable_fill_nton %}
		{{refer.resource.plural_class_name}}: r{{refer.resource.plural_class_name}},{{""-}}
		{%- endif %}
		{%- endfor %}
		{%- if enable_display_index %}
		DisplayIndex: {{var_name}}.DisplayIndex,
		{%- endif %}
		CreatedAt: {{var_name}}.CreatedAt.Format("2006-01-02 15:04:05"),
	}
}

//EncodeMany 对实体对象进行批量编码
func (this *Encode{{class_name}}Service) EncodeMany({{plural_var_name}} []*{{class_name}}) []*R{{class_name}} {
	rDatas := make([]*R{{class_name}}, 0)
	for _, {{var_name}} := range {{plural_var_name}} {
		rDatas = append(rDatas, this.Encode({{var_name}}))
	}
	
	return rDatas
}

func init() {
}
