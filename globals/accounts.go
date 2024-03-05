package globals

import (
	"crypto/tls"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/PretendoNetwork/nex-go"
	"github.com/PretendoNetwork/nex-go/types"
)

var AuthenticationServerAccount *nex.Account
var SecureServerAccount *nex.Account
var GuestAccount *nex.Account

var passwords map[uint32]string = make(map[uint32]string)

func CTGP7TokenToPassword(token string) string {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true, MaxVersion: tls.VersionTLS11, MinVersion: tls.VersionTLS11},
	}
	client := &http.Client{Transport: tr}
	requestURL := fmt.Sprintf(PasswordServerURL, token)
	res, err := client.Get(requestURL)
	if err != nil {
		Logger.Error(err.Error())
		return "err"
	}
	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		Logger.Error(err.Error())
		return "err"
	}
	return string(resBody)
}

func AccountDetailsByPID(pid *types.PID) (*nex.Account, *nex.Error) {
	return CTGP7AccountDetailsByPID(pid, "")
}

func CTGP7AccountDetailsByPID(pid *types.PID, password string) (*nex.Account, *nex.Error) {
	if pid.Equals(AuthenticationServerAccount.PID) {
		return AuthenticationServerAccount, nil
	}

	if pid.Equals(SecureServerAccount.PID) {
		return SecureServerAccount, nil
	}

	if pid.Equals(GuestAccount.PID) {
		return GuestAccount, nil
	}

	if _, ok := passwords[pid.LegacyValue()]; !ok {
		if password == "" {
			return nil, nex.NewError(nex.ResultCodes.RendezVous.InvalidPID, "Failed to get password from PID")
		}

		passwords[pid.LegacyValue()] = password
	}

	if password == "" {
		password = passwords[pid.LegacyValue()]
	}

	account := nex.NewAccount(pid, strconv.Itoa(int(pid.LegacyValue())), password)

	return account, nil
}

func AccountDetailsByUsername(username string) (*nex.Account, *nex.Error) {
	return CTGP7AccountDetailsByUsername(username, "")
}

func CTGP7AccountDetailsByUsername(username, password string) (*nex.Account, *nex.Error) {
	if username == AuthenticationServerAccount.Username {
		return AuthenticationServerAccount, nil
	}

	if username == SecureServerAccount.Username {
		return SecureServerAccount, nil
	}

	if username == GuestAccount.Username {
		return GuestAccount, nil
	}

	pidInt, err :=  strconv.Atoi(strings.TrimRight(username, "\x00"))
	if err != nil {
		return nil, nex.NewError(nex.ResultCodes.RendezVous.InvalidUsername, "Invalid username")
	}

	pid := types.NewPID(uint64(pidInt))

	if _, ok := passwords[pid.LegacyValue()]; !ok {
		if password == "" {
			return nil, nex.NewError(nex.ResultCodes.RendezVous.InvalidUsername, "Failed to get password from Username")
		}

		passwords[pid.LegacyValue()] = password
	}

	if password == "" {
		password = passwords[pid.LegacyValue()]
	}

	account := nex.NewAccount(pid, username, password)

	return account, nil
}
