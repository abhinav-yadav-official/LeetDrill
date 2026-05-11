// Command ingest pulls the LeetCode problem catalog into the local DB.
//
// Phase 0 stub: walks pages and prints counts. Phase 1 wires DB upserts.
package main

import (
	"context"
	"flag"
	"log"
	"time"

	"leetdrill/internal/leetcode"
)

func main() {
	pageSize := flag.Int("page", 50, "page size")
	maxPages := flag.Int("max-pages", 0, "0 = all")
	pause := flag.Duration("pause", 1100*time.Millisecond, "delay between requests")
	flag.Parse()

	ctx := context.Background()
	c := leetcode.New()

	skip := 0
	page := 0
	for {
		items, total, err := c.ListProblems(ctx, skip, *pageSize)
		if err != nil {
			log.Fatalf("ListProblems(skip=%d): %v", skip, err)
		}
		log.Printf("page=%d skip=%d got=%d total=%d", page, skip, len(items), total)
		if len(items) == 0 || skip+len(items) >= total {
			break
		}
		skip += len(items)
		page++
		if *maxPages > 0 && page >= *maxPages {
			break
		}
		time.Sleep(*pause)
	}
	log.Println("ingest stub done")
}
