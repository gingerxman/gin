# -*- coding: utf-8 -*-

import os
import json
import shutil
import sys
import time
import requests
import zipfile
from datetime import datetime

code_zip_path = './codebase.zip'
TEMPLATE_FILE_DIR = '_gofile_template'

def snake2camel(name):
	"""
	将product_count转换为ProductCount
	:param self:
	:param name:
	:return:
	"""
	items = name.split('_')
	for i, item in enumerate(items):
		items[i] = item.capitalize()
	return ''.join(items)

def get_plural_name(name):
	"""
	将单数形式的名字转换为复数形式
	:param name:
	:return:
	"""
	plural_name = name
	if plural_name[-1] == 'y':
		plural_name = '%sies' % plural_name[:-1]
	elif plural_name[-1] == 's':
		plural_name = '%ses' % plural_name
	else:
		plural_name = '%ss' % plural_name

	return plural_name

def get_var_name(name):
	"""
	将ProductCount转换为productCount
	:param name:
	:return:
	"""
	return name[0].lower() + name[1:]

class Field(object):
	type2infos = {
		"str": {
			"go_type": "string",
			"json_type": "string",
			"orm_annatation": '',
			"json_must_func": 'MustString',
			"rest_type": "string"
		},
		"int": {
			"go_type": "int",
			"json_type": "int",
			"orm_annatation": '',
			"json_must_func": 'MustInt',
			"rest_type": "int"
		},
		"float": {
			"go_type": "float64",
			"json_type": "float64",
			"orm_annatation": '',
			"json_must_func": 'MustFloat64',
			"rest_type": "float"
		},
		"bool": {
			"go_type": "bool",
			"json_type": "bool",
			"orm_annatation": '`gorm:"default:false"`',
			"json_must_func": 'MustBool',
			"rest_type": "bool"
		},
		"time": {
			"go_type": "time.Time",
			"json_type": "string",
			"orm_annatation": '`gorm:"type:datetime"`',
			"json_must_func": 'MustString',
			"rest_type": "string"
		},
		"date": {
			"go_type": "time.Time",
			"json_type": "string",
			"orm_annatation": '`gorm:"type:date"`',
			"json_must_func": 'MustString',
			"rest_type": "string"
		},
	}

	def get_snake_name(self, name):
		buf = []
		for index, char in enumerate(name):
			val = ord(char)
			if val >= 65 and val <= 90:
				val = val + 32
				if index != 0:
					buf.append("_")
				buf.append(chr(val))
			else:
				buf.append(char)

		return "".join(buf)

	def get_default_value(self, username, default_index):
		if self.type == "int":
			return 1 + default_index
		elif self.type == "float64":
			return 1.9 + default_index
		elif self.type == "bool":
			return 'true'
		elif self.type == "time.Time":
			return '"%s 12:13:14"' % datetime.today().strftime("%Y-%m-%d")
		else:
			return '"%s_%s_%d"' % (username, self.snake_name, default_index)

	def get_bdd_default_value(self, prefix, default_index, trim_quotes=False):
		if self.type == "int":
			return 0 + default_index
		elif self.type == "float64":
			if default_index == 1:
				return 3
			elif default_index == 2:
				return 1.01
			elif default_index == 3:
				return 9.99
			else:
				return 1.9 + default_index
		elif self.type == "bool":
			return 'True'
		elif self.type == "time.Time":
			return '"%s 12:13:14"' % datetime.today().strftime("%Y-%m-%d")
		else:
			if prefix:
				value = '"%s_%s_%d"' % (prefix, self.snake_name, default_index)
			else:
				value = '"%s_%d"' % (self.snake_name, default_index)

			if trim_quotes:
				return value.replace('"', '')
			else:
				return value

	def get_bdd_default_list_values(self, prefix, default_index):
		if default_index == 1:
			return ''
		elif default_index == 2:
			if prefix:
				value = '"%s_%s_1"' % (prefix, self.snake_name)
			else:
				value = '"%s_1"' % self.snake_name

			return value
		elif default_index == 3:
			if prefix:
				value = '"%s_%s_1", "%s_%s_2"' % (prefix, self.snake_name, prefix, self.snake_name)
			else:
				value = '"%s_1", "%s_2"' % (self.snake_name, self.snake_name)

			return value
		elif default_index == 4:
			if prefix:
				value = '"%s_%s_2", "%s_%s_3"' % (prefix, self.snake_name, prefix, self.snake_name)
			else:
				value = '"%s_2", "%s_3"' % (self.snake_name, self.snake_name)

			return value
		else:
			return ''

	@property
	def py_default_value(self):
		if self.type == "int":
			return '""'
		elif self.type == "float64":
			return '""'
		elif self.type == "bool":
			return True
		elif self.type == "time.Time":
			return '"%s 00:00"' % datetime.today().strftime("%Y-%m-%d")
		else:
			return '""'

	@property
	def reactman_validate(self):
		if self.type == "int":
			return "require-positive-int"
		elif self.type == "float64":
			return "require-float"
		else:
			return "require-notempty"

	@property
	def reactman_control(self):
		if self._reactman_control == 'input':
			if self.type == 'bool':
				return 'bool_radio'
			elif self.type == "time.Time":
				return 'date_picker'
			else:
				return 'input'
		else:
			return self._reactman_control

	def __init__(self, scope, field_info):
		self.scope = scope
		self.name = field_info['name']
		self.snake_name = self.get_snake_name(self.name)
		self.var_name = get_var_name(self.name)
		self.is_name_field = field_info.get('is_name_field', False)
		self.valid_when_update = field_info.get('valid_when_update', True)
		self.valid_when_create = field_info.get('valid_when_create', True)
		self.is_data_owner_fk = field_info.get('is_data_owner_fk', False)
		self.default_value = field_info.get('default', None)

		try:
			filed_type = field_info['type']
		except KeyError:
			print '[ERROR] Field `%s` unspecified type' % self.name
			sys.exit(1)
		type_infos = self.type2infos[filed_type]
		self.type = type_infos["go_type"]
		self.json_type = type_infos["json_type"]
		self.rest_type = type_infos["rest_type"]
		self.json_must_func = type_infos['json_must_func']
		self._reactman_control = field_info.get('reactman_control', 'input')
		self.reactman_label = field_info.get('reactman_label', self.name)
		if 'db_type' in field_info:
			self.orm_annotation = '`gorm:"type:%s"`' % field_info['db_type']
		else:
			self.orm_annotation = field_info.get('orm_annatation', type_infos['orm_annatation'])
		self.meta_data = field_info

	def __repr__(self):
		return self.name

class Resource(object):
	def __init__(self, data):
		self.meta_type = data.get('meta_type', 'resource')
		self.name = data['name']
		self.plural_name = get_plural_name(self.name)
		self.cn_name = data['cn_name']
		self.enable_display_index = data.get('enable_display_index', False)
		self.data_owner = data.get('data_owner', '') # 数据的所有者: corp, user, platform

		class_name = snake2camel(self.name)
		self.class_name = class_name
		self.capital_class_name = class_name.upper()
		self.plural_class_name = get_plural_name(class_name)
		self.capital_plural_class_name = self.plural_class_name.upper()

		self.var_name = get_var_name(self.class_name)
		self.plural_var_name = get_var_name(self.plural_class_name)

		#解析fields
		self.fields = []
		for field_info in data['fields']:
			self.fields.append(Field("resource", field_info))

		#补充完整fields
		if self.belong_to_user:
			self.fields.insert(0, Field("resource", {
				"name": "UserId",
				"type": "int",
				"is_data_owner_fk": True,
				"valid_when_update": False,
				"valid_when_create": False,
				"orm_annatation": '`gorm:"index:userid_isdelete_isenable"`'
			}))
			self.fields.append(Field("resource", {
				"name": "IsEnabled",
				"type": "bool",
				"default": 'true',
				"valid_when_update": False,
				"valid_when_create": False,
				"orm_annatation": '`gorm:"index:userid_isdelete_isenable"`'
			}))
			self.fields.append(Field("resource", {
				"name": "IsDeleted",
				"type": "bool",
				"default": 'false',
				"valid_when_update": False,
				"valid_when_create": False,
				"orm_annatation": '`gorm:"index:userid_isdelete_isenable"`'
			}))
			self.fields[0].iface = {
				"var_name": "user",
				"type": "business.IUser"
			}
		if self.belong_to_corp:
			self.fields.insert(0, Field("resource", {
				"name": "CorpId",
				"type": "int",
				"is_data_owner_fk": True,
				"valid_when_update": False,
				"valid_when_create": False,
				"orm_annatation": '`gorm:"index:corpid_isdelete_isenable"`'
			}))
			self.fields.append(Field("resource", {
				"name": "IsEnabled",
				"type": "bool",
				"valid_when_update": False,
				"valid_when_create": False,
				"orm_annatation": '`gorm:"index:corpid_isdelete_isenable"`'
			}))
			self.fields.append(Field("resource", {
				"name": "IsDeleted",
				"type": "bool",
				"valid_when_update": False,
				"valid_when_create": False,
				"orm_annatation": '`gorm:"index:corpid_isdelete_isenable"`'
			}))
			self.fields[0].iface = {
				"var_name": "corp",
				"type": "business.ICorp"
			}

		#解析refer
		self.refers = data.get('refers')

	@property
	def belong_to_user(self):
		return self.data_owner == 'user'

	@property
	def belong_to_corp(self):
		return self.data_owner == 'corp'

	@property
	def has_name_field(self):
		for field_info in self.fields:
			if field_info.name == 'name':
				return True

		return False

	@property
	def should_import_model(self):
		for refer in self.refers:
			if refer['enable_fill_object'] or refer['enable_fill_objects']:
				return True
		return False

	@property
	def should_select_other_resource_in_reactman(self):
		for refer in self.refers:
			if refer['type'] == 'n-1' and refer['quantity'] == 'n':
				return True
		return False

	@property
	def name_field(self):
		for field in self.fields:
			if field.is_name_field:
				return field

		return None

	def to_dict(self):
		return {
			"meta_type": self.meta_type,
			"name": self.name,
			"plural_name": self.plural_name,
			"cn_name": self.cn_name,
			"var_name": self.var_name,
			"plural_var_name": self.plural_var_name,
			"class_name": self.class_name,
			"capital_class_name": self.capital_class_name,
			"plural_class_name": self.plural_class_name,
			"capital_plural_class_name": self.capital_plural_class_name,
			"fields": self.fields,
			"updatable_fields": filter(lambda field: field.valid_when_update and (not field.is_data_owner_fk), self.fields),
			"creatable_fields": filter(lambda field: field.valid_when_create and (not field.is_data_owner_fk), self.fields),
			"non_creatable_fields": filter(lambda field: (not field.valid_when_create) and (not field.is_data_owner_fk), self.fields),
			"data_owner_field": next( (field for field in self.fields if field.is_data_owner_fk), None),
			"enable_display_index": self.enable_display_index,
			"name_field": self.name_field,
			'has_name_field': self.has_name_field,
			'refers': self.refers,
			'should_import_model': self.should_import_model,
			'should_select_other_resource_in_reactman': self.should_select_other_resource_in_reactman,
			'belong_to_user': self.belong_to_user,
			'belong_to_corp': self.belong_to_corp
		}

	def __repr__(self):
		return '%s' % self.to_dict()


class AppInfo(object):
	def __init__(self):
		self.service_name = ''
		self.app_name = ''
		self.package = ''
		self.app_resource = None
		self.app_activity_resource = None
		self.app_product_resource = None
		self.app_cn_name = ''
		self.product_cn_name = ''
		self.resources = []

	def to_dict(self):
		return {
			"service_name": self.service_name,
			"lint_service_name": self.service_name.split('/')[-1],
			"table_prefix": self.app_name.replace("_", ""),
			"app_name": self.app_name,
			"ignore_app": self.app_resource is None,
			"app_cn_name": self.app_cn_name,
			"product_cn_name": self.product_cn_name,
			"package": self.package,
			"resources": self.resources,
			"app_extra_packages": self.app_extra_packages,
			"app_activity_extra_packages": self.app_activity_extra_packages,
			"app_resource": self.app_resource,
			"app_activity_resource": self.app_activity_resource,
			"app_product_resource": self.app_product_resource
		}

	@property
	def app_extra_packages(self):
		go_packages = []
		if self.app_resource:
			for field in self.app_resource.fields:
				if field.type == 'time.Time':
					go_packages.append("time")

		return go_packages

	@property
	def app_activity_extra_packages(self):
		go_packages = []
		if self.app_activity_resource:
			for field in self.app_activity_resource.fields:
				if field.type == 'time.Time':
					go_packages.append("time")

		return go_packages

	@staticmethod
	def build_resource_relation(resources):
		"""
		建立resource之间的引用关系
		:return:
		"""
		name2resource = {resource.name:resource for resource in resources}

		#将refer中的resource从字符串替换为Resource对象
		for resource in resources:
			if resource.refers:
				for refer_resource_info in resource.refers:
					resource_name = refer_resource_info['resource']

					#检查resource name
					if not resource_name in name2resource:
						print '[ERROR] invalid resource `%s`' % resource_name
						sys.exit(1)

					refer_resource_info['resource'] = name2resource[resource_name]
					refer_resource_info['resource_name'] = resource_name

					refer_type = refer_resource_info['type']
					items = refer_type.split('-')
					if len(items) != 2:
						print '[ERROR] invalid refer_type `%s`' % refer_type
						sys.exit(1)
					self_quantity = items[0]
					other_quantity = items[1]
					refer_resource_info['quantity'] = other_quantity

					if not 'present_when_create' in refer_resource_info:
						refer_resource_info['present_when_create'] = False
					if not 'present_when_update' in refer_resource_info:
						refer_resource_info['present_when_update'] = False

					enable_fill = refer_resource_info.get('enable_fill', False)
					is_nto1 = (refer_resource_info['type'] == '1-n') or (refer_resource_info['type'] == 'n-1')
					is_nton = (refer_resource_info['type'] == 'n-n')
					refer_resource_info['enable_fill_nto1_1'] = enable_fill and is_nto1 and (other_quantity == '1')
					refer_resource_info['enable_fill_nto1_n'] = enable_fill and is_nto1 and (other_quantity == 'n')
					refer_resource_info['enable_fill_nton'] = enable_fill and is_nton
					refer_resource_info['enable_fill_object'] = refer_resource_info['enable_fill_nto1_1']
					refer_resource_info['enable_fill_objects'] = (refer_resource_info['enable_fill_nto1_n'] or refer_resource_info['enable_fill_nton'])

					refer_resource_info['update_nto1_1'] = is_nto1 and (other_quantity == '1') and refer_resource_info['present_when_update']
					refer_resource_info['update_nto1_n'] = is_nto1 and (other_quantity == 'n') and refer_resource_info['present_when_update']
					refer_resource_info['update_nton'] = is_nton and refer_resource_info['present_when_update']

					refer_resource_info['create_nto1_1'] = is_nto1 and (other_quantity == '1') and refer_resource_info['present_when_create']
					refer_resource_info['create_nto1_n'] = is_nto1 and (other_quantity == 'n') and refer_resource_info['present_when_create']
					refer_resource_info['create_nton'] = is_nton and refer_resource_info['present_when_create']

					print '>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>'
					print resource.name, ' -> ', resource_name
					for key in ['type', 'quantity', 'enable_fill_nto1_1', 'enable_fill_nto1_n', 'enable_fill_nton', 'update_nto1_1', 'update_nto1_n', 'update_nton', 'create_nto1_1', 'create_nto1_n', 'create_nton']:
						print '%s: %s' % (key, refer_resource_info.get(key))
					print '>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>'


					if refer_resource_info['type'] == 'n-n':
						refer_resource_info['should_generate_n2n_table'] = True
					else:
						refer_resource_info['should_generate_n2n_table'] = False
			else:
				resource.refers = []

		# #构建resource之间的关系
		# for resource in resources:
		# 	for refer_resource_info in resource.refers:
		# 		#检查循环refer
		# 		resource_name = refer_resource_info['resource_name']
		# 		refer_resource = name2resource[resource_name]
		# 		if refer_resource.refers:
		# 			can_ignore_refer_resource = False
		# 			for refer_2_resource in refer_resource.refers:
		# 				if refer_2_resource['resource'].name == resource.name:
		# 					#发现循环refer，表明不用处理该refer_resource了
		# 					can_ignore_refer_resource = True
		# 					break
		#
		# 			if can_ignore_refer_resource:
		# 				continue
		#
		# 		#处理resource与refer_resource
		# 		refer_type = refer_resource_info['type']
		# 		items = refer_type.split('-')
		# 		if len(items) != 2:
		# 			print '[ERROR] invalid refer_type `%s`' % refer_type
		# 			sys.exit(1)
		# 		self_quantity = items[0]
		# 		other_quantity = items[1]
		# 		refer_resource_info['quantity'] = other_quantity
		# 		refer_resource_info['enable_fill_objects'] = False
		# 		refer_resource_info['enable_fill_object'] = False
		# 		refer_resource_info['present_when_create'] = False
		# 		reverse_present_when_create = False
		# 		refer_resource_info['present_when_update'] = False
		# 		reverse_present_when_update = False
		#
		# 		reverse_fill_n_in_1 = False
		# 		reverse_fill_1_in_n = False
		# 		reverse_fill_object = False
		# 		reverse_fill_objects = False
		# 		if not 'enable_fill_n_in_1' in refer_resource_info:
		# 			refer_resource_info['enable_fill_n_in_1'] = False
		# 		if not 'enable_fill_1_in_n' in refer_resource_info:
		# 			refer_resource_info['enable_fill_1_in_n'] = False
		#
		# 		if refer_type == 'n-1':
		# 			if self_quantity == 'n':
		# 				print '@1'
		# 				refer_resource_info['enable_fill_object'] = refer_resource_info.get('enable_fill_1_in_n', False)
		# 				refer_resource_info['enable_fill_objects'] = False
		# 				reverse_fill_n_in_1 = refer_resource_info.get('enable_fill_n_in_1', False)
		# 				reverse_fill_objects = reverse_fill_n_in_1
		#
		# 				refer_resource_info['present_when_create'] = True
		# 				reverse_present_when_create = refer_resource_info.get('present_when_create_1', False)
		# 				refer_resource_info['present_when_update'] = True
		# 				reverse_present_when_update = refer_resource_info.get('present_when_update_1', False)
		# 			else:
		# 				print '@2'
		# 				refer_resource_info['enable_fill_objects'] = refer_resource_info.get('enable_fill_n_in_1', False)
		# 				refer_resource_info['enable_fill_object'] = False
		# 				reverse_fill_1_in_n = refer_resource_info.get('enable_fill_1_in_n', False)
		# 				reverse_fill_object = reverse_fill_1_in_n
		# 		if refer_type == 'n-n':
		# 			refer_resource_info['enable_fill_object'] = False
		# 			refer_resource_info['enable_fill_objects'] = True
		# 			reverse_fill_objects = True
		#
		# 		if not 'is_relation_master' in refer_resource_info:
		# 			refer_resource_info['is_relation_master'] = False
		# 		refer_resource.refers.append({
		# 			'resource_name': resource.name,
		# 			'resource': name2resource[resource.name],
		# 			'type': refer_type,
		# 			'quantity': self_quantity,
		# 			'present_when_create': reverse_present_when_create,
		# 			'present_when_update': reverse_present_when_update,
		# 			'enable_fill_1_in_n': reverse_fill_1_in_n,
		# 			'enable_fill_n_in_1': reverse_fill_n_in_1,
		# 			'enable_fill_object': reverse_fill_object,
		# 			'enable_fill_objects': reverse_fill_objects,
		# 			'should_generate_n2n_table': False,
		# 			'is_relation_master': False
		# 		})

	@staticmethod
	def parse():
		if not os.path.exists("./app.json"):
			print '[error] You have no app.json!!!'
			return None

		app_json = None
		with open('./app.json', 'rb') as f:
			content = f.read()
			app_json = json.loads(content)

		app_info = AppInfo()
		app_info.app_name = app_json['name']
		app_info.package = app_json['name']

		# for field_info in app_json['app']['details']:
		# 	app_info.app_detail_fields.append(Field("app", field_info))

		# for field_info in app_json['app']['activity_details']:
		# 	app_info.activity_detail_fields.append(Field("activity", field_info))

		resources = []
		for resource_info in app_json['resources']:
			resource = Resource(resource_info)
			if resource.meta_type == "app":
				app_info.app_resource = resource
				app_info.app_cn_name = resource.cn_name
			elif resource.meta_type == "app_activity":
				app_info.app_activity_resource = resource
			elif resource.meta_type == "app_product":
				app_info.app_product_resource = resource
				app_info.product_cn_name = resource.cn_name
			else:
				app_info.resources.append(resource)
			resources.append(resource)

		AppInfo.build_resource_relation(resources)

		#get service name
		with open('./service.json', 'rb') as f:
			content = f.read()
			service_info = json.loads(content)
			app_info.service_name = service_info['name']

		return app_info

class Command(object):
	help = "gen_app"
	args = ''

	def confirm_dir_exists(self, dir):
		if not os.path.exists(dir):
			print 'make dir: ', dir
			os.makedirs(dir)
		else:
			print 'dir(%s) is already exists' % dir

	def check_file_exists(self, file_path):
		return os.path.exists(file_path)

	def generate_file(self, args):
		file_type = args['file_type']
		resource_templates = args.get('resource_templates', {})
		file_suffixs = args['file_suffixs']
		file_map = args['file_map']
		ignore_files = args['ignore_files']
		context = args['context']
		print '\n>>>>>>>>>> generate %s objects <<<<<<<<<<' % file_type
		target_dir_path = '_generate/%s' % file_type
		self.confirm_dir_exists(target_dir_path)

		template_path = os.path.join(TEMPLATE_FILE_DIR, "app", file_type)
		for file_name in os.listdir(template_path):
			if ignore_files and file_name in ignore_files:
				print 'skip: %s' % file_name
				continue

			match_one_file_suffix = False
			for file_suffix in file_suffixs:
				if file_name.endswith(file_suffix):
					match_one_file_suffix = True

			if not match_one_file_suffix:
				continue

			if resource_templates and file_name in resource_templates:
				#render extra resource files
				resources = context['resources']
				if len(resources) == 0:
					continue

				file_path = "app/%s/%s" % (file_type, file_name)
				for resource in resources:
					if not resource.enable_display_index and 'display_index' in file_name:
						continue
					#替换文件名
					file_name_context = {
						"resource.name": resource.name,
						"resource.plural_name": resource.plural_name
					}
					replace_pattern = resource_templates[file_name]
					target_path = os.path.join(target_dir_path, replace_pattern % file_name_context)

					#渲染文件
					resource_context = resource.to_dict()
					resource_context.update({
						"service_name": context['service_name'],
						"lint_service_name": context['lint_service_name'],
						"app_name": context['app_name'],
						"app_cn_name": context['app_cn_name'],
						"package": context['package'],
						"table_prefix": context['table_prefix']
					})
					self.render_file_to(file_path, target_path, resource_context)
			else:
				file_path = "app/%s/%s" % (file_type, file_name)
				if file_map and file_name in file_map:
					#替换文件名
					replace_pattern = file_map[file_name]
					target_path = os.path.join(target_dir_path, replace_pattern % context)
				else:
					target_path = os.path.join(target_dir_path, file_name)
				self.render_file_to(file_path, target_path, context)

	def generate_resource_ui_files(self, args):
		file_suffixs = args['file_suffixs']
		context = args['context']
		resources = context['resources']
		for resource in resources:
			print '\n>>>>>>>>>> generate reactman ui files for resource "%s" <<<<<<<<<<' % resource.name
			path_infos = [{
				"template_path": "app/reactman/ui/resource",
				"target_dir_path": os.path.join('_generate/reactman/ui', resource.name)
			}, {
				"template_path": "app/reactman/ui/resources",
				"target_dir_path": os.path.join('_generate/reactman/ui', resource.plural_name)
			}]
			for path_info in path_infos:
				target_dir_path = path_info['target_dir_path']
				self.confirm_dir_exists(target_dir_path)
				template_path = path_info['template_path']

				for file_name in os.listdir(os.path.join(TEMPLATE_FILE_DIR, template_path)):
					match_one_file_suffix = False
					for file_suffix in file_suffixs:
						if file_name.endswith(file_suffix):
							match_one_file_suffix = True

					if not match_one_file_suffix:
						continue

					file_path = os.path.join(template_path, file_name)
					target_path = os.path.join(target_dir_path, file_name)

					#渲染文件
					resource_context = resource.to_dict()
					resource_context.update({
						"service_name": context['service_name'],
						"lint_service_name": context['lint_service_name'],
						"app_name": context['app_name'],
						"package": context['package'],
					})
					self.render_file_to(file_path, target_path, resource_context)


	def render_file_to(self, template_name, target_dir_path, context):
		from jinja2 import Template
		print '> generate: %s\n\tfrom %s' % (target_dir_path, template_name)

		with open('%s/%s' % (TEMPLATE_FILE_DIR, template_name), 'rb') as f:
			template_content = f.read()

		template = Template(template_content.decode('utf-8'))
		content = template.render(context)

		with open(target_dir_path, 'wb') as f:
			print >> f, content.encode('utf-8')

	def copy_files(self, context):
		copy_infos = [{
			'src': 'models',
			'dst': 'models'
		}, {
			'src': 'business',
			'dst': 'business',
			'ignore': ['entity.go', 'fill_entity_service.go', 'encode_entity_service.go', 'entity_repository.go']
		}, {
			'src': 'rest',
			'dst': 'rest',
			'ignore': ['resource.go', 'corp_resources.go', 'user_resources.go', 'disabled_resource.go']
		}, {
			'src': 'features',
			'dst': 'features',
			'ignore': ['resource.py', 'resources.py']
		}, {
			'src': 'steps',
			'dst': 'features/steps'
		}]
		for copy_info in copy_infos:
			src_dir = os.path.join('_generate', copy_info['src'])
			if not os.path.exists(src_dir):
				print 'copy SKIP: ', src_dir
				continue
			dst_dir = os.path.join(copy_info['dst'], context['package'])
			print '> copy: %s -> %s' % (src_dir, dst_dir)
			if os.path.exists(dst_dir):
				shutil.rmtree(dst_dir)
			shutil.copytree(src_dir, dst_dir)

			ignores = copy_info.get('ignore')
			if ignores:
				for file_name in ignores:
					ignore_file = os.path.join(dst_dir, file_name)
					if os.path.exists(ignore_file):
						os.remove(ignore_file)

	def copy_reactman_files(self, context):
		paths = [{
			'src': './_generate/reactman/ui',
			'dst': ('/Users/chenru/xiaocheng/ceres/static/component/app/%s' % context['app_name']),
			'ignore': ['resource', 'resources']
		}, {
			'src': './_generate/reactman/python',
			'dst': ('/Users/chenru/xiaocheng/ceres/app/%s' % context['app_name'])
		}]

		for path in paths:
			src_dir = path['src']
			dst_dir = path['dst']
			print '> copy: %s -> %s' % (src_dir, dst_dir)
			if os.path.exists(dst_dir):
				shutil.rmtree(dst_dir)
			shutil.copytree(src_dir, dst_dir)

			ignores = path.get('ignore')
			if ignores:
				for dir_name in ignores:
					ignore_dir = os.path.join(dst_dir, dir_name)
					if os.path.exists(ignore_dir):
						shutil.rmtree(ignore_dir)

	def download_code_base(self, url, zipfile):
		total_bytes = 0
		with open(zipfile, 'wb') as handle:
			response = requests.get(url, stream=True)

			for block in response.iter_content(1024):
				total_bytes += len(block)
				print 'download......%sk' % round((total_bytes/1024.0), 2)
				handle.write(block)
				time.sleep(.001)

	def unzip_code_base_to(self, name):
		zfobj = zipfile.ZipFile(code_zip_path)
		zfobj.extractall()

		for dir in os.listdir('.'):
			if not os.path.isdir(dir):
				continue

			if not 'golang-service-resource-template' in dir:
				continue

			if os.path.exists(name):
				shutil.rmtree(name)

			print 'rename %s to %s' % (dir, name)
			os.rename(dir, name)

		zfobj.close()
		os.remove(code_zip_path)

	def download_template_files(self):
		print 'download Golang template file'
		code_base_url = 'https://code.aliyun.com/clubxiaocheng/golang-service-resource-template/repository/archive.zip'
		self.download_code_base(code_base_url, code_zip_path)
		self.unzip_code_base_to(TEMPLATE_FILE_DIR)

	def handle(self):
		generated_dir = './_generate'
		if os.path.exists(generated_dir):
			shutil.rmtree(generated_dir)

		app_info = AppInfo.parse()

		#self.download_template_files()

		app_info_dict = app_info.to_dict()
		self.generate_file({
			"file_type": "models",
			"file_suffixs": [".go",],
			"resource_templates":None,
			"file_map": {'model.go':'%(app_name)s.go'},
			"ignore_files": None,
			"context": app_info_dict
		})
		self.generate_file({
			"file_type": "business",
			"file_suffixs": [".go",],
			"resource_templates": {
				'entity.go': '%(resource.name)s.go',
				'fill_entity_service.go': 'fill_%(resource.name)s_service.go',
				'encode_entity_service.go': 'encode_%(resource.name)s_service.go',
				'entity_repository.go': '%(resource.name)s_repository.go'
			},
			"file_map": None,
			"ignore_files": [
				"app.go", "app_detail.go", "app_repository.go", 
				"fill_app_service.go", "encode_app_service.go",
				"activity.go", "activity_detail.go", "activity_repository.go",
				"fill_activity_service.go", "encode_activity_service.go",
				"create_app_strategy.go", "create_activity_strategy.go",
				"product.go", "product_repository.go", "fill_product_service.go", 
				"encode_product_service.go"] if not app_info.app_resource else None,
			"context": app_info_dict
		})
		self.generate_file({
			"file_type": "rest",
			"file_suffixs": [".go",],
			"resource_templates": {
				'resource.go': '%(resource.name)s.go',
				'resource_display_index.go': '%(resource.name)s_display_index.go',
				'resources.go': '%(resource.plural_name)s.go',
				'corp_resources.go': 'corp_%(resource.plural_name)s.go',
				'user_resources.go': 'user_%(resource.plural_name)s.go',
				'disabled_resource.go': 'disabled_%(resource.name)s.go'
			},
			"file_map": None,
			"ignore_files": [
				"app.go", "corp_apps.go", "app_status.go", "disabled_app.go", 
				"activity.go", "activities.go", "user_activities.go",
				"product.go", "corp_products.go", "products.go", "disabled_product.go"] if not app_info.app_resource else None,
			"context": app_info_dict
		})
		return


		#genearete features & steps
		feature_templates = {
			'create_resource.feature': 'create_%(resource.name)s.feature',
			'update_resource.feature': 'update_%(resource.name)s.feature',
			'update_resource_display_index.feature': 'update_%(resource.name)s_display_index.feature',
			'disable_resource.feature': 'disable_%(resource.name)s.feature',
			'delete_resource.feature': 'delete_%(resource.name)s.feature',
		}
		self.generate_file({
			"file_type": "features",
			"file_suffixs": [".feature",],
			"resource_templates": feature_templates,
			"file_map": None,
			"ignore_files": [
				"change_app_status.feature", "create_app.feature", "update_app.feature",
				"create_product.feature", "update_product.feature", "disable_product.feature",
				"create_activity.feature", "pay_activity.feature"] if not app_info.app_resource else None,
			"context": app_info_dict
		})
		self.generate_file({
			"file_type": "steps",
			"file_suffixs": [".py",],
			"resource_templates": {
				'resource_steps.py': '%(resource.name)s_steps.py'
			},
			"file_map": None,
			"ignore_files": [
				"app_steps.py", "product_steps.py", "activity_steps.py"
			] if not app_info.app_resource else None,
			"context": app_info_dict
		})

		#generate ui files
		# self.generate_file({
		# 	"file_type": "reactman/python",
		# 	"file_suffixs": [".py",],
		# 	"resource_templates": {
		# 		'resource.py': '%(resource.name)s.py',
		# 		'resources.py': '%(resource.plural_name)s.py'
		# 	},
		# 	"file_map": None,
		# 	"ignore_files": [
		# 		"apps.py", "one_app.py", "product.py", "products.py"
		# 	] if not app_info.app_resource else None,
		# 	"context": app_info_dict
		# })
		# if app_info.app_resource:
		# 	self.generate_file({
		# 		"file_type": "reactman/ui/apps",
		# 		"file_suffixs": [".js", ".scss"],
		# 		"resource_templates": None,
		# 		"file_map": None,
		# 		"ignore_files": None,
		# 		"context": app_info_dict
		# 	})
		# 	self.generate_file({
		# 		"file_type": "reactman/ui/one_app",
		# 		"file_suffixs": [".js", ".scss"],
		# 		"resource_templates": None,
		# 		"file_map": None,
		# 		"ignore_files": None,
		# 		"context": app_info_dict
		# 	})
		# 	self.generate_file({
		# 		"file_type": "reactman/ui/products",
		# 		"file_suffixs": [".js", ".scss"],
		# 		"resource_templates": None,
		# 		"file_map": None,
		# 		"ignore_files": None,
		# 		"context": app_info_dict
		# 	})
		# 	self.generate_file({
		# 		"file_type": "reactman/ui/product",
		# 		"file_suffixs": [".js", ".scss"],
		# 		"resource_templates": None,
		# 		"file_map": None,
		# 		"ignore_files": None,
		# 		"context": app_info_dict
		# 	})
		# self.generate_resource_ui_files({
		# 	"file_type": None,
		# 	"file_suffixs": [".js", ".scss"],
		# 	"resource_templates": None,
		# 	"file_map": None,
		# 	"ignore_files": None,
		# 	"context": app_info_dict
		# })
		self.generate_file({
			"file_type": "copy",
			"file_suffixs": [".txt"],
			"resource_templates": None,
			"file_map": None,
			"ignore_files": None,
			"context": app_info_dict
		})

		return

		print '\n******************** Generate File ********************'
		print 'file is generated under ./_generate dir, please copy to real dirs'
		print 'Do you want to copy files now? (y/n): ',
		#input = raw_input().strip()
		input = 'n'

		# if os.path.exists(TEMPLATE_FILE_DIR):
		# 	print "remove %s" % TEMPLATE_FILE_DIR
		# 	shutil.rmtree(TEMPLATE_FILE_DIR)

		if input == 'Y' or input == 'y':
			self.copy_files(app_info_dict)
			#self.copy_reactman_files(app_info_dict)
			print '\n******************** Success ********************'
			print 'Modify `models/init.go`, `routers/router.go` to connect resource into system'
		else:
			print '\n******************** Success ********************'
			print 'NOT COPY FILE. Please copy files manually'
			print 'And modify `models/init.go`, `routers/router.go` to connect resource into system'

def generate_code():
	command = Command()
	command.handle()

if __name__ == '__main__':
	generate_code()