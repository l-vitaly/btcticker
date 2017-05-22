package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"os/signal"
	"runtime"
	"syscall"

	"github.com/BurntSushi/toml"
	"github.com/l-vitaly/btcticker"
	_ "github.com/l-vitaly/btcticker/fetcher/btce"
	_ "github.com/l-vitaly/btcticker/fetcher/coindesk"
	_ "github.com/l-vitaly/btcticker/fetcher/spectrocoin"
)

var (
	configPath = flag.String("c", "", "-c=/path/to/config")
)

func init() {
	flag.Parse()
	runtime.GOMAXPROCS(runtime.NumCPU())
}

func main() {
	var conf btcticker.Config

	if *configPath == "" {
		flag.Usage()
		os.Exit(1)
	}

	data, err := ioutil.ReadFile(*configPath)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	_, err = toml.Decode(string(data), &conf)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	errCh := make(chan error)

	t := btcticker.NewBtcTicker(conf)

	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		errCh <- fmt.Errorf("%s", <-c)
		t.Stop()
	}()

	t.Start()

	fmt.Println(<-errCh)
}
