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


IMP points -

CR - action on CR resource like status update with field addition
in above case we need custom k8s client - klient klientset.Interface

OR

client.Client from "sigs.k8s.io/controller-runtime/pkg/client" [Controller-runtime has in built function to get kubeconfig -  ctrl.GetConfigOrDie() and creating newmanager ctrl.NewManager()]
it's similar to kubectl command
-------------

But, for native resorces we could use k8sclient kubernetes.Interface or *kuberenetes.clientset from client-go/kuberentes package

--------------
one client we use in controller and one we use in utils folders
---------------------------------------------------------------------------
