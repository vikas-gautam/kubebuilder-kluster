# kubebuilder-kluster
https://book.kubebuilder.io/quick-start.html

```
 #it will generate makefile, initialise your go module and main.go 
  kubebuilder init --domain golearning.dev --repo github.com/vikas-gautam/kubebuilder-kluster 

 #create api and implement set  group & version
  kubebuilder create api --group demo --version v1alpha1 --kind Kluster

#modify types.go and then run
 make manifests && make generate && make install && make run

# modify controller to write your business logic

```
