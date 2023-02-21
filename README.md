# stock-service-operator
stock-service-operator is a simple kubernetes operator that retrieves the latest stock price and stores it in a config map.

## Description

stock-service-operator supports the below fields:

1. symbol is the stock symbol eg:- AAPL, GOOGL etc
2. api is the finance api to use to fetch the stock price. The operator currently supports yahoo and vantage finance apis. The operator is written in such a way that new apis can be added without making any changes to the controllers. The backend api code is in the backend package.
3. refreshInterval is the interval at which stock prices are refreshed.

An example config using yahoo finance 
```
apiVersion: stock-service.edb.com/v1
kind: StockReport
metadata:
  name: stockreport-sample
spec:
  symbol: AAPL
  api: yahoo
  refreshInterval: 60s

```
On applying the yaml, if the operator succesfully retrieves the stock price, then it would create a configmap with the stock details. 
```
(base) joeljacob@Joels-MacBook-Pro stock-service-operator % kubectl get stockreport stockreport-sample
NAME                 CONFIG MAP              LAST REFRESHED   STATUS
stockreport-sample   stockreport-sample-cm   13s              Ready
```
This would create a configmap like below
```
(base) joeljacob@Joels-MacBook-Pro stock-service-operator % kubectl get configmap stockreport-sample-cm -o yaml

apiVersion: v1
data:
  api: yahoo
  price: "152.55"
  stock_symbol: AAPL
  updateTime: 2023-02-20 23:46:03.4816821 +0000 UTC m=+505.495988401
kind: ConfigMap
metadata:
  creationTimestamp: "2023-02-20T23:44:03Z"
    ---
```

An example config using vantage finance api is shown below. To use vantage api you must create a free vantage api key which can be done here https://www.alphavantage.co/ . You need to then assign the api key to constant `APIKey` in backend/vantage.go. In a production environment, the api key should be stored somewhere safe and retrieved from there. 

**Note:** Vantage free api service provides only upto 5 api requests per minute and 500 requests per day.

```
apiVersion: stock-service.edb.com/v1
kind: StockReport
metadata:
  name: stockreport-sample
spec:
  symbol: AAPL  
  api: vantage  // remember to create an apikey from https://www.alphavantage.co/ . 
  refreshInterval: 60s

```

**Note:** More examples can be found in config/samples

### Tests

Tests have been setup using ginkgo framework, To run the tests 
```sh
make test
```

### ValidatingWebhookConfiguration
Validation webhooks have been setup to validate the symbol and refreshInterval provided in the config, this will prevent customers from creating incorrect configs. 

**Note:** The operator needs to be deployed on the cluster for the webhooks to work. The steps to do this is discussed in the `Running on the cluster` section below.


Consider the below
test.yaml
```
apiVersion: stock-service.edb.com/v1
kind: StockReport
metadata:
  name: test
spec:
  symbol: test-symbol \\ invalid symbol
  api: yahoo
  refreshInterval: 5m
```
Webhook will block the yaml 

```
(base) joeljacob@Joels-MacBook-Pro stock-service-operator % kubectl apply -f test.yaml
The StockReport "test" is invalid: spec.symbol: Invalid value: "test-symbol": error updating test due to invalid symbol test-symbol
```



## Getting Started
You’ll need a Kubernetes cluster to run against. You can use [KIND](https://sigs.k8s.io/kind) to get a local cluster for testing, or run against a remote cluster.
**Note:** Your controller will automatically use the current context in your kubeconfig file (i.e. whatever cluster `kubectl cluster-info` shows).


### Running locally 
1. Install the CRDs into the cluster:

```sh
make install
```

2. Run your controller (this will run in the foreground, so switch to a new terminal if you want to leave it running):

```sh
make run
```

### Running on the cluster

This will run the operator as a pod. Webhooks are enabled by default, it can be disabled by setting `ENABLE_WEBHOOKS` env variable to false in config/manager/manager.yaml. 

1. Build your image:

```sh
make docker-build  IMG=<some-registry>/stock-service-operator:tag
```

2. Install cert-manager for running webhooks 

```sh
kubectl apply -f https://github.com/cert-manager/cert-manager/releases/download/v1.11.0/cert-manager.yaml
```

3. Load the image into the kind cluster

```sh
kind load docker-image <image-name>
```

4. Deploy the controller to the cluster with the image specified by `IMG`:

```sh
make deploy IMG=<some-registry>/stock-service-operator:tag
```

5. Switch to stock-service-operator-system to see the operator pod

```sh
kubectl config set-context --current --namespace=stock-service-operator-system
```

### Uninstall CRDs
To delete the CRDs from the cluster:

```sh
make uninstall
```

### Undeploy controller
UnDeploy the controller from the cluster:

```sh
make undeploy
```

### Modifying the API definitions
If you are editing the API definitions, generate the manifests such as CRs or CRDs using:

```sh
make manifests
```
