package models

import (
	"errors"
	"fmt"
	"strconv"

	"GoTenancy/backend/config"
	"GoTenancy/backend/database"
	"GoTenancy/backend/validates"
	"github.com/fatih/color"
	"github.com/jinzhu/gorm"
)

/**
*初始化系统 账号 权限 角色
 */
func CreateSystemData(perms []*validates.PermissionRequest) {
	if config.GetAppCreateSysData() {
		CreateSystemAdmin() //后端管理员

		// 租户端公用系统数据
		permIds := CreateSystemAdminPermission(perms) //初始化权限
		role := CreateSystemAdminRole(permIds)        //初始化角色
		if role.ID != 0 {
			CreateSystemUser(role.ID) //初始化管理员
		}
	}
}

/**
*创建系统管理员
*@param role_id uint
*@return   *models.AdminUserTranform api格式化后的数据格式
 */
func CreateSystemUser(roleId uint) {
	aul := &validates.CreateUpdateUserRequest{
		Username: config.GetTestDataUserName(),
		Password: config.GetTestDataPwd(),
		Name:     config.GetTestDataName(),
		RoleIds:  []uint{roleId},
	}

	user := NewUserByStruct(aul)
	user.GetUserByUsername()
	if user.ID == 0 {
		user.CreateUser(aul)
	}
}

/**
*创建系统管理员
*@return   *models.AdminRoleTranform api格式化后的数据格式
 */
func CreateSystemAdminRole(permIds []uint) *Role {
	rr := &validates.RoleRequest{
		Name:        "admin",
		DisplayName: "管理员",
		Description: "管理员",
	}
	role := NewRoleByStruct(rr)
	role.GetRoleByName()
	if role.ID == 0 {
		role.CreateRole(permIds)
	}

	return role
}

/**
 * 创建系统权限
 * @return
 */
func CreateSystemAdminPermission(perms []*validates.PermissionRequest) []uint {
	var permIds []uint
	for _, perm := range perms {
		p := NewPermission(0, perm.Name, perm.Act)
		p.GetPermissionByNameAct()
		if p.ID != 0 {
			continue
		}
		p.CreatePermission()
		permIds = append(permIds, p.ID)
	}
	return permIds
}

/**
 * 后端管理员
 */
func CreateSystemAdmin() {
	cuar := &validates.CreateUpdateAdminRequest{
		Name:     config.GetAdminName(),
		Username: config.GetAdminUserName(),
		Password: config.GetAdminPwd(),
	}

	admin := NewAdminByStruct(cuar)
	admin.GetAdminByUserName()
	if admin.ID == 0 {
		admin.CreateAdmin(cuar)
	}
}

func IsNotFound(err error) {
	if ok := errors.Is(err, gorm.ErrRecordNotFound); !ok && err != nil {
		color.Red(fmt.Sprintf("error :%v \n ", err))
	}
}

/**
 * 获取列表
 * @method GetAll
 * @param  {[type]} string string    [description]
 * @param  {[type]} orderBy string    [description]
 * @param  {[type]} relation string    [description]
 * @param  {[type]} offset int    [description]
 * @param  {[type]} limit int    [description]
 */
func GetAll(string, orderBy string, offset, limit int) *gorm.DB {
	db := database.GetGdb()
	if len(orderBy) > 0 {
		db.Order(orderBy + "desc")
	} else {
		db.Order("created_at desc")
	}
	if len(string) > 0 {
		db.Where("name LIKE  ?", "%"+string+"%")
	}
	if offset > 0 {
		db.Offset((offset - 1) * limit)
	}
	if limit > 0 {
		db.Limit(limit)
	}
	return db
}

// 清除系统基础数据
func DelAllData() {
	database.GetGdb().Unscoped().Delete(&OauthToken{})
	database.GetGdb().Unscoped().Delete(&Permission{})
	database.GetGdb().Unscoped().Delete(&Role{})
	database.GetGdb().Unscoped().Delete(&User{})
	database.GetGdb().Unscoped().Delete(&Admin{})
	database.GetGdb().Exec("DELETE FROM casbin_rule;")
}

// 更新数据
func Update(v, d interface{}) error {
	if err := database.GetGdb().Model(v).Updates(d).Error; err != nil {
		return err
	}
	return nil
}

// 获取用户角色
func GetRolesForUser(uid uint) []string {
	uids, err := database.GetEnforcer().GetRolesForUser(strconv.FormatUint(uint64(uid), 10))
	if err != nil {
		color.Red(fmt.Sprintf("GetRolesForUser 错误: %v", err))
		return []string{}
	}

	return uids
}

// 获取角色权限 （注意不是用户权限）
func GetPermissionsForUser(uid uint) [][]string {
	return database.GetEnforcer().GetPermissionsForUser(strconv.FormatUint(uint64(uid), 10))
}

// 自动创建表结构
func AutoMigrate() {
	database.GetGdb().AutoMigrate(
		&User{},
		&OauthToken{},
		&Role{},
		&Permission{},
		&Admin{},
		&App{},
		&Tenancy{},
	)
}

// 删除数据表
func DropTables() {
	database.GetGdb().DropTable("users", "roles", "permissions", "oauth_tokens", "casbin_rule", "admins", "apps", "tenancies")
}
