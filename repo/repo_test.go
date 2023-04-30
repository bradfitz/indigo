package repo

import (
	"context"
	crand "crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/bluesky-social/indigo/api/bsky"
	cid "github.com/ipfs/go-cid"
	"github.com/ipfs/go-datastore"
	blockstore "github.com/ipfs/go-ipfs-blockstore"
)

func TestRepo(t *testing.T) {
	t.Skip()
	fi, err := os.Open("repo.car")
	if err != nil {
		t.Fatal(err)
	}
	defer fi.Close()

	ctx := context.TODO()
	r, err := ReadRepoFromCar(ctx, fi)
	if err != nil {
		t.Fatal(err)
	}

	if err := r.ForEach(ctx, "app.bsky.feed.post", func(k string, v cid.Cid) error {
		fmt.Println(k, v)
		return nil
	}); err != nil {
		t.Fatal(err)
	}
}

func TestRepoStress(t *testing.T) {
	mkTest := func(adversarial bool) func(*testing.T) {
		return func(t *testing.T) {
			bs := blockstore.NewBlockstore(datastore.NewMapDatastore())
			ctx := context.Background()
			repo := NewRepo(ctx, "did:3:bafyreigv5er7vcxlbikkwedmtd7b3kp7wrcyffep5ogcuxosloxfox5reu", bs)
			val := &bsky.FeedPost{
				Text: strings.Repeat("a", 1<<10),
			}
			lastTime := time.Now()
			for i := 1; i <= 200_000; i++ {
				rpath := genRpath(adversarial)
				if _, err := repo.PutRecord(ctx, rpath, val); err != nil {
					t.Fatal(err)
				}
				if i%10_000 == 0 {
					now := time.Now()
					d := now.Sub(lastTime)
					lastTime = now
					t.Logf("after %v, +%v", i, d)
				}
			}
		}
	}

	t.Run("normal", mkTest(false))
	t.Run("adversarial", mkTest(true))
}

func genRpath(adversarial bool) string {
	const pfx = "foo/"
	buf := make([]byte, len(pfx)+64)
	rbuf := make([]byte, 32)
	for {
		copy(buf, pfx)
		crand.Read(rbuf)
		hex.Encode(buf[len(pfx):], rbuf)
		hv := sha256.Sum256(buf)
		if !adversarial {
			return string(buf)
		}
		if hv[0]&0xc0 != 0x00 {
			return string(buf)
		}
	}
}
