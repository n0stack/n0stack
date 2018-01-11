package repository

import (
	"github.com/n0stack/n0core/message"
	"github.com/n0stack/n0core/model"
	uuid "github.com/satori/go.uuid"
)

type MockRepository struct {
	PatchStoreNotification func(m *message.Notification)
	PatchDigModel          func(i *uuid.UUID, e string, d uint) (*model.Model, error)
}

func (mr MockRepository) StoreNotification(m *message.Notification) {
	mr.PatchStoreNotification(m)
}

func (mr MockRepository) DigModel(i *uuid.UUID, e string, d uint) (*model.Model, error) {
	return mr.PatchDigModel(i, e, d)
}
