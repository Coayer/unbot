#!/bin/bash
mkdir unbot && cd unbot && mkdir config data data/internal data/user
python3 -m venv .
source bin/activate
pip3 install tensorflow
wget -qO- https://storage.googleapis.com/tensorflow/libtensorflow/libtensorflow-cpu-linux-x86_64-1.15.0.tar.gz | tar -xvz
mv lib/libtensorflow.so.1 lib/libtensorflow.so.2
