package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"

	"golang.org/x/sync/errgroup"
)

func run(ctx context.Context, l net.Listener) error {
	// サーバーの組み立て
	s := &http.Server{
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request){
			fmt.Fprintf(w, "Hello, %s!", r.URL.Path[1:])
		}),
	}
	// syncにはerrgroupサブパッケージが含まれており、これを使用するとgo routine並行処理が簡単に実装できる https://pkg.go.dev/golang.org/x/sync/errgroup#Group.Go
	// 派生したContextは、Goに渡された関数が最初にnilでないエラーを返すか、Waitが最初に返すか、どちらか先に発生したときにキャンセルされる
	eg,ctx := errgroup.WithContext(ctx)
	eg.Go(func() error{
		if err := s.Serve(l); err != nil && //サーバー起動してエラーチェック
		err != http.ErrServerClosed {
			log.Printf("failed to close: %+v",err)
			return err
		}
		return nil
	})

	// eg.Goでエラーが発生すると、Doneが入る
	<-ctx.Done()
	if err := s.Shutdown(context.Background()); err != nil {
		log.Printf("failed to shutdown: %+v",err)
	}
	// 別ゴルーチンの終了を待つ
	return eg.Wait()
}

func main(){
	if len(os.Args) != 2 {
		log.Printf("need port num\n")
		os.Exit(1)
	}
	p := os.Args[1]
	l,err := net.Listen("tcp", ":"+p)
	if err != nil{
		log.Fatalf("failed to listen port %s : %v",p,err)
	}
	if err := run(context.Background(),l); err != nil {
		log.Printf("failed to terminate server: %v",err)
		os.Exit(1)
	}
}