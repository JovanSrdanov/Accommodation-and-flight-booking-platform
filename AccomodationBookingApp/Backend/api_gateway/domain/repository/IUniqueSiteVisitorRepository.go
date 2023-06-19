package repository

import (
	"api_gateway/domain/model"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type IUniqueSiteVisitorRepository interface {
	CreateUniqueVisitor(uniqueVisitor *model.UniqueVisitor) (primitive.ObjectID, error)
	GetVisitorByIpAndBrowser(ipAddress, browser string) (*model.UniqueVisitor, error)
}
