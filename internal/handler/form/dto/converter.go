package dto

import (
	"github.com/platonso/hrmate/internal/domain"
	"github.com/platonso/hrmate/internal/service/form/model"
)

func ToFormCreateInput(req FormCreateRequest) model.FormCreateInput {
	return model.FormCreateInput{
		Title:       req.Title,
		Description: req.Description,
		StartDate:   req.StartDate,
		EndDate:     req.EndDate,
	}
}

func ToFormResponse(form *domain.Form) FormResponse {
	return FormResponse{
		ID:          form.ID,
		UserID:      form.UserID,
		Title:       form.Title,
		Description: form.Description,
		StartDate:   form.StartDate,
		EndDate:     form.EndDate,
		CreatedAt:   form.CreatedAt,
		ReviewedAt:  form.ReviewedAt,
		Status:      string(form.Status),
		Comment:     form.Comment,
	}
}

func ToFormResponses(forms []domain.Form) []FormResponse {
	if len(forms) == 0 {
		return []FormResponse{}
	}
	responses := make([]FormResponse, len(forms))
	for i := range forms {
		responses[i] = ToFormResponse(&forms[i])
	}
	return responses
}

func ToUserResponse(user *domain.User) UserResponse {
	return UserResponse{
		ID:        user.ID,
		Role:      string(user.Role),
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Position:  user.Position,
		Email:     user.Email,
		IsActive:  user.IsActive,
	}
}

func ToFormsWithUserResponses(data []model.FormsWithUser) []FormsWithUserResponse {
	if len(data) == 0 {
		return []FormsWithUserResponse{}
	}
	responses := make([]FormsWithUserResponse, len(data))
	for i := range data {
		responses[i] = FormsWithUserResponse{
			User:  ToUserResponse(&data[i].User),
			Forms: ToFormResponses(data[i].Forms),
		}
	}
	return responses
}
