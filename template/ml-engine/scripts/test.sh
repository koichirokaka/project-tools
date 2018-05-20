#!/bin/bash

echo "Test training local ML model"

TRAIN_FILE=data.csv
EVAL_FILE=test.csv

JOB_DIR="keras"
TRAIN_STEPS=200

gcloud ml-engine local train --package-path trainer \
--module-name trainer.task \
-- \
--train-files $TRAIN_FILE \
--eval-files $EVAL_FILE \
--job-dir $JOB_DIR \
--train-steps $TRAIN_STEPS
