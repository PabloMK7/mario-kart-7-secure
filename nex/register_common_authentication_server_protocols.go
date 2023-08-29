package nex

import (
	"os"
	"strconv"

	nex "github.com/PretendoNetwork/nex-go"
	ticket_granting "github.com/PretendoNetwork/nex-protocols-common-go/ticket-granting"
	"github.com/PretendoNetwork/mario-kart-7/globals"

	nex_ticket_granting "github.com/PretendoNetwork/mario-kart-7/nex/ticket-granting"
)

func registerCommonAuthenticationServerProtocols() {
	ticketGrantingProtocol := ticket_granting.NewCommonTicketGrantingProtocol(globals.AuthenticationServer)

	port, _ := strconv.Atoi(os.Getenv("PN_MK7_SECURE_SERVER_PORT"))

	secureStationURL := nex.NewStationURL("")
	secureStationURL.SetScheme("prudps")
	secureStationURL.SetAddress(os.Getenv("PN_MK7_SECURE_SERVER_HOST"))
	secureStationURL.SetPort(uint32(port))
	secureStationURL.SetCID(1)
	secureStationURL.SetPID(2)
	secureStationURL.SetSID(1)
	secureStationURL.SetStream(10)
	secureStationURL.SetType(2)

	ticketGrantingProtocol.SetSecureStationURL(secureStationURL)
	nex_ticket_granting.SecureStationURL = secureStationURL

	ticketGrantingProtocol.SetBuildName(serverBuildString)
	nex_ticket_granting.BuildName = serverBuildString

	ticketGrantingProtocol.LoginEx(nex_ticket_granting.LoginEx)
	ticketGrantingProtocol.RequestTicket(nex_ticket_granting.RequestTicket)
}
