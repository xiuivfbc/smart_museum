package dao

import (
	"group_ten_server/model"

	"gorm.io/gorm"
)

// CreateTicket 新建工单
func CreateTicket(ticket *model.Ticket) error {
	return db.Create(ticket).Error
}

// ListTickets 查询所有工单
func ListTickets() ([]model.Ticket, error) {
	var tickets []model.Ticket
	err := db.Find(&tickets).Error
	return tickets, err
}

// DeleteTicketByID 删除工单
func DeleteTicketByID(id int) error {
	return db.Delete(&model.Ticket{}, id).Error
}

// UpdateTicket 更新工单
func UpdateTicket(ticket *model.Ticket) error {
	return db.Save(ticket).Error
}

// GetTicketByID 查询单个工单
func GetTicketByID(db *gorm.DB, id int) (*model.Ticket, error) {
	var ticket model.Ticket
	err := db.First(&ticket, id).Error
	return &ticket, err
}
