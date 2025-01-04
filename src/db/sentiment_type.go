package db

import "fmt"

type SentimentType string

const (
	NeutralSentimentType     SentimentType = "neutral"
	PositiveSentimentType    SentimentType = "positive"
	NegativeSentimentType    SentimentType = "negative"
	ExcitingSentimentType    SentimentType = "exctiting"
	SupportiveSentimentType  SentimentType = "supportive"
	InformativeSentimentType SentimentType = "informative"
	ProvocativeSentimentType SentimentType = "provocative"
	ViolentSentimentType     SentimentType = "violent"
	AbuseSentimentType       SentimentType = "abuse"
	ScamSentimentType        SentimentType = "scam"
)

var sentimentMap = map[string]SentimentType{
	string(NeutralSentimentType):     NeutralSentimentType,
	string(PositiveSentimentType):    PositiveSentimentType,
	string(NegativeSentimentType):    NegativeSentimentType,
	string(ExcitingSentimentType):    ExcitingSentimentType,
	string(SupportiveSentimentType):  SupportiveSentimentType,
	string(InformativeSentimentType): InformativeSentimentType,
	string(ProvocativeSentimentType): ProvocativeSentimentType,
	string(ViolentSentimentType):     ViolentSentimentType,
	string(AbuseSentimentType):       AbuseSentimentType,
	string(ScamSentimentType):        ScamSentimentType,
}

func ParseSentimentType(s string) (SentimentType, error) {
	if sentiment, ok := sentimentMap[s]; ok {
		return sentiment, nil
	}
	return "", fmt.Errorf("invalid sentiment type: %q", s)
}
