#!/bin/bash

# Simple test script for the Pomodoro API
BASE_URL="http://localhost:8080/api/v1"

echo "Testing Pomodoro API..."

# Test health check
echo "1. Testing health check..."
curl -s "$BASE_URL/healthcheck" | jq .

echo -e "\n2. Testing user creation..."
USER_RESPONSE=$(curl -s -X POST "$BASE_URL/users" \
  -H "Content-Type: application/json" \
  -d '{
    "firebase_uid": "test_uid_123",
    "email": "test@example.com",
    "name": "Test User"
  }')

echo $USER_RESPONSE | jq .

# Extract user ID (assuming the API returns the full object)
USER_ID=$(echo $USER_RESPONSE | jq -r '.data.inserted_id' 2>/dev/null)

if [ "$USER_ID" != "null" ] && [ "$USER_ID" != "" ]; then
  echo -e "\n3. Testing task creation..."
  TASK_RESPONSE=$(curl -s -X POST "$BASE_URL/tasks" \
    -H "Content-Type: application/json" \
    -d "{
      \"user_id\": \"$USER_ID\",
      \"title\": \"Test Task\",
      \"description\": \"This is a test task\",
      \"estimated_pomodoros\": 3
    }")
  
  echo $TASK_RESPONSE | jq .
  
  echo -e "\n4. Testing get tasks by user..."
  curl -s "$BASE_URL/tasks/user/$USER_ID" | jq .
else
  echo "User creation failed, skipping task tests..."
fi

echo -e "\n\nAPI testing completed!"
