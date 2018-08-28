## Running URL shortener on kubernetes cluster

### Prerequisites:
1. Kubernetes cluster (GKE/AWS/AKS/self-hosted)
2. [kubectl](https://kubernetes.io/docs/reference/kubectl/overview/) installed.

### Steps:
1. Head over to `deploy/k8s` directory of the repository.
2. Create a persistent volume and persistent volume claim with `kubectl apply -f disk-pv.yml`
3. Create a deployment for URL shortener service with `kubectl apply -f go-url-shortener.yml`
4. Verify if the service is up and running by `kubectl get svc | grep 'url-shortener-service'`
