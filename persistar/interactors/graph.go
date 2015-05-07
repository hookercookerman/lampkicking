package interactors

import (
	"encoding/json"

	"github.com/hookercookerman/lampkicking/persistar"
)

type GraphInteractor struct {
	Adaptor persistar.GraphAdaptor
}

type RelationResult struct {
	Key   string                 `json:"path"`
	Value map[string]interface{} `json:"value"`
}

type RelationResults struct {
	Count   int               `json:"count"`
	Results []*RelationResult `json:"results"`
}

func (interactor *GraphInteractor) GetRelation(collection, key, relation string) ([]byte, error) {
	members, result, err := interactor.Adaptor.GetRelation(collection, key, relation)
	if err != nil {
		return nil, err
	}

	var results []*RelationResult

	for i, m := range members {
		var egg map[string]interface{}
		json.Unmarshal(result[i], &egg)
		results = append(results, &RelationResult{Key: m, Value: egg})
	}

	relationResults := &RelationResults{
		Count:   len(members),
		Results: results,
	}
	output, err := json.Marshal(relationResults)
	if err != nil {
		return nil, err
	}
	return output, nil
}

func (interactor *GraphInteractor) AddRelation(collection, key, relation, relatedCollection, relatedKey string) (bool, error) {
	return interactor.Adaptor.AddRelation(collection, key, relation, relatedCollection, relatedKey)
}
