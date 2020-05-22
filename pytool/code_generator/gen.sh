rm -rf ./__gin_master.zip
cp ./app.json ./_gen_workspace
cp ./service.json ./_gen_workspace/service.json
cp ./_gen_workspace/gin-master/pytool/code_generator/gen_code_tool.py ./_gen_workspace/gen_code_tool.py
cp -rf ./_gen_workspace/gin-master/pytool/code_generator/template ./_gen_workspace/_gofile_template
cd ./_gen_workspace
python gen_code_tool.py
