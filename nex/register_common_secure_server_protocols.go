package nex

import (
	database "github.com/ItzSwirlz/angry-birds-star-wars/database"
	globals "github.com/ItzSwirlz/angry-birds-star-wars/globals"
	common_ranking "github.com/PretendoNetwork/nex-protocols-common-go/v2/ranking"
	common_secure "github.com/PretendoNetwork/nex-protocols-common-go/v2/secure-connection"
	ranking "github.com/PretendoNetwork/nex-protocols-go/v2/ranking"
	secure "github.com/PretendoNetwork/nex-protocols-go/v2/secure-connection"
)

func registerCommonSecureServerProtocols() {
	secureProtocol := secure.NewProtocol()
	globals.SecureEndpoint.RegisterServiceProtocol(secureProtocol)
	common_secure.NewCommonProtocol(secureProtocol)

	rankingProtocol := ranking.NewProtocol()
	globals.SecureEndpoint.RegisterServiceProtocol(rankingProtocol)
	ranking_protocol := common_ranking.NewCommonProtocol(rankingProtocol)
	ranking_protocol.UploadCommonData = database.UploadCommonData
	// TODO: UploadScore
}
