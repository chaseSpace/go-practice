#!/usr/bin/env bash

__exitWith() {
    tips=$1
    if [ ! -z "$tips" ]; then
      echo "********** $tips ***********"
    fi
    exit 1
}