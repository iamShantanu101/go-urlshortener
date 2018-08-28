## Stress Testing with Vegeta

### Attack overview
1. File `attack.go` contains rules written for performing an attack. This file uses [vegeta](https://github.com/tsenart/vegeta) as a library.
2. For `POST` requests, a file `reqBody.json` contains the JSON payload which will be injected in the attack.
3. Run `go build attack.go` to build the executable file.
4. Run `./attack` for running the attack. This will give the basic stats in `stdout` as well as output an HTML formatted report.

### Kubernetes Deployment
1. Vegeta can also be deployed on a kubernetes cluster.
2. Head over to `benchmark/load_test_vegeta/deploy_k8s` inside the repository.
3. For simply deploying vegeta on a kubernetes cluster run `kubectl apply -f vegeta-k8s.yml` This will pull image `ishantanu16/vegeta-docker-alpine` from dockerhub repository. If you don't want to use this, you can build your own image.
