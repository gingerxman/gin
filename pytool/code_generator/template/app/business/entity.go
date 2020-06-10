package {{package}}

import (
	"fmt"
	"errors"
	"context"
	"{{service_name}}/business"
	m_{{package}} "{{service_name}}/models/{{package}}"
	"time"

	"github.com/gingerxman/eel"
	"github.com/gingerxman/gorm"
)

type {{class_name}} struct {
	eel.EntityBase
	Id int
	{%- for field in fields %}
	{{ field.name }} {{ field.type -}}
	{% endfor %}
	{%- if enable_display_index %}
	DisplayIndex int //显示时的排序
	OriginalDisplayIndex int //置顶置底前的排序
	{%- endif %}
	CreatedAt time.Time

	//foreign key
	{%- for refer in refers %}
	{%- if refer.enable_fill_object %}
	{{refer.resource.class_name}}Id int //refer to {{refer.resource.name}}
	{{refer.resource.class_name}} *{{refer.resource.class_name}}
	{%- endif %}

	{%- if refer.enable_fill_objects %}
	{{refer.resource.plural_class_name}} []*{{refer.resource.class_name}}
	{%- endif %}
	{%- endfor %}
}

//Update 更新对象
func (this *{{class_name}}) Update(
	{%- for field in updatable_fields %}
	{{ field.var_name }} {{ field.type -}},{{""-}}
	{% endfor %}

	{%- for refer in refers %}
	{%- if refer.enable_fill_object %}
	{{refer.resource.var_name}}Id int,{{""-}}
	{%- endif %}
	{%- endfor %}

	{%- for refer in refers %}
	{%- if refer.enable_fill_objects %}
	{{refer.resource.var_name}}Ids []int,{{""-}}
	{%- endif %}
	{%- endfor %}
) error {
	o := eel.GetOrmFromContext(this.Ctx)

	db := o.Model(&m_{{package}}.{{class_name}}).Where("id", this.Id).Update(gorm.Params{
		{%- for field in updatable_fields %}
		"{{ field.snake_name }}": {{ field.var_name }},{{""-}}
		{% endfor %}

		{%- for refer in refers %}
		{%- if refer.enable_fill_object %}
		"{{ refer.resource.name }}_id": {{ refer.resource.var_name }}Id,{{""-}}
		{%- endif %}
		{%- endfor %}
	})

	if db.Error != nil {
		eel.Logger.Error(db.Error)
		return errors.New("{{name}}:update_fail")
	}

	{%- for refer in refers %}
	{%- if refer.enable_fill_objects %}

	//删除{{class_name}}Has{{refer.resource.class_name}}中的老数据
	db = o.Where("{{name}}_id", this.Id).Delete(&m_{{package}}.{{class_name}}Has{{refer.resource.class_name}}{})
	if db.Error != nil {
		eel.Logger.Error(db.Error)
	}

	//创建{{class_name}}Has{{refer.resource.class_name}}记录
	for _, {{refer.resource.var_name}}Id := range {{refer.resource.var_name}}Ids {
		relationModel := m_{{package}}.{{class_name}}Has{{refer.resource.class_name}}{}
		relationModel.{{class_name}}Id = this.Id
		relationModel.{{refer.resource.class_name}}Id = {{refer.resource.var_name}}Id
		db = o.Create(&relationModel)
		if db.Error != nil {
			eel.Logger.Error(db.Error)
			panic(eel.NewBusinessError("{{name}}:create_relation_fail", fmt.Sprintf("创建失败")))
		}
	}
	{%- endif %}
	{%- endfor %}

	return nil
}

func (this *{{class_name}}) enable(isEnabled bool) {
	o := eel.GetOrmFromContext(this.Ctx)
	
	db := o.Model(&m_{{package}}.{{class_name}}{}).Where(eel.Map{
		"id": this.Id,
	}).Update(gorm.Params{
		"is_enabled": isEnabled,
	})

	if db.Error != nil {
		eel.Logger.Error(db.Error)
	}
}

func (this *{{class_name}}) Enable() {
	this.enable(true)
}

func (this *{{class_name}}) Disable() {
	this.enable(false)
}

func (this *{{class_name}}) Delete() error {
	o := eel.GetOrmFromContext(this.Ctx)
	
	db := o.Model(&m_{{package}}.{{class_name}}{}).Where(eel.Map{
		"id", this.Id,
	}).Update(gorm.Params{
		"is_deleted": true,
	})

	if db.Error != nil {
		eel.Logger.Error(db.Error)
		return db.Error
	}

	return nil
}

{% if enable_display_index %}
func (this *{{class_name}}) UpdateDisplayIndex(action string) error {
	model := m_{{package}}.{{class_name}}{}
	item := itemPos{
		Id: this.Id,
		DisplayIndex: this.DisplayIndex,
		OriginalDisplayIndex: this.OriginalDisplayIndex,
		Table: model.TableName(),
	}
	err := NewUpdateDisplayIndexService(this.Ctx, DISPLAY_INDEX_ORDER_ASC).Update(&item, action)
	if err != nil {
		eel.Logger.Error(err)
		return err
	}
	
	return nil
}

{%- endif %}

//工厂方法
func New{{class_name}}(
	ctx context.Context,
	{%- if data_owner_field %}
	{{data_owner_field.iface.var_name}} {{data_owner_field.iface.type}},
	{%- endif %}
	{%- for field in creatable_fields %}
	{{ field.var_name }} {{ field.type -}},{{""-}}
	{% endfor %}

	{%- for refer in refers %}
	{%- if refer.enable_fill_object %}
	{{refer.resource.var_name}}Id int,{{""-}}
	{%- endif %}
	{%- endfor %}

	{%- for refer in refers %}
	{%- if refer.enable_fill_objects %}
	{{refer.resource.var_name}}Ids []int,{{""-}}
	{%- endif %}
	{%- endfor %}
) *{{class_name}} {
	o := eel.GetOrmFromContext(ctx)

	//保存数据
	model := m_{{package}}.{{class_name}}{}
	{%- if data_owner_field %}
	model.{{ data_owner_field.name }} = {{data_owner_field.iface.var_name}}.GetId()
	{%- endif -%}
	{% for field in creatable_fields %}
	model.{{ field.name }} = {{ field.var_name -}}
	{% endfor %}
	{%- for field in non_creatable_fields %}
	model.{{ field.name }} = {{ field.default_value -}}
	{% endfor %}

	{%- for refer in refers %}
	{%- if refer.enable_fill_object %}
	model.{{refer.resource.class_name}}Id = {{refer.resource.var_name}}Id
	{%- endif %}
	{%- endfor %}
	db := o.Create(&model)
	if db.Error != nil {
		eel.Logger.Error(db.Error)
		panic(eel.NewBusinessError("{{name}}:create_fail", fmt.Sprintf("创建失败")))
	}
	{% if enable_display_index %}
	//更新display_index
	db = o.Model(&model).Where("id", model.Id).Update(gorm.Params{
		"display_index": model.Id,
	})
	if db.Error != nil {
		eel.Logger.Error(db.Error)
		panic(eel.NewBusinessError("{{name}}:update_display_index_fail", fmt.Sprintf("创建失败")))
	}
	{%- endif %}

	{%- for refer in refers %}
	{%- if refer.enable_fill_objects %}

	//创建{{class_name}}Has{{refer.resource.class_name}}记录
	for _, {{refer.resource.var_name}}Id := range {{refer.resource.var_name}}Ids {
		relationModel := m_{{package}}.{{class_name}}Has{{refer.resource.class_name}}{}
		relationModel.{{class_name}}Id = model.Id
		relationModel.{{refer.resource.class_name}}Id = {{refer.resource.var_name}}Id
		db = o.Create(&relationModel)
		if db.Error != nil {
			eel.Logger.Error(db.Error)
			panic(eel.NewBusinessError("{{name}}:create_relation_fail", fmt.Sprintf("创建失败")))
		}
	}
	{%- endif %}
	{%- endfor %}

	return New{{class_name}}FromModel(ctx, &model)
}

//根据model构建对象
func New{{class_name}}FromModel(ctx context.Context, model *m_{{package}}.{{class_name}}) *{{class_name}} {
	instance := new({{class_name}})
	instance.Ctx = ctx
	instance.Model = model
	instance.Id = model.Id
	{%- if enable_display_index %}
	instance.DisplayIndex = model.DisplayIndex
	instance.OriginalDisplayIndex = model.OriginalDisplayIndex
	{%- endif %}
	{%- for field in fields %}
	instance.{{ field.name }} = model.{{ field.name -}}
	{% endfor %}

	{%- for refer in refers %}
	{%- if refer.enable_fill_object %}
	instance.{{refer.resource.class_name}}Id = model.{{refer.resource.class_name}}Id
	{%- endif %}
	{%- endfor %}
	instance.CreatedAt = model.CreatedAt

	return instance
}

func init() {
}
