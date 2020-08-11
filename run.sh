#!/bin/bash
touch data/memories.json
touch data/reminders.json
source bin/activate
export LIBRARY_PATH=$LIBRARY_PATH:~/unbot/lib
export LD_LIBRARY_PATH=$LD_LIBRARY_PATH:~/unbot/lib
nohup ./unbot &
