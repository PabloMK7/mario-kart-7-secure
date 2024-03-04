package nex_ticket_granting

import (
	"github.com/PretendoNetwork/nex-go"
	"github.com/PretendoNetwork/nex-go/types"
	common_globals "github.com/PretendoNetwork/nex-protocols-common-go/globals"
	ticket_granting "github.com/PretendoNetwork/nex-protocols-go/ticket-granting"
	ticket_granting_types "github.com/PretendoNetwork/nex-protocols-go/ticket-granting/types"

	"github.com/PretendoNetwork/mario-kart-7/globals"
)

func LoginEx(err error, packet nex.PacketInterface, callID uint32, strUserName *types.String, oExtraData *types.AnyDataHolder) (*nex.RMCMessage, *nex.Error) {
	if err != nil {
		common_globals.Logger.Error(err.Error())
		return nil, nex.NewError(nex.ResultCodes.Core.InvalidArgument, "change_error")
	}

	if !oExtraData.TypeName.Equals(types.NewString("AuthenticationInfo")) {
		return nil, nex.NewError(nex.ResultCodes.Core.InvalidArgument, "change_error")
	}

	authenticationInfo := oExtraData.ObjectData.(*ticket_granting_types.AuthenticationInfo)

	connection := packet.Sender().(*nex.PRUDPConnection)
	endpoint := connection.Endpoint().(*nex.PRUDPEndPoint)

	sourceAccount, errorCode := globals.CTGP7AccountDetailsByUsername(strUserName.Value, authenticationInfo.Token.Value)
	if errorCode != nil && errorCode.ResultCode != nex.ResultCodes.RendezVous.InvalidUsername {
		// * Some other error happened
		return nil, errorCode
	}

	targetAccount, errorCode := globals.CTGP7AccountDetailsByUsername(globals.SecureServerAccount.Username, authenticationInfo.Token.Value)
	if errorCode != nil && errorCode.ResultCode != nex.ResultCodes.RendezVous.InvalidUsername {
		// * Some other error happened
		return nil, errorCode
	}

	encryptedTicket, errorCode := generateTicket(sourceAccount, targetAccount, globals.AuthenticationServer.SessionKeyLength, endpoint)

	if errorCode != nil && errorCode.ResultCode != nex.ResultCodes.RendezVous.InvalidUsername {
		// * Some other error happened
		return nil, errorCode
	}

	var retval *types.QResult
	pidPrincipal := types.NewPID(0)
	pbufResponse := types.NewBuffer([]byte{})
	pConnectionData := types.NewRVConnectionData()
	strReturnMsg := types.NewString("")

	// * From the wiki:
	// *
	// * "If the username does not exist, the %retval% field is set to
	// * RendezVous::InvalidUsername and the other fields are left blank."
	if errorCode != nil && errorCode.ResultCode == nex.ResultCodes.RendezVous.InvalidUsername {
		retval = types.NewQResultError(errorCode.ResultCode)
	} else {
		retval = types.NewQResultSuccess(nex.ResultCodes.Core.Unknown)
		pidPrincipal = sourceAccount.PID
		pbufResponse = types.NewBuffer(encryptedTicket)
		strReturnMsg = globals.CommonTicketGrantingProtocol.BuildName.Copy().(*types.String)

		specialProtocols := types.NewList[*types.PrimitiveU8]()

		specialProtocols.Type = types.NewPrimitiveU8(0)
		specialProtocols.SetFromData(globals.CommonTicketGrantingProtocol.SpecialProtocols)

		pConnectionData.StationURL = globals.CommonTicketGrantingProtocol.SecureStationURL
		pConnectionData.SpecialProtocols = specialProtocols
		pConnectionData.StationURLSpecialProtocols = globals.CommonTicketGrantingProtocol.StationURLSpecialProtocols
		pConnectionData.Time = types.NewDateTime(0).Now()
	}

	rmcResponseStream := nex.NewByteStreamOut(endpoint.LibraryVersions(), endpoint.ByteStreamSettings())

	retval.WriteTo(rmcResponseStream)
	pidPrincipal.WriteTo(rmcResponseStream)
	pbufResponse.WriteTo(rmcResponseStream)
	pConnectionData.WriteTo(rmcResponseStream)
	strReturnMsg.WriteTo(rmcResponseStream)

	rmcResponseBody := rmcResponseStream.Bytes()

	rmcResponse := nex.NewRMCSuccess(endpoint, rmcResponseBody)
	rmcResponse.ProtocolID = ticket_granting.ProtocolID
	rmcResponse.MethodID = ticket_granting.MethodLoginEx
	rmcResponse.CallID = callID

	return rmcResponse, nil
}
