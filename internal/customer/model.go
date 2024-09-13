package customer

type Data struct {
	BrandName   string `json:"brandName"`
	CivilName   string `json:"civilName"`
	CPF         string `json:"cpfNumber"`
	CompanyInfo struct {
		CNPJ string `json:"cnpjNumber"`
		Name string `json:"name"`
	} `json:"companyInfo"`
}
