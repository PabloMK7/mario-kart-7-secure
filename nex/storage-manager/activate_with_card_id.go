package nex_storage_manager

import (
	"github.com/PretendoNetwork/mario-kart-7/globals"
	nex "github.com/PretendoNetwork/nex-go/v2"
	"github.com/PretendoNetwork/nex-go/v2/types"
	storage_manager "github.com/PretendoNetwork/nex-protocols-go/v2/storage-manager"
)

func ActivateWithCardID(err error, packet nex.PacketInterface, callID uint32, unknown *types.PrimitiveU8, cardID *types.PrimitiveU64) (*nex.RMCMessage, *nex.Error) {
	if err != nil {
		globals.Logger.Error(err.Error())
		return nil, nex.NewError(nex.ResultCodes.Core.InvalidArgument, err.Error())
	}

	uniqueID := types.NewPrimitiveU32(1)       // Unique ID
	firstTime := types.NewPrimitiveBool(false) // First time

	rmcResponseStream := nex.NewByteStreamOut(globals.SecureServer.LibraryVersions, globals.SecureServer.ByteStreamSettings)

	uniqueID.WriteTo(rmcResponseStream)
	firstTime.WriteTo(rmcResponseStream)

	rmcResponseBody := rmcResponseStream.Bytes()

	rmcResponse := nex.NewRMCSuccess(globals.SecureEndpoint, rmcResponseBody)
	rmcResponse.ProtocolID = storage_manager.ProtocolID
	rmcResponse.MethodID = storage_manager.MethodActivateWithCardID
	rmcResponse.CallID = callID

	return rmcResponse, nil
}
