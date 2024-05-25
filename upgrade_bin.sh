#!/bin/bash
date
if [ "$#" -ne 1 ]; then
    echo "usage: ./upgrade_bin.sh CLIENT; eg: ./upgrade_bin.sh nuc-1"
    exit 0;
fi
nuc=$1
echo "machine is:$nuc"

echo "uploading request-matcher-openai"

ssh ${nuc} "[ ! -d 'upgrade.request-matcher-openai' ] &&  mkdir upgrade.request-matcher-openai"
rsync -P bin/linux/request-matcher-openai ${nuc}:upgrade.request-matcher-openai/request-matcher-openai
rsync -P config/config.json ${nuc}:upgrade.request-matcher-openai/config.json
rsync -P upgrade.sh ${nuc}:upgrade.request-matcher-openai/upgrade.sh
rsync -P upgrade.sh ${nuc}:upgrade.sh.request-matcher-openai

