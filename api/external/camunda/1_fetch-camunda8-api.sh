TAG=8.8.110
curl -L "https://github.com/camunda/camunda-docs/archive/refs/tags/${TAG}.tar.gz" \
  | tar -xz --strip-components=1 "camunda-docs-${TAG}/api"