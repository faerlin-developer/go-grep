package grep

type Config struct {
	numberWorkers int
	bufferSize    int
}

func (c *Config) SetNumberWorkers(numberWorkers int) {
	c.numberWorkers = numberWorkers
}

func (c *Config) SetBufferSize(bufferSize int) {
	c.bufferSize = bufferSize
}

func (c *Config) GetNumberWorkers() int {
	return c.numberWorkers
}

func (c *Config) GetBufferSize() int {
	return c.bufferSize
}
