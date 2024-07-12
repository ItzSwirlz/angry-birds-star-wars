package database

import (
	"github.com/PretendoNetwork/nex-go/v2/types"
)

func UploadCommonData(pid *types.PID, uniqueID *types.PrimitiveU64, commonData *types.Buffer) error {
	_, err := Postgres.Exec(`
		UPDATE common_datas SET common_data = $1, unique_id = $2, owner_pid = $3
	`,
		commonData.Value,
		uniqueID.Value,
		pid.Value(),
	)
	return err
}
