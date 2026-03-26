package domain

// DefaultPayrollConcepts retorna los conceptos base para nómina
func DefaultPayrollConcepts() []PayrollConcept {
	return []PayrollConcept{
		{Code: ConceptBaseSalary, Name: "Salario Base", Type: PayrollTypeEarning, IsMandatory: true},
		{Code: ConceptTransport, Name: "Auxilio Transporte", Type: PayrollTypeEarning, IsMandatory: false},
		{Code: ConceptHousing, Name: "Auxilio Vivienda", Type: PayrollTypeEarning, IsMandatory: false},
		{Code: ConceptOvertime, Name: "Horas Extra", Type: PayrollTypeEarning, IsMandatory: false},
		{Code: ConceptBonus, Name: "Bonificación", Type: PayrollTypeEarning, IsMandatory: false},
		{Code: ConceptHealth, Name: "Aporte Salud", Type: PayrollTypeDeduction, IsMandatory: true},
		{Code: ConceptPension, Name: "Aporte Pensión", Type: PayrollTypeDeduction, IsMandatory: true},
		{Code: ConceptTax, Name: "Retención Impuesto", Type: PayrollTypeDeduction, IsMandatory: false},
		{Code: ConceptOtherDeduction, Name: "Otra Deducción", Type: PayrollTypeDeduction, IsMandatory: false},
		{Code: ConceptHealthEmployer, Name: "Aporte Salud Empleador", Type: PayrollTypeEmployerContribution, IsMandatory: true},
		{Code: ConceptPensionEmployer, Name: "Aporte Pensión Empleador", Type: PayrollTypeEmployerContribution, IsMandatory: true},
		{Code: ConceptParafiscales, Name: "Parafiscales", Type: PayrollTypeEmployerContribution, IsMandatory: true},
	}
}
