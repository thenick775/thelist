#!/usr/bin/env bash

cd ./src;go fmt; cd ..;
appname=TheList
GO_PATH=$(go env | grep "GOPATH" | grep -o "\".*\"" | tr -d '"')
CGO_ENABLED=1 $GO_PATH/bin/fyne package -os darwin --src ./src -icon ./images/TheList-large.png --name $appname

if [ $? -eq 0 ]; then
   echo "compilation/packaging succeeded"
else
   echo "compilation/packaging failed"
   exit 1
fi

FILE=./$appname.app/Contents/MacOS/conf.json
if [[ -f "$FILE" ]]; then
    echo "$FILE configuration file exists."
else
	echo "Creating default configuration file at $FILE"
	echo '{
 "configuration": {
  "default branches open": [
   "Quick Actions"
  ],
  "default list": "",
  "default selected": "Inquire",
  "default theme": "Light",
  "local item file": ""
 }
}' > $FILE
fi
