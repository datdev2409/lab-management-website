#!/bin/bash

set -e

# Configuration
AWS_REGION="${AWS_REGION:-ap-southeast-1}"
APP_NAME="lab-admin-go"

ENVIRONMENT="$1"
OUTPUT_FILE=".env"

PREFIX="/${APP_NAME}/${ENVIRONMENT}/"

# Check AWS CLI
if ! command -v aws &> /dev/null; then
    echo "Error: AWS CLI not found"
    exit 1
fi

# Fetch parameters
echo "Fetching parameters from ${PREFIX}..."
aws ssm get-parameters-by-path \
    --path "${PREFIX}" \
    --recursive \
    --with-decryption \
    --region "${AWS_REGION}" \
    --output text \
    --query 'Parameters[*].[Name,Value]' | \
while read -r name value; do
    # Remove prefix to get environment variable name
    env_var=$(echo "$name" | sed "s|${PREFIX}||")
    echo "${env_var}=${value}"
done > "$OUTPUT_FILE"

echo "Parameters saved to ${OUTPUT_FILE}"
