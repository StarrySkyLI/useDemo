code:
	goctl api go -api ./dsl/main.api -dir .  --style=go_zero

swg:
	goctl api plugin -plugin goctl-swagger="swagger -filename api.json" -api dsl/main.api -dir swagger/


up:
	go get gitlab.leadsea.cn/go-template/rpc-client
	go get gitlab.leadsea.cn/go-template/base-common
