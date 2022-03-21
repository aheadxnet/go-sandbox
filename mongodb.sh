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

MONGO_DB_NAME="go-sandbox-mongodb"
MONGO_DB_PID_FILE="${BASE_DIR}/.mongodb.pid"

function showUsage() {
  echo "Script to help starting and stopping Mongo-DB Docker container"
  echo "Usage: $0 [start|stop|clean]"
}

if [ $# -lt 1 ]; then
  showUsage
  exit 1
fi

function startContainer() {
  if [[ -f "${MONGO_DB_PID_FILE}" ]]
  then
    echo "There is a (stale?) PID file, or the container is already running."
    exit 2
  fi
  # Container alread exists?
  CONTAINER="$(docker container ls --all | grep -o ${MONGO_DB_NAME})"
  if [[ "${CONTAINER}x" == "x" ]]
  then
    # Start mit fixem Passwort:
    MONGO_DB_PASSWORD="$(cat ${BASE_DIR}/run/secrets/mongodb_password)"

    MONGO_DB_PID="$(docker run -d --name ${MONGO_DB_NAME} -v "${BASE_DIR}/data" -e MONGO_INITDB_ROOT_USERNAME=admin -e MONGO_INITDB_ROOT_PASSWORD=${MONGO_DB_PASSWORD} -p 27017:27017 mongo:5.0.6)"

    # Start mit Passwort aus docker secret file
    # Anlegen des Passwortes mit
    # openssl rand -base64 12 | docker secret create mongodb_password -
    # Redis-Service anlegen
    # docker service  create --name redis --secret mongodb_password redis:alpine
    # docker run -d --name mongodb -e MONGO_INITDB_ROOT_USERNAME=admin -e MONGO_INITDB_ROOT_PASSWORD_FILE=/run/secrets/mongodb_password -p 27017:27017 mongo:latest
  else
    docker start "${MONGO_DB_NAME}"
    MONGO_DB_PID="$(docker ps -l | grep ${MONGO_DB_NAME} | sed -e 's/\([0-9a-f]*\)\(.*\)/\1/g')"
  fi
  echo ${MONGO_DB_PID} >> "${BASE_DIR}/.mongodb.pid"

  echo "MongoDB container ${MONGO_DB_NAME} startet with PID ${MONGO_DB_PID}"
}

function stopContainer() {
  if [[ ! -f "${MONGO_DB_PID_FILE}" ]]
  then
    echo "There is no PID file, the container seems to be not running."
    exit 2
  fi
  docker stop "${MONGO_DB_NAME}"
  echo "Stopped ${MONGO_DB_NAME} ..."
  rm "${MONGO_DB_PID_FILE}"
  echo "Removed PID file."
}

function cleanContainer() {
  if [[ -f "${MONGO_DB_PID_FILE}" ]]
  then
    echo "There is a (stale?) PID file, or the container is still running."
    exit 2
  fi
  docker rm "${MONGO_DB_NAME}"
  echo "Removed container ${MONGO_DB_NAME}."
}

case "$1" in
  "start")
    startContainer
    ;;
  "stop")
    stopContainer
    ;;
  "clean")
    cleanContainer
    ;;
  *)
    showUsage
    ;;
esac