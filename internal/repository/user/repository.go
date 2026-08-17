package repository

type IUserRepository interface {
	List()
	Get(column string, value any)
	Create()
	Update()
	Delete()
}
