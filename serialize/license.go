package serialize

import uuid "github.com/satori/go.uuid"

type LicenseInfo struct {
	Process           uuid.UUID //process            : 94232f94-be78-4acd-a11e-09911bd4f4ed
	Session           uuid.UUID //session            : e45c1c2b-b3ac-4fea-9f0c-0583ad65d117
	UserName          string    //user-name          : User
	Host              string    //host               : host
	AppId             string    //app-id             : 1CV8
	FullName          string    //full-name          :
	Series            string    //series             : "ORG8A"
	IssuedByServer    bool      //issued-by-server   : yes
	LicenseType       string    //license-type       : HASP
	Net               bool      //net                : yes
	MaxUsersAll       int32     //max-users-all      : 300
	MaxUsersCur       int32     //max-users-cur      : 300
	RmngrAddress      string    //rmngr-address      : "app"
	RmngrPort         int       //rmngr-port         : 1541
	RmngrPid          int32     //rmngr-pid          : 2300
	ShortPresentation string    //short-presentation : "Сервер, ORG8A Сет 300"
	FullPresentation  string    //full-presentation  : "Сервер, 2300, app, 1541, ORG8A Сетевой 300"
}

//
//final String fullName = decoder.decodeString(buffer);
//final String fullPresentation = decoder.decodeString(buffer);
//final boolean issuedByServer = decoder.decodeBoolean(buffer);
//final int licenseType = decoder.decodeInt(buffer);
//final int maxUsersAll = decoder.decodeInt(buffer);
//final int maxUsersCur = decoder.decodeInt(buffer);
//final boolean net = decoder.decodeBoolean(buffer);
//final String rmngrAddress = decoder.decodeString(buffer);
//final String rmngrPid = decoder.decodeString(buffer);
//final int rmngrPort = decoder.decodeInt(buffer);
//final String series = decoder.decodeString(buffer);
//final String shortPresentation = decoder.decodeString(buffer);
