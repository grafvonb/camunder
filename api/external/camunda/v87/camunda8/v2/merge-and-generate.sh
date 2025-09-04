#!/bin/zsh
set -euo pipefail
setopt null_glob    # if no overlays match, the glob expands to nothing (not a literal)

SCHEMA="camunda8-openapi"
SCHEMA_FILE="$SCHEMA.json"
OVERLAY_DIR="overlays"
OUT="$SCHEMA-merged.json"

# Collect and sort overlays (so merges are deterministic)
OVERLAYS=(${OVERLAY_DIR}/*.json)
OVERLAYS=(${(on)OVERLAYS})  # sort by name

# Merge base schema + all overlays using yq v4
# The *+ operator does a deep merge with overlays overriding base where needed.
yq eval-all '. as $item ireduce ({}; . *+ $item)' "$SCHEMA_FILE" ${OVERLAYS:+${OVERLAYS[@]}} > "$OUT"

# Generate code from the merged spec
oapi-codegen --config oapi-codegen-config.yaml "$OUT"