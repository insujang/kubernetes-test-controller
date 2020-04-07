# Kubernetes Test Custom Controller

This is an implementation of a Kubernetes CRD and a custom controller using [code-generator](https://github.com/kubernetes/code-generator). This repository contains one Go module (`insujang.github.io/kubernetes-test-controller`) with a generated code regarding the CRD (`/lib/testresource`) and a custom controller based on this (`/cmd/controller`).

The code generates and uses the following CRD type:

```
Name:         testresources.insujang.github.io
Namespace:    
Labels:       <none>
Annotations:  <none>
API Version:  apiextensions.k8s.io/v1
Kind:         CustomResourceDefinition
Metadata: <OMITTED>
    Manager:         kube-apiserver
    Operation:       Update
    Time:            2020-04-06T10:07:43Z
  Resource Version:  21349980
  Self Link:         /apis/apiextensions.k8s.io/v1/customresourcedefinitions/testresources.insujang.github.io
  UID:               
Spec:
  Conversion:
    Strategy:  None
  Group:       insujang.github.io
  Names:
    Kind:       TestResource
    List Kind:  TestResourceList
    Plural:     testresources
    Short Names:
      tr
    Singular:               testresource
  Preserve Unknown Fields:  true
  Scope:                    Namespaced
  Versions:
    Additional Printer Columns:
      Json Path:  .spec.command
      Name:       command
      Type:       string
    Name:         v1beta1
    Schema:
      openAPIV3Schema:
        Properties:
          Spec:
            Properties:
              Command:
                Pattern:  ^(echo).*
                Type:     string
              Custom Property:
                Type:  string
            Required:
              command
              customProperty
            Type:  object
        Type:      object
    Served:        true
    Storage:       true
Status:
  Accepted Names:
    Kind:       TestResource
    List Kind:  TestResourceList
    Plural:     testresources
    Short Names:
      tr
    Singular:  testresource
  Conditions:
    Last Transition Time:  2020-04-06T10:07:43Z
    Message:               no conflicts found
    Reason:                NoConflicts
    Status:                True
    Type:                  NamesAccepted
    Last Transition Time:  2020-04-06T10:07:43Z
    Message:               the initial names have been accepted
    Reason:                InitialNamesAccepted
    Status:                True
    Type:                  Established
  Stored Versions:
    v1beta1
Events:  <none>
```

# Directory Structure

## 1. `/lib/testresource`

This directory, at first, contains three files for code generation; `v1beta1/doc.go`, `v1beta1/register.go`, and `v1beta1/types.go`. After generating code with code-generator, `v1beta1/zz_generated.deepcopy.go` and `generated` are generated.

## 2. `/cmd/controller`

This directory contains basic controller logics that creates a CRD type itself, creates a custom resource with the CRD type, and watches events regarding the CRD type.

# Usage

Use Dockerfile to build a custom controller. All operations are done by the controller (creating a CRD, creating a custom resource, and watching events).

```
$ git clone https://github.com/insujang/kubernetes-test-controller
$ cd kubernetes-test-controller
$ docker build -t controller-test .
$ docker run -it --rm --net=host -v $KUBECONFIG:$KUBECONFIG -e KUBECONFIG=$KUBECONFIG controller-test
```

## Expected output

```shell
~/kubernetes-test-controller$ docker run -it --rm --net=host -v $KUBECONFIG:$KUBECONFIG -e KUBECONFIG=$KUBECONFIG controller-test
go: downloading k8s.io/api v0.17.0
go: downloading k8s.io/apiextensions-apiserver v0.17.0
go: downloading k8s.io/apimachinery v0.17.0
...
I0407 05:13:44.500578    2240 handler.go:21] Waiting cache to be synced.
I0407 05:13:44.500610    2240 handler.go:33] Starting custom controller.
I0407 05:13:44.500423    2240 handler.go:40] Added: &{{ } {example-tr2  default /apis/insujang.github.io/v1beta1/namespaces/default/testresources/example-tr2 5890795a-81fb-4fe6-85d8-3ec8db89145d 21501514 1 2020-04-07 05:13:44 +0000 UTC <nil> <nil> map[] map[] [] []  [{controller Update insujang.github.io/v1beta1 2020-04-07 05:13:44 +0000 UTC FieldsV1 &FieldsV1{Raw:*[123 34 102 58 115 112 101 99 34 58 123 34 46 34 58 123 125 44 34 102 58 99 111 109 109 97 110 100 34 58 123 125 44 34 102 58 99 117 115 116 111 109 80 114 111 112 101 114 116 121 34 58 123 125 125 44 34 102 58 115 116 97 116 117 115 34 58 123 125 125],}}]} {echo Hello World! asdasd=1234} }
I0407 05:13:44.503252    2240 types.go:72] Updated: &{{ } {example-tr2  default /apis/insujang.github.io/v1beta1/namespaces/default/testresources/example-tr2 5890795a-81fb-4fe6-85d8-3ec8db89145d 21501515 2 2020-04-07 05:13:44 +0000 UTC <nil> <nil> map[] map[] [] []  [{controller Update insujang.github.io/v1beta1 2020-04-07 05:13:44 +0000 UTC FieldsV1 &FieldsV1{Raw:*[123 34 102 58 115 112 101 99 34 58 123 34 46 34 58 123 125 44 34 102 58 99 111 109 109 97 110 100 34 58 123 125 44 34 102 58 99 117 115 116 111 109 80 114 111 112 101 114 116 121 34 58 123 125 125 44 34 102 58 115 116 97 116 117 115 34 58 123 125 125],}}]} {echo Hello World! asdasd=1234} HANDLED}
I0407 05:13:44.503321    2240 event.go:281] Event(v1.ObjectReference{Kind:"TestResource", Namespace:"default", Name:"example-tr2", UID:"5890795a-81fb-4fe6-85d8-3ec8db89145d", APIVersion:"insujang.github.io/v1beta1", ResourceVersion:"21501514", FieldPath:""}): type: 'Normal' reason: 'ObjectHandled' Object is handled by custom controller.
```