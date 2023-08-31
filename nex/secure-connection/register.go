package nex_secure_connection

import (
	"github.com/PretendoNetwork/mario-kart-7/globals"
	nex "github.com/PretendoNetwork/nex-go"
	secure_connection "github.com/PretendoNetwork/nex-protocols-go/secure-connection"
)

func Register(err error, client *nex.Client, callID uint32, stationUrls []*nex.StationURL) uint32 {
	if err != nil {
		globals.Logger.Error(err.Error())
		return nex.Errors.Core.InvalidArgument
	}

	nextConnectionID := uint32(globals.SecureServer.ConnectionIDCounter().Increment())
	client.SetConnectionID(nextConnectionID)

	localStation := stationUrls[0]

	// Mario Kart 7 already sets the public station for us
	publicStation := stationUrls[1]

	localStation.SetLocal()
	publicStation.SetPublic()

	client.AddStationURL(localStation)
	client.AddStationURL(publicStation)

	retval := nex.NewResultSuccess(nex.Errors.Core.Unknown)

	rmcResponseStream := nex.NewStreamOut(globals.SecureServer)

	rmcResponseStream.WriteResult(retval) // Success
	rmcResponseStream.WriteUInt32LE(client.ConnectionID())
	rmcResponseStream.WriteString(publicStation.EncodeToString())

	rmcResponseBody := rmcResponseStream.Bytes()

	// Build response packet
	rmcResponse := nex.NewRMCResponse(secure_connection.ProtocolID, callID)
	rmcResponse.SetSuccess(secure_connection.MethodRegister, rmcResponseBody)

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

	globals.SecureServer.Send(responsePacket)

	return 0
}
