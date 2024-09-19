package utils

var allowedImageFormats = []string{
	"image/svg+xml",
	"image/jpeg",
	"image/gif",
	"image/png",
}

func IsAllowedImageFormat(mime string) bool {
	for _, allowed := range allowedImageFormats {
		if allowed == mime {
			return true
		}
	}
	return false
}
