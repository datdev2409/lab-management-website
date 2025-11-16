#!/bin/bash

BACKUP_DIR=${1:-"./backup"}

S3_BUCKET="anhquanlab-mongodb-backup"
ENV="stg"

# Fetch latest backup from S3 for stg environment
latest_backup=$(aws s3api list-objects-v2 --bucket $S3_BUCKET --prefix $ENV/ --query 'Contents[?LastModified>=`2024-10-01`].[Key, LastModified]' --output text | sort -k2 -r | head -n 1)

if [ -z "$latest_backup" ]; then
  echo "No backups found for environment: $ENV"
  exit 1
fi

env_prefix="$ENV/"
backup_key=$(echo $latest_backup | awk '{print $1}' | sed "s|$env_prefix||")
backup_date=$(echo $latest_backup | awk '{print $2}')

echo "Downloading backup from $backup_date..... With key: $backup_key"

backup_file="$BACKUP_DIR/$ENV-$backup_key.tar.gz"
aws s3 cp "s3://$S3_BUCKET/$ENV/$backup_key" "$backup_file"

tar xzf "$backup_file" -C "$BACKUP_DIR"
