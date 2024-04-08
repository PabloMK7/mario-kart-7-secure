package globals

import (
	"github.com/PretendoNetwork/nex-go/v2"
	common_ticket_granting "github.com/PretendoNetwork/nex-protocols-common-go/v2/ticket-granting"
	"github.com/PretendoNetwork/plogger-go"
)

var Logger *plogger.Logger
var KerberosPassword = "password" // * Default password
var AuthenticationServer *nex.PRUDPServer
var AuthenticationEndpoint *nex.PRUDPEndPoint
var SecureServer *nex.PRUDPServer
var SecureEndpoint *nex.PRUDPEndPoint
var CommonTicketGrantingProtocol *common_ticket_granting.CommonProtocol
var PasswordServerURL = "https://localhost:80/t/%s" //%s will be replaced by the token
