package profile

const maxAvatarSize = 2 << 20 // 2MB

var allowedAvatarTypes = map[string]bool{
	"image/jpeg": true,
	"image/jpg":  true,
	"image/png":  true,
}
