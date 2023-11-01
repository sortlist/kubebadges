/*
Copyright 2023 Your Company.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

// Code generated by informer-gen. DO NOT EDIT.

package v1

import (
	"context"
	time "time"

	kubebadgesv1 "github.com/kubebadges/kubebadges/pkg/apis/kubebadges/v1"
	versioned "github.com/kubebadges/kubebadges/pkg/generated/clientset/versioned"
	internalinterfaces "github.com/kubebadges/kubebadges/pkg/generated/informers/externalversions/internalinterfaces"
	v1 "github.com/kubebadges/kubebadges/pkg/generated/listers/kubebadges/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	runtime "k8s.io/apimachinery/pkg/runtime"
	watch "k8s.io/apimachinery/pkg/watch"
	cache "k8s.io/client-go/tools/cache"
)

// KubeBadgeInformer provides access to a shared informer and lister for
// KubeBadges.
type KubeBadgeInformer interface {
	Informer() cache.SharedIndexInformer
	Lister() v1.KubeBadgeLister
}

type kubeBadgeInformer struct {
	factory          internalinterfaces.SharedInformerFactory
	tweakListOptions internalinterfaces.TweakListOptionsFunc
	namespace        string
}

// NewKubeBadgeInformer constructs a new informer for KubeBadge type.
// Always prefer using an informer factory to get a shared informer instead of getting an independent
// one. This reduces memory footprint and number of connections to the server.
func NewKubeBadgeInformer(client versioned.Interface, namespace string, resyncPeriod time.Duration, indexers cache.Indexers) cache.SharedIndexInformer {
	return NewFilteredKubeBadgeInformer(client, namespace, resyncPeriod, indexers, nil)
}

// NewFilteredKubeBadgeInformer constructs a new informer for KubeBadge type.
// Always prefer using an informer factory to get a shared informer instead of getting an independent
// one. This reduces memory footprint and number of connections to the server.
func NewFilteredKubeBadgeInformer(client versioned.Interface, namespace string, resyncPeriod time.Duration, indexers cache.Indexers, tweakListOptions internalinterfaces.TweakListOptionsFunc) cache.SharedIndexInformer {
	return cache.NewSharedIndexInformer(
		&cache.ListWatch{
			ListFunc: func(options metav1.ListOptions) (runtime.Object, error) {
				if tweakListOptions != nil {
					tweakListOptions(&options)
				}
				return client.KubebadgesV1().KubeBadges(namespace).List(context.TODO(), options)
			},
			WatchFunc: func(options metav1.ListOptions) (watch.Interface, error) {
				if tweakListOptions != nil {
					tweakListOptions(&options)
				}
				return client.KubebadgesV1().KubeBadges(namespace).Watch(context.TODO(), options)
			},
		},
		&kubebadgesv1.KubeBadge{},
		resyncPeriod,
		indexers,
	)
}

func (f *kubeBadgeInformer) defaultInformer(client versioned.Interface, resyncPeriod time.Duration) cache.SharedIndexInformer {
	return NewFilteredKubeBadgeInformer(client, f.namespace, resyncPeriod, cache.Indexers{cache.NamespaceIndex: cache.MetaNamespaceIndexFunc}, f.tweakListOptions)
}

func (f *kubeBadgeInformer) Informer() cache.SharedIndexInformer {
	return f.factory.InformerFor(&kubebadgesv1.KubeBadge{}, f.defaultInformer)
}

func (f *kubeBadgeInformer) Lister() v1.KubeBadgeLister {
	return v1.NewKubeBadgeLister(f.Informer().GetIndexer())
}
