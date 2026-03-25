package dto

import "github.com/platonso/hrmate/internal/service/form/model"

func ToFormCreateInput(req *FormCreateRequest) *model.FormCreateInput {
	return &model.FormCreateInput{
		Title:       req.Title,
		Description: req.Description,
		StartDate:   req.StartDate,
		EndDate:     req.EndDate,
	}
}
