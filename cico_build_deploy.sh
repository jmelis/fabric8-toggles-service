#!/bin/bash
#
# Build script for CI builds on CentOS CI https://ci.centos.org/view/Devtools/job/devtools-fabric8-jenkins-proxy-build-master/

set -e

###################################################################################
# Installs all requires build tools to compile, test and build the container image
# Arguments:
#   Nore
# Returns:
#   None
###################################################################################
function setup_build_environment() {
    [ -f jenkins-env ] && cat jenkins-env | grep -e GIT -e DEVSHIFT -e JOB_NAME > inherit-env
    [ -f inherit-env ] && . inherit-env

    # We need to disable selinux for now, XXX
    /usr/sbin/setenforce 0

    yum -y install docker make golang git
    service docker start

    echo 'CICO: Build environment created.'
}

###################################################################################
# Setup the environment for Go, aka the GOPATH
# Arguments:
#   Nore
# Returns:
#   None
###################################################################################
function setup_golang() {
  # Show Go version
  go version
  # Setup GOPATH
  mkdir $HOME/go $HOME/go/src $HOME/go/bin $HOME/go/pkg
  export GOPATH=$HOME/go
  export PATH=$GOPATH/bin:$PATH
}

###################################################################################
# Make sure the Go sources are at their proper location within GOPATH.
# See https://golang.org/doc/code.html
# Arguments:
#   Nore
# Returns:
#   None
###################################################################################
function setup_workspace() {
  mkdir -p $GOPATH/src/github.com/fabric8-services
  cp -r $HOME/payload $GOPATH/src/github.com/fabric8-services/fabric8-toggles-service
}

setup_build_environment
setup_golang
setup_workspace

cd $GOPATH/src/github.com/fabric8-services/fabric8-toggles-service
echo "HEAD of repository `git rev-parse --short HEAD`"
make all

bash <(curl -s https://codecov.io/bash) -f coverage.txt -t cbdff99f-9158-4128-8dec-ef6afb6d78ab

if [[ "$JOB_NAME" = "devtools-fabric8-toggles-service-build-master" ]]; then
    TAG=$(echo ${GIT_COMMIT} | cut -c1-${DEVSHIFT_TAG_LEN})
    make push-openshift REGISTRY_USER=${DEVSHIFT_USERNAME} REGISTRY_PASSWORD=${DEVSHIFT_PASSWORD} IMAGE_TAG=${TAG}
fi