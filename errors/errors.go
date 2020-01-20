package customerrors



type InvalidCredentialsError struct {
	S string
}


func (err * InvalidCredentialsError) Error() string {
	return err.S
}


type TokenExpired struct {
	S string
}

func (err * TokenExpired) Error() string {
	return err.S
}