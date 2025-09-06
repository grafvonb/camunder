#!/bin/zsh

if [ -z "$1" ]; then
  echo "Usage: source $0 <token>"
  # If sourced, return; if executed, exit:
  return 1 2>/dev/null || exit 1
fi

token="$1"

export CAMUNDER_CAMUNDA8_API_TOKEN="$token"
export CAMUNDER_OPERATE_API_TOKEN="$token"
export CAMUNDER_TASKLIST_API_TOKEN="$token"