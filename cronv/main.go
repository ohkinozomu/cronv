package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/signal"

	"github.com/jessevdk/go-flags"
	"github.com/ohkinozomu/cronv"
	"github.com/ohkinozomu/cronv/server"
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
		os.Exit(0)
	}

	ctx, err := cronv.NewCtx(opts)
	if err != nil {
		panic(err)
	}

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		if _, err := ctx.AppendNewLine(scanner.Text()); err != nil {
			panic(err)
		}
	}
	if err := scanner.Err(); err != nil {
		panic(err)
	}

	err = ctx.Dump()
	if err != nil {
		panic(err)
	}

	log.Printf("[%s] %d tasks.\n", opts.Title, len(ctx.CronEntries))

	go server.Serve()
	log.Println("server start http://localhost:8080")

	open.Run("http://localhost:8080/index.html")

	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<-quit
}
