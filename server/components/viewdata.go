package components

import (
	"jinovatka/entities"
	"jinovatka/utils"
)

type AdminViewData struct {
	Pagination utils.Pagination
}

func NewAdminViewData(records []*entities.Seed, requstedPage, noPages int) *AdminViewData {
	pagination := utils.NewPagination(requstedPage, len(records)/utils.DefaultLinesPerPage, len(records))
	return &AdminViewData{
		Pagination: pagination,
	}
}

type GroupViewData struct {
	Group *entities.SeedsGroup
}

func NewGroupViewData(seedsGroup *entities.SeedsGroup) *GroupViewData {
	return &GroupViewData{
		Group: seedsGroup,
	}
}
