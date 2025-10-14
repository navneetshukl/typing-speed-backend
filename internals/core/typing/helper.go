package typing

func TypingDataValid(data *TypingData)error{
	if data.UserId==""{
		return ErrInvalidUser
	}
	return nil
	
}