package slur

import (
	"context"
	"encoding/csv"
	"errors"
	"os"
	"slices"
	"strconv"
)

var ErrExport = errors.New("failed to export")

type MessageSource interface {
	Next(ctx context.Context, count uint64) ([]Message, error)
}

type Message interface {
	UserID() string
	MessageID() int64
	Text() string
}

type SlurCounts struct {
	Common string
	Total  int
}

type PlayerCounts struct {
	SteamID string
	Total   int
	Rank    int
}

func Calc(ctx context.Context, source MessageSource) ([]PlayerCounts, []int64, error) {
	var (
		steamIDTotals     = make(map[string]int)
		matchTotals       = make(map[string]int)
		offset            = uint64(0)
		limit             = uint64(10_000_000)
		flaggedMessageIDs []int64
		total             = 0
		playerResults     []PlayerCounts
		slurResults       []SlurCounts //nolint:prealloc
	)

	for {
		results, err := source.Next(ctx, limit)
		if err != nil || len(results) == 0 {
			break
		}

		for _, result := range results {
			if match, found := Check(result.Text()); found {
				total++

				flaggedMessageIDs = append(flaggedMessageIDs, result.MessageID())
				if _, ok := steamIDTotals[result.UserID()]; !ok {
					steamIDTotals[result.UserID()] = 1
				} else {
					steamIDTotals[result.UserID()]++
				}

				if _, ok := matchTotals[match.Common]; !ok {
					matchTotals[match.Common] = 1
				} else {
					matchTotals[match.Common]++
				}
			}
		}

		offset += limit
	}

	for text, total := range matchTotals {
		slurResults = append(slurResults, SlurCounts{Common: text, Total: total})
	}

	slices.SortStableFunc(slurResults, func(a, b SlurCounts) int {
		if a.Total > b.Total {
			return -1
		} else if a.Total < b.Total {
			return 1
		}

		return 0
	})

	//nolint:prealloc
	for sid, total := range steamIDTotals {
		playerResults = append(playerResults, PlayerCounts{SteamID: sid, Total: total})
	}

	slices.SortStableFunc(playerResults, func(a, b PlayerCounts) int {
		if a.Total > b.Total {
			return -1
		} else if a.Total < b.Total {
			return 1
		}

		return 0
	})

	if errPlayers := savePlayerSlursCSV(playerResults, "players.csv"); errPlayers != nil {
		return nil, nil, errPlayers
	}

	if errSlurs := saveSlursCSV(slurResults, "slurs.csv"); errSlurs != nil {
		return nil, nil, errSlurs
	}

	return playerResults, nil, nil
}

func savePlayerSlursCSV(slurs []PlayerCounts, path string) error {
	file, err := os.Create(path)
	if err != nil {
		return errors.Join(err, ErrExport)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	for _, slur := range slurs {
		err := writer.Write([]string{slur.SteamID, strconv.Itoa(slur.Total)})
		if err != nil {
			return errors.Join(err, ErrExport)
		}
	}

	return nil
}

func saveSlursCSV(slurs []SlurCounts, path string) error {
	file, err := os.Create(path)
	if err != nil {
		return errors.Join(err, ErrExport)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	for _, slur := range slurs {
		err := writer.Write([]string{slur.Common, strconv.Itoa(slur.Total)})
		if err != nil {
			return errors.Join(err, ErrExport)
		}
	}

	return nil
}
