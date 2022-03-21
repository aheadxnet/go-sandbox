#! /bin/bash
# set -e

SOURCE="${BASH_SOURCE[0]}"
while [ -h "$SOURCE" ]; do # resolve $SOURCE until the file is no longer a symlink
  DIR="$(cd -P "$(dirname "$SOURCE")" && pwd)"
  SOURCE="$(readlink "$SOURCE")"
  [[ $SOURCE != /* ]] && SOURCE="$DIR/$SOURCE" # if $SOURCE was a relative symlink, we need to resolve it relative to the path where the symlink file was located
done
DIR="$(cd -P "$(dirname "$SOURCE")" && pwd)"

BASE_DIR="$(cd -P "$DIR"  && pwd)"

MONGO_DB_PASSWORD="$(cat ${BASE_DIR}/run/secrets/mongodb_password)"

export MONGO_URI="mongodb://admin:${MONGO_DB_PASSWORD}@localhost:27017/test?authSource=admin"
export MONGO_DATABASE="demo"
go run main.go