# Tokenizer

## Prerequisite

* macOS 10.11+

* Docker

    ```
    brew install docker
    ```

* Hyperkit
    
    HyperKit is an open-source hypervisor for macOS hypervisor, optimized for lightweight virtual machines and container deployment.

    ```
    brew install hyperkit
    ```

* Minikube

    Minikube is a tool that makes it easy to run Kubernetes locally. Minikube runs a single-node Kubernetes cluster inside a Virtual Machine (VM) on your laptop for users looking to try out Kubernetes or develop with it day-to-day.

    ```
    brew install minikube
    ```
* Helm

    Helm is the package manager for Kubernetes.

    ```
    brew install helm
    ```

* GNU Make  

    If `make` is not installed on your machine please run the following command:
    ```
    xcode-select --install
    ```
## Build the Tokenizer container in the Minikube docker environment

```
make build
```

## Deploy Tokenizer in Minikube
```
make deploy
```

## Testing Tokenizer Scaling
```
make load-test
```

## Tear Down Environment 
```
make clean-all
```