#!/bin/bash
rm -rf ./_gen_workspace
mkdir ./_gen_workspace
cp ./app.json ./_gen_workspace
cp ./service.json ./_gen_workspace/service.json
cp ./pytool/code_generator/gen_code_tool.py ./_gen_workspace/gen_code_tool.py
cp -rf ./pytool/code_generator/template ./_gen_workspace/_gofile_template
cd ./_gen_workspace
python gen_code_tool.py

cd ..
bash ./polaris.sh
