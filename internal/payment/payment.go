package payment

type Provider interface {
	Pay(amount float64) error
}

type CreditCardProvider struct {
	CardNumber string
	ExpiryDate string
	CVV        string
}

func (c *CreditCardProvider) Pay(amount float64) error {
	// Implement credit card payment logic here
	return nil
}

type PayPalProvider struct {
	Email    string
	Password string
}

func (p *PayPalProvider) Pay(amount float64) error {
	// Implement PayPal payment logic here
	return nil
}
