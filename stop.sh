#!/bin/bash

ps -ef | grep request-matcher-openai  | grep env |grep uat | awk '{print $2}' | xargs kill
