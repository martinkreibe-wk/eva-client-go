####
# Docker Commands
gen-docker:
	docker build \
		--build-arg GIT_SSH_KEY \
		-f workivabuild.Dockerfile \
		-t drydock.workiva.net/workiva/eva-client-go:latest-release .

update-tocs:
	./.circleci/scripts/update-tocs.sh
