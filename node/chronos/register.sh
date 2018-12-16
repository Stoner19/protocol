#!/bin/bash

if [ -z "$1" ]
then
	echo "Usage: register.sh {uniqueName}"
	exit -1
fi

olclient update -c Chronos --account "LocalTest"

olclient register -c Chronos --identity "$1" --account "LocalTest" --node Chronos-Node

