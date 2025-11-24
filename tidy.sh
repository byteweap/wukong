#!/bin/bash

readonly directory=$(cd "$(dirname "$0")" && pwd)
readonly modules=(
    "./"
    "game"
    "gate"
    "contrib/locator/redis"
    "contrib/logger/logrus"
    "contrib/logger/zap"
    "contrib/logger/zerolog"
    "contrib/network/websocket"
    "contrib/network/tcp"
    "contrib/network/kcp"
    "contrib/registry/consul"
    "contrib/registry/etcd"
    "contrib/registry/nacos"
    "contrib/registry/zookeeper"
)

for module in ${modules[@]}
do
  cd "${module}"
  go mod tidy
  cd "${directory}"
done