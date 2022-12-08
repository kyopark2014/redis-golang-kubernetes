# redis-golang-kubernetes

This repository shows how to create and use redis based on golang in kubernetes.

## Preperation
### Install eksctl

For linux, use these commends as [[Installing or updating eksctl]](https://docs.aws.amazon.com/eks/latest/userguide/eksctl.html)

```java
curl --silent --location "https://github.com/weaveworks/eksctl/releases/latest/download/eksctl_$(uname -s)_amd64.tar.gz" | tar xz -C /tmp

sudo mv /tmp/eksctl /usr/local/bin

eksctl version
```

### SSH Key Gnenration

In order to access node, generate ssh key as bellow. Then copy into "./ssh"

```java
ssh-keygen
```

### Kubectl

To install kubectl, use these commends as [[Install and Set Up kubectl on Linux]](https://kubernetes.io/docs/tasks/tools/install-kubectl-linux/)

```java
curl -LO "https://dl.k8s.io/release/$(curl -L -s https://dl.k8s.io/release/stable.txt)/bin/linux/amd64/kubectl"

sudo install -o root -g root -m 0755 kubectl /usr/local/bin/kubectl

kubectl cluster-info
```

If connection is failed, use this

```java
aws eks --region ap-northeast-2 update-kubeconfig --name eks-ksdyb
```

## Main

### Create EKS cluster
```c
eksctl create cluster -f k8s/cluster-redis-golang-kubernetes.yaml
```

### Deploy Redis server
```c
kubectl create -f k8s/redis-master-deployment.yaml

kubectl create -f k8s/redis-master-service.yaml 
```

### Change the type from ClusterIP to LoadBalancer in order to easily use 

```c
$ kubectl edit service/redis-master

$ kubectl get service/redis-master
NAME           TYPE           CLUSTER-IP     EXTERNAL-IP                                                               PORT(S)          AGE
redis-master   LoadBalancer   10.100.22.49   a2a61bdc0208d11eaaabc0a6b8228ff9-2077444820.eu-west-2.elb.amazonaws.com   6379:32502/TCP   28m
```

### Input these commands using redis-cli 

```c
$ redis-cli -h a2a61bdc0208d11eaaabc0a6b8228ff9-2077444820.eu-west-2.elb.amazonaws.com -p 6379
a2a61bdc0208d11eaaabc0a6b8228ff9-2077444820.eu-west-2.elb.amazonaws.com:6379> lpush comments "hello"
(integer) 1
a2a61bdc0208d11eaaabc0a6b8228ff9-2077444820.eu-west-2.elb.amazonaws.com:6379> lpush comments "world"
(integer) 2
a2a61bdc0208d11eaaabc0a6b8228ff9-2077444820.eu-west-2.elb.amazonaws.com:6379> lpush comments "testing"
(integer) 3
a2a61bdc0208d11eaaabc0a6b8228ff9-2077444820.eu-west-2.elb.amazonaws.com:6379> lpush comments "1 2 3"
(integer) 4
a2a61bdc0208d11eaaabc0a6b8228ff9-2077444820.eu-west-2.elb.amazonaws.com:6379> lpush comments "1 2 3"
(integer) 5
a2a61bdc0208d11eaaabc0a6b8228ff9-2077444820.eu-west-2.elb.amazonaws.com:6379> exit
```

### Prepare main.go and index.html as bellow

#### main.go
```c
package main

import (
	"net/http"
	"github.com/gorilla/mux"
	"github.com/go-redis/redis"
	"html/template"
)

var client *redis.Client
var templates *template.Template

func main() {
	client = redis.NewClient(&redis.Options{
		Addr: "a2a61bdc0208d11eaaabc0a6b8228ff9-2077444820.eu-west-2.elb.amazonaws.com:6379",
	})
	templates = template.Must(template.ParseGlob("templates/*.html"))
	r := mux.NewRouter()
	r.HandleFunc("/", indexHandler).Methods("GET")
	http.Handle("/", r)
	http.ListenAndServe(":8080", nil)
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	comments, err := client.LRange("comments", 0, 10).Result()
	if err != nil {
		return
	}
	templates.ExecuteTemplate(w, "index.html", comments)
}
```

#### templates/index.html

```c
<html>
	<head>
		<title>Comments</title>
	</head>
	<body>
		<h1>Comments</h1>
		{{ range . }}
		<div>{{ . }}</div>
		{{ end }}
	</body>
</html>
```

### Run it for checking the operation.
```c
$ go run main.go

$ curl -i localhost:8080
```

```c
HTTP/1.1 200 OK
Date: Tue, 17 Dec 2019 13:00:51 GMT
Content-Length: 210
Content-Type: text/html; charset=utf-8

<html>
	<head>
		<title>Comments</title>
	</head>
	<body>
		<h1>Comments</h1>
		
		<div>1 2 3</div>
		
		<div>1 2 3</div>
		
		<div>testing</div>
		
		<div>hello</div>
		
		<div>hello</div>
		
	</body>
</html>
```

### Make docker image
```c
$ docker build -t redis-golang-kubernetes:v1 .

$ docker run -d -p 8080:8080 redis-golang-kubernetes:v1
```

### Tagging
```c
$ docker tag redis-golang-kubernetes:v1 884942771862.dkr.ecr.eu-west-2.amazonaws.com/repository-redis-golang
```

### Create repository if required
```c
$ aws ecr create-repository --region eu-west-2 --repository-name repository-redis-golang
```

### Push the image to ECR
```c
$ docker push 994942771862.dkr.ecr.eu-west-2.amazonaws.com/repository-redis-golang
```

### Deploy and run
```c
$ kubectl create -f redis-golang-kubernetes-deployment.yaml
$ kubectl create -f redis-golang-kubernetes-server-service.yaml 
```

### check the url in order to show the operation of radis using golang.
```c
$ kubectl get svc
NAME                            TYPE           CLUSTER-IP       EXTERNAL-IP                                                               PORT(S)                      AGE
redis-golang                    LoadBalancer   10.100.244.82    a1cf9cffe20cf11ea9c1206f47235bc6-737905859.eu-west-2.elb.amazonaws.com    8080:32085/TCP               12s
redis-master                    LoadBalancer   10.100.22.49     a2a61bdc0208d11eaaabc0a6b8228ff9-2077444820.eu-west-2.elb.amazonaws.com   6379:32502/TCP               7h52m
```


### Make sure the opration using the url.
```c
$ curl -i a1cf9cffe20cf11ea9c1206f47235bc6-737905859.eu-west-2.elb.amazonaws.com:8080
```

```c
HTTP/1.1 200 OK
Date: Tue, 17 Dec 2019 13:25:35 GMT
Content-Length: 210
Content-Type: text/html; charset=utf-8

<html>
	<head>
		<title>Comments</title>
	</head>
	<body>
		<h1>Comments</h1>
		
		<div>1 2 3</div>
		
		<div>1 2 3</div>
		
		<div>testing</div>
		
		<div>hello</div>
		
		<div>hello</div>
		
	</body>
</html>
```

## Delete Resouces 

```java
eksctl delete cluster --region=ap-northeast-2 --name=eks-ksdyb
```

## Reference

https://kubernetes.io/ko/docs/tutorials/stateless-application/guestbook/

https://www.youtube.com/watch?v=Hbt56gFj998&feature=youtu.be
