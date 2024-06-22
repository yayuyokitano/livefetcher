package util

import (
	"io"
	"os/exec"
)

func spawnMecab() (cmd *exec.Cmd, stdin io.WriteCloser, stdout io.ReadCloser, err error) {
	cmd = exec.Command("mecab", `--node-format=%pS%f[7]`, `--unk-format=%M`, `--eos-format=`)

	stdin, err = cmd.StdinPipe()
	if err != nil {
		return
	}

	stdout, err = cmd.StdoutPipe()
	if err != nil {
		return
	}
	err = cmd.Start()
	return
}

func Mecab(s string) (kana string, err error) {
	cmd, stdin, stdout, err := spawnMecab()
	if err != nil {
		return
	}

	_, err = stdin.Write([]byte(s + "\n"))
	if err != nil {
		return
	}
	stdin.Close()

	out, err := io.ReadAll(stdout)
	if err != nil {
		return
	}

	err = cmd.Wait()
	if err != nil {
		return
	}
	kana = string(out)
	return
}
