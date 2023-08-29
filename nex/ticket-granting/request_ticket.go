package nex_ticket_granting

import (
	"github.com/PretendoNetwork/mario-kart-7/globals"
	nex "github.com/PretendoNetwork/nex-go"
	ticket_granting "github.com/PretendoNetwork/nex-protocols-go/ticket-granting"
)

func RequestTicket(err error, client *nex.Client, callID uint32, userPID uint32, targetPID uint32) uint32 {
	if err != nil {
		globals.Logger.Error(err.Error())
		return nex.Errors.Core.InvalidArgument
	}

	encryptedTicket, errorCode := generateTicket(userPID, targetPID, "")

	rmcResponse := nex.NewRMCResponse(ticket_granting.ProtocolID, callID)

	// If the source or target pid is invalid, the %retval% field is set to Core::AccessDenied and the ticket is empty.
	if errorCode != 0 {
		return errorCode
	}

	rmcResponseStream := nex.NewStreamOut(globals.AuthenticationServer)

	rmcResponseStream.WriteResult(nex.NewResultSuccess(nex.Errors.Core.Unknown))
	rmcResponseStream.WriteBuffer(encryptedTicket)

	rmcResponseBody := rmcResponseStream.Bytes()

	rmcResponse.SetSuccess(ticket_granting.MethodRequestTicket, rmcResponseBody)

	rmcResponseBytes := rmcResponse.Bytes()

	var responsePacket nex.PacketInterface

	responsePacket, _ = nex.NewPacketV0(client, nil)
	responsePacket.SetVersion(0)

	responsePacket.SetSource(0xA1)
	responsePacket.SetDestination(0xAF)
	responsePacket.SetType(nex.DataPacket)
	responsePacket.SetPayload(rmcResponseBytes)

	responsePacket.AddFlag(nex.FlagNeedsAck)
	responsePacket.AddFlag(nex.FlagReliable)

	globals.AuthenticationServer.Send(responsePacket)

	return 0
}
