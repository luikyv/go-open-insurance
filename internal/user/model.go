package user

type User struct {
	UserName  string
	CPF       string
	Name      string
	Companies []string
}

type Company struct {
	Name string
	CNPJ string
}

type Document struct {
	Identification string `json:"identification" binding:"required"`
	Rel            string `json:"rel" binding:"required"`
}

type Logged struct {
	Document Document `json:"document" binding:"required"`
}

type BusinessEntity struct {
	Document Document `json:"document" binding:"required"`
}
