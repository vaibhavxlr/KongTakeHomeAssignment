package DTOs

type ListServicesResp struct {
	Services    []Service   `json:"services"`
	SortOrder   SortOrder   `json:"sortOrder"`
	PageDetails PageDetails `json:"pageDetails"`
}
type Service struct {
	ID            string    `json:"id"`
	Name          string    `json:"name"`
	Info          string    `json:"info"`
	VersionsCount int       `json:"versionsCount,omitempty"`
	Versions      []Version `json:"versions,omitempty"`
}
type SortOrder struct {
	AZ int `json:"A-Z"`
	ZA int `json:"Z-A"`
}
type PageDetails struct {
	Curr  int `json:"curr"`
	Total int `json:"total"`
	Count int `json:"count"`
}
type Version struct {
	VerName string `json:"verName"`
	VerInfo string `json:"verInfo"`
	Changes string `json:"changes"`
}
