package core

type FanInResult[T any] struct {
	i    int
	data T
}

func FanIn[T any](channels []chan T) chan FanInResult[T] {
	out := make(chan FanInResult[T])

	for i, channel := range channels {
		go func() {
			for {
				data, ok := <-channel
				if !ok {
					break
				}

				out <- FanInResult[T]{i, data}
			}
		}()
	}

	return out
}

func FanOut[T any](output []chan T) chan T {
	input := make(chan T)

	go func() {
		for {
			data, ok := <-input
			if !ok {
				break
			}

			for _, out := range output {
				out <- data
			}
		}

		for _, out := range output {
			close(out)
		}
	}()

	return input
}
