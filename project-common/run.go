package common

import (
	"context"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func Run(r *gin.Engine, srvName string, addr string, stop func()) {
	srv := &http.Server{
		Addr:    addr,
		Handler: r,
	}
	go func() {
		log.Printf("%s web Server running in %s \n", srvName, srv.Addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalln(err)
		}
	}()
	quit := make(chan os.Signal)
	//sigin 用户发送intr字符(CTRL+C)触发kill -2
	//sigterm 结束程序(可以被捕获，阻塞或者忽略)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Printf("Shutting Down %s web Server...\n", srvName)
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	//在此可以同时关闭其他服务
	if stop != nil {
		stop()
	}
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("%s web Server Shutdown,cause by:%s\n", srvName, err)
	}
	select {
	case <-ctx.Done():
		log.Println("wait timeout...")
	}
	log.Printf("%s web Server stop success...\n", srvName)
}
