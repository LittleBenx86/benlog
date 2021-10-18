package controller

import (
	"github.com/gofiber/fiber/v2"
)

type IPostController interface {
	Create() fiber.Handler
}

type IGetDetailController interface {
	Detail() fiber.Handler
}

type IGetListController interface {
	List() fiber.Handler
}

type IGetPageController interface {
	Page() fiber.Handler
}

type IGetController interface {
	IGetDetailController
	IGetListController
	IGetPageController
}

type IPatchController interface {
	Modify() fiber.Handler
}

type IDeleteOneController interface {
	Delete() fiber.Handler
}

type IDeleteByController interface {
	DeleteByIdentify() fiber.Handler
}

type IDeleteBatchController interface {
	BatchDelete() fiber.Handler
}

type IDeleteController interface {
	IDeleteOneController
	IDeleteByController
	IDeleteBatchController
}
