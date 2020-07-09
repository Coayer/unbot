#!/bin/bash
model_path="../internal/wikiQA/model"

if ! [ -d $model_path ]; then
  wget -O $model_path.tar.gz "https://cdn.huggingface.co/distilbert-base-uncased-distilled-squad-saved_model.tar.gz"
  mkdir $model_path
  tar -xzf $model_path.tar.gz -C $model_path
  rm $model_path.tar.gz
fi
