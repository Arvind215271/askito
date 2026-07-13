package python

import (
	"encoding/binary"
	"encoding/json"
	"io"
	"os/exec"
)

type PythonWorker struct {
	cmd    *exec.Cmd
	stdin  io.WriteCloser
	stdout io.ReadCloser
	stderr io.ReadCloser
}

func NewWorker(scriptPath string) (*PythonWorker, error) {
	cmd := exec.Command("python3", "-u", scriptPath)
	stdin, err := cmd.StdinPipe()
	if err != nil {
		return nil, err
	}
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, err
	}
	stderr, err := cmd.StderrPipe()
	if err != nil {
		return nil, err
	}

	if err := cmd.Start(); err != nil {
		return nil, err
	}

	return &PythonWorker{
		cmd:    cmd,
		stdin:  stdin,
		stdout: stdout,
		stderr: stderr,
	}, nil
}

func (w *PythonWorker) SendCommand(req any) error {
	reqData, err := json.Marshal(req)
	if err != nil {
		return err
	}

	if err := binary.Write(w.stdin, binary.BigEndian, uint32(len(reqData))); err != nil {
		return err
	}
	if _, err := w.stdin.Write(reqData); err != nil {
		return err
	}
	return nil
}

func (w *PythonWorker) ReceiveResponse() ([]byte, error) {
	var length uint32
	if err := binary.Read(w.stdout, binary.BigEndian, &length); err != nil {
		return nil, err
	}

	respData := make([]byte, length)
	if _, err := io.ReadFull(w.stdout, respData); err != nil {
		return nil, err
	}
	return respData, nil
}

func (w *PythonWorker) Close() error {
	w.stdin.Close()
	return w.cmd.Wait()
}
