#!/bin/bash
# Seed script to create initial test user and data for E2E tests

set -e

MONGODB_URI=${MONGODB_URI:-"mongodb://root:password123@localhost:27017"}
DB_NAME=${DB_NAME:-"labadmin"}

echo "Waiting for MongoDB to be ready..."
timeout 30 bash -c 'until mongosh "$MONGODB_URI" --eval "db.adminCommand(\"ping\")" > /dev/null 2>&1; do sleep 1; done'

echo "MongoDB is ready. Seeding test data..."

# Create test user using mongosh
mongosh "$MONGODB_URI/$DB_NAME" <<EOF
// Create test admin user
db.users.insertOne({
  _id: "user_admin_test",
  username: "admin",
  password: "\$2a\$10\$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy", // password: admin123
  created_at: new Date(),
  updated_at: new Date()
});

print("✓ Test user created: username=admin, password=admin123");
print("✓ Database seeded successfully!");
EOF

echo "Seeding completed!"
