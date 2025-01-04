package consts

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

var possibleSentiments []string

func PossibleSentiments() []string {
	if possibleSentiments != nil {
		return possibleSentiments
	}

	possible := make([]string, len(sentimentMap))
	i := 0
	for k := range sentimentMap {
		possible[i] = k
		i += 1
	}

	possibleSentiments = possible
	return possibleSentiments
}

func ParseSentimentType(s string) (SentimentType, error) {
	if sentiment, ok := sentimentMap[s]; ok {
		return sentiment, nil
	}
	return "", fmt.Errorf("invalid sentiment type: %q", s)
}
