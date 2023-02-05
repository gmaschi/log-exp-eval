package mock_generator

import (
	_ "github.com/golang/mock/mockgen/model"
)

//go:generate mockgen -package mockedexpstore -destination ../../services/datastore/postgresql/exp/mocks/mock_store.go github.com/gmaschi/log-exp-eval/internal/services/datastore/postgresql/exp Store
//go:generate mockgen -package mockedeval -destination ../../services/eval/mocks/mock_eval.go github.com/gmaschi/log-exp-eval/internal/services/eval Evaluator
// internal/services/eval/eval.go
