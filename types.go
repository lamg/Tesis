package tesis

import (
	"github.com/dgrijalva/jwt-go"
	"time"
)

// This interface is an abstract program specification of the server.
// The rest of types and procedures are defined for using them inside
// it, or its implementation.
type DBManager interface {
	Authenticate(user, password string) (bool, error)
	UserInfo(string) (*UserInfo, error)
	Record(string, int) (*PageC, error)
	Propose(string, []string) error
	Proposed(string, int) (*PageD, error)
	Pending(int) (*PageD, error)
	RevertProp(string, []string) error
	Close() error
}

type RecordProvider interface {
	Records() ([]DBRecord, error)
	Name() string
}

type Activity struct {
	Record   []Change `json:"record"`
	Proposed []Diff   `json:"proposed"`
}

type StateSys struct {
	Pending []Diff               `json:"pending"`
	UsrAct  map[string]*Activity `json:"usrAct"`
}

type RecordReceptor interface {
	Create(id string, d *DBRecord) (e error)
	Update(id string, d *DBRecord) (e error)
	Delete(id string) (e error)
}

type UserDB interface {
	Authenticate(string, string) (bool, error)
	UserInfo(string) (*UserInfo, error)
}

//TODO construct adequate interfaces

type PageN struct {
	PageN int `json:"pageN"`
}

type PageC struct {
	Total   int      `json:"total"`
	PageN   int      `json:"pageN"`
	ChangeP []Change `json:"changeP"`
}

type PageD struct {
	Total int    `json:"total"`
	PageN int    `json:"pageN"`
	DiffP []Diff `json:"diffP"`
}

type UserInfo struct {
	Name     string `json:"name"`
	UserName string `json:"username"`
	Token    string `json:"token"`
}

type Error struct {
	Message string `json:"error"`
}

type Change struct {
	Time time.Time `json:"time"`
	SRec []Diff    `json:"srec"` //success
	FRec []Diff    `json:"frec"` //failed
}

type Credentials struct {
	User string `json:"user"`
	Pass string `json:"pass"`
}

type User struct {
	UserName string `json:"user"`
	jwt.StandardClaims
}

type Diff struct {
	LDAPRec DBRecord `json:"ldapRec"`
	//LDAPRec.Id is distinguishedName in LDAP
	//LDAPRec.Name is user DN in LDAP
	//LDAPRec.IN is employeeID in LDAP
	//LDAPRec.Addr is streetAddress in LDAP
	//LDAPRec.Tel is telephoneNumber in LDAP
	DBRec DBRecord `json:"dbRec"`
	//DBRec.Id is identification in SIGENU
	//DBRec.Name is name+middle_name+last_name in SIGENU
	//DBRec.IN is id_student in SIGENU
	//DBRec.Addr is address in SIGENU
	//DBRec.Tel is phone in SIGENU
	Src      string `json:"src"`
	Exists   bool   `json:"exists"`
	Mismatch bool   `json:"mismatch"`
}

type DBRecord struct {
	//database key field
	//DN in LDAP
	Id string `json:"id"`
	//identity number
	//employeeID in LDAP
	IN string `json:"in"`
	//person name
	//displayName in LDAP
	Name string `json:"name"`
	//address
	//streetAddress in LDAP
	Addr string `json:"addr"`
	//telephone number
	//telephoneNumber in LDAP
	Tel string `json:"tel"`
}

// isCandidate ≡ ¬hasId ∨ existsSimilar
// existsSimilar ≡ toLowerEq ∨ unAccentEq

type Eq interface {
	Equals(interface{}) bool
}

type Sim interface {
	Similar(interface{}) bool
	Eq
}

type Reporter interface {
	Progress(float32)
}

type Logger interface {
	Logf(string, ...interface{})
}
