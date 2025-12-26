#!/bin/bash

set -o nounset
set -o errexit
set -o pipefail

BIN_DIR="../bin"
SCHEDULER_BINARY_NAME=scheduler
MANAGER_BINARY_NAME=manager
INJECTOR_BINARY_NAME=injector

curDir=$(cd "$(dirname "$0")" && pwd)
cd "${curDir}" || return

. ./env.sh

install() {
    case "${1-}" in
    scheduler)
        install-scheduler
        ;;
    manager)
        install-manager
        ;;
    injector)
        install-injector
    esac
}

install-scheduler() {
    local bin="${INSTALL_HOME}/${INSTALL_BIN_PATH}"
    echo "install: ${bin}"
    mkdir -p "${bin}"

    cp "${BIN_DIR}/${GOOS}_${GOARCH}/${SCHEDULER_BINARY_NAME}" "${bin}"

    createLink "${bin}/${SCHEDULER_BINARY_NAME}" /usr/local/bin/scheduler
}

uninstall-scheduler() {
    echo "unlink /usr/local/bin/scheduler"
    test -e /usr/local/bin/scheduler && unlink /usr/local/bin/scheduler
}

install-manager() {
    local bin="${INSTALL_HOME}/${INSTALL_BIN_PATH}"
    echo "install: ${bin}"
    mkdir -p "${bin}"

    cp "${BIN_DIR}/${GOOS}_${GOARCH}/${MANAGER_BINARY_NAME}" "${bin}"

    createLink "${bin}/${MANAGER_BINARY_NAME}" /usr/local/bin/manager
}

uninstall-manager() {
    echo "unlink /usr/local/bin/manager"
    test -e /usr/local/bin/manager && unlink /usr/local/bin/manager
}

install-injector() {
    local bin="${INSTALL_HOME}/${INSTALL_BIN_PATH}"
    echo "install: ${bin}"
    mkdir -p "${bin}"

    cp "${BIN_DIR}/${GOOS}_${GOARCH}/${INJECTOR_BINARY_NAME}" "${bin}"

    createLink "${bin}/${INJECTOR_BINARY_NAME}" /usr/local/bin/injector
}

uninstall-injector() {
    echo "unlink /usr/local/bin/injector"
    test -e /usr/local/bin/injector && unlink /usr/local/bin/injector
}

createLink() {
    srcPath="$1"
    linkPath="$2"

    echo "create link ${linkPath} to ${srcPath}"
    test -e "${linkPath}" && unlink "${linkPath}"
    ln -s "${srcPath}" "${linkPath}"
}

createDir() {
    test -e "$1" && rm -rf "$1"
    mkdir -p "$1"
}

main() {
    case "${1-}" in
    install)
        install "${2-}"
        ;;
    uninstall)
        uninstall "${2-}"
        ;;
    *)
        echo "You must specify the subcommand 'install' or 'uninstall'."
        ;;
    esac
}

main "$@"
