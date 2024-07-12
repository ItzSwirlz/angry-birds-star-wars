package database

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/ItzSwirlz/angry-birds-star-wars/globals"
	"github.com/PretendoNetwork/nex-go/v2/types"
	ranking_types "github.com/PretendoNetwork/nex-protocols-go/v2/ranking/types"
	"github.com/lib/pq"
)

func GetRankingsAndCountByCategoryAndRankingOrderParam(category *types.PrimitiveU32, rankingOrderParam *ranking_types.RankingOrderParam) (*types.List[*ranking_types.RankingRankData], uint32, error) {
	var pqErr *pq.Error
	globals.Logger.Info(rankingOrderParam.FormatToString(1))

	rows, err := Postgres.Query(`
		SELECT * FROM "rankings"
	`,
	)
	if errors.Is(err, sql.ErrNoRows) || (errors.As(err, &pqErr) && pqErr.SQLState() == "42P01") {
		fmt.Printf("No rows")
		return nil, 0, err
	}
	if err != nil {
		fmt.Printf("Error selecting rankings data.")
		return nil, 0, err
	}

	results, err := parseRankingDataList(rows)
	if err != nil {
		fmt.Printf("Error parsing ranking data")
		return nil, 0, err
	}

	globals.Logger.Info(results.String())

	return results, uint32(results.Length()), nil
}
