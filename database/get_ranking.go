package database

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/ItzSwirlz/angry-birds-star-wars/globals"
	"github.com/PretendoNetwork/nex-go/v2/types"
	ranking_types "github.com/PretendoNetwork/nex-protocols-go/v2/ranking/types"
	"github.com/lib/pq"
)

func parseRankingDataList(rows *sql.Rows) (*types.List[*ranking_types.RankingRankData], error) {
	results := types.NewList[*ranking_types.RankingRankData]()

	for rows.Next() {
		result := ranking_types.NewRankingRankData()
		var updateDate uint64
		var userPid uint64
		var updateMode uint32
		var idk byte

		err := rows.Scan(
			&userPid,
			&result.UniqueID.Value,
			&result.Category.Value,
			&result.Score.Value,
			&result.Order.Value,
			&updateMode,
			&result.Groups.Value,
			&idk,
			&updateDate,
		)
		if err != nil {
			fmt.Printf("err line 33")
			return nil, err
		}

		result.PrincipalID = types.NewPID(userPid)
		result.UpdateTime.FromTimestamp(time.Now()) // todo: fix
		var buf []byte
		err = Postgres.QueryRow(`SELECT "common_data" FROM "common_datas" WHERE "owner_pid" = $1`, userPid).Scan(&buf)
		if err != nil {
			fmt.Printf("couldnt find common data")
			return nil, err
		}
		result.CommonData.Value = buf
		results.Append(result)
	}

	return results, nil
}

// FIXME: do later
func GetRankingsAndCountByCategoryAndRankingOrderParam(category *types.PrimitiveU32, rankingOrderParam *ranking_types.RankingOrderParam) (*types.List[*ranking_types.RankingRankData], uint32, error) {
	var pqErr *pq.Error
	globals.Logger.Info(rankingOrderParam.FormatToString(1))
	// rankingTable := `ranking.ranks_` + strconv.Itoa(int(category.Value))
	rows, err := Postgres.Query(`
		SELECT * FROM "rankings"
	`,
	)
	if errors.Is(err, sql.ErrNoRows) || (errors.As(err, &pqErr) && pqErr.SQLState() == "42P01") {
		fmt.Printf("no rows")
	}
	if err != nil {
		fmt.Printf("err line 63")
		return nil, 0, err
	}

	results, err := parseRankingDataList(rows)
	if err != nil {
		fmt.Printf("err line 60")
		return nil, 0, err
	}

	globals.Logger.Info(results.String())

	return results, uint32(results.Length()), nil
}

func InsertRankingByPIDAndRankingScoreData(pid *types.PID, rankingScoreData *ranking_types.RankingScoreData, uniqueID *types.PrimitiveU64) error {
	globals.Logger.Info(rankingScoreData.FormatToString(1))
	globals.Logger.Info(uniqueID.String())
	now := time.Now().Unix()

	var exists bool
	e := Postgres.QueryRow(`SELECT EXISTS(SELECT 1 FROM rankings WHERE owner_pid = $1)`, pid.Value()).Scan(&exists)
	if e != nil {
		return e
	}
	if !exists {
		_, err := Postgres.Exec(`
			INSERT INTO rankings (score, groups, param, updated_at, category, owner_pid, unique_id, order_by, update_mode)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		`,
			rankingScoreData.Score.Value,
			rankingScoreData.Groups.Value,
			rankingScoreData.Param.Value,
			now,
			rankingScoreData.Category.Value,
			pid.Value(),
			uniqueID.Value,
			rankingScoreData.OrderBy.Value,
			rankingScoreData.UpdateMode.Value,
		)
		return err
	}

	_, err := Postgres.Exec(`
		UPDATE rankings SET score = $1, groups = $2, param = $3, updated_at = $4, order_by = $8, update_mode = $9
		WHERE category = $5 AND owner_pid = $6 AND unique_id = $7
	`,
		rankingScoreData.Score.Value,
		rankingScoreData.Groups.Value,
		rankingScoreData.Param.Value,
		now,
		rankingScoreData.Category.Value,
		pid.Value(),
		uniqueID.Value,
		rankingScoreData.OrderBy.Value,
		rankingScoreData.UpdateMode.Value,
	)

	return err
}
