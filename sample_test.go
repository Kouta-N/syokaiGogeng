package main

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"testing"

	"golang.org/x/sync/errgroup"
)

func TestRun(t *testing.T){
	ctx, cancel := context.WithCancel(context.Background())
	eg,ctx := errgroup.WithContext(ctx)
	eg.Go(func() error {
		return run(ctx)
	})
	in := "message"
	// http.Getでresponseが帰ってくる(defer rsp.Body.Close() を忘れない)
	rsp,err := http.Get("http://localhost:18080/" + in)
	if err != nil {
		t.Error("failed to get:", err)
	}
	// rsp bodyをio.ReadAllで読み込むことが可能
	got, err := io.ReadAll(rsp.Body)
	defer rsp.Body.Close()
	if err != nil {
		t.Fatalf("failed to read body: %v", err)
	}

	want := fmt.Sprintf("Hello, %s!", in)
	if string(got) != want {
		t.Error("want",want,"but got is ",string(got))
	}
	cancel()
	if err := eg.Wait(); err != nil {
		t.Fatal(err)
	}
}