package user

// TODO: remove from here to main.
var companyA = Company{
	Name: "A Business",
	CNPJ: "27737785000136",
}

var userBob = User{
	UserName:  "bob@mail.com",
	CPF:       "78628584099",
	Name:      "Mr. Bob",
	Companies: []string{companyA.CNPJ},
}
