#!/usr/bin/env bash
set -euo pipefail

ROOT=/host-root
ROOT_VOL="--mount type=bind,src=/,dst=${ROOT},options=rbind:ro"
CONFIG=/tmp/kube/config
USER_VERBOSE=""

if [ -z ${IMG+x} ]; then
    IMG="ghcr.io/onmetal/inventory:latest"
fi

if [ -z ${NAMESPACE+x} ]; then
    NAMESPACE="default"
fi

function main() {
    echo "Inventorize and test started"
    clean_before_start;
    pull_image;
    inventorization;
    benchmarks;
}

function clean_before_start() {
    echo "Cleaning up old containers if they are exists"
    for container in inventorization benchmarks; do
        if ctr c ls | grep $container; then
        ctr c rm $container
    fi
    done
}

function pull_image() {
    echo "Pulling image: ${IMG}"
    ctr image pull "${IMG}"
}

function inventorization() {
    echo "Inventarization process started"
    CONTAINER_KUBECONFIG="--mount type=bind,src=${HOST_CONFIG},dst=${CONFIG},options=rbind:ro"
    NAME="inventorization"
    ctr run -d --privileged --net-host ${ROOT_VOL} ${CONTAINER_KUBECONFIG} ${IMG} ${NAME} /app/inventory -r ${ROOT} -k ${CONFIG} -n ${NAMESPACE} ${USER_VERBOSE}
}

function benchmarks() {
    echo "Benchmark process started"
    NAME="benchmarks"
    ctr run --env ROOT=${ROOT} -d --privileged --net-host ${ROOT_VOL} ${IMG} ${NAME} /app/benchmark-scheduler run -n ${NAMESPACE}
}

function print_help() {
    echo "Script used to run docker container with inventory tool."
    echo "Container runs once and will be deleted after completion."
    echo "There are two use cases - when k8s config file contains "
    echo "paths to the certificates and keys and when it contains "
    echo "certificates and keys data."
    echo ""
    echo "Usage: ./inventory.sh [options]"
    echo ""
    echo "  -h|--help to print this message"
    echo ""
    echo "When k8s config file contains paths to certificates and "
    echo "keys, mandatory options are:"
    echo ""
    echo "  -k|--kubeconfig   string absolute path to the k8s config file "
    echo "                      stored on the host"
    echo ""
    echo "Optional parameter is:"
    echo "  -i|--image  string  repository/image:tag if not set, will"
    echo "                      use default repository: 'ghcr.io/onmetal/inventory'"
    echo ""
    echo "  -v|--verbose        enables verbose output"
    echo "                      may be used to troubleshoot the process if data is not collected for some reason"
    echo ""
    echo "  -n|--namespace  string  resource will be pushed to selected namespace"
    echo "                          default value: 'default'"
    echo ""
    exit 0
}

POSITIONAL=()
while [[ $# -gt 0 ]]; do
    key="$1"

    case $key in
        -v|--verbose)
            USER_VERBOSE="-v"
            shift
            ;;
        -k|--kubeconfig)
            HOST_CONFIG="$2"
            shift
            ;;
        -i|--image)
            IMG="$2"
            shift
            ;;
        -n|--namespace)
            NAMESPACE="$2"
            shift
            ;;
        -h|--help)
            print_help
            ;;
        *)
            POSITIONAL+=("$1")
            shift
            ;;
    esac
done

set -- "${POSITIONAL[@]}"

if [ -z ${HOST_CONFIG+x} ]; then
    echo "-k|--kubeconfig parameter is mandatory"
    exit 1
fi

main
