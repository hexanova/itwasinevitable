package main

import (
	"context"
	"log"
	"time"

	atp "github.com/bluesky-social/indigo/api/atproto"
	bsky "github.com/bluesky-social/indigo/api/bsky"
	lexutil "github.com/bluesky-social/indigo/lex/util"
	xrpc "github.com/bluesky-social/indigo/xrpc"
	"github.com/urfave/cli/v2"
)

type Config struct {
	BlueskyURL      string `env:"BLUESKY_URL,required"`
	BlueskyUsername string `env:"BLUESKY_USERNAME,required"`
	BlueskyPassword string `env:"BLUESKY_PASSWORD,required"`
}

func postToBluesky(cctx *cli.Context, xrpcc *xrpc.Client, message string, cfg *Config) {

	for attempts := 0; attempts < 5; attempts++ {
		post := &bsky.FeedPost{
			CreatedAt: time.Now().Local().Format(time.RFC3339),
			Text:      message,
		}
		resp, err := atp.RepoCreateRecord(context.TODO(), xrpcc, &atp.RepoCreateRecord_Input{
			Collection: "app.bsky.feed.post",
			Repo:       xrpcc.Auth.Did,
			Record: &lexutil.LexiconTypeDecoder{
				Val: post,
			},
		})
		if err != nil {
			log.Fatalf("Error posting to Bluesky: %v\n, Response: %v\n", err, resp)
			log.Println("Attempt", attempts+1, "of 5.")
		}
		log.Println("Giving up on post:", message)
	}
}

func authenticateSession(xrpcc *xrpc.Client, cfg *Config) error {
	ses, err := atp.ServerCreateSession(context.TODO(), xrpcc, &atp.ServerCreateSession_Input{
		Identifier: cfg.BlueskyUsername,
		Password:   cfg.BlueskyPassword,
	})
	xrpcc.Auth = &xrpc.AuthInfo{
		AccessJwt:  ses.AccessJwt,
		RefreshJwt: ses.RefreshJwt,
		Handle:     ses.Handle,
		Did:        ses.Did,
	}
	if err != nil {
		log.Fatal("Error creating session: ", err)
	}
	return nil
}
