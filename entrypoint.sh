#!/bin/sh

if pidof $CLI_NAME; then
    export PATH=/proc/$(pidof $CLI_NAME)/root/usr/local/bin:$PATH
fi

exec "$@"
