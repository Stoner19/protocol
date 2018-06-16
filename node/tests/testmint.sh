#!/bin/bash

#
# Test creating a single send transaction in a 1-node chain, reset each time
#
CMD=$GOPATH/src/github.com/Oneledger/protocol/node/scripts
TEST=$GOPATH/src/github.com/Oneledger/protocol/node/tests

# Clear out the existing chains
$CMD/resetOneLedger

# Add in or update users
$TEST/register.sh

# Startup the chains
$CMD/startOneLedger

addrAlice=`$CMD/lookup Alice RPCAddress tcp://127.0.0.1:`
addrBob=`$CMD/lookup Bob RPCAddress tcp://127.0.0.1:`

# olclient wait --initialized
#sleep 2 

# Put some money in the user accounts
SEQ=`$CMD/nextSeq`
olclient testmint -s $SEQ -a $addrAlice --party Alice --amount 100000 --currency OLT 

SEQ=`$CMD/nextSeq`
olclient testmint -s $SEQ -a $addrBob --party Bob --amount 100000 --currency OLT 

olclient account -a $addrAlice --identity Alice
olclient account -a $addrBob --identity Bob

olclient account -a $addrBob --identity Zer


$CMD/stopOneLedger
