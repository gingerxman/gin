package {{package}}

import (
	"context"
	"github.com/gingerxman/eel"
	{% if should_import_model %}
	m_{{package}} "{{service_name}}/models/{{app_name}}"
	{% endif %}
)

type Fill{{class_name}}Service struct {
	eel.ServiceBase
}

func NewFill{{class_name}}Service(ctx context.Context) *Fill{{class_name}}Service {
	service := new(Fill{{class_name}}Service)
	service.Ctx = ctx
	return service
}

func (this *Fill{{class_name}}Service) FillOne({{var_name}} *{{class_name}}, option eel.FillOption) {
	this.Fill([]*{{class_name}}{ {{var_name}} }, option)
}

func (this *Fill{{class_name}}Service) Fill({{plural_var_name}} []*{{class_name}}, option eel.FillOption) {
	if len({{plural_var_name}}) == 0 {
		return
	}
	
	ids := make([]int, 0)
	for _, {{var_name}} := range {{plural_var_name}} {
		ids = append(ids, {{var_name}}.Id)
	}

	{%- for refer in refers %}
	{%- if refer.enable_fill_nto1_1 or refer.enable_fill_nto1_n or refer.enable_fill_nton %}

	if enableOption, ok := option["with_{{refer.resource_name}}"]; ok && enableOption {
		this.fill{{refer.resource.class_name}}({{plural_var_name}}, ids)
	}

	{%- endif %}
	{%- endfor %}
	return
}

{%- for refer in refers %}
{%- if refer.enable_fill_nto1_1 %}


func (this *Fill{{class_name}}Service) fill{{refer.resource.class_name}}({{plural_var_name}} []*{{class_name}}, ids []int) {
	//获取关联的id集合
	{{refer.resource.var_name}}Ids := make([]int, 0)
	for _, {{var_name}} := range {{plural_var_name}} {
		{{refer.resource.var_name}}Ids = append({{refer.resource.var_name}}Ids, {{var_name}}.{{refer.resource.class_name}}Id)
	}

	//获取{{refer.resource.plural_var_name}}, 构建<id, {{refer.resource.var_name}}>
	{{refer.resource.plural_var_name}} := New{{refer.resource.class_name}}Repository(this.Ctx).Get{{refer.resource.plural_class_name}}ByIds({{refer.resource.var_name}}Ids)
	id2{{refer.resource.var_name}} := make(map[int]*{{refer.resource.class_name}})
	for _, {{refer.resource.var_name}} := range {{refer.resource.plural_var_name}} {
		id2{{refer.resource.var_name}}[{{refer.resource.var_name}}.Id] = {{ refer.resource.var_name }}
	}

	//填充{{name}}的{{refer.resource.class_name}}对象
	for _, {{var_name}} := range {{plural_var_name}} {
		if {{refer.resource.var_name}}, ok := id2{{refer.resource.var_name}}[{{var_name}}.{{refer.resource.class_name}}Id]; ok {
			{{var_name}}.{{refer.resource.class_name}} = {{refer.resource.var_name}}
		}
	}
}
{%- endif %}

{%- if refer.enable_fill_nto1_n %}


func (this *Fill{{class_name}}Service) fill{{refer.resource.class_name}}({{plural_var_name}} []*{{class_name}}, ids []int) {
	//从db中获取{{refer.resource.resource_class_name}}数据集合
	var models []*m_{{package}}.{{refer.resource.class_name}}
	o := eel.GetOrmFromContext(this.Ctx)
	db := o.Model(&m_{{package}}.{{refer.resource.class_name}}{}).Where("{{name}}_id__in", ids).Find(&models)
	if db.Error != nil {
		eel.Logger.Error(db.Error)
		return
	}

	//构建<id, {{name}}>
	id2entity := make(map[int]*{{class_name}})
	for _, {{var_name}} := range {{plural_var_name}} {
		id2entity[{{var_name}}.Id] = {{var_name}}
	}

	for _, model := range models {
		if entity, ok := id2entity[model.{{class_name}}Id]; ok {
			entity.{{refer.resource.plural_class_name}} = append(entity.{{refer.resource.plural_class_name}}, New{{refer.resource.class_name}}FromModel(this.Ctx, model))
		}
	}
}
{%- endif %}

{%- if refer.enable_fill_nton %}


func (this *Fill{{class_name}}Service) fill{{refer.resource.class_name}}({{plural_var_name}} []*{{class_name}}, ids []int) {
	//构建<id, {{name}}>
	id2entity := make(map[int]*{{class_name}})
	for _, {{var_name}} := range {{plural_var_name}} {
		id2entity[{{var_name}}.Id] = {{var_name}}
	}
	
	o := eel.GetOrmFromContext(this.Ctx)
	
	//从db中获取relation models
	var relationModels []*m_{{package}}.{{class_name}}Has{{refer.resource.class_name}}
	db := o.Model(&m_{{package}}.{{class_name}}Has{{refer.resource.class_name}}{}).Where("{{name}}_id__in", ids).Find(&relationModels)
	if db.Error != nil {
		eel.Logger.Error(db.Error)
		return
	}
	
	if len(relationModels) == 0 {
		return
	}
	
	//获取关联的id集合
	{{refer.resource.var_name}}Ids := make([]int, 0)
	for _, relationModel := range relationModels {
		{{refer.resource.var_name}}Ids = append({{refer.resource.var_name}}Ids, relationModel.{{refer.resource.class_name}}Id)
	}
	//从db中获取数据集合
	var models []*m_{{package}}.{{refer.resource.class_name}}
	db = o.Model(&m_{{package}}.{{refer.resource.class_name}}{}).Where("id__in", {{refer.resource.var_name}}Ids).Find(&models)
	if db.Error != nil {
		eel.Logger.Error(db.Error)
		return
	}
	//构建<id, model>
	id2model := make(map[int]*m_{{package}}.{{refer.resource.class_name}})
	for _, model := range models {
		id2model[model.Id] = model
	}
	
	//填充{{name}}的{{refer.resource.plural_class_name}}对象
	for _, relationModel := range relationModels {
		{{var_name}}Id := relationModel.{{class_name}}Id
		{{refer.resource.var_name}}Id := relationModel.{{refer.resource.class_name}}Id
		
		if {{var_name}}, ok := id2entity[{{var_name}}Id]; ok {
			if model, ok2 := id2model[{{refer.resource.var_name}}Id]; ok2 {
				{{var_name}}.{{refer.resource.plural_class_name}} = append({{var_name}}.{{refer.resource.plural_class_name}}, New{{refer.resource.class_name}}FromModel(this.Ctx, model))
			}
		}
	}
}
{%- endif %}
{% endfor %}

func init() {
}
