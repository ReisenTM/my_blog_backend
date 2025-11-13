package enum

type RoleType uint8

const (
	RoleUserType RoleType = iota + 1
	RoleAdminType
	RoleGuestType
)
