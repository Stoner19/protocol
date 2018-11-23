#!/bin/bash

#
# Test creating a single send transaction in a 1-node chain, reset each time
#
CMD=$GOPATH/src/github.com/Oneledger/protocol/node/scripts


echo "=================== Send Transactions =================="
olclient install -c David --owner David-OneLedger --name Test --version v0.0.1 -f test.old
olclient execute -c David --owner David-OneLedger --name TestExecute --version v0.10.1

