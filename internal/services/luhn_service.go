package services

type luhnService struct{}

func NewLuhnService() Validator {
	return &luhnService{}
}

func (s *luhnService) IsValid(number string) bool {
	for _, char := range number {
		if char < '0' || char > '9' {
			return false
		}
	}

	sum := 0
	isEven := false
	for i := len(number) - 1; i >= 0; i-- {
		digit := int(number[i] - '0')
		if isEven {
			digit *= 2
			if digit > 9 {
				digit -= 9
			}
		}

		sum += digit
		isEven = !isEven
	}

	return sum%10 == 0
}
