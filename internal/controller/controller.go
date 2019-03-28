package controller

type Controller struct {
	Storage Storage
}

func New(storage Storage) *Controller {
	return &Controller{Storage: storage}
}
