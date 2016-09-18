package licensecode

const TypeIndividual string = "A"
const TypeBusiness string = "B"
const TypeSite string = "E"

func IsValidType(t string) bool {
	return t == TypeIndividual || t == TypeBusiness || t == TypeSite
}
