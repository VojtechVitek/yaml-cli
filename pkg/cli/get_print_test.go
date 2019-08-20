package cli_test

import (
	"strings"
	"testing"
)

func TestGetPrint(t *testing.T) {
	var kubectlGetPod = readFile("_testfiles/kubectl-get-pod.yml")

	// On single input object/document, get and print should behave the same.

	tt := []*cliTestCase{
		{
			in:  kubectlGetPod,
			cmd: []string{"status.containerStatuses[0].name"},
			out: "goose-metrixdb",
		},
		{
			in:  kubectlGetPod,
			cmd: []string{"status.containerStatuses[0].state.terminated.finishedAt"},
			out: "\"2019-08-18T12:23:29Z\"",
		},
		{
			in:  kubectlGetPod,
			cmd: []string{"status.containerStatuses[1].name"},
			out: "linkerd-proxy",
		},
		{
			in:  kubectlGetPod,
			cmd: []string{"status.containerStatuses[1].state.running.startedAt"},
			out: "\"2019-08-18T12:23:30Z\"",
		},
		{
			in:  kubectlGetPod,
			cmd: []string{"status.containerStatuses[2].name"},
			err: true,
		},
		{
			in:  kubectlGetPod,
			cmd: []string{"status.containerStatuses[2].state.terminated.finishedAt"},
			err: true,
		},
	}

	for _, tc := range tt {
		get := &cliTestCase{
			in:  join(tc.in, "\n---\n"),
			cmd: append([]string{"yaml", "get"}, tc.cmd...),
			out: join(tc.out, "\n---\n"), // Get returns multiple YAML objects separated by `---`.
			err: tc.err,
		}
		get.runTest(t)

		print := &cliTestCase{
			in:  join(tc.in, "\n---\n"),
			cmd: append([]string{"yaml", "print"}, tc.cmd...),
			out: join(tc.out, "\n"), // Print returns naked lines. Not necessary a valid YAML.
			err: tc.err,
		}
		print.runTest(t)
	}
}

func join(str string, sep string) string {
	if str == "" {
		return str
	}
	return strings.Join([]string{str, str}, sep) + "\n"
}
