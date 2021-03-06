package {{package}}

import (
	"fmt"
	"errors"
	"context"
	m_{{package}} "{{service_name}}/models/{{package}}"
	"time"
	{% if should_create_batch_json_factory %}
	"encoding/json"
	{%- endif -%}

	{% if enable_display_index %}
	"{{service_name}}/business/common"
	{%- endif %}

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

	//refer object
	{%- for refer in refers %}
	{%- if refer.enable_fill_nto1_1 %}
	{{refer.resource.class_name}}Id int //refer to {{refer.resource.name}}
	{{refer.resource.class_name}} *{{refer.resource.class_name}}
	{%- endif %}

	{%- if refer.enable_fill_nto1_n or refer.enable_fill_nton %}
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
	{%- if refer.update_nto1_1 %}
	{{refer.resource.var_name}}Id int,{{""-}}
	{%- endif %}
	{%- endfor %}

	{%- for refer in refers %}
	{%- if refer.update_nto1_n %}
	{{refer.resource.plural_var_name}}JsonStr string,{{""-}}
	{%- endif %}
	{%- endfor %}

	{%- for refer in refers %}
	{%- if refer.update_nton %}
	{{refer.resource.var_name}}Ids []int,{{""-}}
	{%- endif %}
	{%- endfor %}
) error {
	o := eel.GetOrmFromContext(this.Ctx)

	db := o.Model(&m_{{package}}.{{class_name}}{}).Where("id", this.Id).Update(gorm.Params{
		{%- for field in updatable_fields %}
		"{{ field.snake_name }}": {{ field.var_name }},{{""-}}
		{% endfor %}

		{%- for refer in refers %}
		{%- if refer.update_nto1_1 %}
		"{{ refer.resource.name }}_id": {{ refer.resource.var_name }}Id,{{""-}}
		{%- endif %}
		{%- endfor %}
	})

	if db.Error != nil {
		eel.Logger.Error(db.Error)
		return errors.New("{{name}}:update_fail")
	}

	{%- for refer in refers %}
	{%- if refer.update_nton %}

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

	{%- for refer in refers %}
	{%- if refer.create_nto1_n %}
	
	//更新{{refer.resource.class_name}}记录
	NewBatch{{refer.resource.plural_class_name}}FromJSON(
		this.Ctx,
		this,
		{{refer.resource.plural_var_name}}JsonStr,
		true,
	)
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
		"id": this.Id,
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
	item := common.ItemPos{
		Id: this.Id,
		DisplayIndex: this.DisplayIndex,
		OriginalDisplayIndex: this.OriginalDisplayIndex,
		Table: model.TableName(),
	}
	err := common.NewUpdateDisplayIndexService(this.Ctx, common.DISPLAY_INDEX_ORDER_ASC).Update(&item, action, eel.Map{})
	if err != nil {
		eel.Logger.Error(err)
		return err
	}
	
	return nil
}

{%- endif %}

//工厂方法
{%- if should_create_batch_json_factory %}
type _{{class_name}}InputData struct {
	{%- for field in creatable_fields %}
	{{ field.name }} {{ field.json_type }} `json:"{{field.snake_name}}"`{{""-}}
	{%- endfor %}
}

func NewBatch{{plural_class_name}}FromJSON(
	ctx context.Context,
	{%- if data_owner_field %}
	{{data_owner_field.iface.var_name}} {{data_owner_field.iface.type}},
	{%- endif %}

	{%- for refer in refers %}
	{%- if refer.create_nto1_1 %}
	{{refer.resource.var_name}} *{{refer.resource.class_name}},{{""-}}
	{%- endif %}
	{%- endfor %}
	
	jsonStr string,
	removeOldData bool,
) {
	inputDatas := make([]*_{{class_name}}InputData, 0)
	err := json.Unmarshal([]byte(jsonStr), &inputDatas)
	if err != nil {
		eel.Logger.Error(err)
		panic(eel.NewBusinessError("{{var_name}}:parse_batch_json_data_fail", "解析batch json data出错"))
	}

	for _, inputData := range inputDatas {
		New{{class_name}}(
			ctx,
			{%- if data_owner_field %}
			{{data_owner_field.iface.var_name}},
			{%- endif %}

			{%- for field in creatable_fields %}
			inputData.{{ field.name }},{{""-}}
			{% endfor %}

			{%- for refer in refers %}
			{%- if refer.create_nto1_1 %}
			{{refer.resource.var_name}},{{""-}}
			{%- endif %}
			{%- endfor %}
		)
	}
}
{%- endif %}

func New{{class_name}}(
	ctx context.Context,
	{%- if data_owner_field %}
	{{data_owner_field.iface.var_name}} {{data_owner_field.iface.type}},
	{%- endif %}
	{%- for field in creatable_fields %}
	{{ field.var_name }} {{ field.type -}},{{""-}}
	{% endfor %}

	{%- for refer in refers %}
	{%- if refer.create_nto1_1 %}
	{{refer.resource.var_name}} *{{refer.resource.class_name}},{{""-}}
	{%- endif %}
	{%- endfor %}

	{%- for refer in refers %}
	{%- if refer.update_nto1_n %}
	{{refer.resource.plural_var_name}}JsonStr string,{{""-}}
	{%- endif %}
	{%- endfor %}

	{%- for refer in refers %}
	{%- if refer.create_nton %}
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
	{%- if refer.create_nto1_1 %}
	model.{{refer.resource.class_name}}Id = {{refer.resource.var_name}}.Id
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
	{%- if refer.create_nton %}

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

	instance := New{{class_name}}FromModel(ctx, &model)

	{%- for refer in refers %}
	{%- if refer.create_nto1_n %}
	//创建{{refer.resource.class_name}}记录
	NewBatch{{refer.resource.plural_class_name}}FromJSON(
		ctx,
		instance,
		{{refer.resource.plural_var_name}}JsonStr,
		false,
	)
	{%- endif %}
	{%- endfor %}
	return instance
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
