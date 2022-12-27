package constants

// 编译期指定
var (
	VERSION       string
	BUILD_COMMIT  string
	BUILD_TIME    string
	BUILD_STATIC  string
	BUILD_LINKAGE string
)

func init() {
	if BUILD_STATIC != "" {
		BUILD_LINKAGE = "statically linked"
	} else {
		BUILD_LINKAGE = "dynamically linked"
	}
}
