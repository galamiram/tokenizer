app_name 				 	= tokenizer
monitoring_namespce      	= monitoring
custom_metrics_namespace 	= custom-metrics
prometheus_operator_version = 0.40.0
tag		 					= $(shell cat VERSION)


.PHONY: add-helm-repos
add-helm-repos:
		@helm repo add stable https://kubernetes-charts.storage.googleapis.com	
		@helm repo add bitnami https://charts.bitnami.com/bitnami

.PHONY: start-minikube
start-minikube:
		@echo "Starting Minikube..."
		@minikube start --driver=hyperkit > /dev/null

.PHONY: install-prometheus-operator
install-prometheus-operator: start-minikube add-helm-repos
		helm upgrade --install monitoring \
			bitnami/prometheus-operator \
			-f helm-values/prometheus-operator-values.yaml \
			-n $(monitoring_namespce) \
			--create-namespace

.PHONY: install-prometheus-adapter
install-prometheus-adapter: start-minikube add-helm-repos
		helm upgrade --install monitoring \
			stable/prometheus-adapter \
			-f helm-values/prometheus-adapter-values.yaml \
			--set rbac.create="true" \
			--version 2.4.0 \
			--namespace $(custom_metrics_namespace) \
			--create-namespace 			

.PHONY: vendor
vendor:
		@go mod vendor


.PHONY: build
build: start-minikube vendor
		@eval $$(minikube docker-env); \
		docker image build -t $(app_name):$(tag) -f Dockerfile .

.PHONY: deploy
deploy: build install-prometheus-operator install-prometheus-adapter
		helm upgrade --install tokenizer \
			./charts/tokenizer \
			--set image.tag=$(tag) \
			--namespace $(app_name) \
			--create-namespace

.PHONY:
load-test: deploy
		eval $(minikube docker-env --unset); \
			docker build -t test-tokenizer . -f tests/Dockerfile && \
			docker run -it --network host \
			test-tokenizer ./tests/test \
			--dir ./tests/test-cases/ \
			--url "$$(minikube service tokenizer --url --namespace $(app_name))/tokenize" \
			--parallelism 10
				
.PHONE: clean-all
clean-all:
		minikube delete
		@rm -rf dist

