#!/bin/bash

seeds="6f724aaeb0a9232b91060fc98af9f034fa7970d0@35.203.120.202:26615,f2815fcb2deb6209c26d86e1cbb6f7708f21e458@35.203.59.66:26615"
peers="f2815fcb2deb6209c26d86e1cbb6f7708f21e458@35.203.59.66:26615,6f724aaeb0a9232b91060fc98af9f034fa7970d0@35.203.120.202:26615"

olfullnode node -c Chronos \
	--root ./Chronos-Node/olfullnode \
	--tendermintRoot ./Chronos-Node/consensus \
        --seeds $seeds \
	--persistent_peers $peers \
       >> ./Chronos-Node/olfullnode.log 2>&1 &

sleep 2

