package docmodels

type ResponseSuccess struct {
	Message string `json:"messsage" example:"success"`
	Status  int    `json:"status" example:"200"`
}

type ResponseBadRequest struct {
	Message string `json:"messsage" example:"Bad request message"`
	Status  int    `json:"status" example:"400"`
}
type ResponseUnauthorized struct {
	Message string `json:"messsage" example:"Request Unauthorized message"`
	Status  int    `json:"status" example:"401"`
	Expired bool   `json:"token_expired" example:"true"`
}

type ResponseForbidden struct {
	Message string `json:"messsage" example:"Forbidden message"`
	Status  int    `json:"status" example:"403"`
}

type ResponseNotFound struct {
	Message string `json:"messsage" example:"Not found message"`
	Status  int    `json:"status" example:"404"`
}

type ResponseInternalServerError struct {
	Message string `json:"messsage" example:"Internal Server Error"`
	Status  int    `json:"status" example:"500"`
}

type ResponseImageElement struct {
	Id         string `json:"id" example:"6006d3cc95f8c8e32d660c04"`
	UploadDate string `json:"uploadDate" swaggertype:"string" format:"date-time"`
}

type ResponseLoLHistory struct {
	Id        string `json:"id" example:"6006d3cc95f8c8e32d660c04"`
	Timestamp string `json:"timestamp" swaggertype:"integer" example:"1610811479544"`
	Win       bool   `json:"win" swaggertype:"boolean" example:"true"`
	Participates       []string   `json:"participates"`
}

type ResponseDateLog struct {
	Count    int      `json:"count" example:"4"`
	Year     int      `json:"year" example:"2021"`
	Month    int      `json:"month" example:"2"`
	Day      int      `json:"day" example:"12"`
	QueueId  int      `json:"queueid" example:"430"`
	Win      int      `json:"win" example:"2"`
}


type RequestBodyAnnouncePost struct {
	Title       string `json:"title" example:"Title"`
	Description string `json:"description" example:"This field is optional"`
	TargetDate  string `json:"target_date" swaggertype:"string" format:"date-time"`
}

type RequestBodyImagePost struct {
	Image string `json:"img" swaggertype:"string" format:"base64" example:"data:image/gif;base64,R0lGODlhAQABAAAAACH5BAEKAAEALAAAAAABAAEAAAICTAEAOw=="`
}

type RequestBodyStatePost struct {
	State string `json:"state" swaggertype:"string" example:"Working"`
}

type RequestBodyToken struct {
	Token string `json:"token" example:"riot-api-access"`
}

type RequestBodyPatchUser struct {
	Username     string `json:"username" example:"devsdk"`
	ProfileImage string `json:"profile_image_id" example:"6006d3cc95f8c8e32d660c04"`
}

type RequestEmpty struct {
}
