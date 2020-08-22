#!/bin/bash
touch data/user/memories.json data/user/reminders.json
source bin/activate
export LIBRARY_PATH=$LIBRARY_PATH:$PWD/lib LD_LIBRARY_PATH=$LD_LIBRARY_PATH:$PWD/lib
nohup ./unbot &
