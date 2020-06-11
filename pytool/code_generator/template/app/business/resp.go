package {{package}}

{%- if not ignore_app %}

import (
	"{{service_name}}/business/account"
)

type RProduct struct {
	Id int    `json:"id"`
	Name string `json:"name"`
	Image string `json:"image"`
	Price float64 `json:"price"`
	IsEnabled bool   `json:"is_enabled"`
	IsDeleted bool   `json:"is_deleted"`
	{%- for field in app_product_resource.fields %}
	{{ field.name }} {{field.json_type}} `json:"{{field.snake_name}}"`{{""-}}
	{% endfor %}
	CreatedAt string `json:"created_at"`
}

type RActivityDetail struct {
	{%- for field in app_activity_resource.fields %}
	{{ field.name }} {{field.json_type}} `json:"{{field.snake_name}}"`{{""-}}
	{% endfor %}
	Product *RProduct `json:"product"`
}

type RActivityOrder struct {
	Bid string `json:"bid"`
	Status string `json:"status"`
}

type RActivity struct {
	Id        int    `json:"id"'`
	IsDeleted bool   `json:"is_deleted"`
	User    *account.RUser `json:"user"`
	Order *RActivityOrder `json:"order"`
	Detail *RActivityDetail `json:"detail"`
	Score int `json:"score"`
	Ranking int `json:"ranking"`
	IsCashed bool `json:"is_cashed"`
	{%- for refer in app_activity_resource.refers %}
	{%- if refer.enable_fill_nto1_n or refer.enable_fill_nton %}
	{{refer.resource.plural_class_name}} []*R{{refer.resource.class_name}} `json:"{{refer.resource.plural_name}}"`{{""-}}
	{%- endif -%}
	{%- endfor %}
	CreatedAt string `json:"created_at"`
}

type RActivityScore struct {
	Id int `json:"id"`
	UserId int `json:"user_id"`
	Name string `json:"name"`
	Thumbnail string `json:"thumbnail"`
	Score int `json:"score"`
	Rank int `json:"rank"`
}

type RAppDetail struct {
	{%- for field in app_resource.fields %}
	{{ field.name }} {{field.json_type}} `json:"{{field.snake_name}}"`{{""-}}
	{% endfor %}
}

type RApp struct {
	Id int `json:"id"'`
	Name string `json:"name"`
	Status string `json:"status"`
	IsDeleted bool `json:"is_deleted"`
	IsEnabled bool `json:"is_enabled"`
	Detail *RAppDetail `json:"detail"`
	Products []*RProduct `json:"products"`
	{%- for refer in app_resource.refers %}
	{%- if refer.enable_fill_nto1_n or refer.enable_fill_nton %}
	{{refer.resource.plural_class_name}} []*R{{refer.resource.class_name}} `json:"{{refer.resource.plural_name}}"`{{""-}}
	{%- endif -%}
	{%- endfor %}
	CreatedAt string `json:"created_at"`
}
{%- endif %}
{% for resource in resources %}
type R{{resource.class_name}} struct {
	Id int `json:"id"`
	{%- for field in resource.fields %}
	{{ field.name }} {{ field.json_type }} `json:"{{field.snake_name}}"`{{""-}}
	{%- endfor %}

	{%- for refer in resource.refers %}
	{%- if refer.enable_fill_nto1_1 %}
	{{refer.resource.class_name}} *R{{refer.resource.class_name}} `json:"{{refer.resource.name}}"`{{""-}}
	{%- endif %}
	{%- if refer.enable_fill_objects %}
	{{refer.resource.plural_class_name}} []*R{{refer.resource.class_name}} `json:"{{refer.resource.plural_name}}"`{{""-}}
	{%- endif -%}
	{%- endfor %}
	{%- if resource.enable_display_index %}
	DisplayIndex int `json:"display_index"`
	{%- endif %}
	CreatedAt string `json:"created_at"`
}
{% endfor %}

func init() {
}
