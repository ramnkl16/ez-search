package global

type ReferenceType string

const (
	DefaultMenuPermission  = 31
	DefaultUserPassword    = "welcome@123"
	DefaultUserFirstName   = "Admin"
	DefaultPlayerFirstName = "Guest"
	DefaultLastName        = "User"
	RootMenuID             = "ROOT"

	RefTypeNamespace ReferenceType = "NS"
	RefTypeGroup     ReferenceType = "GR"
	RefTypeUser      ReferenceType = "UR"
)

func GetAllGlobalAdminNS() []string {
	return []string{"PLATFORM"}
}
