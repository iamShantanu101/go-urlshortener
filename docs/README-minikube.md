## Running URL shortener with minikube

### Steps:
1. Clone the repository
2. Head over to `deploy/minikube` directory in the repository.
3. Run `bash deploy_minikube.sh`. 

### What does this script do?
1. Checks if prerequisites are installed. Prerequisites checked during the process:
   1. Docker
   2. kubectl
   3. virtualbox
   3. minikube
2. Creates a local cluster of Kubernetes version `v1.10.0`.
3. Deploys the URL shortener application on the cluster.
4. At last, outputs the external ip where the service is accessible.
