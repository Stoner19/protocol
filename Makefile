#
# Docker targets
#
#

# Prepare the base image for docker
docker:
	docker build -t oneledger/chronos -f ./DOCKER/Dockerfile .

# Ensure Docker doesn't use its cache, useful for making small edits to scripts
docker-nocache:
	docker build -t oneledger/chronos -f ./DOCKER/Dockerfile --no-cache .

run:
	docker run -it --volume /tmp/VOLUME:/home/oneledger/.olfullnode --env NODE_NAME=Zebra oneledger/chronos

.PHONY: docker-fullnode run-singlenode local-testnet local-testnet-down prepare-volume reset-volume build run
