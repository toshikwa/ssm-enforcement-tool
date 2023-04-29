SHELL := /bin/bash
LAMBDA_DIR := assets/functions
FUNCTIONS := list_unmanaged_instances attach_ssm_permissions send_alert_to_sns

define build
$(1):
	GOARCH=amd64 GOOS=linux go build -C ${LAMBDA_DIR}/cmd/${1} -o ../../bin/${1}
	zip -j ${LAMBDA_DIR}/archive/${1}.zip ${LAMBDA_DIR}/bin/${1}
endef

all: $(FUNCTIONS)
$(foreach func,$(FUNCTIONS),$(eval $(call build,$(func))))
