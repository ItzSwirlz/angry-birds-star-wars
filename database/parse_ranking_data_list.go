package database

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/PretendoNetwork/nex-go/v2/types"
	ranking_types "github.com/PretendoNetwork/nex-protocols-go/v2/ranking/types"
)

func parseRankingDataList(rows *sql.Rows) (*types.List[*ranking_types.RankingRankData], error) {
	results := types.NewList[*ranking_types.RankingRankData]()

	for rows.Next() {
		result := ranking_types.NewRankingRankData()
		var updateDate int64
		var userPid uint64
		var updateMode uint32

		err := rows.Scan(
			&userPid,
			&result.UniqueID.Value,
			&result.Category.Value,
			&result.Score.Value,
			&result.Order.Value,
			&updateMode,
			&result.Groups.Value,
			&result.Param.Value,
			&updateDate,
		)
		if err != nil {
			fmt.Printf("Rrror scanning the row")
			return nil, err
		}

		result.PrincipalID = types.NewPID(userPid)
		result.UpdateTime.FromTimestamp(time.Unix(updateDate, 0))

		// Grab our common data otherwise the game will crash
		var commonData []byte
		err = Postgres.QueryRow(`SELECT "common_data" FROM "common_datas" WHERE "owner_pid" = $1`, userPid).Scan(&commonData)
		if err != nil {
			fmt.Printf("couldnt find common data")
			return nil, err
		}

		result.CommonData.Value = commonData
		results.Append(result)
	}

	return results, nil
}
