package database

import (
	"fmt"
	"time"

	"github.com/ItzSwirlz/angry-birds-star-wars/globals"
	"github.com/PretendoNetwork/nex-go/v2/types"
	ranking_types "github.com/PretendoNetwork/nex-protocols-go/v2/ranking/types"
)

func InsertRankingByPIDAndRankingScoreData(pid *types.PID, rankingScoreData *ranking_types.RankingScoreData, uniqueID *types.PrimitiveU64) error {
	globals.Logger.Info(rankingScoreData.FormatToString(1))
	now := time.Now().Unix()

	// If the data does not exist, UDPATE wont work on it.
	var exists bool
	err := Postgres.QueryRow(`SELECT EXISTS(SELECT 1 FROM rankings WHERE owner_pid = $1 category = $2)`, pid.Value(), rankingScoreData.Category.Value).Scan(&exists)
	if err != nil {
		fmt.Printf(err.Error())
		return err
	}
	if !exists {
		// TODO: See if we need all these fields
		_, err := Postgres.Exec(`
			INSERT INTO rankings (owner_pid, unique_id, category, score, order_by, update_mode, groups, param, updated_at)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		`,
			pid.Value(),
			uniqueID.Value,
			rankingScoreData.Category.Value,
			rankingScoreData.Score.Value,
			rankingScoreData.OrderBy.Value,
			rankingScoreData.UpdateMode.Value,
			rankingScoreData.Groups.Value,
			rankingScoreData.Param.Value,
			now,
		)
		return err
	}

	_, err = Postgres.Exec(`
		UPDATE rankings SET score = $1, groups = $2, param = $3, updated_at = $4, order_by = $5, update_mode = $6
		WHERE category = $7 AND owner_pid = $8 AND unique_id = $9
	`,
		rankingScoreData.Score.Value,
		rankingScoreData.Groups.Value,
		rankingScoreData.Param.Value,
		now,
		rankingScoreData.OrderBy.Value,
		rankingScoreData.UpdateMode.Value,
		rankingScoreData.Category.Value,
		pid.Value(),
		uniqueID.Value,
	)

	return err
}
