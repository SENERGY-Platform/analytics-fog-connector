package dependencies

import (
	"context"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"os"
	"path/filepath"
)


type VerneMQ struct {
	container testcontainers.Container;
}

func NewVerneMQ(ctx context.Context) (*VerneMQ, error) {
	absPath, err := filepath.Abs(filepath.Join("..", "dependencies", "acl.txt"))
	if err != nil {
		return &VerneMQ{}, err
	}
	r, err := os.Open(absPath)
	if err != nil {
		return &VerneMQ{}, err 
	}

	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: testcontainers.ContainerRequest{
			Image:           "vernemq/vernemq",
			Tmpfs:           map[string]string{},
			ExposedPorts:    []string{"1883/tcp"},
			WaitingFor:      wait.ForListeningPort("1883/tcp"),
			AlwaysPullImage: true,
			Env: map[string]string{
				"DOCKER_VERNEMQ_LOG__CONSOLE__LEVEL": "debug",
				"DOCKER_VERNEMQ_ALLOW_ANONYMOUS": "on",
				"DOCKER_VERNEMQ_ACCEPT_EULA": "yes",
			},
			Files: []testcontainers.ContainerFile{
				{
					Reader:            r,
					HostFilePath:      "./tests/dependencies/acl.txt", // will be discarded internally
					ContainerFilePath: "/etc/vernemq/vmq.acl",
					FileMode:          0o777,
				},
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