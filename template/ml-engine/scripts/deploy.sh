#!/bin/bash

echo "Deploy trained ML model to Cloud ML Engine"

REGION="us-central1"
BUCKET="your-bucket-name" # change to your bucket name

MODEL_NAME="you_model_name" # change to your estimator name
MODEL_VERSION="your.model.version" # change to your model version

JOB_DIR="keras"

MODEL_BINARIES=$JOB_DIR/export

# delete model version
gcloud ml-engine versions delete ${MODEL_VERSION} --model=${MODEL_NAME}

# delete model
gcloud ml-engine models delete ${MODEL_NAME}

# deploy model to GCP
gcloud ml-engine models create $MODEL_NAME --regions $REGION

# deploy model version
gcloud ml-engine versions create $MODEL_VERSION --model $MODEL_NAME --origin $MODEL_BINARIES
