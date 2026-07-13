package python

import (
	"context"
	"encoding/json"
	"fmt"
	"path/filepath"
	"runtime"
)

type SingleClient struct {
	worker *PythonWorker
}

func NewSingleWorker() (*SingleClient, error) {
	_, file, _, _ := runtime.Caller(0)
	script := filepath.Join(filepath.Dir(file), "python_worker_single.py")
	worker, err := NewWorker(script)
	if err != nil {
		return nil, err
	}
	return &SingleClient{worker: worker}, nil
}

func (c *SingleClient) GetVideo(ctx context.Context, id string) (map[string]any, error) {
	if err := c.worker.SendCommand(map[string]string{"id": id}); err != nil {
		return nil, err
	}
	respData, err := c.worker.ReceiveResponse()
	if err != nil {
		return nil, err
	}
	var resp struct {
		Ok   bool           `json:"ok"`
		Data map[string]any `json:"data"`
		Err  string         `json:"error"`
	}
	if err := json.Unmarshal(respData, &resp); err != nil {
		return nil, err
	}
	if !resp.Ok {
		return nil, fmt.Errorf("worker error: %s", resp.Err)
	}
	return resp.Data, nil
}

func (c *SingleClient) WarmUp(ctx context.Context) error {
	if err := c.worker.SendCommand(map[string]string{"cmd": "warmup"}); err != nil {
		return err
	}
	respData, err := c.worker.ReceiveResponse()
	if err != nil {
		return err
	}
	var resp struct {
		Ok  bool   `json:"ok"`
		Err string `json:"error"`
	}
	if err := json.Unmarshal(respData, &resp); err != nil {
		return err
	}
	if !resp.Ok {
		return fmt.Errorf("worker warmup error: %s", resp.Err)
	}
	return nil
}

func (c *SingleClient) Close() error {
	return c.worker.Close()
}

type BatchClient struct {
	worker *PythonWorker
}

func NewBatchWorker() (*BatchClient, error) {
	_, file, _, _ := runtime.Caller(0)
	script := filepath.Join(filepath.Dir(file), "python_worker_batch.py")
	worker, err := NewWorker(script)
	if err != nil {
		return nil, err
	}
	return &BatchClient{worker: worker}, nil
}

func (c *BatchClient) GetVideos(ctx context.Context, ids []string) ([]map[string]any, error) {
	if err := c.worker.SendCommand(map[string][]string{"ids": ids}); err != nil {
		return nil, err
	}
	respData, err := c.worker.ReceiveResponse()
	if err != nil {
		return nil, err
	}
	var results []struct {
		Ok   bool           `json:"ok"`
		Data map[string]any `json:"data"`
		Err  string         `json:"error"`
	}
	if err := json.Unmarshal(respData, &results); err != nil {
		return nil, err
	}
	
	final := make([]map[string]any, len(results))
	for i, r := range results {
		if !r.Ok {
			return nil, fmt.Errorf("batch item %d error: %s", i, r.Err)
		}
		final[i] = r.Data
	}
	return final, nil
}

func (c *BatchClient) Close() error {
	return c.worker.Close()
}
