package user

type UserRole string

const (
	ADMIN    UserRole = "ADMIN"
	EMPLOYEE UserRole = "EMPLOYEE"
)
const maxAvatarSize = 2 << 20 // 2MB

var AllowedAvatarTypes = map[string]bool{
	"image/jpeg": true,
	"image/jpg":  true,
	"image/png":  true,
}
