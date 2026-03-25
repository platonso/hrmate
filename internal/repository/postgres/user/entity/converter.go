package entity

import "github.com/platonso/hrmate/internal/domain"

func ToUserRecord(u domain.User) UserRecord {
	record := UserRecord{
		ID:             u.ID,
		Role:           string(u.Role),
		FirstName:      u.FirstName,
		LastName:       u.LastName,
		Position:       u.Position,
		Email:          u.Email,
		HashedPassword: u.HashedPassword,
		IsActive:       u.IsActive,
	}
	return record
}

func ToDomainUser(record UserRecord) domain.User {
	user := domain.User{
		ID:             record.ID,
		Role:           domain.Role(record.Role),
		FirstName:      record.FirstName,
		LastName:       record.LastName,
		Position:       record.Position,
		Email:          record.Email,
		HashedPassword: record.HashedPassword,
		IsActive:       record.IsActive,
	}
	return user
}

func ToDomainUsers(records []UserRecord) []domain.User {
	if len(records) == 0 {
		return []domain.User{}
	}

	users := make([]domain.User, 0, len(records))
	for _, rec := range records {
		user := ToDomainUser(rec)
		users = append(users, user)
	}
	return users
}
