package model

type Advertisement struct {
	Name  string
	Price string
	Url   string
}

func NewAdvertisement(name string, price string, url string) *Advertisement {
	ad := new(Advertisement)
	ad.Name = name
	ad.Price = price
	ad.Url = url
	return ad
}
