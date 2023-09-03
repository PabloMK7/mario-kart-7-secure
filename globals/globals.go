package globals

import (
	"github.com/PretendoNetwork/nex-go"
	"github.com/PretendoNetwork/plogger-go"
)

var Logger *plogger.Logger
var KerberosPassword = "password" // * Default password
var AuthenticationServer *nex.Server
var SecureServer *nex.Server
var PasswordServerURL = "https://localhost:80/t/%s" //%s will be replaced by the token
