Feature: 删除{% if app_cn_name %}{{app_cn_name}}活动的{% endif %}{{cn_name}}
	Background:
		Given ginger登录系统
		When ginger创建公司
		"""
		[{
			"name": "Apple",
			"username": "jobs"
		}, {
			"name": "Microsoft",
			"username": "bill"
		}, {
			"name": "Facebook",
			"username": "tom"
		}]
		"""

	@{{lint_service_name}} @{{app_name}}
	Scenario: 管理员删除{% if app_cn_name %}{{app_cn_name}}活动的{% endif %}{{cn_name}}

		# 创建{{cn_name}}
		Given jobs登录系统
		{%- for refer in refers %}
		{%- if refer.enable_fill_object or refer.enable_fill_objects %}
		When jobs创建{% if app_cn_name %}{{app_cn_name}}活动的{% endif %}{{refer.resource.cn_name}}
			"""
			[{
				{%- for field in refer.resource.fields %}
				"{{ field.snake_name }}": {{ field.get_bdd_default_value(refer.resource.name, 1) }}{{ "," if not loop.last -}}
				{%- endfor %}
			}]
			"""
		{%- endif -%}
		{%- endfor %}
		When jobs创建{% if app_cn_name %}{{app_cn_name}}活动的{% endif %}{{cn_name}}
		"""
		[
		{%- for index in [1, 2, 3] -%}
		{
			{%- for refer in refers %}
			{%- if refer.enable_fill_object %}
			"{{ refer.resource.name }}": {{ refer.resource.name_field.get_bdd_default_value(refer.resource.name, 1) }},
			{%- endif -%}
			{%- if refer.enable_fill_objects %}
			"{{ refer.resource.plural_name }}": [{{ refer.resource.name_field.get_bdd_default_list_values(refer.resource.name, 1) }}],
			{%- endif -%}
			{%- endfor %}
			{%- for field in fields %}
			"{{ field.snake_name }}": {{ field.get_bdd_default_value("", index) }}{{ "," if not loop.last -}}
			{%- endfor %}
		}{{ "," if not loop.last -}}
		{%- endfor -%}
		]
		"""
		Then jobs能看到{% if app_cn_name %}{{app_cn_name}}活动的{% endif %}{{cn_name}}列表
		"""
		[
		{%- for index in [1, 2, 3] -%}
		{
			"{{ name_field.snake_name }}": {{ name_field.get_bdd_default_value("", index) }}
		}{{ "," if not loop.last -}}
		{%- endfor -%}
		]
		"""

		#删除
		When jobs删除{% if app_cn_name %}{{app_cn_name}}活动的{% endif %}{{cn_name}}'{{ name_field.get_bdd_default_value("", 2, true) }}'
		When jobs删除{% if app_cn_name %}{{app_cn_name}}活动的{% endif %}{{cn_name}}'{{ name_field.get_bdd_default_value("", 3, true) }}'
		Then jobs能看到{% if app_cn_name %}{{app_cn_name}}活动的{% endif %}{{cn_name}}列表
		"""
		[
		{%- for index in [1] -%}
		{
			"{{ name_field.snake_name }}": {{ name_field.get_bdd_default_value("", index) }}
		}{{ "," if not loop.last -}}
		{%- endfor -%}
		]
		"""