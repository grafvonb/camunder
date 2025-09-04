#!/usr/bin/env zsh
set -e
yq eval-all '. as $item ireduce ({}; . *+ $item)' operate-openapi.json overlays/process-definitions-overlay.json > operate-openapi-merged.json
oapi-codegen --config oapi-codegen-config.yaml operate-openapi-merged.json