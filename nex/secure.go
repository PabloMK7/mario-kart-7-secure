package nex

import (
	"os"
	"strconv"

	"github.com/PretendoNetwork/mario-kart-7/ctgp7"
	"github.com/PretendoNetwork/mario-kart-7/globals"
	nex "github.com/PretendoNetwork/nex-go/v2"
	common_globals "github.com/PretendoNetwork/nex-protocols-common-go/v2/globals"
)

func StartSecureServer() {
	globals.SecureServer = nex.NewPRUDPServer()

	globals.SecureEndpoint = nex.NewPRUDPEndPoint(1)
	globals.SecureEndpoint.IsSecureEndPoint = true
	globals.SecureEndpoint.ServerAccount = globals.SecureServerAccount
	globals.SecureEndpoint.AccountDetailsByPID = ctgp7.AccountDetailsByPID
	globals.SecureEndpoint.AccountDetailsByUsername = ctgp7.AccountDetailsByUsername
	globals.SecureEndpoint.DefaultStreamSettings.MaxSilenceTime = 90000 / 2
	globals.SecureEndpoint.DefaultStreamSettings.KeepAliveTimeout = 500
	globals.SecureEndpoint.DefaultStreamSettings.ExtraRetransmitTimeoutTrigger = 0xFFFFFFFF
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

	globals.SecureEndpoint.OnConnectionEnded(ctgp7.OnConnectionEnded)
	globals.SecureEndpoint.OnConnectionEnded(ctgp7.OnConnectionEndedVRHandler)
	common_globals.FilterFoundCandidateSessions(ctgp7.FilterFoundCandidateSessions)

	registerCommonSecureServerProtocols()
	registerSecureServerNEXProtocols()

	port, _ := strconv.Atoi(os.Getenv("PN_MK7_SECURE_SERVER_PORT"))

	ctgp7.CTGP7HTTPReadyToServe = true
	globals.SecureServer.Listen(port)
}
