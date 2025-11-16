# BACKUP_FILE=

DEFAULT_BACKUP_FOLDER=$(ls -1 -d ./tmp/*/ | head -n 1)

BACKUP_FOLDER=${1:-"$DEFAULT_BACKUP_FOLDER"}
MONGODB_URI=${2:-"mongodb://root:password123@localhost:27017"}
DB_NAME="labadmin"

echo "Restoring database from backup folder: $BACKUP_FOLDER to MongoDB URI: $MONGODB_URI"

read -p "Please confirm the restore operation (y/n): " confirmation

if [ "$confirmation" != "y" ]; then
  echo "Restore operation cancelled."
  exit 0
fi

mongorestore --uri=$MONGODB_URI --nsInclude $DB_NAME.* --drop $BACKUP_FOLDER