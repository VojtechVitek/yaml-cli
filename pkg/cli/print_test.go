package cli_test

import (
	"testing"
)

func TestPrint(t *testing.T) {
	var kubectlGetPod = readFile("_testfiles/kubectl-get-pod.yml")

	// On single input object/document, get and print should behave the same.

	tt := []*cliTestCase{
		{
			in:   kubectlGetPod,
			args: []string{"print", "status.containerStatuses[0].name"},
			out:  "name: goose-metrixdb\n",
		},
		{
			in:   kubectlGetPod,
			args: []string{"print", "status.containerStatuses[0].state.terminated.finishedAt"},
			out:  "finishedAt: \"2019-08-18T12:23:29Z\"\n",
		},
		{
			in:   kubectlGetPod,
			args: []string{"print", "status.containerStatuses[1].name"},
			out:  "name: linkerd-proxy\n",
		},
		{
			in:   kubectlGetPod,
			args: []string{"print", "status.containerStatuses[1].state.running.startedAt"},
			out:  "startedAt: \"2019-08-18T12:23:30Z\"\n",
		},
		{
			in:   kubectlGetPod,
			args: []string{"print", "status.containerStatuses[2].name"},
			err:  true,
		},
		{
			in:   kubectlGetPod,
			args: []string{"print", "status.containerStatuses[2].state.terminated.finishedAt"},
			err:  true,
		},
	}

	for _, tc := range tt {
		tc.runTest(t)
	}
}
