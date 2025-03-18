package helpers

import (
	"os/exec"
	"strings"
	"testing"
	"time"

	"github.com/agglayer/go_signer/log"
	"github.com/stretchr/testify/require"
)

type DockerCompose struct {
}

func NewDockerCompose() *DockerCompose {
	return &DockerCompose{}
}

func (d *DockerCompose) Up(t *testing.T) {
	log.Debug("starting docker")
	msg, err := exec.Command("bash", "-l", "-c", "docker compose up -d").CombinedOutput()
	require.NoError(t, err, string(msg))
	log.Debug(string(msg))
	log.Debug("docker started")
}

func (d *DockerCompose) Down(t *testing.T) {
	msg, err := exec.Command("bash", "-l", "-c", "docker compose down").CombinedOutput()
	require.NoError(t, err, string(msg))
}

func (d *DockerCompose) WaitHealthy(t *testing.T, timeout time.Duration) {
	end := time.Now().Add(timeout)
	for time.Now().Before(end) {
		msg, err := exec.Command("bash", "-l", "-c", `docker compose ps -q \
| xargs -I {} docker inspect \
    --format='{{if .Config.Healthcheck}}{{.Name}} => {{.State.Health.Status}}{{end}}' {}
`).CombinedOutput()
		require.NoError(t, err, string(msg))
		lines := strings.Split(string(msg), `\n`)
		allHealthy := true
		for _, line := range lines {
			if !strings.Contains(line, "healthy") {
				log.Debugf("docker not healthy: %s", line)
				allHealthy = false
			}
		}
		if allHealthy {
			log.Infof("dockers healty OK")
			return
		}

		time.Sleep(1 * time.Second)
	}
}
