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
	t.Helper()
	log.Debug("starting docker")
	runCommand(t, true, "docker compose down")
	log.Debug("docker started")
}

func (d *DockerCompose) Down(t *testing.T) {
	t.Helper()
	runCommand(t, true, "docker compose down")
	log.Debug("docker down")
}

func (d *DockerCompose) WaitHealthy(t *testing.T, timeout time.Duration) {
	t.Helper()
	end := time.Now().Add(timeout)
	for time.Now().Before(end) {
		msg, err := exec.Command(bash, "-l", "-c", `docker compose ps -q \
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
