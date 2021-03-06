#!/usr/bin/env bats

load ../helpers

function teardown() {
	# if the state of the container is paused, it can't be removed(rm -f)
	run docker_swarm unpause test_container
	swarm_manage_cleanup
	stop_docker
}

@test "docker pause" {
	start_docker_with_busybox 2
	swarm_manage

	docker_swarm run -d --name test_container busybox sleep 1000

	# make sure container is up
	# FIXME(#748): Retry required because of race condition.
	retry 5 0.5 eval "[ -n $(docker_swarm ps -q --filter=name=test_container --filter=status=running) ]"

	docker_swarm pause test_container

	# verify
	# FIXME(#748): Retry required because of race condition.
	retry 5 0.5 eval "[ -n $(docker_swarm ps -q --filter=name=test_container --filter=status=paused) ]"
}
