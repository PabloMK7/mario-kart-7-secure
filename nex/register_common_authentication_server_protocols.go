package nex

import (
	"os"

	"github.com/PretendoNetwork/mario-kart-7/globals"
	nex_ticket_granting "github.com/PretendoNetwork/mario-kart-7/nex/ticket-granting"
	"github.com/PretendoNetwork/nex-go/types"
	common_ticket_granting "github.com/PretendoNetwork/nex-protocols-common-go/ticket-granting"
	ticket_granting "github.com/PretendoNetwork/nex-protocols-go/ticket-granting"
)

func registerCommonAuthenticationServerProtocols() {
	ticketGrantingProtocol := ticket_granting.NewProtocol()
	globals.AuthenticationEndpoint.RegisterServiceProtocol(ticketGrantingProtocol)
	globals.CommonTicketGrantingProtocol = common_ticket_granting.NewCommonProtocol(ticketGrantingProtocol)
	ticketGrantingProtocol.SetHandlerLoginEx(nex_ticket_granting.LoginEx)
	ticketGrantingProtocol.SetHandlerRequestTicket(nex_ticket_granting.RequestTicket)

	secureStationURL := types.NewStationURL("")
	secureStationURL.Scheme = "prudps"
	secureStationURL.Fields["address"] = os.Getenv("PN_MK7_SECURE_SERVER_HOST")
	secureStationURL.Fields["port"] = os.Getenv("PN_MK7_SECURE_SERVER_PORT")
	secureStationURL.Fields["CID"] = "1"
	secureStationURL.Fields["PID"] = "2"
	secureStationURL.Fields["sid"] = "1"
	secureStationURL.Fields["stream"] = "10"
	secureStationURL.Fields["type"] = "2"

	globals.CommonTicketGrantingProtocol.SecureStationURL = secureStationURL
	globals.CommonTicketGrantingProtocol.BuildName = types.NewString(serverBuildString)
	globals.CommonTicketGrantingProtocol.SecureServerAccount = globals.SecureServerAccount
}
