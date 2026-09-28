package simgos

type antrolTokenResponse struct {
	Metadata struct {
		Code int `json:"code"`
	} `json:"metadata"`

	Response struct {
		Token string `json:"token"`
	} `json:"response"`
}

type antrolCheckInRequest struct {
	KodeBooking string `json:"kodebooking"`
	Waktu       string `json:"waktu"`
}
