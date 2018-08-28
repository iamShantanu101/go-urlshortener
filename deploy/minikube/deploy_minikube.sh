#!/bin/bash

red=$( tput setaf 1 )
green=$( tput setaf 2 )
NC=$( tput sgr0 )

check_minikube_install() {
    type minikube >/dev/null 2>&1 || { error_exit "${red}Minikube is not installed. Installing minikube.${NC}" && install_minikube; }
}

check_kubectl_install() {
    type kubectl >/dev/null 2>&1 || { error_exit "${red}Kubectl is not installed. Installing kubectl.${NC}" && install_kubectl; }
}

check_docker_install() {
    type docker >/dev/null 2>&1 || { error_exit "${red}Docker is not installed. Installing Docker.${NC}" && install_docker; }
}

install_docker() {
    echo -e "Starting Docker installation.\n"
    curl -fsSL get.docker.com -o get-docker.sh
    sh get-docker.sh || { error_exit "${red}Error install Docker. Exiting.${NC}" exit; }
    echo -e "Successfully installed Docker.\n"
    rm get-docker.sh
}

install_minikube() {

    echo -e "Minikube installation started.\n"

    # Update apt repo
    sudo apt update -y

    # Install virtualbox
    sudo apt install -y virtualbox >/dev/null 2>&1 || { error_exit "${red}Error install virtualbox. Exiting.${NC}" exit; }

    # Make sure no prior copy of minikube exists.
    sudo rm -rf .minikube/

    # Install minikube
    curl -Lo minikube https://storage.googleapis.com/minikube/releases/latest/minikube-linux-amd64 && chmod +x minikube && sudo cp minikube /usr/local/bin/ && rm minikube || { error_exit "${red}Error installing minikube. Exiting.${NC}" exit; }

    echo -e "${green}Minikube installed successfully.${NC}\n"
}

install_kubectl() {

    echo -e "kubectl installation started.\n"
    sudo apt update && sudo apt install -y apt-transport-https
    curl -s https://packages.cloud.google.com/apt/doc/apt-key.gpg | sudo apt-key add -
    sudo touch /etc/apt/sources.list.d/kubernetes.list
    echo "deb http://apt.kubernetes.io/ kubernetes-xenial main" | sudo tee -a /etc/apt/sources.list.d/kubernetes.list
    sudo apt update
    sudo apt install -y kubectl >/dev/null 2>&1 || { error_exit "${red}Error installing kubectl.${NC}" exit; }

    echo -e "${green}Kubectl installed successfully.${NC}\n"
}


setup_url_shortener() {

    echo -e "==================Starting minikube v1.10.0 cluster.===================\n"
    minikube start --kubernetes-version v1.10.0 >/dev/null 2>&1 || { error_exit "${red}Error starting minikube kubernetes cluster. Exiting.${NC}" exit; }
    echo -e "====================Starting URL shortener service.====================\n"
    echo -e "${green}Creating PV and PV claim.${NC}\n"
    kubectl apply -f ../k8s/disk-pv.yml >/dev/null 2>&1 || { error_exit "${red}Error creating PV and PVC${NC}" exit; }
    echo -e "${green}Creating deployment for URL shortener service.${NC}\n"
    kubectl apply -f ../k8s/go-url-shortener.yml >/dev/null 2>&1 || { error_exit "${red}Error creating deployment for URL shortener${NC}" exit; }
    echo -e "Wait a couple of seconds until URL shortening service gets an external IP to access. You can get the external IP with:\n"
    echo -e "minikube service url-shortener-service --url\n"
    echo -e "Getting URL shortener service URL.\n"
    sleep 20
    minikube service url-shortener-service --url

}

error_exit() {
    echo "$1" 1>&2
}

check_minikube_install

check_kubectl_install

check_docker_install

setup_url_shortener
