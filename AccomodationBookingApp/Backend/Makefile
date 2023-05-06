proto-authorization:
	protoc -I ./common/proto/authorization_service -I ./common/proto \
	--go_out ./common/proto/authorization_service/generated --go_opt paths=source_relative \
	--go-grpc_out ./common/proto/authorization_service/generated --go-grpc_opt paths=source_relative \
	--grpc-gateway_out ./common/proto/authorization_service/generated --grpc-gateway_opt paths=source_relative \
	--swagger_out=logtostderr=true,allow_merge=true:./common/swagger/ \
	./common/proto/authorization_service/authorization_service.proto

proto-user-profile:
	protoc -I ./common/proto/user_profile_service -I ./common/proto \
	--go_out ./common/proto/user_profile_service/generated --go_opt paths=source_relative \
	--go-grpc_out ./common/proto/user_profile_service/generated --go-grpc_opt paths=source_relative \
	--grpc-gateway_out ./common/proto/user_profile_service/generated --grpc-gateway_opt paths=source_relative \
	--swagger_out=logtostderr=true,allow_merge=true:./common/swagger/ \
	./common/proto/user_profile_service/user_profile_service.proto

swagger-merge:
	swagger-codegen generate -i ./common/swagger/generated/*  -o ./common/swagger/apiDocumentation.json


proto-all: proto-authorization proto-user-profile
