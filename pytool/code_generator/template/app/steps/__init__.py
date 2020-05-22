# -*- coding: utf-8 -*-
{%- if app_resource %}
import app_steps
import product_steps
import activity_steps
{%- endif %}
{% for resource in resources %}
import {{resource.name}}_steps
{% endfor %}

