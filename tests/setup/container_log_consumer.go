package setup

import (
	"fmt"

	"github.com/testcontainers/testcontainers-go"
)

var _ testcontainers.LogConsumer = (*containerLogsConsumer)(nil)

// containerLogsConsumer collects logs from a container.
type containerLogsConsumer struct {
	containerName string
	logs          []byte
}

func newContainerLogsConsumer(containerName string) *containerLogsConsumer {
	return &containerLogsConsumer{containerName: containerName}
}

// Accept records a log message from the container.
// It implements [testcontainers.LogConsumer] interface.
func (c *containerLogsConsumer) Accept(log testcontainers.Log) {
	c.logs = append(c.logs, log.Content...)
}

// Print returns prints all logs.
func (c *containerLogsConsumer) Print() {
	fmt.Printf(`### Start of %[1]s container logs
%[2]s
### End of %[1]s container logs
`, c.containerName, c.logs)
}
