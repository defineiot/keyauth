package user

// Department user's department
type Department struct {
	ID       string `json:"id"`
	Name     string `json:"name"`      // 部门名称
	Grade    string `json:"grade"`     // 第几级部门
	Path     string `json:"path"`      // 部门访问路径
	CreateAt int64  `json:"create_at"` // 部门创建时间
	DomainID string `json:"domain_id"` // 部门所属域

	ParentID  string `json:"parent_id"`  // 上级部门ID
	ManagerID string `json:"manager_id"` // 部门管理者ID
}
