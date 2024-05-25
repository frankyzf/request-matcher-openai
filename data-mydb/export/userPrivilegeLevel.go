package export

var GroupPrivilege = map[string]int{
	"unknown":       0,
	"user":          100,
	"user-operator": 110,
	"operator":      1000, //account operator
	"account":       1000, //account operator
	"admin":         2000,
}
