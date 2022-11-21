package visit

type Clinic struct {
	Id   int    `json:"Id"`
	Name string `json:"Name"`
}

type Doctor struct {
	Id   int    `json:"Id"`
	Name string `json:"Name"`
}

type Impediment struct {
	IsImpediment   bool   `json:"IsImpediment"`
	ImpedimentText string `json:"ImpedimentText,omitempty"`
}

type VisitDate struct {
	StartDateTime string `json:"StartDateTime"`
	FormattedDate string `json:"FormattedDate"`
	EndDateTime   string `json:"EndDateTime"`
}

type VisitTerm struct {
	ServiceId                 int           `json:"ServiceId"`
	Clinic                    Clinic        `json:"Clinic"`
	Impediment                Impediment    `json:"Impediment,omitempty"`
	VisitDate                 VisitDate     `json:"VisitDate"`
	IsFree                    bool          `json:"IsFree"`
	RoomId                    int           `json:"RoomId"`
	ScheduleId                int           `json:"ScheduleId"`
	ReferralRequiredByService bool          `json:"ReferralRequiredByService"`
	ReferralRequiredByProduct bool          `json:"ReferralRequiredByProduct"`
	PayerDetailsList          []interface{} `json:"PayerDetailsList"`
}

type DateRange struct {
	FromDate string `json:"FromDate"`
	ToDate   string `json:"ToDate"`
}

type VisitTermsResponse struct {
	VisitTerms      []VisitTerm `json:"AvailableVisitsTermPresentation"`
	SearchDateRange DateRange         `json:"SearchDateRange"`
}