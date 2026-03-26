package utils

// Source - https://stackoverflow.com/a/73750401
// Posted by JAMOLIDDIN BAKHRIDDINOV, modified by community. See post 'Timeline' for change history
// Retrieved 2026-03-15, License - CC BY-SA 4.0

func IsUpper(s string) bool {
	for _, charNumber := range s {
		if charNumber > 90 || charNumber < 65 {
			return false
		}
	}
	return true
}

func IsLower(s string) bool {
	for _, charNumber := range s {
		if charNumber > 122 || charNumber < 97 {
			return false
		}
	}
	return true
}

// =============================================
