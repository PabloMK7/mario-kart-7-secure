package nex

import (
	"os"
	"strconv"

	"github.com/PretendoNetwork/mario-kart-7/globals"
	nex "github.com/PretendoNetwork/nex-go/v2"
)

func StartSecureServer() {
	globals.SecureServer = nex.NewPRUDPServer()

	globals.SecureEndpoint = nex.NewPRUDPEndPoint(1)
	globals.SecureEndpoint.IsSecureEndPoint = true
	globals.SecureEndpoint.ServerAccount = globals.SecureServerAccount
	globals.SecureEndpoint.AccountDetailsByPID = globals.AccountDetailsByPID
	globals.SecureEndpoint.AccountDetailsByUsername = globals.AccountDetailsByUsername
	globals.SecureEndpoint.DefaultStreamSettings.MaxSilenceTime = 90000
	globals.SecureEndpoint.DefaultStreamSettings.KeepAliveTimeout = 500
	globals.SecureEndpoint.DefaultStreamSettings.ExtraRestransmitTimeoutTrigger = 0xFFFFFFFF
	globals.SecureEndpoint.DefaultStreamSettings.RetransmitTimeoutMultiplier = 1.0
	globals.SecureEndpoint.DefaultStreamSettings.MaxPacketRetransmissions = 15
	globals.SecureServer.BindPRUDPEndPoint(globals.SecureEndpoint)

	globals.SecureServer.LibraryVersions.SetDefault(nex.NewLibraryVersion(2, 4, 3))
	globals.SecureServer.AccessKey = "6181dff1"

	globals.SecureEndpoint.OnData(func(packet nex.PacketInterface) {
		//request := packet.RMCMessage()

		//globals.Logger.Info("=== MK7 - Secure ===")
		//globals.Logger.Infof("Protocol ID: %#v", request.ProtocolID)
		//globals.Logger.Infof("Method ID: %#v", request.MethodID)
		//globals.Logger.Info("====================")
	})

	registerCommonSecureServerProtocols()
	registerSecureServerNEXProtocols()

	port, _ := strconv.Atoi(os.Getenv("PN_MK7_SECURE_SERVER_PORT"))

	globals.SecureServer.Listen(port)
}
