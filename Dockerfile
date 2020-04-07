FROM golang:1.13 as clone
WORKDIR /
RUN git clone --branch v0.17.0 https://github.com/kubernetes/code-generator

# Generate code based on the input Go code.
FROM golang:1.13 as generatecode
COPY --from=clone /code-generator /code-generator/
COPY ./ /kubernetes-test-controller/
WORKDIR /code-generator
RUN go mod edit -replace=insujang.github.io/kubernetes-test-controller=../kubernetes-test-controller
RUN ./generate-groups.sh all insujang.github.io/kubernetes-test-controller/lib/testresource/generated insujang.github.io/kubernetes-test-controller/lib testresource:v1beta1 --go-header-file ./hack/boilerplate.go.txt --output-base ..
RUN cp -r /insujang.github.io/kubernetes-test-controller/lib/testresource/* /kubernetes-test-controller/lib/testresource/ && rm -r /insujang.github.io

# Run custom resource controller. Requires KUBECONFIG env and a running Kubernetes cluster.
FROM golang:1.13
COPY --from=generatecode /kubernetes-test-controller/ /kubernetes-test-controller/
WORKDIR /kubernetes-test-controller/cmd/controller
ENTRYPOINT [ "go", "run", "." ]
