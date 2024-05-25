#!/bin/bash

dd=$(date +%Y-%m-%d.%H%M%S)
[ ! -d log_archive ] && mkdir log_archive
[ ! -d backup ] && mkdir backup

cp nohup.out  log_archive/nohup-$dd.out
cp -r config backup/config.$dd

mv request-matcher-openai backup/request-matcher-openai.$dd

./stop.sh
cp ~/upgrade.request-matcher-openai/upgrade.sh ./upgrade.sh
cp ~/upgrade.request-matcher-openai/request-matcher-openai .
cp ~/upgrade.request-matcher-openai/config.json config/config.json
#./start.sh
