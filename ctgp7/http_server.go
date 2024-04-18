package ctgp7

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"io"
	"net/http"
	"os"
	"sync"

	"github.com/PretendoNetwork/mario-kart-7/globals"
	nex "github.com/PretendoNetwork/nex-go/v2"
	"github.com/PretendoNetwork/nex-go/v2/constants"
	"github.com/PretendoNetwork/nex-go/v2/types"
	common_globals "github.com/PretendoNetwork/nex-protocols-common-go/v2/globals"
	notifications "github.com/PretendoNetwork/nex-protocols-go/v2/notifications"
	notifications_types "github.com/PretendoNetwork/nex-protocols-go/v2/notifications/types"
)

var outAESKey []byte
var CTGP7HTTPReadyToServe bool

func doPKCS5Padding(ciphertext []byte) []byte {
	blockSize := aes.BlockSize
	padding := (blockSize - len(ciphertext) % blockSize)
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

func doPKCS5Trimming(plaintext []byte) []byte {
	blockSize := aes.BlockSize
	padding := plaintext[len(plaintext)-1]
	cut := len(plaintext)-int(padding)
	if cut < 0 || int(padding) >= blockSize {
		// Invalid padding
		return []byte{}
	}
	return plaintext[:cut]
}

func aes256ENC(plaintext []byte) []byte {
	bIV := make([]byte, 16)
	rand.Read(bIV)
	bPlaintext := doPKCS5Padding(plaintext)
	block, _ := aes.NewCipher(outAESKey)
	ciphertext := make([]byte, len(bPlaintext))
	mode := cipher.NewCBCEncrypter(block, bIV)
	mode.CryptBlocks(ciphertext, bPlaintext)
	return append(bIV, ciphertext...)
}

func aes256DEC(ciphertext []byte) []byte {
	if len(ciphertext) <= 16 {
		return []byte{}
	}
	bIV := ciphertext[:16]
	bCipher := ciphertext[16:]
	block, _ := aes.NewCipher(outAESKey)
	mode := cipher.NewCBCDecrypter(block, bIV)
	plaintext := make([]byte, len(bCipher))
	mode.CryptBlocks(plaintext, bCipher)
	return doPKCS5Trimming(plaintext)
}

type userData struct {
	PID uint32
	CID uint32
	NatReportMyself natReport // * Reported by myself
	NatReportOthers natReport // * Reported by others
}

type sessionData struct {
	GID uint32
	HostPID uint32
	OwnerPID uint32
	Connections []uint32
}

type reportData struct {
	Users []userData
	Sessions []sessionData
}

func usersHandler(rw http.ResponseWriter, r *http.Request) {
	report := reportData{
		Users: make([]userData, 0),
		Sessions: make([]sessionData, 0),
	}

	if CTGP7HTTPReadyToServe {
		globals.SecureEndpoint.Connections.Each(func(key string, connection *nex.PRUDPConnection) bool {
			if connection.PID().LegacyValue() == 0 || connection.ID == 0 {
				return false
			}
			ud := userData{
				PID: connection.PID().LegacyValue(), 
				CID: connection.ID,
				NatReportMyself: natReport{},
				NatReportOthers: natReport{},
			}
			if report, ok := playerNATRepotsMyself.Get(connection.ID); ok {
				report.lock.RLock()
				ud.NatReportMyself.Results = report.Results
				report.lock.RUnlock()
			}
			if report, ok := playerNATRepotsOther.Get(connection.ID); ok {
				report.lock.RLock()
				ud.NatReportOthers.Results = report.Results
				report.lock.RUnlock()
			}
			report.Users = append(report.Users, ud)
			return false
		})
	
		common_globals.EachSession(func(index uint32, value *common_globals.CommonMatchmakeSession) bool {
			sd := sessionData{
				GID: value.GameMatchmakeSession.ID.Value,
				HostPID: value.GameMatchmakeSession.HostPID.LegacyValue(),
				OwnerPID: value.GameMatchmakeSession.OwnerPID.LegacyValue(),
				Connections: make([]uint32, 0),
			}
			value.ConnectionIDs.Each(func(index int, conID uint32) bool {
				sd.Connections = append(sd.Connections, conID)
				return false
			})
			report.Sessions = append(report.Sessions, sd)
			return false
		})
	}

	b, _ := json.Marshal(report)
	rw.Write(aes256ENC(b))
}

type natReport struct {
	Results [10]int; // 0 -> Not initialized, 1 -> False, 2 -> True
	resultIndex int;
	lock *sync.RWMutex;
}
var playerNATRepotsMyself *nex.MutexMap[uint32, *natReport]
var playerNATRepotsOther *nex.MutexMap[uint32, *natReport]

func OnAfterReportNATTraversalResult(packet nex.PacketInterface, cid *types.PrimitiveU32, result *types.PrimitiveBool, rtt *types.PrimitiveU32) {
	if cid == nil || result == nil {
		return
	}
	// * Other
	natreportOther, ok := playerNATRepotsOther.Get(cid.Value)
	if !ok {
		natreportOther = &natReport{
			resultIndex: 0,
			lock: &sync.RWMutex{},
		}
		playerNATRepotsOther.Set(cid.Value, natreportOther)
	}
	natreportOther.lock.Lock()
	defer natreportOther.lock.Unlock()
	if result.Value {
		natreportOther.Results[natreportOther.resultIndex] = 2
	} else {
		natreportOther.Results[natreportOther.resultIndex] = 1
	}
	natreportOther.resultIndex++
	if (natreportOther.resultIndex >= 10) {
		natreportOther.resultIndex = 0
	}
	// * Myself
	natreportMyself, ok := playerNATRepotsMyself.Get(packet.Sender().(*nex.PRUDPConnection).ID)
	if !ok {
		natreportMyself = &natReport{
			resultIndex: 0,
			lock: &sync.RWMutex{},
		}
		playerNATRepotsMyself.Set(packet.Sender().(*nex.PRUDPConnection).ID, natreportMyself)
	}
	natreportMyself.lock.Lock()
	defer natreportMyself.lock.Unlock()
	if result.Value {
		natreportMyself.Results[natreportMyself.resultIndex] = 2
	} else {
		natreportMyself.Results[natreportMyself.resultIndex] = 1
	}
	natreportMyself.resultIndex++
	if (natreportMyself.resultIndex >= 10) {
		natreportMyself.resultIndex = 0
	}
}

func OnPlayerJoinLeaveSession(gid uint32, cid uint32) {
	// playerNATRepots.Delete(cid)
}

type KickUsers struct {
	Connections []uint32
}

func kickHandler(rw http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost || r.ContentLength <= 0 || r.ContentLength > 1000 {
		rw.Write([]byte("NACK"))
		return
	}
	body, error := io.ReadAll(r.Body)
	if error != nil {
		rw.Write([]byte("NACK"))
		return
	}

	bodyDec := aes256DEC(body)
	kUsers := KickUsers{}

	if err := json.Unmarshal(bodyDec, &kUsers); err != nil {
		rw.Write([]byte("NACK"))
		return
	}

	for _, cID := range kUsers.Connections {
		connection := globals.SecureEndpoint.FindConnectionByID(cID)
		if connection == nil {
			continue
		}

		endpoint := connection.Endpoint().(*nex.PRUDPEndPoint)
		var category uint32 = 106
		var subtype uint32 = 0

		oEvent := notifications_types.NewNotificationEvent()
		oEvent.PIDSource = connection.PID()
		oEvent.Type = types.NewPrimitiveU32(notifications.BuildNotificationType(category, subtype))

		stream := nex.NewByteStreamOut(endpoint.LibraryVersions(), endpoint.ByteStreamSettings())

		oEvent.WriteTo(stream)

		rmcRequest := nex.NewRMCRequest(endpoint)
		rmcRequest.ProtocolID = notifications.ProtocolID
		rmcRequest.CallID = common_globals.CurrentMatchmakingCallID.Next()
		rmcRequest.MethodID = notifications.MethodProcessNotificationEvent
		rmcRequest.Parameters = stream.Bytes()

		rmcRequestBytes := rmcRequest.Bytes()

		var messagePacket nex.PRUDPPacketInterface

		if connection.DefaultPRUDPVersion == 0 {
			messagePacket, _ = nex.NewPRUDPPacketV0(endpoint.Server, connection, nil)
		} else {
			messagePacket, _ = nex.NewPRUDPPacketV1(endpoint.Server, connection, nil)
		}

		messagePacket.SetType(constants.DataPacket)
		messagePacket.AddFlag(constants.PacketFlagNeedsAck)
		messagePacket.AddFlag(constants.PacketFlagReliable)
		messagePacket.SetSourceVirtualPortStreamType(connection.StreamType)
		messagePacket.SetSourceVirtualPortStreamID(endpoint.StreamID)
		messagePacket.SetDestinationVirtualPortStreamType(connection.StreamType)
		messagePacket.SetDestinationVirtualPortStreamID(connection.StreamID)
		messagePacket.SetPayload(rmcRequestBytes)

		endpoint.Server.Send(messagePacket)
	}

	rw.Write([]byte("ACK"))
}

func StartHttpServer() {
	outAESKey, _ = hex.DecodeString(os.Getenv("PN_HTTP_SERVER_AES_KEY_OUT"))
	playerNATRepotsMyself = nex.NewMutexMap[uint32, *natReport]()
	playerNATRepotsOther = nex.NewMutexMap[uint32, *natReport]()

	common_globals.OnPlayerJoinSession(OnPlayerJoinLeaveSession)
	common_globals.OnPlayerLeaveSession(func (gid uint32, cid uint32, gracefully bool) {OnPlayerJoinLeaveSession(gid, cid)})

	m := http.NewServeMux()
	m.HandleFunc("/stats", usersHandler)
	m.HandleFunc("/kick", kickHandler)

	server := http.Server{
		Addr: os.Getenv("PN_HTTP_SERVER_LISTEN_ADDR"), // :9000
		Handler:  m,
	}
	server.ListenAndServe()
}