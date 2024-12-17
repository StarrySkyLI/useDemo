
code:
	#sh ./dsl/main.sh
	goctl rpc protoc dsl/rpc_demo.proto --go_out=./  --go-grpc_out=./  --zrpc_out=./ --client=true --multiple=true --style=go_zero

swg:
	# goctl api plugin -plugin goctl-swagger="swagger -filename rpc.json" -api rpc.proto -dir ../swagger

up:
	#go get gitlab.leadsea.cn/go-template/base-common
	#go mod tidy