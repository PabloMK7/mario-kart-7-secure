package nex

import (
	"fmt"
	"os"

	nex "github.com/PretendoNetwork/nex-go"
	"github.com/PretendoNetwork/mario-kart-7/globals"
)

var serverBuildString string

func StartAuthenticationServer() {
	globals.AuthenticationServer = nex.NewServer()
	globals.AuthenticationServer.SetPRUDPVersion(0)
	globals.SecureServer.SetDefaultNEXVersion(&nex.NEXVersion{
		Major: 2,
		Minor: 4,
		Patch: 3,
	})

	globals.AuthenticationServer.SetKerberosPassword(globals.KerberosPassword)
	globals.AuthenticationServer.SetAccessKey("6181dff1")

	globals.AuthenticationServer.On("Data", func(packet *nex.PacketV0) {
		_ := packet.RMCRequest()

		// fmt.Println("=== MK7 - Auth ===")
		// fmt.Printf("Protocol ID: %#v\n", request.ProtocolID())
		// fmt.Printf("Method ID: %#v\n", request.MethodID())
		// fmt.Println("==================")
	})

	registerCommonAuthenticationServerProtocols()

	globals.AuthenticationServer.Listen(fmt.Sprintf(":%s", os.Getenv("PN_MK7_AUTHENTICATION_SERVER_PORT")))
}
