#!/bin/zsh

# Default values
# https://www.postman.com/:publicHandle/collection/20317927-361a5fc9-69a2-4dd9-b5da-3e98f6095ab9
DEFAULT_COLLECTION_ID="20317927-361a5fc9-69a2-4dd9-b5da-3e98f6095ab9"
DEFAULT_OUTPUT_FILE="operate-openapi.json"

# Use environment variables if set, otherwise fall back to defaults
COLLECTION_ID="${POSTMAN_COLLECTION_ID:-$DEFAULT_COLLECTION_ID}"
OUTPUT_FILE="${POSTMAN_OUTPUT_FILE:-$DEFAULT_OUTPUT_FILE}"

# Ensure the API key is set
if [[ -z "$POSTMAN_API_KEY" ]]; then
  echo "Error: POSTMAN_API_KEY environment variable is not set."
  exit 1
fi

# Perform the API call
curl -s "https://api.getpostman.com/collections/${COLLECTION_ID}/transformations" \
  --header 'Content-Type: application/json' \
  --header "X-Api-Key: $POSTMAN_API_KEY" | jq '.output | fromjson' > "$OUTPUT_FILE"
