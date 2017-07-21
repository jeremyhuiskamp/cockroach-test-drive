#!/bin/bash
#
# Wrapper around go's dep tool to manipulate $GOPATH.
# Needed so that you can avoid having to put this entire
# repository in your central $GOPATH.

GOPATH=$(dirname $(dirname $(pwd))) dep "$@"
