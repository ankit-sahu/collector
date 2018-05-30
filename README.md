# collector
Collector - Collects resource information from Kubernetes cluster and passes it on to another service.

# Setup:
- To run from source code, needs golang installed in the system.
- git clone <this repo>
- switch to collector director
- execute : go run main.go

# How it works?
- Provide the url of the cluster from where you want to scrape the metric in clusterURL variable.
- Provide the kubernetes cluster admin token in clusterToken variable.
  Eg -> token: "Bearer your_cluster_admin_token"
- Provide the destination service where this data has to be sent in destinationURL variable.
- collector -> config -> appConfig.yaml consists of list of metricResources you want to capture from kubernetes.
  Add the list of {Name(resource name) and path to access the resource} as key value pair.
- The provides key value pair of resources would be picked from appConfig.yaml file metric would be collected and sent to destination service.  
