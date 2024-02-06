package main

import (
	"bufio"
	"context"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/bluesky-social/indigo/util/cliutil"
	xrpc "github.com/bluesky-social/indigo/xrpc"
	"github.com/caarlos0/env"
	"github.com/urfave/cli/v2"
)

var cctx *cli.Context
var xrpcc *xrpc.Client

func main() {
	parseFlags()

	cleanup()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	cfg, err := readConfig()
	if err != nil {
		log.Fatal("Error reading .env variables:", err)
	}

	xrpcc, err := cliutil.GetXrpcClient(cctx, false)
	if err != nil {
		log.Fatal("Error getting XRPC client:", err)
	}
	xrpcc.Host = cfg.BlueskyURL

	err = authenticateSession(xrpcc, cfg)
	if err != nil {
		log.Fatal("Error authenticating:", err)
	}

	go func() {
		sigch := make(chan os.Signal, 1)
		defer signal.Stop(sigch)

		signal.Notify(sigch, syscall.SIGINT, syscall.SIGTERM)

		select {
		case sig := <-sigch:
			log.Println(sig, "- cleaning up")
		case <-ctx.Done():
		}

		cancel()
	}()

	buffer := &dataBuffer{
		recent: make([]string, minLinesBeforeDuplicate),
		fuzzy:  make([][]string, fuzzyDuplicateWindow),
		queue:  make([]string, 0, maxQueuedLines),
	}

	ch := make(chan string)
	go dwarfFortress(ctx, buffer, ch)

	initialDelay := nextDelay()
	log.Println("Waiting", initialDelay, "before making first post.")
	time.Sleep(initialDelay)

	for {
		var line string
		select {
		case <-ctx.Done():
			return
		case line = <-ch:
		default:
			log.Println("Warning: no posts are ready")
			select {
			case <-ctx.Done():
				return
			case line = <-ch:
			}
		}
		authenticateSession(xrpcc, cfg)
		postToBluesky(cctx, xrpcc, line, cfg)
		time.Sleep(nextDelay())
	}
}

func nextDelay() time.Duration {
	return postInterval - time.Duration(time.Now().UnixNano())%postInterval
}

func loadEnvFile(filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		log.Fatal("Error reading .env file. ", err)
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.SplitN(line, "=", 2)
		if len(parts) == 2 {
			key := parts[0]
			value := parts[1]
			os.Setenv(key, value)
		}
	}

	if err := scanner.Err(); err != nil {
		return err
	}

	return nil
}

func readConfig() (*Config, error) {
	err := loadEnvFile(".env")
	if err != nil {
		log.Fatal("No .env file found. ", err)
		return nil, err
	}

	config := &Config{}
	err = env.Parse(config)
	if err != nil {
		return nil, err
	}

	return config, nil
}
