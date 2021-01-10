package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/signal"

	"github.com/jessevdk/go-flags"
	"github.com/ohkinozomu/cronv"
	"github.com/ohkinozomu/cronv/pkg/server"
	"github.com/skratchdot/open-golang/open"
)

const (
	version = "0.4.2"
	name    = "Cronv"
)

func main() {
	opts := cronv.NewCronvCommand()

	parser := flags.NewParser(opts, flags.Default)
	parser.Name = fmt.Sprintf("%s v%s", name, version)
	if _, err := parser.Parse(); err != nil {
		fmt.Println("ここ")
		fmt.Println(err)
		os.Exit(1)
	}

	ctx, err := cronv.NewCtx(opts)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		if _, err := ctx.AppendNewLine(scanner.Text()); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	err = ctx.Dump()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	log.Printf("[%s] %d tasks.\n", opts.Title, len(ctx.CronEntries))

	go server.Serve(ctx)
	log.Println("server start http://localhost:8080")

	open.Run("http://localhost:8080/index.html")

	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<-quit
}
