#!/bin/bash

if [ -z "$1" ]
then
	echo "Usage: register.sh {uniqueName}"
	exit -1
fi

olclient update -c Chronos --account "$1-Account"

olclient register -c Chronos --identity "$1" --account "$1-Account" --node Chronos-Node

