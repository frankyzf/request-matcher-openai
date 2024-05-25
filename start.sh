#!/bin/bash

nohup ./request-matcher-openai --env="uat" >> nohup.out 2>&1 &
