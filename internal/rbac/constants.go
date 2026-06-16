package rbac

const UserPermissionsCacheKey = "user_permissions:%s"
const UserRoleCacheKey = "auth:role:%s"

type RoleType string

const (
	ADMIN      RoleType = "ADMIN"
	EMPLOYEE   RoleType = "EMPLOYEE"
	SUPERADMIN RoleType = "SUPERADMIN"
)

const (
	PermissionCategoryView   = "category:view"
	PermissionCategoryCreate = "category:create"
	PermissionCategoryUpdate = "category:update"
	PermissionCategoryDelete = "category:delete"

	PermissionDivisionView   = "division:view"
	PermissionDivisionCreate = "division:create"
	PermissionDivisionUpdate = "division:update"
	PermissionDivisionDelete = "division:delete"

	PermissionUserView   = "user:view"
	PermissionUserCreate = "user:create"
	PermissionUserUpdate = "user:update"
	PermissionUserDelete = "user:delete"

	PermissionTicketView       = "ticket:view"
	PermissionTicketCreate     = "ticket:create"
	PermissionTicketUpdate     = "ticket:update"
	PermissionTicketDelete     = "ticket:delete"
	PermissionTicketAssign     = "ticket:assign"
	PermissionTicketPriority   = "ticket:priority"
	PermissionTicketResolution = "ticket:resolution"
	PermissionTicketClose      = "ticket:close"

	PermissionFeedbackView   = "feedback:view"
	PermissionFeedbackCreate = "feedback:create"
	PermissionFeedbackUpdate = "feedback:update"
	PermissionFeedbackDelete = "feedback:delete"

	PermissionDashboardView = "dashboard:view"

	PermissionRBACView   = "rbac:view"
	PermissionRBACCreate = "rbac:create"
	PermissionRBACUpdate = "rbac:update"
	PermissionRBACDelete = "rbac:delete"
)
