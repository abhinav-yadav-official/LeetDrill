// Command ingest pulls the LeetCode problem catalog into the local DB.
//
// Walks problemsetQuestionList in pages, upserts into problems, patterns,
// problem_patterns. Idempotent — re-run anytime to refresh.
package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"leetdrill/internal/leetcode"
	"leetdrill/internal/models"
	"leetdrill/internal/store"
)

func main() {
	pageSize := flag.Int("page", 50, "page size")
	maxPages := flag.Int("max-pages", 0, "0 = walk to end")
	pause := flag.Duration("pause", 1100*time.Millisecond, "delay between requests; keep >= 1s")
	dryRun := flag.Bool("dry-run", false, "fetch but don't write to DB")
	flag.Parse()

	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" && !*dryRun {
		log.Fatal("DATABASE_URL not set (or pass -dry-run)")
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go func() {
		ch := make(chan os.Signal, 1)
		signal.Notify(ch, os.Interrupt, syscall.SIGTERM)
		<-ch
		log.Println("interrupt: shutting down")
		cancel()
	}()

	var st *store.Store
	if !*dryRun {
		var err error
		st, err = store.Open(ctx, dsn)
		if err != nil {
			log.Fatalf("open store: %v", err)
		}
		defer st.Close()
	}

	c := leetcode.New()
	skip := 0
	page := 0
	totalSeen := 0
	upserts := 0

	for {
		items, total, err := c.ListProblems(ctx, skip, *pageSize)
		if err != nil {
			log.Fatalf("ListProblems(skip=%d): %v", skip, err)
		}
		log.Printf("page=%d skip=%d got=%d total=%d", page, skip, len(items), total)
		if len(items) == 0 {
			break
		}

		if !*dryRun {
			if err := persistPage(ctx, st, items); err != nil {
				log.Fatalf("persist page %d: %v", page, err)
			}
			upserts += len(items)
		}

		totalSeen += len(items)
		if skip+len(items) >= total {
			break
		}
		skip += len(items)
		page++
		if *maxPages > 0 && page >= *maxPages {
			log.Printf("max-pages=%d reached", *maxPages)
			break
		}
		select {
		case <-ctx.Done():
			log.Println("context canceled")
			return
		case <-time.After(*pause):
		}
	}
	log.Printf("ingest done: seen=%d upserts=%d", totalSeen, upserts)
}

func persistPage(ctx context.Context, st *store.Store, items []leetcode.ProblemListItem) error {
	for _, it := range items {
		tags := make([]models.Tag, 0, len(it.TopicTags))
		for _, t := range it.TopicTags {
			tags = append(tags, models.Tag{Name: t.Name, Slug: t.Slug})
		}
		p := models.Problem{
			LeetcodeSlug:       it.TitleSlug,
			LeetcodeQuestionID: it.QuestionID,
			LeetcodeFrontendID: it.QuestionFrontendID,
			Title:              it.Title,
			Difficulty:         models.Difficulty(it.Difficulty),
			TopicTags:          tags,
			ACRate:             it.ACRate,
			PaidOnly:           it.IsPaidOnly,
		}
		if p.Difficulty == "" {
			p.Difficulty = models.DifficultyMedium // shouldn't happen, but be defensive
		}

		problemID, err := store.UpsertProblem(ctx, st.DB(), p)
		if err != nil {
			return err
		}

		for _, t := range it.TopicTags {
			patID, err := store.UpsertPattern(ctx, st.DB(), t.Slug, t.Name)
			if err != nil {
				return err
			}
			if err := store.LinkProblemPattern(ctx, st.DB(), problemID, patID); err != nil {
				return err
			}
		}
	}
	return nil
}

