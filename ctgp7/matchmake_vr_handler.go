package ctgp7

import (
	"math"
	"sort"

	"github.com/PretendoNetwork/mario-kart-7/globals"
	nex "github.com/PretendoNetwork/nex-go/v2"
	common_globals "github.com/PretendoNetwork/nex-protocols-common-go/v2/globals"
	match_making_types "github.com/PretendoNetwork/nex-protocols-go/v2/match-making/types"
)

var playerVRs *nex.MutexMap[uint32, int32]

func getSessionVRMean(sessionID uint32) int32 {
	var vrMean int32 = 0
	var vrCount int32 = 0

	session, ok := common_globals.GetSession(sessionID)
	if !ok {
		return -1000000
	}

	session.ConnectionIDs.Each(func(index int, conID uint32) bool {
		vr, ok := playerVRs.Get(conID)
		if ok {
			vrMean += vr
			vrCount++
		}
		return false
	})

	if vrCount == 0 {
		return -1000000
	}

	return vrMean / vrCount
}

func FilterFoundCandidateSessions(sessions []uint32, connection *nex.PRUDPConnection, searchMatchmakeSession *match_making_types.MatchmakeSession) []uint32 {
	vrPrimitive, _ := searchMatchmakeSession.Attributes.Get(1)
	vr := int32(vrPrimitive.Value)
	playerVRs.Set(connection.ID, vr)

	transformVR := func(vr int32) int32 {
		vrFloat := float64(vr)
		return int32(math.Log(vrFloat) + (2000000.0/math.Sqrt(vrFloat)))
	}

	sort.Slice(sessions, func(i, j int) bool {
		transformedVR := transformVR(vr)

		var vrDifi int32 = transformVR(getSessionVRMean(sessions[i])) - transformedVR
		if vrDifi < 0 {vrDifi = -vrDifi}

		var vrDifj int32 = transformVR(getSessionVRMean(sessions[j])) - transformedVR
		if vrDifj < 0 {vrDifj = -vrDifj}

		return vrDifi < vrDifj
	})
	return sessions
}

func OnConnectionEndedVRHandler(connection *nex.PRUDPConnection) {
	playerVRs.Delete(connection.ID)
}

func InitMatchMakeVRHandler() {
	playerVRs = nex.NewMutexMap[uint32, int32]()
	common_globals.FilterFoundCandidateSessions(FilterFoundCandidateSessions)
	globals.SecureEndpoint.OnConnectionEnded(OnConnectionEndedVRHandler)
}