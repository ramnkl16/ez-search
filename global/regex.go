package global

import (
	"regexp"
)

const (
	Email          = `^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`
	Mobile         = `^\s*(?:\+?(\d{1,3}))?[-. (]*(\d{3})[-. )]*(\d{3})[-. ]*(\d{4})(?: *x(\d+))?\s*$`
	SlotID         = `^([0-9]{12}-[0-9]{1,3})$`
	OperationHours = `^[0-9]{1,4}-[0-9]{1,4}([\|]?[0-9]{1}=([0-9]{1,4}-[0-9]{1,4},?){1,24}){0,7}$`
	Price          = `^([a-zA-z0-9]{3,36}-[0-9]{1,4},?){1,}$`
)

var (
	RegexParseDate             = regexp.MustCompile(`\{([^}]+)\}`)                                      //to parse this pattern indexName{2006-01-02}
	RegExkeyWords              = regexp.MustCompile("select |from |where |limit |since |facets |sort ") //for query parsing
	RegexParseHasCapitalLetter = regexp.MustCompile("[[:upper:]]+")
)

func IsValidEmail(email string) bool {
	//fmt.Println("IsValidEmail", email)
	re := regexp.MustCompile(Email)
	return re.MatchString(email)
}

func IsValidPhoneNumber(phone string) bool {
	//fmt.Println("IsValidPhoneNumber", phone)
	//re := regexp.MustCompile(Mobile)
	return true //re.MatchString(phone)
}

func IsValidSlotID(id string) bool {
	re := regexp.MustCompile(SlotID)
	return re.MatchString(id)
}

func IsValidOperationHours(opHour string) bool {
	re := regexp.MustCompile(OperationHours)
	return re.MatchString(opHour)
}

func IsValidPrice(price string) bool {
	re := regexp.MustCompile(Price)
	return re.MatchString(price)
}
