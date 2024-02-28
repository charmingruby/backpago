package queue

type MockQueue struct {
	q []*QueueDto
}

func (mq *MockQueue) Publish(msg []byte) error {
	dto := new(QueueDto)
	dto.Unmarshal(msg)

	mq.q = append(mq.q, dto)

	return nil
}

func (mq *MockQueue) Consume(c chan<- QueueDto) error {
	return nil
}
