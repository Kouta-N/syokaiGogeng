package main

import (
	"context"
	"fmt"
	"io"
	"net"
	"net/http"
	"testing"

	"golang.org/x/sync/errgroup"
)

func TestRun(t *testing.T){
	l,err := net.Listen("tcp","localhost:0")
	if err != nil{
		t.Fatalf("failed to listen port %v", err)
	}
	ctx, cancel := context.WithCancel(context.Background())
	eg,ctx := errgroup.WithContext(ctx)
	eg.Go(func() error {
		return run(ctx, l)
	})
	in := "message"
	url := fmt.Sprintf("http://%s/%s", l.Addr().String(), in)
	// http.Getでresponseが帰ってくる(defer rsp.Body.Close() を忘れない)
	t.Logf("try yo %q", url)
	rsp,err := http.Get(url)
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