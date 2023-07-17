package terms

type Info struct {
	TermsForService `json:"termsForService"`
}

type TermsForService struct {
	ServiceVariantID int            `json:"serviceVariantId"`
	TermsForDays     []TermsForDays `json:"termsForDays"`
}

type TermsForDays struct {
	Day   string `json:"day"`
	Terms []Term `json:"terms"`
}

type Term struct {
	Doctor   `json:"doctor"`
	From     string `json:"dateTimeFrom"`
	To       string `json:"dateTimeTo"`
	ClinicID int    `json:"clinicId"`
	Clinic   string `json:"clinic"`
}

type Doctor struct {
	ID            int    `json:"id"`
	GenderID      int    `json:"genderId"`
	AcademicTitle string `json:"academicTitle"`
	FirstName     string `json:"firstName"`
	LastName      string `json:"lastName"`
}
