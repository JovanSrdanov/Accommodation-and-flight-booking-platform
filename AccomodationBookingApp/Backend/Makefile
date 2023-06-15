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

proto-accommodation:
	protoc -I ./common/proto/accommodation_service -I ./common/proto \
	--go_out ./common/proto/accommodation_service/generated --go_opt paths=source_relative \
	--go-grpc_out ./common/proto/accommodation_service/generated --go-grpc_opt paths=source_relative \
	--grpc-gateway_out ./common/proto/accommodation_service/generated --grpc-gateway_opt paths=source_relative \
	./common/proto/accommodation_service/accommodation_service.proto

proto-reservation:
	protoc -I ./common/proto/reservation_service -I ./common/proto \
	--go_out ./common/proto/reservation_service/generated --go_opt paths=source_relative \
	--go-grpc_out ./common/proto/reservation_service/generated --go-grpc_opt paths=source_relative \
	--grpc-gateway_out ./common/proto/reservation_service/generated --grpc-gateway_opt paths=source_relative \
	./common/proto/reservation_service/reservation_service.proto

proto-notification:
	protoc -I ./common/proto/notification_service -I ./common/proto \
	--go_out ./common/proto/notification_service/generated --go_opt paths=source_relative \
	--go-grpc_out ./common/proto/notification_service/generated --go-grpc_opt paths=source_relative \
	--grpc-gateway_out ./common/proto/notification_service/generated --grpc-gateway_opt paths=source_relative \
	./common/proto/notification_service/notification_service.proto

proto-rating:
	protoc -I ./common/proto/rating_service -I ./common/proto \
    	--go_out ./common/proto/rating_service/generated --go_opt paths=source_relative \
    	--go-grpc_out ./common/proto/rating_service/generated --go-grpc_opt paths=source_relative \
    	--grpc-gateway_out ./common/proto/rating_service/generated --grpc-gateway_opt paths=source_relative \
    	./common/proto/rating_service/rating_service.proto

proto-all: proto-authorization proto-user-profile proto-accommodation proto-reservation notification_service proto-rating
