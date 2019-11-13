#!/bin/sh

set -eu

ROOT_DIR=$(dirname "$0")
cd $ROOT_DIR


color_off="\033[0m"
color_cyan="\033[1;36m"
color_green="\033[1;32m"

echo "${color_cyan}Starting venom tests:${color_off}"
docker-compose build venom

VENOM_CONTAINER_NAME="venom-container-$(date +%s)"
docker-compose run -d --name $VENOM_CONTAINER_NAME venom tail -f /dev/null

docker exec $VENOM_CONTAINER_NAME venom run --var-from-file variables.yml --parallel 5 --format=xml --output-dir="." --strict

docker exec  $VENOM_CONTAINER_NAME cat test_results.xml > test_results.xml
docker-compose rm --force --stop -v venom || true

echo "${color_green}Venom integration tests completed!${color_off}"
