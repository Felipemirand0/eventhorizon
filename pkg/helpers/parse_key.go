package eventhorizon

import (
	"strings"

	"k8s.io/client-go/tools/cache"
)

func ParseKey(key string) (string, string, string) {
	strs := strings.Split(key, "#")

	ns, name, _ := cache.SplitMetaNamespaceKey(strs[1])

	if "" == ns {
		ns = "default"
	}

	return strs[0], ns, name
}
