package client

import (
	"time"

	"k8s.io/client-go/informers"
	"k8s.io/client-go/tools/cache"
)

type AddObjectFunc func(obj interface{})
type UpdateObjectFunc func(oldObj, newObj interface{})
type DeleteObjectFunc func(obj interface{})

func (k *K8SClient) RegisterConfigMapChangeHandler(addFunc AddObjectFunc,
	updateFn UpdateObjectFunc, deleteFunc DeleteObjectFunc) {
	informerFactory := informers.NewSharedInformerFactory(k.client, time.Second*30)
	configMapInformer := informerFactory.Core().V1().ConfigMaps().Informer()
	configMapInformer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc:    addFunc,
		UpdateFunc: updateFn,
		DeleteFunc: deleteFunc,
	})
	k.configMapInformer = configMapInformer
}

func (k *K8SClient) StartObjectChangeInformer() {
	stopCh := make(chan struct{})
	defer close(stopCh)
	go k.informerFactory.Start(stopCh)
	k.informerFactory.WaitForCacheSync(stopCh)
	<-stopCh
}
