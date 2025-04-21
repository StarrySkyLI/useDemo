# 设置安装目录
CONSUL_DIR = D:/consul_1.17.0_windows_386
PROMETHEUS_DIR = D:/prometheus-2.45.0.windows-amd64
GRAFANA_DIR = D:/grafana-10.0.2
JAEGER_DIR = D:/jaeger-1.47.0-windows-amd64

# 设置可执行文件路径
PROMETHEUS = $(PROMETHEUS_DIR)/prometheus
CONSUL = $(CONSUL_DIR)/consul
GRAFANA = $(GRAFANA_DIR)/bin/grafana server
JAEGER = $(JAEGER_DIR)/jaeger-all-in-one.exe

# 设置启动参数
DEV_MODE_ARGS = agent -dev
SERVER_MODE_ARGS = agent -server -bootstrap-expect=1 -data-dir="D:/consul_data" -node="node1" -bind="127.0.0.1" -client="0.0.0.0"
PROMETHEUS_ARGS = --config.file=$(PROMETHEUS_DIR)/prometheus.yml --storage.tsdb.path=D:/prometheus_data
GRAFANA_ARGS = --homepath=D:/grafana-10.0.2
JAEGER_ARGS =

# 启动开发模式
consul_dev:
	@echo "Starting Consul in dev mode..."
	$(CONSUL) $(DEV_MODE_ARGS)

# 启动 Consul 服务器模式
consul_server:
	@echo "Starting Consul in server mode..."
	$(CONSUL) $(SERVER_MODE_ARGS)

# 启动 Prometheus
prom:
	@echo "Starting Prometheus..."
	$(PROMETHEUS) $(PROMETHEUS_ARGS)

# 启动 Grafana
gra:
	@echo "Starting Grafana..."
	$(GRAFANA) $(GRAFANA_ARGS)

# 启动 Jaeger
ja:
	@echo "Starting Jaeger..."
	$(JAEGER) $(JAEGER_ARGS)


# 清理命令（如果需要）
clean:
	@echo "Cleaning up..."
	# 可以根据需要添加清理操作
abp:
	ab -p POST.json -T application/json -c 3 -n 20 "http://127.0.0.1:8888/v1/demo/find_id"

abg:
	ab -c 3 -n 20 "http://127.0.0.1:8888/v1/demo/breaker_test?userId=1"

	ab -c 3 -n 2000 "http://43.135.42.236:8080/test/test"