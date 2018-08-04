BASE_DIR:=$(shell dirname $(realpath $(lastword $(MAKEFILE_LIST))))
SUBDIRS := $(wildcard $(BASE_DIR)/*/.)
SUBDIRS := $(BASE_DIR)/services/. $(BASE_DIR)/gateway/. $(BASE_DIR)/common/.

TARGETS := all test
SUBDIRS_TARGETS := $(foreach target,$(TARGETS),$(addsuffix $(target),$(SUBDIRS)))
DEPLOY_GATEWAY_TARGETS := gateway
DEPLOY_SERVICE_TARGETS := auth user registration decision rsvp checkin upload mail event stat
DEPLOY_TARGETS := $(DEPLOY_GATEWAY_TARGETS) $(DEPLOY_SERVICE_TARGETS)

.PHONY: $(TARGETS) $(SUBDIRS_TARGETS)

$(TARGETS): % : $(addsuffix %,$(SUBDIRS))
	@echo 'Finished running target "$*"'

$(SUBDIRS_TARGETS):
	$(MAKE) -C $(@D) $(@F:.%=%)

.PHONY: run
run:
	@$(BASE_DIR)/scripts/run.sh

.PHONY: fmt
fmt:
	go fmt ./...

.PHONY: deploy
deploy:
	mkdir -p $(BASE_DIR)/deploy/
	$(foreach target,$(DEPLOY_TARGETS),mkdir -p $(BASE_DIR)/deploy/$(target)/;cp $(GOPATH)/bin/hackillinois-api-$(target) $(BASE_DIR)/deploy/$(target)/hackillinois-api-$(target);)
	$(foreach target,$(DEPLOY_GATEWAY_TARGETS),cp $(BASE_DIR)/$(target)/Dockerfile $(BASE_DIR)/deploy/$(target)/Dockerfile;)
	$(foreach target,$(DEPLOY_SERVICE_TARGETS),cp $(BASE_DIR)/services/$(target)/Dockerfile $(BASE_DIR)/deploy/$(target)/Dockerfile;)
	cp $(BASE_DIR)/buildspec.yml $(BASE_DIR)/deploy/buildspec.yml