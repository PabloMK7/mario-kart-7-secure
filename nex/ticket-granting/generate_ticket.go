package nex_ticket_granting

import (
	"crypto/rand"

	"github.com/PretendoNetwork/mario-kart-7/globals"
	"github.com/PretendoNetwork/nex-go"
)

var passwords map[uint32]string = make(map[uint32]string)

func generateTicket(userPID uint32, targetPID uint32, password string) ([]byte, uint32) {
	var userPassword string
	var targetPassword string

	// TODO - Maybe we should error out if the user PID is the server account?
	switch userPID {
	case 2: // "Quazal Rendez-Vous" (the server user) account
		userPassword = globals.KerberosPassword
	case 100: // guest user account
		userPassword = "MMQea3n!fsik"
	default:
		userPassword = password
	}

	if _, ok := passwords[userPID]; !ok {
		if userPassword == "" {
			return []byte{}, nex.Errors.RendezVous.InvalidUsername
		}

		passwords[userPID] = userPassword
	}

	if userPassword == "" {
		userPassword = passwords[userPID]
	}

	switch targetPID {
	case 2: // "Quazal Rendez-Vous" (the server user) account
		targetPassword = globals.KerberosPassword
	case 100: // guest user account
		targetPassword = "MMQea3n!fsik"
	default:
		targetPassword = password
	}

	if _, ok := passwords[targetPID]; !ok {
		if targetPassword == "" {
			return []byte{}, nex.Errors.RendezVous.InvalidUsername
		}

		passwords[targetPID] = targetPassword
	}

	if targetPassword == "" {
		targetPassword = passwords[targetPID]
	}

	userKey := nex.DeriveKerberosKey(userPID, []byte(userPassword))
	targetKey := nex.DeriveKerberosKey(targetPID, []byte(targetPassword))
	sessionKey := make([]byte, globals.AuthenticationServer.KerberosKeySize())
	_, err := rand.Read(sessionKey)
	if err != nil {
		globals.Logger.Error(err.Error())
		return []byte{}, nex.Errors.Authentication.Unknown
	}

	ticketInternalData := nex.NewKerberosTicketInternalData()
	serverTime := nex.NewDateTime(0)
	serverTime.UTC()
	ticketInternalData.SetTimestamp(serverTime)
	ticketInternalData.SetUserPID(userPID)
	ticketInternalData.SetSessionKey(sessionKey)

	encryptedTicketInternalData, err := ticketInternalData.Encrypt(targetKey, nex.NewStreamOut(globals.AuthenticationServer))
	if err != nil {
		globals.Logger.Error(err.Error())
		return []byte{}, nex.Errors.Authentication.Unknown
	}

	ticket := nex.NewKerberosTicket()
	ticket.SetSessionKey(sessionKey)
	ticket.SetTargetPID(targetPID)
	ticket.SetInternalData(encryptedTicketInternalData)

	encryptedTicket, err := ticket.Encrypt(userKey, nex.NewStreamOut(globals.AuthenticationServer))
	if err != nil {
		globals.Logger.Error(err.Error())
		return []byte{}, nex.Errors.Authentication.Unknown
	}

	return encryptedTicket, 0
}
