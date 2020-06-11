package {{package}}

import (
	"time"

	"github.com/gingerxman/eel"
)

{%- if app_resource %}

const APP_STATUS_NOT_START = 0
const APP_STATUS_ONGOING = 1
const APP_STATUS_FINISHED = 2
var APPSTATUS2STR = map[int]string{
	APP_STATUS_NOT_START: "not_start",
	APP_STATUS_ONGOING: "ongoing",
	APP_STATUS_FINISHED: "finished",
}
var STR2APPSTATUS = map[string]int{
	"not_start": APP_STATUS_NOT_START,
	"ongoing": APP_STATUS_ONGOING,
	"finished": APP_STATUS_FINISHED,
}
//App 应用
type App struct {
	//内建field，谨慎修改
	eel.Model
	CorpId int `gorm:"index"`
	Status int `gorm:"index"`
	Name string
	IsDeleted bool //是否删除
	IsEnabled bool //是否开启
	CreatedAt  time.Time
	FinishedAt  time.Time

	//detail字段
	{%- for field in app_resource.fields %}
	{{ field.name }} {{ field.type }} {{ field.orm_annotaton -}}
	{% endfor %}
}
func (this *App) TableName() string {
	return "{{table_prefix}}_app"
}


{%- for refer in app_resource.refers %}
{%- if refer.type == 'n-n' and refer.should_generate_n2n_table %}


type AppHas{{refer.resource.class_name}} struct {
	eel.Model
	AppId int `gorm:"index"`//foreign key app
	{{refer.resource.class_name}}Id int //foreign key {{refer.resource.name}}
}
func (self *AppHas{{refer.resource.class_name}}) TableName() string {
	return "{{table_prefix}}_app_has_{{refer.resource.name}}"
}
{%- endif %}
{%- endfor %}


//Product 应用中的商品
type Product struct {
	eel.Model
	CorpId int `gorm:"index:corpid_isenable"`
	Name string `gorm:"size:125;index"`
	Image string
	Price float64
	{%- for field in app_product_resource.fields %}
	{{ field.name }} {{ field.type }} {{ field.orm_annotaton -}}
	{% endfor %}
	IsEnabled bool `gorm:"default:true;index:corpid_isenable"`
	IsDeleted bool `gorm:"default:false"`
	DisplayIndex int `gorm:"index"`//显示时的排序
	OriginalDisplayIndex int //置顶置底前的原始排序
}
func (this *Product) TableName() string {
	return "{{table_prefix}}_product"
}
//func (this *Product) TableIndex() [][]string {
//	return [][]string{
//		[]string{"CorpId", "IsDeleted", "IsEnabled"},
//	}
//}


const ACTIVITY_ORDER_NOT_PAY = 0
const ACTIVITY_ORDER_PULLING = 1
const ACTIVITY_ORDER_PAY_SUCCESS = 2
const ACTIVITY_ORDER_PAY_FAIL = 3
//Activity 参与应用
type Activity struct {
	//内建field，谨慎修改
	eel.Model
	CorpId int `gorm:"index"`
	AppId int `gorm:"index:appid_user"`
	UserId int `gorm:"index:appid_user"`
	OrderBid string `gorm:"index;default:''"`
	OrderStatus int `gorm:"index"`
	PullOrderTimes int
	PullFailReason string `gorm:"size:1024"`
	IsSendDeadMessage bool //是否发送了dead activity的钉钉通知
	LastPullOrderAt time.Time
	ProductId int `gorm:"index"`
	ProductPrice float64
	Score int //用户参与的分数（一般为摇动手机的次数）
	Ranking int //用户参与的最终排名
	IsCashed bool //如果活动有奖品，标记是否已经兑付了该笔参与
	IsDeleted bool `gorm:"default:false"`
	DeletedAt time.Time
	
	//detail字段
	{%- for field in app_activity_resource.fields %}
	{{ field.name }} {{ field.type }} {{ field.orm_annotation -}}
	{% endfor %}
}
func (this *Activity) TableName() string {
	return "{{table_prefix}}_activity"
}
//func (this *Activity) TableIndex() [][]string {
//	return [][]string{
//		[]string{"AppId", "IsDeleted", "CreatedAt"},
//		[]string{"AppId", "UserId", "IsDeleted"},
//		[]string{"CorpId", "AppId", "IsDeleted"},
//	}
//}


{%- for refer in app_activity_resource.refers %}
{%- if refer.type == 'n-n' and refer.should_generate_n2n_table %}


type ActivityHas{{refer.resource.class_name}} struct {
	eel.Model
	ActivityId int `gorm:"index"` //foreign key activity
	{{refer.resource.class_name}}Id int //foreign key {{refer.resource.name}}
}
func (self *ActivityHas{{refer.resource.class_name}}) TableName() string {
	return "{{table_prefix}}_activity_has_{{refer.resource.name}}"
}
{%- endif %}
{%- endfor %}
{%- endif %}

{% for resource in resources %}
//{{resource.class_name}} Model
type {{resource.class_name}} struct {
	eel.Model
	{%- for refer in resource.refers %}
	{%- if refer.type == 'n-1' and refer.quantity == '1' %}
	{{ refer.resource.class_name }}Id int //foreign key {{refer.resource.name}}
	{%- endif %}
	{%- endfor %}
	{%- for field in resource.fields %}
	{{ field.name }} {{ field.type }} {{ field.orm_annotation -}}
	{%- endfor %}
	{%- if resource.enable_display_index %}
	DisplayIndex int `gorm:"index"`//显示时的排序
	OriginalDisplayIndex int //置顶置底前的原始排序
	{%- endif %}
}
func (self *{{resource.class_name}}) TableName() string {
	return "{{table_prefix}}_{{resource.name}}"
}

{%- for refer in resource.refers %}
{%- if refer.type == 'n-n' and refer.should_generate_n2n_table %}


type {{resource.class_name}}Has{{refer.resource.class_name}} struct {
	eel.Model
	{{resource.class_name}}Id int //foreign key {{resource.name}}
	{{refer.resource.class_name}}Id int //foreign key {{refer.resource.name}}
}
func (self *{{resource.class_name}}Has{{refer.resource.class_name}}) TableName() string {
	return "{{table_prefix}}_{{resource.name}}_has_{{refer.resource.name}}"
}
{%- endif %}
{%- endfor %}

{% endfor %}
func init() {
	{% if not ignore_app -%}
	eel.RegisterModel(new(App))
	{%- for refer in app_resource.refers %}
	{%- if refer.type == 'n-n' and refer.should_generate_n2n_table %}
	eel.RegisterModel(new(AppHas{{refer.resource.class_name}}))
	{%- endif %}
	{%- endfor %}
	eel.RegisterModel(new(Product))
	eel.RegisterModel(new(Activity))
	{%- for refer in app_activity_resource.refers %}
	{%- if refer.type == 'n-n' and refer.should_generate_n2n_table %}
	eel.RegisterModel(new(ActivityHas{{refer.resource.class_name}}))
	{%- endif %}
	{%- endfor %}
	{% endif -%}
	{% for resource in resources -%}
	eel.RegisterModel(new({{resource.class_name}}))
	{%- for refer in resource.refers %}
	{%- if refer.type == 'n-n' and refer.should_generate_n2n_table %}
	eel.RegisterModel(new({{resource.class_name}}Has{{refer.resource.class_name}}))
	{%- endif %}
	{%- endfor %}
	{% endfor %}
}
