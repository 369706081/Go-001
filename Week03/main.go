package main

import (
	"context"
	"fmt"
	"github.com/pkg/errors"
	"golang.org/x/sync/errgroup"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
)
func test(w http.ResponseWriter, r *http.Request) {

	fmt.Fprintf(w, "Hello	word test!")
}
func main(){
	addr1 := "0.0.0.0:9090"
	addr2 := "0.0.0.0:9091"
	signalChan := make(chan os.Signal)
	signal.Notify(signalChan, os.Interrupt)
	http.HandleFunc("/test", test)
	server := &http.Server{
		Addr:         addr1,
		Handler:      nil,

	}
	server2 := &http.Server{
		Addr:         addr2,
		Handler:      nil,

	}
	//errgro
	ctx, cancel := context.WithCancel(context.Background())
	group, _ := errgroup.WithContext(ctx)

	group.Go(func() error {

		err := server.ListenAndServe()
		if err != nil{
			return err
		}
		return nil

	})
	group.Go(func() error {

		err := server2.ListenAndServe()
		if err != nil{
			return err
		}
		return nil

	})
	group.Go(func() error {
		for{
			select {
			case sig := <-signalChan:
				log.Println("Get Signal:", sig)
				log.Println("Shutdown Server ...")
				cancel()

			case <-ctx.Done():

				log.Println("context.Canceled ...")
				ctx2, cancel := context.WithTimeout(context.Background(), 5*time.Second)
				defer cancel()
				log.Println("start stop Server")
				if err := server.Shutdown(ctx2); err != nil {
					errors.Wrap(err,"Shutdown Server :")
					return  err
				}
				log.Println(" stop Serve sucess")
				log.Println("start stop Server2")
				if err := server2.Shutdown(ctx2); err != nil {
					errors.Wrap(err,"Shutdown Server2 :")
					return  err
				}
				log.Println(" stop Serve sucess")
				return context.Canceled
			}
		}




		return nil
	})
	err :=group.Wait()
	if err != nil {
		log.Println()
	}

    defer cancel()
}
