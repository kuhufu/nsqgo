package watcher

func sub(left, right []string) []string {
	var res []string
	for _, v := range left {
		if !contains(right, v) {
			res = append(res, v)
		}
	}

	return res
}

func contains(src []string, s string) bool {
	for _, e := range src {
		if e == s {
			return true
		}
	}

	return false
}
