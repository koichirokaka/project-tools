#!/bin/bash

echo "Submitting a Cloud ML Engine job"

REGION="us-central1"
BUCKET="you-bucket-name" # change to your bucket name

MODEL_NAME="your_model_name" # change to your model name

GCS_TRAIN_FILE=gs://${BUCKET}/path/to/data/data.csv
GCS_EVAL_FILE=gs://${BUCKET}/path/to/data/test.csv
MODEL_DIR=gs://${BUCKET}/path/to/models/${MODEL_NAME}

CURRENT_DATE=`date +%Y%m%d_%H%M%S`
JOB_NAME=train_${MODEL_NAME}_${TIER}_${CURRENT_DATE}
#JOB_NAME=tune_${MODEL_NAME}_${CURRENT_DATE} # for hyper-parameter tuning jobs

TRAIN_STEPS=200

gcloud ml-engine jobs submit training $JOB_NAME \
--stream-logs \
--package-path trainer \
--module-name trainer.task \
-- \
--train-files $GCS_TRAIN_FILE \
--eval-files $GCS_EVAL_FILE \
--train-steps $TRAIN_STEPS
