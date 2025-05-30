package serializer

type RegisterRequestSerializer struct {
	Username     string `json:"username"`
	Password     string `json:"password"`
	CorePassword string `json:"core_password"`
}

type LoginRequestSerializer struct {
	Username     string `json:"username"`
	Password     string `json:"password"`
	CorePassword string `json:"core_password"`
}

type InsertAccountRequestSerializer struct {
	Website  string `json:"website" binding:"required"`
	Account  string `json:"account"`
	Password string `json:"password" binding:"required"`
}
