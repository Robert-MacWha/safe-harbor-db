package user

// User represents a user document in Firestore.
type User struct {
	FirstName       string `firestore:"firstName"`
	LastName        string `firestore:"lastName"`
	Email           string `firestore:"email"`
	Title           string `firestore:"title"`
	CompanyName     string `firestore:"companyName"`
	Street          string `firestore:"street"`
	City            string `firestore:"city"`
	ZipCode         string `firestore:"zipCode"`
	Country         string `firestore:"country"`
	Region          string `firestore:"region"`
	IndividualKYCId string `firestore:"individualKYCId"`
	BusinessKYCId   string `firestore:"businessKYCId"`
}
