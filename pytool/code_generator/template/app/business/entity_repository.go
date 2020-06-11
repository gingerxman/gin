package {{package}}

import (
	"context"
	"{{service_name}}/business"
	m_{{package}} "{{service_name}}/models/{{package}}"

	"github.com/gingerxman/eel"
)

type {{class_name}}Repository struct {
	eel.RepositoryBase
}

func New{{class_name}}Repository(ctx context.Context) *{{class_name}}Repository {
	repository := new({{class_name}}Repository)
	repository.Ctx = ctx
	return repository
}

func (this *{{class_name}}Repository) Get{{plural_class_name}}(filters eel.Map, orderExprs ...string) []*{{class_name}} {
	o := eel.GetOrmFromContext(this.Ctx)
	db := o.Model(&m_{{package}}.{{class_name}}{})
	
	var models []*m_{{package}}.{{class_name}}
	if len(filters) > 0 {
		db = db.Where(filters)
	}
	for _, expr := range orderExprs {
		db = db.Order(expr)
	}
	db = db.Find(&models)
	if db.Error != nil {
		eel.Logger.Error(db.Error)
		return nil
	}
	
	instances := make([]*{{class_name}}, 0)
	for _, model := range models {
		instances = append(instances, New{{class_name}}FromModel(this.Ctx, model))
	}
	return instances
}

func (this *{{class_name}}Repository) GetPaged{{plural_class_name}}(filters eel.Map, page *eel.PageInfo, orderExprs ...string) ([]*{{class_name}}, eel.INextPageInfo) {
	o := eel.GetOrmFromContext(this.Ctx)
	db := o.Model(&m_{{package}}.{{class_name}}{})
	
	var models []*m_{{package}}.{{class_name}}
	if len(filters) > 0 {
		db = db.Where(filters)
	}	
	for _, expr := range orderExprs {
		db = db.Order(expr)
	}
	paginateResult, db := eel.Paginate(db, page, &models)
	
	if db.Error != nil {
		eel.Logger.Error(db.Error)
		return nil, paginateResult
	}
	
	instances := make([]*{{class_name}}, 0)
	for _, model := range models {
		instances = append(instances, New{{class_name}}FromModel(this.Ctx, model))
	}
	return instances, paginateResult
}

{% if belong_to_corp -%}
//GetEnabled{{plural_class_name}}ForCorp 获得启用的{{class_name}}对象集合
func (this *{{class_name}}Repository) GetEnabled{{plural_class_name}}ForCorp(corp business.ICorp, filters eel.Map, page *eel.PageInfo) ([]*{{class_name}}, eel.INextPageInfo) {
	filters["corp_id"] = corp.GetId()
	filters["is_enabled"] = true
	filters["is_deleted"] = false
	{% if enable_display_index %}
	return this.GetPaged{{plural_class_name}}(filters, page, "display_index")
	{% else %}
	return this.GetPaged{{plural_class_name}}(filters, page, "id")
	{% endif %}
}

//GetAll{{plural_class_name}}ForCorp 获得所有{{class_name}}对象集合
func (this *{{class_name}}Repository) GetAll{{plural_class_name}}ForCorp(corp business.ICorp, filters eel.Map, page *eel.PageInfo) ([]*{{class_name}}, eel.INextPageInfo) {
	filters["corp_id"] = corp.GetId()
	filters["is_deleted"] = false
	{% if enable_display_index %}
	return this.GetPaged{{plural_class_name}}(filters, page, "display_index")
	{% else %}
	return this.GetPaged{{plural_class_name}}(filters, page, "id")
	{% endif %}
}

//Get{{class_name}}InCorp 根据id和corp获得{{class_name}}对象
func (this *{{class_name}}Repository) Get{{class_name}}InCorp(corp business.ICorp, id int) *{{class_name}} {
	filters := eel.Map{
		"corp_id": corp.GetId(),
		"id": id,
	}
	
	{{plural_var_name}} := this.Get{{plural_class_name}}(filters)
	
	if len({{plural_var_name}}) == 0 {
		return nil
	} else {
		return {{plural_var_name}}[0]
	}
}
{%- endif %}

{%- if belong_to_user -%}
func (this *{{class_name}}Repository) GetEnabled{{plural_class_name}}ForUser(user business.IUser, filters eel.Map, page *eel.PageInfo) ([]*{{class_name}}, eel.INextPageInfo) {
	filters["user_id"] = user.GetId()
	filters["is_enabled"] = true
	filters["is_deleted"] = false
	{% if enable_display_index %}
	return this.GetPaged{{plural_class_name}}(filters, page, "display_index")
	{% else %}
	return this.GetPaged{{plural_class_name}}(filters, page, "id")
	{% endif %}
}

func (this *{{class_name}}Repository) GetAll{{plural_class_name}}ForUser(user business.IUser, filters eel.Map, page *eel.PageInfo) ([]*{{class_name}}, eel.INextPageInfo) {
	filters["user_id"] = user.GetId()
	filters["is_deleted"] = false
	{% if enable_display_index %}
	return this.GetPaged{{plural_class_name}}(filters, page, "display_index")
	{% else %}
	return this.GetPaged{{plural_class_name}}(filters, page, "id")
	{% endif %}
}

func (this *{{class_name}}Repository) Get{{class_name}}ForUser(user business.IUser, id int) *{{class_name}} {
	filters := eel.Map{
		"user_id": user.GetId(),
		"id": id,
	}

	{{plural_var_name}} := this.Get{{plural_class_name}}(filters)
	
	if len({{plural_var_name}}) == 0 {
		return nil
	} else {
		return {{plural_var_name}}[0]
	}
}
{%- endif %}

func (this *{{class_name}}Repository) GetAll{{plural_class_name}}(filters eel.Map, page *eel.PageInfo) ([]*{{class_name}}, eel.INextPageInfo) {
	filters["is_deleted"] = false
	{% if enable_display_index %}
	return this.GetPaged{{plural_class_name}}(filters, page, "display_index")
	{% else %}
	return this.GetPaged{{plural_class_name}}(filters, page, "id")
	{% endif %}
}

func (this *{{class_name}}Repository) Get{{plural_class_name}}ByIds(ids []int) []*{{class_name}} {
	filters := eel.Map{
		"id__in": ids,
	}
	{% if enable_display_index %}
	return this.Get{{plural_class_name}}(filters, "display_index")
	{% else %}
	return this.Get{{plural_class_name}}(filters, "id")
	{% endif %}
}

func (this *{{class_name}}Repository) Get{{class_name}}ById(id int) *{{class_name}} {
	filters := eel.Map{
		"id": id,
	}
	
	{{plural_var_name}} := this.Get{{plural_class_name}}(filters)
	
	if len({{plural_var_name}}) == 0 {
		return nil
	} else {
		return {{plural_var_name}}[0]
	}
}

func init() {
}
