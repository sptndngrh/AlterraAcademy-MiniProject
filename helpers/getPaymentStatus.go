package helpers

func GetPaymentStatusString(paymentStatus bool) string {
	if paymentStatus {
		return "Lunas"
	}
	return "Belum Lunas"
}