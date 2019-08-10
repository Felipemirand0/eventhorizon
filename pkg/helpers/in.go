package eventhorizon

func In(key string, list []string) bool {
	for _, val := range list {
		if val == key {
			return true
		}
	}

	return false
}
