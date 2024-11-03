package api

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/bsontype"
)

func (q QuoteCustomerData) MarshalBSONValue() (bsontype.Type, []byte, error) {
	return bson.MarshalValue(q.union)
}

func (q *QuoteCustomerData) UnmarshalBSONValue(t bsontype.Type, data []byte) error {
	rv := bson.RawValue{Type: t, Value: data}
	return rv.Unmarshal(&q.union)
}

// TODO: Change for business.
func (q QuoteCustomerData) ToQuoteCustomer() QuoteCustomer {
	personalCustomerData, _ := q.AsQuotePersonalCustomerData()
	personalCustomer := QuotePersonalCustomer{
		Identification:           personalCustomerData.IdentificationData,
		Qualification:            personalCustomerData.QualificationData,
		ComplimentaryInformation: personalCustomerData.ComplimentaryInformationData,
	}

	quoteCustomer := QuoteCustomer{}
	_ = quoteCustomer.FromQuotePersonalCustomer(personalCustomer)
	return quoteCustomer
}
