package sm

import "github.com/zhikh23/sm-instruction/internal/common/commonerrs"

type Role struct {
	s string
}

var (
	Participant   = Role{s: "participant"}
	Administrator = Role{s: "administrator"}
)

func NewRoleFromString(s string) (Role, error) {
	switch s {
	case "participant":
		return Participant, nil
	case "administrator":
		return Administrator, nil
	}
	return Role{}, commonerrs.NewInvalidInputErrorf(
		"invalid user role: %s; expected one of ['participant', 'administrator']", s,
	)
}

func (r Role) String() string {
	return r.s
}

func (r Role) IsZero() bool {
	return r.s == ""
}
