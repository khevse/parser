package engine

type buffer struct {
	Queue chan interface{}
}

func newBuffer(count int) *buffer {
	return &buffer{
		Queue: make(chan interface{}, count),
	}
}

func (b *buffer) Insert(val interface{}) {
	b.Queue <- val
}

func (b *buffer) GetData() <-chan interface{} {

	lenght := cap(b.Queue)
	retval := make(chan interface{})

	if lenght == 0 {
		close(retval)
	} else {
		go func() {
			var count int
			for item := range b.Queue {
				retval <- item
				count += 1

				if count == lenght {
					close(retval)
					return
				}
			}
		}()
	}

	return retval
}
