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

export GOHOME="${BASE_DIR}"

#  with go >= 1.16 disable Go modules
export GO111MODULE=off

export GOPATH="${HOME}/go"
export PATH="${PATH}":"${GOPATH}"/bin
