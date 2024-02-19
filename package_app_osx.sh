#!/usr/bin/env bash

cd ./src;go fmt; cd ..;
appname=TheList
CGO_ENABLED=1 fyne package -os darwin --src ./src

if [ $? -eq 0 ]; then
   echo "compilation/packaging succeeded"
else
   echo "compilation/packaging failed"
   exit 1
fi

echo "Copying fonts to executable dir..."
cp -r ./fonts ./$appname.app/Contents/MacOS/

FILE=./$appname.app/Contents/MacOS/conf.json
if [[ -f "$FILE" ]]; then
    echo "$FILE configuration file exists."
else
	echo "Creating default configuration file at $FILE"
	echo '{
 "configuration": {
  "default list": "",
  "default selected": "Inquire",
  "default theme": "Dark",
  "local item file": ""
 }
}' > $FILE
fi
