#!/bin/bash

fetch_votes() {
    local id=$1
    if [ -z "$id" ]; then
        echo "Usage: fetch_votes <id>"
        return 1
    fi
    
    gnokey query vm/qrender --data "gno.land/r/gov/dao:${id}/votes" -remote https://rpc.test7.testnets.gno.land:443
}

for id in {0..6}; do
    fetch_votes $id  |
    # Keep only lines containing g1 or @
    sed '/g1\|@/!d' |
    # Remove all "- " characters
    sed 's/^- //'  |
    # Keep only username
    sed -E 's/^\[@([A-Za-z0-9_]+)\]\(.*/\1/' |
    # Convert to json array
    jq -R . | jq -s . > data/proposal_$id
done
