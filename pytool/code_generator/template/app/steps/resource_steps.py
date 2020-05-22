# -*- coding: utf-8 -*-
import json

from behave import *

from features.bdd import util as bdd_util

def get_{{name}}_id_by_name(name):
	objs = bdd_util.exec_sql("select * from {{table_prefix}}_{{name}} where name = %s", [name])
	return objs[0]['id']


@Then(u"{corp_user}能看到{% if app_cn_name %}{{app_cn_name}}活动的{% endif %}{{cn_name}}列表")
def step_impl(context, corp_user):
	expected = json.loads(context.text)
	resp = context.client.get("{{app_name}}.corp_{{plural_name}}")
	actual = resp.data["{{plural_name}}"]

	{%- for refer in refers %}
	{%- if refer.enable_fill_object %}

	for item in actual:
		item['{{ refer.resource.name }}'] = item['{{ refer.resource.name }}']['{{ refer.resource.name_field.snake_name }}']
	{%- endif -%}

	{%- if refer.enable_fill_objects %}

	for item in actual:
		item['{{ refer.resource.plural_name }}'] = map(lambda x: x['{{ refer.resource.name_field.snake_name }}'], item['{{ refer.resource.plural_name }}'])
	{%- endif -%}
	{%- endfor %}

	bdd_util.assert_api_call_success(resp)
	bdd_util.assert_list(expected, actual)

@When(u"{corp_user}创建{% if app_cn_name %}{{app_cn_name}}活动的{% endif %}{{cn_name}}")
def step_impl(context, corp_user):
	datas = json.loads(context.text)
	for data in datas:
		{%- for refer in refers %}
		import {{ refer.resource.name }}_steps
		{%- if refer.enable_fill_object %}
		data["{{ refer.resource.name }}_id"] = {{ refer.resource.name }}_steps.get_{{ refer.resource.name }}_id_by_name(data["{{ refer.resource.name }}"])
		{%- endif %}

		{%- if refer.enable_fill_objects %}
		{{refer.resource.name}}_ids = []
		for {{refer.resource.name}}_name in data['{{refer.resource.plural_name}}']:
			{{refer.resource.name}}_ids.append({{ refer.resource.name }}_steps.get_{{ refer.resource.name }}_id_by_name({{refer.resource.name}}_name))
		data['{{refer.resource.name}}_ids'] = json.dumps({{refer.resource.name}}_ids)
		{%- endif %}
		{% endfor %}
		resp = context.client.put("{{app_name}}.{{name}}", data)
		bdd_util.assert_api_call_success(resp)

@When(u"{corp_user}删除{% if app_cn_name %}{{app_cn_name}}活动的{% endif %}{{cn_name}}'{name}'")
def step_impl(context, corp_user, name):
	id = get_{{name}}_id_by_name(name)
	resp = context.client.delete("{{app_name}}.{{name}}", {"id": id})
	bdd_util.assert_api_call_success(resp)

@When(u"{corp_user}修改{% if app_cn_name %}{{app_cn_name}}活动的{% endif %}{{cn_name}}'{name}'的信息")
def step_impl(context, corp_user, name):
	params = json.loads(context.text)
	id = get_{{name}}_id_by_name(name)
	params['id'] = id
	{%- for refer in refers %}

	import {{ refer.resource.name }}_steps 
	{%- if refer.enable_fill_object %}
	params["{{ refer.resource.name }}_id"] = {{ refer.resource.name }}_steps.get_{{ refer.resource.name }}_id_by_name(params["{{ refer.resource.name }}"])
	{%- endif %}

	{%- if refer.enable_fill_objects %}
	{{refer.resource.name}}_ids = []
	for {{refer.resource.name}}_name in params['{{refer.resource.plural_name}}']:
		{{refer.resource.name}}_ids.append({{ refer.resource.name }}_steps.get_{{ refer.resource.name }}_id_by_name({{refer.resource.name}}_name))
	params['{{refer.resource.name}}_ids'] = json.dumps({{refer.resource.name}}_ids)
	{%- endif %}
	{%- endfor %}

	resp = context.client.post("{{app_name}}.{{name}}", params)
	bdd_util.assert_api_call_success(resp)

@When(u"{corp_user}修改{% if app_cn_name %}{{app_cn_name}}活动的{% endif %}{{cn_name}}'{name}'的排序")
def step_impl(context, corp_user, name):
	params = json.loads(context.text)
	id = get_{{name}}_id_by_name(name)
	params['id'] = id
	params['action'] = json.loads(context.text)['action']

	resp = context.client.post("{{app_name}}.{{name}}_display_index", params)
	bdd_util.assert_api_call_success(resp)

@When(u"{corp_user}启用{% if app_cn_name %}{{app_cn_name}}活动的{% endif %}{{cn_name}}'{name}'")
def step_impl(context, corp_user, name):
	id = get_{{name}}_id_by_name(name)
	resp = context.client.delete("{{app_name}}.disabled_{{name}}", {"id": id})
	bdd_util.assert_api_call_success(resp)

@When(u"{corp_user}禁用{% if app_cn_name %}{{app_cn_name}}活动的{% endif %}{{cn_name}}'{name}'")
def step_impl(context, corp_user, name):
	id = get_{{name}}_id_by_name(name)
	resp = context.client.put("{{app_name}}.disabled_{{name}}", {"id": id})
	bdd_util.assert_api_call_success(resp)