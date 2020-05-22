Feature: 更新{% if app_cn_name %}{{app_cn_name}}活动的{% endif %}{{cn_name}}的排序

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

@{{lint_service_name}} @{{app_name}} @di
Scenario: 向上向下调整{% if app_cn_name %}{{app_cn_name}}活动的{% endif %}{{cn_name}}的排序
	Given jobs登录系统
	{%- for refer in refers %}
	{%- if refer.enable_fill_object or refer.enable_fill_objects %}
	When jobs创建{% if app_cn_name %}{{app_cn_name}}活动的{% endif %}{{refer.resource.cn_name}}
		"""
		[
		{%- for index in [1, 2, 3] -%}
		{
			{%- for field in refer.resource.fields %}
			"{{ field.snake_name }}": {{ field.get_bdd_default_value(refer.resource.name, index) }}{{ "," if not loop.last -}}
			{%- endfor %}
		}{{ "," if not loop.last -}}
		{%- endfor -%}
		]
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
			"{{ refer.resource.name }}": {{ refer.resource.name_field.get_bdd_default_value(refer.resource.name, index) }},
			{%- endif -%}
			{%- if refer.enable_fill_objects %}
			"{{ refer.resource.plural_name }}": [{{ refer.resource.name_field.get_bdd_default_list_values(refer.resource.name, index) }}],
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
		}{{ "," if not loop.last }}
		{%- endfor -%}
		]
		"""

	#向上调整
	When jobs修改{% if app_cn_name %}{{app_cn_name}}活动的{% endif %}{{cn_name}}'{{ name_field.get_bdd_default_value("", 2, true) }}'的排序
		"""
		{
			"action": "up"
		}
		"""	
	Then jobs能看到{% if app_cn_name %}{{app_cn_name}}活动的{% endif %}{{cn_name}}列表
		"""
		[
		{%- for index in [2, 1, 3] -%}
		{
			"{{ name_field.snake_name }}": {{ name_field.get_bdd_default_value("", index) }}
		}{{ "," if not loop.last }}
		{%- endfor -%}
		]
		"""
	#向上调整顶部数据，无变化
	When jobs修改{% if app_cn_name %}{{app_cn_name}}活动的{% endif %}{{cn_name}}'{{ name_field.get_bdd_default_value("", 2, true) }}'的排序
		"""
		{
			"action": "up"
		}
		"""	
	Then jobs能看到{% if app_cn_name %}{{app_cn_name}}活动的{% endif %}{{cn_name}}列表
		"""
		[
		{%- for index in [2, 1, 3] -%}
		{
			"{{ name_field.snake_name }}": {{ name_field.get_bdd_default_value("", index) }}
		}{{ "," if not loop.last }}
		{%- endfor -%}
		]
		"""	

	#向下调整
	When jobs修改{% if app_cn_name %}{{app_cn_name}}活动的{% endif %}{{cn_name}}'{{ name_field.get_bdd_default_value("", 1, true) }}'的排序
		"""
		{
			"action": "down"
		}
		"""	
	Then jobs能看到{% if app_cn_name %}{{app_cn_name}}活动的{% endif %}{{cn_name}}列表
		"""
		[
		{%- for index in [2, 3, 1] -%}
		{
			"{{ name_field.snake_name }}": {{ name_field.get_bdd_default_value("", index) }}
		}{{ "," if not loop.last }}
		{%- endfor -%}
		]
		"""
	#向下调整底部数据，无变化
	When jobs修改{% if app_cn_name %}{{app_cn_name}}活动的{% endif %}{{cn_name}}'{{ name_field.get_bdd_default_value("", 1, true) }}'的排序
		"""
		{
			"action": "down"
		}
		"""	
	Then jobs能看到{% if app_cn_name %}{{app_cn_name}}活动的{% endif %}{{cn_name}}列表
		"""
		[
		{%- for index in [2, 3, 1] -%}
		{
			"{{ name_field.snake_name }}": {{ name_field.get_bdd_default_value("", index) }}
		}{{ "," if not loop.last }}
		{%- endfor -%}
		]
		"""	

@{{lint_service_name}} @{{app_name}} @di
Scenario: 将{% if app_cn_name %}{{app_cn_name}}活动的{% endif %}{{cn_name}}放到顶部
	Given jobs登录系统
	{%- for refer in refers %}
	{%- if refer.enable_fill_object or refer.enable_fill_objects %}
	When jobs创建{% if app_cn_name %}{{app_cn_name}}活动的{% endif %}{{refer.resource.cn_name}}
		"""
		[
		{%- for index in [1, 2, 3] -%}
		{
			{%- for field in refer.resource.fields %}
			"{{ field.snake_name }}": {{ field.get_bdd_default_value(refer.resource.name, index) }}{{ "," if not loop.last -}}
			{%- endfor %}
		}{{ "," if not loop.last -}}
		{%- endfor -%}
		]
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
			"{{ refer.resource.name }}": {{ refer.resource.name_field.get_bdd_default_value(refer.resource.name, index) }},
			{%- endif -%}
			{%- if refer.enable_fill_objects %}
			"{{ refer.resource.plural_name }}": [{{ refer.resource.name_field.get_bdd_default_list_values(refer.resource.name, index) }}],
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
		}{{ "," if not loop.last }}
		{%- endfor -%}
		]
		"""

	When jobs修改{% if app_cn_name %}{{app_cn_name}}活动的{% endif %}{{cn_name}}'{{ name_field.get_bdd_default_value("", 3, true) }}'的排序
		"""
		{
			"action": "top"
		}
		"""	
	Then jobs能看到{% if app_cn_name %}{{app_cn_name}}活动的{% endif %}{{cn_name}}列表
		"""
		[
		{%- for index in [3, 1, 2] -%}
		{
			"{{ name_field.snake_name }}": {{ name_field.get_bdd_default_value("", index) }}
		}{{ "," if not loop.last }}
		{%- endfor -%}
		]
		"""	

	When jobs修改{% if app_cn_name %}{{app_cn_name}}活动的{% endif %}{{cn_name}}'{{ name_field.get_bdd_default_value("", 2, true) }}'的排序
		"""
		{
			"action": "top"
		}
		"""	
	Then jobs能看到{% if app_cn_name %}{{app_cn_name}}活动的{% endif %}{{cn_name}}列表
		"""
		[
		{%- for index in [2, 3, 1] -%}
		{
			"{{ name_field.snake_name }}": {{ name_field.get_bdd_default_value("", index) }}
		}{{ "," if not loop.last }}
		{%- endfor -%}
		]
		"""

	#向上调整会成功
	When jobs修改{% if app_cn_name %}{{app_cn_name}}活动的{% endif %}{{cn_name}}'{{ name_field.get_bdd_default_value("", 1, true) }}'的排序
		"""
		{
			"action": "up"
		}
		"""	
	Then jobs能看到{% if app_cn_name %}{{app_cn_name}}活动的{% endif %}{{cn_name}}列表
		"""
		[
		{%- for index in [2, 1, 3] -%}
		{
			"{{ name_field.snake_name }}": {{ name_field.get_bdd_default_value("", index) }}
		}{{ "," if not loop.last }}
		{%- endfor -%}
		]
		"""

@{{lint_service_name}} @{{app_name}} @di
Scenario: 将{% if app_cn_name %}{{app_cn_name}}活动的{% endif %}{{cn_name}}放到底部
	Given jobs登录系统
	{%- for refer in refers %}
	{%- if refer.enable_fill_object or refer.enable_fill_objects %}
	When jobs创建{% if app_cn_name %}{{app_cn_name}}活动的{% endif %}{{refer.resource.cn_name}}
		"""
		[
		{%- for index in [1, 2, 3] -%}
		{
			{%- for field in refer.resource.fields %}
			"{{ field.snake_name }}": {{ field.get_bdd_default_value(refer.resource.name, index) }}{{ "," if not loop.last -}}
			{%- endfor %}
		}{{ "," if not loop.last -}}
		{%- endfor -%}
		]
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
			"{{ refer.resource.name }}": {{ refer.resource.name_field.get_bdd_default_value(refer.resource.name, index) }},
			{%- endif -%}
			{%- if refer.enable_fill_objects %}
			"{{ refer.resource.plural_name }}": [{{ refer.resource.name_field.get_bdd_default_list_values(refer.resource.name, index) }}],
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
		}{{ "," if not loop.last }}
		{%- endfor -%}
		]
		"""

	When jobs修改{% if app_cn_name %}{{app_cn_name}}活动的{% endif %}{{cn_name}}'{{ name_field.get_bdd_default_value("", 1, true) }}'的排序
		"""
		{
			"action": "bottom"
		}
		"""	
	Then jobs能看到{% if app_cn_name %}{{app_cn_name}}活动的{% endif %}{{cn_name}}列表
		"""
		[
		{%- for index in [2, 3, 1] -%}
		{
			"{{ name_field.snake_name }}": {{ name_field.get_bdd_default_value("", index) }}
		}{{ "," if not loop.last }}
		{%- endfor -%}
		]
		"""	

	When jobs修改{% if app_cn_name %}{{app_cn_name}}活动的{% endif %}{{cn_name}}'{{ name_field.get_bdd_default_value("", 2, true) }}'的排序
		"""
		{
			"action": "bottom"
		}
		"""	
	Then jobs能看到{% if app_cn_name %}{{app_cn_name}}活动的{% endif %}{{cn_name}}列表
		"""
		[
		{%- for index in [3, 1, 2] -%}
		{
			"{{ name_field.snake_name }}": {{ name_field.get_bdd_default_value("", index) }}
		}{{ "," if not loop.last }}
		{%- endfor -%}
		]
		"""

	#向下调整会成功
	When jobs修改{% if app_cn_name %}{{app_cn_name}}活动的{% endif %}{{cn_name}}'{{ name_field.get_bdd_default_value("", 3, true) }}'的排序
		"""
		{
			"action": "down"
		}
		"""	
	Then jobs能看到{% if app_cn_name %}{{app_cn_name}}活动的{% endif %}{{cn_name}}列表
		"""
		[
		{%- for index in [1, 3, 2] -%}
		{
			"{{ name_field.snake_name }}": {{ name_field.get_bdd_default_value("", index) }}
		}{{ "," if not loop.last }}
		{%- endfor -%}
		]
		"""

@{{lint_service_name}} @{{app_name}} @di
Scenario: 将{% if app_cn_name %}{{app_cn_name}}活动的{% endif %}{{cn_name}}置顶
	Given jobs登录系统
	{%- for refer in refers %}
	{%- if refer.enable_fill_object or refer.enable_fill_objects %}
	When jobs创建{% if app_cn_name %}{{app_cn_name}}活动的{% endif %}{{refer.resource.cn_name}}
		"""
		[
		{%- for index in [1, 2, 3, 4] -%}
		{
			{%- for field in refer.resource.fields %}
			"{{ field.snake_name }}": {{ field.get_bdd_default_value(refer.resource.name, index) }}{{ "," if not loop.last -}}
			{%- endfor %}
		}{{ "," if not loop.last -}}
		{%- endfor -%}
		]
		"""
	{%- endif -%}
	{%- endfor %}
	When jobs创建{% if app_cn_name %}{{app_cn_name}}活动的{% endif %}{{cn_name}}
		"""
		[
		{%- for index in [1, 2, 3, 4] -%}
		{
			{%- for refer in refers %}
			{%- if refer.enable_fill_object %}
			"{{ refer.resource.name }}": {{ refer.resource.name_field.get_bdd_default_value(refer.resource.name, index) }},
			{%- endif -%}
			{%- if refer.enable_fill_objects %}
			"{{ refer.resource.plural_name }}": [{{ refer.resource.name_field.get_bdd_default_list_values(refer.resource.name, index) }}],
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
		{%- for index in [1, 2, 3, 4] -%}
		{
			"{{ name_field.snake_name }}": {{ name_field.get_bdd_default_value("", index) }}
		}{{ "," if not loop.last }}
		{%- endfor -%}
		]
		"""

	#验证置顶
	When jobs修改{% if app_cn_name %}{{app_cn_name}}活动的{% endif %}{{cn_name}}'{{ name_field.get_bdd_default_value("", 4, true) }}'的排序
		"""
		{
			"action": "stick_top"
		}
		"""	
	Then jobs能看到{% if app_cn_name %}{{app_cn_name}}活动的{% endif %}{{cn_name}}列表
		"""
		[
		{%- for index in [4, 1, 2, 3] -%}
		{
			"{{ name_field.snake_name }}": {{ name_field.get_bdd_default_value("", index) }}
		}{{ "," if not loop.last }}
		{%- endfor -%}
		]
		"""	

	#放到顶部不会影响置顶
	When jobs修改{% if app_cn_name %}{{app_cn_name}}活动的{% endif %}{{cn_name}}'{{ name_field.get_bdd_default_value("", 3, true) }}'的排序
		"""
		{
			"action": "top"
		}
		"""	
	Then jobs能看到{% if app_cn_name %}{{app_cn_name}}活动的{% endif %}{{cn_name}}列表
		"""
		[
		{%- for index in [4, 3, 1, 2] -%}
		{
			"{{ name_field.snake_name }}": {{ name_field.get_bdd_default_value("", index) }}
		}{{ "," if not loop.last }}
		{%- endfor -%}
		]
		"""

@{{lint_service_name}} @{{app_name}} @di
Scenario: 将{% if app_cn_name %}{{app_cn_name}}活动的{% endif %}{{cn_name}}置底
	Given jobs登录系统
	{%- for refer in refers %}
	{%- if refer.enable_fill_object or refer.enable_fill_objects %}
	When jobs创建{% if app_cn_name %}{{app_cn_name}}活动的{% endif %}{{refer.resource.cn_name}}
		"""
		[
		{%- for index in [1, 2, 3, 4] -%}
		{
			{%- for field in refer.resource.fields %}
			"{{ field.snake_name }}": {{ field.get_bdd_default_value(refer.resource.name, index) }}{{ "," if not loop.last -}}
			{%- endfor %}
		}{{ "," if not loop.last -}}
		{%- endfor -%}
		]
		"""
	{%- endif -%}
	{%- endfor %}
	When jobs创建{% if app_cn_name %}{{app_cn_name}}活动的{% endif %}{{cn_name}}
		"""
		[
		{%- for index in [1, 2, 3, 4] -%}
		{
			{%- for refer in refers %}
			{%- if refer.enable_fill_object %}
			"{{ refer.resource.name }}": {{ refer.resource.name_field.get_bdd_default_value(refer.resource.name, index) }},
			{%- endif -%}
			{%- if refer.enable_fill_objects %}
			"{{ refer.resource.plural_name }}": [{{ refer.resource.name_field.get_bdd_default_list_values(refer.resource.name, index) }}],
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
		{%- for index in [1, 2, 3, 4] -%}
		{
			"{{ name_field.snake_name }}": {{ name_field.get_bdd_default_value("", index) }}
		}{{ "," if not loop.last }}
		{%- endfor -%}
		]
		"""

	#验证置底
	When jobs修改{% if app_cn_name %}{{app_cn_name}}活动的{% endif %}{{cn_name}}'{{ name_field.get_bdd_default_value("", 1, true) }}'的排序
		"""
		{
			"action": "stick_bottom"
		}
		"""	
	Then jobs能看到{% if app_cn_name %}{{app_cn_name}}活动的{% endif %}{{cn_name}}列表
		"""
		[
		{%- for index in [2, 3, 4, 1] -%}
		{
			"{{ name_field.snake_name }}": {{ name_field.get_bdd_default_value("", index) }}
		}{{ "," if not loop.last }}
		{%- endfor -%}
		]
		"""	

	#放到底部不会影响置底
	When jobs修改{% if app_cn_name %}{{app_cn_name}}活动的{% endif %}{{cn_name}}'{{ name_field.get_bdd_default_value("", 2, true) }}'的排序
		"""
		{
			"action": "bottom"
		}
		"""	
	Then jobs能看到{% if app_cn_name %}{{app_cn_name}}活动的{% endif %}{{cn_name}}列表
		"""
		[
		{%- for index in [3, 4, 2, 1] -%}
		{
			"{{ name_field.snake_name }}": {{ name_field.get_bdd_default_value("", index) }}
		}{{ "," if not loop.last }}
		{%- endfor -%}
		]
		"""

@{{lint_service_name}} @{{app_name}} @di
Scenario: 将{% if app_cn_name %}{{app_cn_name}}活动的{% endif %}{{cn_name}}取消置顶或取消置底
	Given jobs登录系统
	{%- for refer in refers %}
	{%- if refer.enable_fill_object or refer.enable_fill_objects %}
	When jobs创建{% if app_cn_name %}{{app_cn_name}}活动的{% endif %}{{refer.resource.cn_name}}
		"""
		[
		{%- for index in [1, 2, 3, 4] -%}
		{
			{%- for field in refer.resource.fields %}
			"{{ field.snake_name }}": {{ field.get_bdd_default_value(refer.resource.name, index) }}{{ "," if not loop.last -}}
			{%- endfor %}
		}{{ "," if not loop.last -}}
		{%- endfor -%}
		]
		"""
	{%- endif -%}
	{%- endfor %}
	When jobs创建{% if app_cn_name %}{{app_cn_name}}活动的{% endif %}{{cn_name}}
		"""
		[
		{%- for index in [1, 2, 3, 4] -%}
		{
			{%- for refer in refers %}
			{%- if refer.enable_fill_object %}
			"{{ refer.resource.name }}": {{ refer.resource.name_field.get_bdd_default_value(refer.resource.name, index) }},
			{%- endif -%}
			{%- if refer.enable_fill_objects %}
			"{{ refer.resource.plural_name }}": [{{ refer.resource.name_field.get_bdd_default_list_values(refer.resource.name, index) }}],
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
		{%- for index in [1, 2, 3, 4] -%}
		{
			"{{ name_field.snake_name }}": {{ name_field.get_bdd_default_value("", index) }}
		}{{ "," if not loop.last }}
		{%- endfor -%}
		]
		"""

	#置顶
	When jobs修改{% if app_cn_name %}{{app_cn_name}}活动的{% endif %}{{cn_name}}'{{ name_field.get_bdd_default_value("", 3, true) }}'的排序
		"""
		{
			"action": "stick_top"
		}
		"""	
	Then jobs能看到{% if app_cn_name %}{{app_cn_name}}活动的{% endif %}{{cn_name}}列表
		"""
		[
		{%- for index in [3, 1, 2, 4] -%}
		{
			"{{ name_field.snake_name }}": {{ name_field.get_bdd_default_value("", index) }}
		}{{ "," if not loop.last }}
		{%- endfor -%}
		]
		"""	
	#取消置顶
	When jobs修改{% if app_cn_name %}{{app_cn_name}}活动的{% endif %}{{cn_name}}'{{ name_field.get_bdd_default_value("", 3, true) }}'的排序
		"""
		{
			"action": "unstick"
		}
		"""	
	Then jobs能看到{% if app_cn_name %}{{app_cn_name}}活动的{% endif %}{{cn_name}}列表
		"""
		[
		{%- for index in [1, 2, 3, 4] -%}
		{
			"{{ name_field.snake_name }}": {{ name_field.get_bdd_default_value("", index) }}
		}{{ "," if not loop.last }}
		{%- endfor -%}
		]
		"""	

	#置底
	When jobs修改{% if app_cn_name %}{{app_cn_name}}活动的{% endif %}{{cn_name}}'{{ name_field.get_bdd_default_value("", 1, true) }}'的排序
		"""
		{
			"action": "stick_bottom"
		}
		"""	
	Then jobs能看到{% if app_cn_name %}{{app_cn_name}}活动的{% endif %}{{cn_name}}列表
		"""
		[
		{%- for index in [2, 3, 4, 1] -%}
		{
			"{{ name_field.snake_name }}": {{ name_field.get_bdd_default_value("", index) }}
		}{{ "," if not loop.last }}
		{%- endfor -%}
		]
		"""	
	#取消置底
	When jobs修改{% if app_cn_name %}{{app_cn_name}}活动的{% endif %}{{cn_name}}'{{ name_field.get_bdd_default_value("", 1, true) }}'的排序
		"""
		{
			"action": "unstick"
		}
		"""	
	Then jobs能看到{% if app_cn_name %}{{app_cn_name}}活动的{% endif %}{{cn_name}}列表
		"""
		[
		{%- for index in [1, 2, 3, 4] -%}
		{
			"{{ name_field.snake_name }}": {{ name_field.get_bdd_default_value("", index) }}
		}{{ "," if not loop.last }}
		{%- endfor -%}
		]
		"""	

@{{lint_service_name}} @{{app_name}}
Scenario: 删除{% if app_cn_name %}{{app_cn_name}}活动的{% endif %}{{product_cn_name}}不影响向上调整
	Given jobs登录系统
	{%- for refer in refers %}
	{%- if refer.enable_fill_object or refer.enable_fill_objects %}
	When jobs创建{% if app_cn_name %}{{app_cn_name}}活动的{% endif %}{{refer.resource.cn_name}}
		"""
		[
		{%- for index in [1, 2, 3] -%}
		{
			{%- for field in refer.resource.fields %}
			"{{ field.snake_name }}": {{ field.get_bdd_default_value(refer.resource.name, index) }}{{ "," if not loop.last -}}
			{%- endfor %}
		}{{ "," if not loop.last -}}
		{%- endfor -%}
		]
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
			"{{ refer.resource.name }}": {{ refer.resource.name_field.get_bdd_default_value(refer.resource.name, index) }},
			{%- endif -%}
			{%- if refer.enable_fill_objects %}
			"{{ refer.resource.plural_name }}": [{{ refer.resource.name_field.get_bdd_default_list_values(refer.resource.name, index) }}],
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
		}{{ "," if not loop.last }}
		{%- endfor -%}
		]
		"""

	#删除
	When jobs删除{% if app_cn_name %}{{app_cn_name}}活动的{% endif %}{{cn_name}}'{{ name_field.get_bdd_default_value("", 2, true) }}'
	#向上调整
	When jobs修改{% if app_cn_name %}{{app_cn_name}}活动的{% endif %}{{cn_name}}'{{ name_field.get_bdd_default_value("", 3, true) }}'的排序
		"""
		{
			"action": "up"
		}
		"""	
	#验证
	Then jobs能看到{% if app_cn_name %}{{app_cn_name}}活动的{% endif %}{{cn_name}}列表
		"""
		[
		{%- for index in [3, 1] -%}
		{
			"{{ name_field.snake_name }}": {{ name_field.get_bdd_default_value("", index) }}
		}{{ "," if not loop.last }}
		{%- endfor -%}
		]
		"""

@{{lint_service_name}} @{{app_name}}
Scenario: 删除{% if app_cn_name %}{{app_cn_name}}活动的{% endif %}{{product_cn_name}}不影响向下调整
	Given jobs登录系统
	{%- for refer in refers %}
	{%- if refer.enable_fill_object or refer.enable_fill_objects %}
	When jobs创建{% if app_cn_name %}{{app_cn_name}}活动的{% endif %}{{refer.resource.cn_name}}
		"""
		[
		{%- for index in [1, 2, 3] -%}
		{
			{%- for field in refer.resource.fields %}
			"{{ field.snake_name }}": {{ field.get_bdd_default_value(refer.resource.name, index) }}{{ "," if not loop.last -}}
			{%- endfor %}
		}{{ "," if not loop.last -}}
		{%- endfor -%}
		]
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
			"{{ refer.resource.name }}": {{ refer.resource.name_field.get_bdd_default_value(refer.resource.name, index) }},
			{%- endif -%}
			{%- if refer.enable_fill_objects %}
			"{{ refer.resource.plural_name }}": [{{ refer.resource.name_field.get_bdd_default_list_values(refer.resource.name, index) }}],
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
		}{{ "," if not loop.last }}
		{%- endfor -%}
		]
		"""

	#删除
	When jobs删除{% if app_cn_name %}{{app_cn_name}}活动的{% endif %}{{cn_name}}'{{ name_field.get_bdd_default_value("", 2, true) }}'
	#向上调整
	When jobs修改{% if app_cn_name %}{{app_cn_name}}活动的{% endif %}{{cn_name}}'{{ name_field.get_bdd_default_value("", 1, true) }}'的排序
		"""
		{
			"action": "down"
		}
		"""	
	#验证
	Then jobs能看到{% if app_cn_name %}{{app_cn_name}}活动的{% endif %}{{cn_name}}列表
		"""
		[
		{%- for index in [3, 1] -%}
		{
			"{{ name_field.snake_name }}": {{ name_field.get_bdd_default_value("", index) }}
		}{{ "," if not loop.last }}
		{%- endfor -%}
		]
		"""

@{{lint_service_name}} @{{app_name}}
Scenario: 多次置顶不影响新建{{product_cn_name}}的排序
	Given jobs登录系统
	{%- for refer in refers %}
	{%- if refer.enable_fill_object or refer.enable_fill_objects %}
	When jobs创建{% if app_cn_name %}{{app_cn_name}}活动的{% endif %}{{refer.resource.cn_name}}
		"""
		[
		{%- for index in [1, 2, 3] -%}
		{
			{%- for field in refer.resource.fields %}
			"{{ field.snake_name }}": {{ field.get_bdd_default_value(refer.resource.name, index) }}{{ "," if not loop.last -}}
			{%- endfor %}
		}{{ "," if not loop.last -}}
		{%- endfor -%}
		]
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
			"{{ refer.resource.name }}": {{ refer.resource.name_field.get_bdd_default_value(refer.resource.name, index) }},
			{%- endif -%}
			{%- if refer.enable_fill_objects %}
			"{{ refer.resource.plural_name }}": [{{ refer.resource.name_field.get_bdd_default_list_values(refer.resource.name, index) }}],
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
		}{{ "," if not loop.last }}
		{%- endfor -%}
		]
		"""

	#置顶
	When jobs修改{% if app_cn_name %}{{app_cn_name}}活动的{% endif %}{{cn_name}}'{{ name_field.get_bdd_default_value("", 2, true) }}'的排序
		"""
		{
			"action": "top"
		}
		"""	
	When jobs修改{% if app_cn_name %}{{app_cn_name}}活动的{% endif %}{{cn_name}}'{{ name_field.get_bdd_default_value("", 1, true) }}'的排序
		"""
		{
			"action": "top"
		}
		"""	
	When jobs修改{% if app_cn_name %}{{app_cn_name}}活动的{% endif %}{{cn_name}}'{{ name_field.get_bdd_default_value("", 2, true) }}'的排序
		"""
		{
			"action": "top"
		}
		"""	
	#新建
	When jobs创建{% if app_cn_name %}{{app_cn_name}}活动的{% endif %}{{cn_name}}
		"""
		[
		{%- for index in [4] -%}
		{
			{%- for refer in refers %}
			{%- if refer.enable_fill_object %}
			"{{ refer.resource.name }}": {{ refer.resource.name_field.get_bdd_default_value(refer.resource.name, index) }},
			{%- endif -%}
			{%- if refer.enable_fill_objects %}
			"{{ refer.resource.plural_name }}": [{{ refer.resource.name_field.get_bdd_default_list_values(refer.resource.name, index) }}],
			{%- endif -%}
			{%- endfor %}
			{%- for field in fields %}
			"{{ field.snake_name }}": {{ field.get_bdd_default_value("", index) }}{{ "," if not loop.last -}}
			{%- endfor %}
		}{{ "," if not loop.last -}}
		{%- endfor -%}
		]
		"""
	#验证
	Then jobs能看到{% if app_cn_name %}{{app_cn_name}}活动的{% endif %}{{cn_name}}列表
		"""
		[
		{%- for index in [2, 1, 3, 4] -%}
		{
			"{{ name_field.snake_name }}": {{ name_field.get_bdd_default_value("", index) }}
		}{{ "," if not loop.last }}
		{%- endfor -%}
		]
		"""

@{{lint_service_name}} @{{app_name}}
Scenario: 多次置底不影响新建{{product_cn_name}}的排序
	Given jobs登录系统
	{%- for refer in refers %}
	{%- if refer.enable_fill_object or refer.enable_fill_objects %}
	When jobs创建{% if app_cn_name %}{{app_cn_name}}活动的{% endif %}{{refer.resource.cn_name}}
		"""
		[
		{%- for index in [1, 2, 3] -%}
		{
			{%- for field in refer.resource.fields %}
			"{{ field.snake_name }}": {{ field.get_bdd_default_value(refer.resource.name, index) }}{{ "," if not loop.last -}}
			{%- endfor %}
		}{{ "," if not loop.last -}}
		{%- endfor -%}
		]
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
			"{{ refer.resource.name }}": {{ refer.resource.name_field.get_bdd_default_value(refer.resource.name, index) }},
			{%- endif -%}
			{%- if refer.enable_fill_objects %}
			"{{ refer.resource.plural_name }}": [{{ refer.resource.name_field.get_bdd_default_list_values(refer.resource.name, index) }}],
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
		}{{ "," if not loop.last }}
		{%- endfor -%}
		]
		"""

	#置顶
	When jobs修改{% if app_cn_name %}{{app_cn_name}}活动的{% endif %}{{cn_name}}'{{ name_field.get_bdd_default_value("", 2, true) }}'的排序
		"""
		{
			"action": "bottom"
		}
		"""	
	When jobs修改{% if app_cn_name %}{{app_cn_name}}活动的{% endif %}{{cn_name}}'{{ name_field.get_bdd_default_value("", 3, true) }}'的排序
		"""
		{
			"action": "bottom"
		}
		"""	
	When jobs修改{% if app_cn_name %}{{app_cn_name}}活动的{% endif %}{{cn_name}}'{{ name_field.get_bdd_default_value("", 2, true) }}'的排序
		"""
		{
			"action": "bottom"
		}
		"""	
	#新建
	When jobs创建{% if app_cn_name %}{{app_cn_name}}活动的{% endif %}{{cn_name}}
		"""
		[
		{%- for index in [4] -%}
		{
			{%- for refer in refers %}
			{%- if refer.enable_fill_object %}
			"{{ refer.resource.name }}": {{ refer.resource.name_field.get_bdd_default_value(refer.resource.name, index) }},
			{%- endif -%}
			{%- if refer.enable_fill_objects %}
			"{{ refer.resource.plural_name }}": [{{ refer.resource.name_field.get_bdd_default_list_values(refer.resource.name, index) }}],
			{%- endif -%}
			{%- endfor %}
			{%- for field in fields %}
			"{{ field.snake_name }}": {{ field.get_bdd_default_value("", index) }}{{ "," if not loop.last -}}
			{%- endfor %}
		}{{ "," if not loop.last -}}
		{%- endfor -%}
		]
		"""
	#验证
	Then jobs能看到{% if app_cn_name %}{{app_cn_name}}活动的{% endif %}{{cn_name}}列表
		"""
		[
		{%- for index in [1, 3, 2, 4] -%}
		{
			"{{ name_field.snake_name }}": {{ name_field.get_bdd_default_value("", index) }}
		}{{ "," if not loop.last }}
		{%- endfor -%}
		]
		"""		