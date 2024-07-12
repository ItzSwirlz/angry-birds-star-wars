package nex

import (
	"fmt"
	"os"
	"strconv"

	"github.com/ItzSwirlz/angry-birds-star-wars/globals"
	nex "github.com/PretendoNetwork/nex-go/v2"
)

func StartSecureServer() {
	globals.SecureServer = nex.NewPRUDPServer()

	globals.SecureEndpoint = nex.NewPRUDPEndPoint(1)
	globals.SecureEndpoint.IsSecureEndPoint = true
	globals.SecureEndpoint.ServerAccount = globals.SecureServerAccount
	globals.SecureEndpoint.AccountDetailsByPID = globals.AccountDetailsByPID
	globals.SecureEndpoint.AccountDetailsByUsername = globals.AccountDetailsByUsername
	globals.SecureServer.BindPRUDPEndPoint(globals.SecureEndpoint)
	// globals.SecureServer.ByteStreamSettings.UseStructureHeader = true

	globals.SecureServer.LibraryVersions.SetDefault(nex.NewLibraryVersion(3, 4, 0))
	globals.SecureServer.AccessKey = "38dd2953"

	globals.SecureEndpoint.OnData(func(packet nex.PacketInterface) {
		request := packet.RMCMessage()

		fmt.Println("======= ABSW-U - Secure =======")
		fmt.Printf("Protocol ID: %#v\n", request.ProtocolID)
		fmt.Printf("Method ID: %#v\n", request.MethodID)
		fmt.Println("==================================================")
	})

	registerCommonSecureServerProtocols()

	port, _ := strconv.Atoi(os.Getenv("PN_SECURE_SERVER_PORT"))

	globals.SecureServer.Listen(port)
}
