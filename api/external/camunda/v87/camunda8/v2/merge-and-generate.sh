#!/usr/bin/env zsh
set -e
yq eval-all '. as $item ireduce ({}; . *+ $item)' camunda-8-openapi.json overlays/cluster-topology-overlay.json > camunda-8-openapi-merged.json
oapi-codegen --config oapi-codegen-config.yaml camunda-8-openapi-merged.json