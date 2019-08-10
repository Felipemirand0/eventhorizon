package eventhorizon

import (
	"k8s.io/client-go/tools/cache"
)

func ValidSynced(items []cache.InformerSynced) []cache.InformerSynced {
	list := []cache.InformerSynced{}

	for _, item := range items {
		if nil != item {
			list = append(list, item)
		}
	}

	return list
}
