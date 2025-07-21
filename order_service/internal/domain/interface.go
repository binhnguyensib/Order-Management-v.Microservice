package domain

type OrderUsecase interface {
	Creat()
	GetByID()
	GetByCustomerID()
	Update()
	UpdateStatus()
}

type OrderRepository interface {
}
