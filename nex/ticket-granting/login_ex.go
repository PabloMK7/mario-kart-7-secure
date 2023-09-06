package nex_ticket_granting

import (
	"strconv"
	"strings"
	"fmt"
	"net/http"
    "crypto/tls"
	"io"

	"github.com/PretendoNetwork/mario-kart-7/globals"
	nex "github.com/PretendoNetwork/nex-go"
	ticket_granting "github.com/PretendoNetwork/nex-protocols-go/ticket-granting"
	ticket_granting_types "github.com/PretendoNetwork/nex-protocols-go/ticket-granting/types"
)

var SecureStationURL *nex.StationURL
var BuildName string

func tokenToPassword(token string) string {
	tr := &http.Transport{
        TLSClientConfig: &tls.Config{InsecureSkipVerify: true, MaxVersion: tls.VersionTLS10, MinVersion: tls.VersionTLS10},
    }
	client := &http.Client{Transport: tr}
	requestURL := fmt.Sprintf(globals.PasswordServerURL, token)
	res, err := client.Get(requestURL)
    if err != nil {
		globals.Logger.Error(err.Error())
		return "err"
	}
	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		globals.Logger.Error(err.Error())
        return "err"
    }
	return string(resBody)
}

func LoginEx(err error, client *nex.Client, callID uint32, username string, oExtraData *nex.DataHolder) uint32 {
	if err != nil {
		globals.Logger.Error(err.Error())
		return nex.Errors.Core.InvalidArgument
	}

	if oExtraData.TypeName() != "AuthenticationInfo" {
		return nex.Errors.Core.InvalidArgument
	}

	authenticationInfo := oExtraData.ObjectData().(*ticket_granting_types.AuthenticationInfo)

	var userPID uint32

	if username == "guest" {
		userPID = 100
	} else {
		converted, err := strconv.Atoi(strings.TrimRight(username, "\x00"))
		if err != nil {
			panic(err)
		}

		userPID = uint32(converted)
	}

	var targetPID uint32 = 2 // "Quazal Rendez-Vous" (the server user) account PID

	encryptedTicket, errorCode := generateTicket(userPID, targetPID, tokenToPassword(authenticationInfo.Token))

	rmcResponse := nex.NewRMCResponse(ticket_granting.ProtocolID, callID)

	if errorCode != 0 && errorCode != nex.Errors.RendezVous.InvalidUsername {
		// Some other error happened
		return errorCode
	}

	var retval *nex.Result
	var pidPrincipal uint32
	var pbufResponse []byte
	var pConnectionData *nex.RVConnectionData
	var strReturnMsg string

	pConnectionData = nex.NewRVConnectionData()
	pConnectionData.SetStationURL(SecureStationURL.EncodeToString())
	pConnectionData.SetSpecialProtocols([]byte{})
	pConnectionData.SetStationURLSpecialProtocols("")
	serverTime := nex.NewDateTime(0)
	pConnectionData.SetTime(nex.NewDateTime(serverTime.UTC()))

	/*
		From the wiki:

		"If the username does not exist, the %retval% field is set to
		RendezVous::InvalidUsername and the other fields are left blank."
	*/
	if errorCode == nex.Errors.RendezVous.InvalidUsername {
		retval = nex.NewResultError(errorCode)
	} else {
		retval = nex.NewResultSuccess(nex.Errors.Core.Unknown)
		pidPrincipal = userPID
		pbufResponse = encryptedTicket
		strReturnMsg = BuildName
	}

	rmcResponseStream := nex.NewStreamOut(globals.AuthenticationServer)

	rmcResponseStream.WriteResult(retval)
	rmcResponseStream.WriteUInt32LE(pidPrincipal)
	rmcResponseStream.WriteBuffer(pbufResponse)
	rmcResponseStream.WriteStructure(pConnectionData)
	rmcResponseStream.WriteString(strReturnMsg)

	rmcResponseBody := rmcResponseStream.Bytes()

	rmcResponse.SetSuccess(ticket_granting.MethodLoginEx, rmcResponseBody)

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
