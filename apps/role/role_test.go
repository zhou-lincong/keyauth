package role_test

import (
	"testing"

	"github.com/zhou-lincong/keyauth/apps/role"
)

func TestHasPermission(t *testing.T) {
	set := role.NewRoleSet()

	r := &role.Role{
		Spec: &role.CreateRoleRequest{
			Permissions: []*role.Permission{
				{
					Service: "cmdb",
					Featrues: []*role.Featrue{
						{
							Resource: "secret",
							Action:   "list",
						},
						{
							Resource: "secret",
							Action:   "get",
						},
						{
							Resource: "secret",
							Action:   "create",
						},
					},
				},
			},
		},
	}

	set.Add(r)

	perm, role := set.HasPermission(&role.PermissionRequest{
		Service:  "cmdb",
		Resource: "secret",
		Action:   "create",
	})

	t.Log(role)
	if perm != true {
		t.Fatal("has perm error")
	}
}

// === RUN   TestHasPermission
//     e:\goproject\keyauth\apps\role\role_test.go:44: spec:{permissions:{service:"cmdb" featrues:{resource:"secret" action:"list"} featrues:{resource:"secret" action:"get"} featrues:{resource:"secret" action:"create"}}}
// --- PASS: TestHasPermission (0.00s)
// PASS
// ok  	github.com/zhou-lincong/keyauth/apps/role	1.016s

//将Action:   "create",修改成Action:   "delete",测试
// === RUN   TestHasPermission
//     e:\goproject\keyauth\apps\role\role_test.go:44: <nil>
//     e:\goproject\keyauth\apps\role\role_test.go:46: has perm error
// --- FAIL: TestHasPermission (0.00s)
// FAIL
// FAIL	github.com/zhou-lincong/keyauth/apps/role	0.787s

func TestHasPermission2(t *testing.T) {
	set := role.NewRoleSet()

	r := &role.Role{
		Spec: &role.CreateRoleRequest{
			Permissions: []*role.Permission{
				{
					Service:  "cmdb",
					AllowAll: true,
				},
			},
		},
	}

	set.Add(r)

	perm, role := set.HasPermission(&role.PermissionRequest{
		Service:  "cmdb",
		Resource: "secret",
		Action:   "delete",
	})

	t.Log(role)
	if perm != true {
		t.Fatal("has perm error")
	}
}

// === RUN   TestHasPermission2
//     e:\goproject\keyauth\apps\role\role_test.go:86: spec:{permissions:{service:"cmdb" allow_all:true}}
// --- PASS: TestHasPermission2 (0.00s)
// PASS
// ok  	github.com/zhou-lincong/keyauth/apps/role	0.573s
