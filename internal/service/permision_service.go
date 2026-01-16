package service

import (
	"context"
	"errors"

	"github.com/arrase21/crm-users/internal/domain"
)

type PermissionService struct {
	userRoleRepo domain.UserRoleRepo
	roleRepo     domain.RoleRepo
}

func NewPermissionService(
	userRoleRepo domain.UserRoleRepo,
	roleRepo domain.RoleRepo,
) *PermissionService {
	return &PermissionService{
		userRoleRepo: userRoleRepo,
		roleRepo:     roleRepo,
	}
}

func (s *PermissionService) UserHasPermission(ctx context.Context, userID uint, resource, action string) (bool, error) {
	if userID == 0 {
		return false, errors.New("invalid user error")
	}
	if resource == "" || action == "" {
		return false, errors.New("resource and action are required")
	}
	roles, err := s.userRoleRepo.GetUserRoles(ctx, userID)
	if err != nil {
		return false, err
	}
	for _, role := range roles {
		if !role.IsActive {
			continue
		}
		for _, rp := range role.RolePermissions {
			if !rp.Action.IsActive || !rp.Action.Resource.IsActive {
				continue
			}
			if rp.Action.Resource.Name == resource && rp.Action.Action == action {
				return true, nil
			}
		}
	}
	return false, nil
}

func (s *PermissionService) UserHasRole(ctx context.Context, userID uint, roleName string) (bool, error) {
	if userID == 0 {
		return false, errors.New("invalid user id")
	}
	if roleName == "" {
		return false, errors.New("role name is required")
	}
	roles, err := s.userRoleRepo.GetUserRoles(ctx, userID)
	if err != nil {
		return false, err
	}
	for _, role := range roles {
		if role.Name == roleName && role.IsActive {
			return true, nil
		}
	}
	return false, nil
}

func (s *PermissionService) UserHasAnyPermission(ctx context.Context, userID uint, permissions map[string]string) (bool, error) {
	for resource, action := range permissions {
		hasPermission, err := s.UserHasPermission(ctx, userID, resource, action)
		if err != nil {
			return false, err
		}
		if hasPermission {
			return true, nil
		}
	}
	return false, nil
}

func (s *PermissionService) GetUserPermissions(ctx context.Context, userID uint) ([]string, error) {
	if userID == 0 {
		return nil, errors.New("invalid user id")
	}
	roles, err := s.userRoleRepo.GetUserRoles(ctx, userID)
	if err != nil {
		return nil, err
	}
	permissionsMap := make(map[string]bool)

	for _, role := range roles {
		if !role.IsActive {
			continue
		}
		for _, rp := range role.RolePermissions {
			if rp.Action.IsActive && rp.Action.Resource.IsActive {
				slug := rp.Action.Resource.Name + "." + rp.Action.Action
				permissionsMap[slug] = true
			}
		}
	}
	permissions := make([]string, 0, len(permissionsMap))
	for perm := range permissionsMap {
		permissions = append(permissions, perm)
	}
	return permissions, nil
}

func (s *PermissionService) AssignRoleToUser(ctx context.Context, userID, roleID uint) error {
	return s.userRoleRepo.AssignRole(ctx, userID, roleID)
}

// RevokeRoleFromUser revoca un rol de un usuario
func (s *PermissionService) RevokeRoleFromUser(ctx context.Context, userID, roleID uint) error {
	return s.userRoleRepo.RevokeRole(ctx, userID, roleID)
}

// AssignPermissionToRole asigna un permiso a un rol
func (s *PermissionService) AssignPermissionToRole(ctx context.Context, roleID, actionID uint) error {
	return s.roleRepo.AssignPermission(ctx, roleID, actionID)
}

// RevokePermissionFromRole revoca un permiso de un rol
func (s *PermissionService) RevokePermissionFromRole(ctx context.Context, roleID, actionID uint) error {
	return s.roleRepo.RevokePermission(ctx, roleID, actionID)
}
