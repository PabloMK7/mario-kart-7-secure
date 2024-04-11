package main

import (
	"os"
	"strconv"
	"strings"

	"github.com/PretendoNetwork/mario-kart-7/globals"
	"github.com/PretendoNetwork/nex-go/v2"
	"github.com/PretendoNetwork/nex-go/v2/types"
	"github.com/PretendoNetwork/plogger-go"
	"github.com/joho/godotenv"
	common_globals "github.com/PretendoNetwork/nex-protocols-common-go/v2/globals"
)

func init() {
	globals.Logger = plogger.NewLogger()
	plogger.LogToFile = false
	plogger.LogToStdOut = false
	common_globals.SessionManagementDebugLog = false

	var err error

	err = godotenv.Load()
	if err != nil {
		globals.Logger.Warning("Error loading .env file")
	}

	kerberosPassword := os.Getenv("PN_MK7_KERBEROS_PASSWORD")
	authenticationServerPort := os.Getenv("PN_MK7_AUTHENTICATION_SERVER_PORT")
	secureServerHost := os.Getenv("PN_MK7_SECURE_SERVER_HOST")
	secureServerPort := os.Getenv("PN_MK7_SECURE_SERVER_PORT")
	globals.PasswordServerURL = os.Getenv("PN_MK7_PASSWORD_SERVER_URL")

	if strings.TrimSpace(kerberosPassword) == "" {
		globals.Logger.Warningf("PN_MK7_KERBEROS_PASSWORD environment variable not set. Using default password: %q", globals.KerberosPassword)
	} else {
		globals.KerberosPassword = kerberosPassword
	}

	globals.AuthenticationServerAccount = nex.NewAccount(types.NewPID(1), "Quazal Authentication", globals.KerberosPassword)
	globals.SecureServerAccount = nex.NewAccount(types.NewPID(2), "Quazal Rendez-Vous", globals.KerberosPassword)
	globals.GuestAccount = nex.NewAccount(types.NewPID(100), "guest", "MMQea3n!fsik")

	if strings.TrimSpace(authenticationServerPort) == "" {
		globals.Logger.Error("PN_MK7_AUTHENTICATION_SERVER_PORT environment variable not set")
		os.Exit(0)
	}

	if port, err := strconv.Atoi(authenticationServerPort); err != nil {
		globals.Logger.Errorf("PN_MK7_AUTHENTICATION_SERVER_PORT is not a valid port. Expected 0-65535, got %s", authenticationServerPort)
		os.Exit(0)
	} else if port < 0 || port > 65535 {
		globals.Logger.Errorf("PN_MK7_AUTHENTICATION_SERVER_PORT is not a valid port. Expected 0-65535, got %s", authenticationServerPort)
		os.Exit(0)
	}

	if strings.TrimSpace(secureServerHost) == "" {
		globals.Logger.Error("PN_MK7_SECURE_SERVER_HOST environment variable not set")
		os.Exit(0)
	}

	if strings.TrimSpace(secureServerPort) == "" {
		globals.Logger.Error("PN_MK7_SECURE_SERVER_PORT environment variable not set")
		os.Exit(0)
	}

	if port, err := strconv.Atoi(secureServerPort); err != nil {
		globals.Logger.Errorf("PN_MK7_SECURE_SERVER_PORT is not a valid port. Expected 0-65535, got %s", secureServerPort)
		os.Exit(0)
	} else if port < 0 || port > 65535 {
		globals.Logger.Errorf("PN_MK7_SECURE_SERVER_PORT is not a valid port. Expected 0-65535, got %s", secureServerPort)
		os.Exit(0)
	}

	// database.ConnectPostgres()
}
