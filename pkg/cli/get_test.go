package cli_test

import (
	"bytes"
	"os"
	"testing"

	"github.com/VojtechVitek/yaml/pkg/cli"
	"github.com/google/go-cmp/cmp"
)

func TestGet(t *testing.T) {
	tt := []struct {
		in  *os.File
		cmd []string
		out string
		err bool
	}{
		{
			in:  openFile("_testfiles/kubectl-get-pod.yml"),
			cmd: []string{"yaml", "get", "status.containerStatuses[0].name"},
			out: "goose-metrixdb\n",
		},
		{
			in:  openFile("_testfiles/kubectl-get-pod.yml"),
			cmd: []string{"yaml", "get", "status.containerStatuses[0].state.terminated.finishedAt"},
			out: "\"2019-08-18T12:23:29Z\"\n",
		},
		{
			in:  openFile("_testfiles/kubectl-get-pod.yml"),
			cmd: []string{"yaml", "get", "status.containerStatuses[1].name"},
			out: "linkerd-proxy\n",
		},
		{
			in:  openFile("_testfiles/kubectl-get-pod.yml"),
			cmd: []string{"yaml", "get", "status.containerStatuses[1].state.running.startedAt"},
			out: "\"2019-08-18T12:23:30Z\"\n",
		},
		{
			in:  openFile("_testfiles/kubectl-get-pod.yml"),
			cmd: []string{"yaml", "get", "status.containerStatuses[2].name"},
			err: true,
		},
		{
			in:  openFile("_testfiles/kubectl-get-pod.yml"),
			cmd: []string{"yaml", "get", "status.containerStatuses[2].state.terminated.finishedAt"},
			err: true,
		},
	}

	for i, tc := range tt {
		var b bytes.Buffer

		err := cli.Run(&b, tc.in, tc.cmd)
		if err != nil && !tc.err {
			t.Errorf("tc[%v]: %v", i, err)
		} else if err == nil && tc.err {
			t.Errorf("tc[%v]: expected error", i)
		}

		if diff := cmp.Diff(tc.out, b.String()); diff != "" {
			t.Errorf("tc[%v] mismatch (-want +got):\n%s", i, diff)
		}
	}
}
