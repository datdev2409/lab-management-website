set -e
MONGODB_URI=${MONGODB_URI:-"mongodb://root:password123@localhost:27017"}
DB_NAME=${DB_NAME:-"labadmin"}

CONTAINER_ID=$(docker ps -qf "name=mongo_db" 2>/dev/null || echo "")
if [ ! -z "$CONTAINER_ID" ]; then
  echo "Using Docker exec to seed database..."
  docker exec "$CONTAINER_ID" mongosh -u root -p password123 --authenticationDatabase admin "$DB_NAME" --eval '
    db.users.deleteMany({_id: "user_admin_test"});
    db.users.insertOne({
      _id: "user_admin_test",
      username: "admin",
      password: "$2a$10$a0bURtaDVOIbG/vOE8WeiunLmPw.WfbqSBURLM3AfvT61uKtyIxKu",
      created_at: new Date(),
      updated_at: new Date()
    });
    print("✓ Test user created: username=admin, password=admin123");
  ' && echo "✓ Database seeded successfully!" && exit 0
fi
