package cartesian

// https://github.com/schwarmco/go-cartesian-product
// Modified in various ways, mostly to add goroutines (poorly)

import "sync"

func Iter(params ...[]string) chan []string {
	c := make(chan []string)
	if len(params) == 0 {
		close(c)
		return c
	}

	go func() {
		var wg sync.WaitGroup
		topLevel := params[0]
		rest := params[1:]

		for _, p := range topLevel {
			wg.Add(1)
			go func(prefix string) {
				defer wg.Done()
				partial := []string{prefix}
				iterate(c, partial, rest...)
			}(p)
		}

		wg.Wait()
		close(c)
	}()

	return c
}

func iterate(channel chan []string, result []string, params ...[]string) {
	if len(params) == 0 {
		channel <- append([]string{}, result...)
		return
	}

	for _, p := range params[0] {
		iterate(channel, append(result, p), params[1:]...)
	}
}
