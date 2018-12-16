#!/bin/bash
  
rm -r Chronos-Node/*.log
rm -r Chronos-Node/olfullnode*
rm -r Chronos-Node/consensus/data
rm -r Chronos-Node/consensus/config/addrbook.json

sed -i '/\"last\_/d' Chronos-Node/consensus/config/priv_validator.json
