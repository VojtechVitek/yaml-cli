package cli_test

import (
	"testing"
)

func TestGet(t *testing.T) {
	var kubectlGetPod = readFile("_testfiles/kubectl-get-pod.yml")

	tt := []*cliTestCase{
		{
			in:  kubectlGetPod,
			cmd: []string{"yaml", "get", "status.containerStatuses[0].name"},
			out: "goose-metrixdb\n",
		},
		{
			in:  kubectlGetPod,
			cmd: []string{"yaml", "get", "status.containerStatuses[0].state.terminated.finishedAt"},
			out: "\"2019-08-18T12:23:29Z\"\n",
		},
		{
			in:  kubectlGetPod,
			cmd: []string{"yaml", "get", "status.containerStatuses[1].name"},
			out: "linkerd-proxy\n",
		},
		{
			in:  kubectlGetPod,
			cmd: []string{"yaml", "get", "status.containerStatuses[1].state.running.startedAt"},
			out: "\"2019-08-18T12:23:30Z\"\n",
		},
		{
			in:  kubectlGetPod,
			cmd: []string{"yaml", "get", "status.containerStatuses[2].name"},
			err: true,
		},
		{
			in:  kubectlGetPod,
			cmd: []string{"yaml", "get", "status.containerStatuses[2].state.terminated.finishedAt"},
			err: true,
		},
	}

	for _, tc := range tt {
		tc.runTest(t)
	}
}
