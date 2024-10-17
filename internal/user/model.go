package user

type User struct {
	UserName     string
	Email        string
	CPF          string
	Name         string
	CompanyCNPJs []string
}

type Company struct {
	Name string
	CNPJ string
}
