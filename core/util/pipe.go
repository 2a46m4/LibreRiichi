package core

type FanInResult[T any] struct {
	I    int
	Data T
}

// Returns a channel that collects messages from channels with the same type.
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

// Returns a channel that accepts a single message and sends them to all channels
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
