package dependencies

import (
	"context"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)


type VerneMQ struct {
	container testcontainers.Container;
}

func NewVerneMQ(ctx context.Context) (*VerneMQ, error) {
	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: testcontainers.ContainerRequest{
			Image:           "ghcr.io/senergy-platform/vernemq:prod",
			Tmpfs:           map[string]string{},
			ExposedPorts:    []string{"1883/tcp"},
			WaitingFor:      wait.ForListeningPort("1883/tcp"),
			AlwaysPullImage: true,
			Env: map[string]string{
				"DOCKER_VERNEMQ_ALLOW_ANONYMOUS": "on",
			},
		},
		Started: false,
		
	})
	if err != nil {
		return &VerneMQ{}, err
	}
	return &VerneMQ{
		container: container,
	}, nil
}

func (m *VerneMQ) StartAndWait(ctx context.Context) (error, string) {
	err := m.container.Start(ctx)
	if err != nil {
		return err, ""
	}
	localhostPort, err := m.container.MappedPort(ctx, "1883")
	if err != nil {
		return err, ""
	}
	return nil, localhostPort.Port()
}