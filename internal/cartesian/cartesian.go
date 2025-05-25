package cartesian

// https://github.com/schwarmco/go-cartesian-product
// Modified to work with just strings

func Iter(params ...[]string) chan []string {
	c := make(chan []string)
	if len(params) == 0 {
		close(c)
		return c
	}
	go func() {
		iterate(c, params[0], []string{}, params[1:]...)
		close(c)
	}()
	return c
}

func iterate(channel chan []string, topLevel, result []string, needUnpacking ...[]string) {
	if len(needUnpacking) == 0 {
		for _, p := range topLevel {
			channel <- append(append([]string{}, result...), p)
		}
		return
	}
	for _, p := range topLevel {
		iterate(channel, needUnpacking[0], append(result, p), needUnpacking[1:]...)
	}
}
