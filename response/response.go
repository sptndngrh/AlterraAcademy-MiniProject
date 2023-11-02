package response

type UserResponse struct {
	Username string `json:"username"`
	Nama     string `json:"nama"`
}

func NewUserResponse(username, nama string) UserResponse {
	return UserResponse{
		Username: username,
		Nama:     nama,
	}
}

type OwnerResponse struct {
	Username string `json:"username"`
	Nama     string `json:"nama"`
}

func NewOwnerResponse(username, nama string) OwnerResponse {
	return OwnerResponse{
		Username: username,
		Nama:     nama,
	}
}
