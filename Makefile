# proto 코드 생성 경로
PROTO_GEN_PATH := ./gen
PROTO_GEN_GO := $(PROTO_GEN_PATH)/go

# proto 옵션
PROTO_OPT_PROTO_PATH := "--proto_path=./proto"
PROTO_OPT_GO_OUT := "--go_out=$(PROTO_GEN_GO)"
PROTO_OPT_GO_PATH := "--go_opt=paths=source_relative"
PROTO_OPT_GO_GRPC_OUT := "--go-grpc_out=$(PROTO_GEN_GO)"
PROTO_OPT_GO_GRPC_PATH := "--go-grpc_opt=paths=source_relative"

# proto 타겟
PROTO_TARGETS := \
	aigateway/v1 \
	channelconfig/v1 \
	channelconfig/v2 \
	dataportal/v1 \
	jarvis/v1


.PHONY: proto
proto:
	@mkdir -p $(PROTO_GEN_GO)
	@for target in $(PROTO_TARGETS); do \
		echo generating $$target; \
		for file in $$(ls ./proto/$$target); do \
			protoc $(PROTO_OPT_PROTO_PATH) $(PROTO_OPT_GO_OUT) $(PROTO_OPT_GO_PATH) $(PROTO_OPT_GO_GRPC_OUT) $(PROTO_OPT_GO_GRPC_PATH) "./proto/$$target/$$file"; \
		done; \
	done

.PHONY: clean
clean:
	@rm -rf $(GEN_PATH)
