package nex

import (
	"os"
	"strconv"

	"github.com/PretendoNetwork/mario-kart-7/ctgp7"
	"github.com/PretendoNetwork/mario-kart-7/globals"
	"github.com/PretendoNetwork/nex-go/v2"
)

var serverBuildString string

func StartAuthenticationServer() {
	globals.AuthenticationServer = nex.NewPRUDPServer()

	globals.AuthenticationEndpoint = nex.NewPRUDPEndPoint(1)
	globals.AuthenticationEndpoint.ServerAccount = globals.AuthenticationServerAccount
	globals.AuthenticationEndpoint.AccountDetailsByPID = ctgp7.AccountDetailsByPID
	globals.AuthenticationEndpoint.AccountDetailsByUsername = ctgp7.AccountDetailsByUsername
	globals.AuthenticationEndpoint.DefaultStreamSettings.MaxSilenceTime = 90000 / 2
	globals.AuthenticationEndpoint.DefaultStreamSettings.KeepAliveTimeout = 500
	globals.AuthenticationEndpoint.DefaultStreamSettings.ExtraRestransmitTimeoutTrigger = 0xFFFFFFFF
	globals.AuthenticationEndpoint.DefaultStreamSettings.RetransmitTimeoutMultiplier = 1.0
	globals.AuthenticationEndpoint.DefaultStreamSettings.MaxPacketRetransmissions = 15
	globals.AuthenticationServer.BindPRUDPEndPoint(globals.AuthenticationEndpoint)

	globals.AuthenticationServer.LibraryVersions.SetDefault(nex.NewLibraryVersion(2, 4, 3))
	globals.AuthenticationServer.AccessKey = "6181dff1"

	globals.AuthenticationEndpoint.OnData(func(packet nex.PacketInterface) {
		//request := packet.RMCMessage()

		//globals.Logger.Info("=== MK7 - Auth ===")
		//globals.Logger.Infof("Protocol ID: %#v", request.ProtocolID)
		//globals.Logger.Infof("Method ID: %#v", request.MethodID)
		//globals.Logger.Info("==================")
	})

	registerCommonAuthenticationServerProtocols()

	port, _ := strconv.Atoi(os.Getenv("PN_MK7_AUTHENTICATION_SERVER_PORT"))

	globals.AuthenticationServer.Listen(port)
}
