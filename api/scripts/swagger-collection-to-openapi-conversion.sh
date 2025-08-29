#!/bin/zsh
# https://www.postman.com/camundateam

# Camunda 8 API (REST) Postman Collection to OpenAPI conversion script
curl -s https://api.getpostman.com/collections/24684262-e4cb02e0-4f3e-42c2-b983-3644feef3565/transformations \
--header 'Content-Type: application/json' \
--header 'X-Api-Key: <Postman API Key>' | jq '.output | fromjson'

