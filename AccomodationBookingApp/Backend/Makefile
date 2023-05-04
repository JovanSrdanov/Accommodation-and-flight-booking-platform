proto-authorization:
	protoc -I ./common/proto/authorization_service -I ./common/proto \
	--go_out ./common/proto/authorization_service/generated --go_opt paths=source_relative \
	--go-grpc_out ./common/proto/authorization_service/generated --go-grpc_opt paths=source_relative \
	--grpc-gateway_out ./common/proto/authorization_service/generated --grpc-gateway_opt paths=source_relative \
	./common/proto/authorization_service/authorization_service.proto

proto-user-profile:
	protoc -I ./common/proto/user_profile_service -I ./common/proto \
	--go_out ./common/proto/user_profile_service/generated --go_opt paths=source_relative \
	--go-grpc_out ./common/proto/user_profile_service/generated --go-grpc_opt paths=source_relative \
	--grpc-gateway_out ./common/proto/user_profile_service/generated --grpc-gateway_opt paths=source_relative \
	./common/proto/user_profile_service/user_profile_service.proto


proto-all: proto-authorization proto-user-profile
