package expserver

import (
	"github.com/gin-gonic/gin"
	expcontroller "github.com/gmaschi/log-exp-eval/internal/controllers/expressions"
	authmid "github.com/gmaschi/log-exp-eval/internal/controllers/middlewares/auth-mid"
	expstore "github.com/gmaschi/log-exp-eval/internal/services/datastore/postgresql/exp"
	"github.com/gmaschi/log-exp-eval/internal/services/eval"
	"github.com/gmaschi/log-exp-eval/pkg/tools/config/env"
	_ "github.com/lib/pq"
)

type (
	// Server holds all the required fields regarding the expression server.
	Server struct {
		store         expstore.Store
		evaluator     eval.Evaluator
		expController *expcontroller.Controller
		Config        env.Config
		Router        *gin.Engine
	}
)

// New instantiates a Server. While initializing, it calls the setupRoutes method.
func New(config env.Config, store expstore.Store, ev eval.Evaluator) (*Server, error) {
	// TODO: implement any required validations.
	srv := &Server{
		store:         store,
		evaluator:     ev,
		expController: expcontroller.New(store, ev),
		Config:        config,
	}
	router := gin.Default()

	srv.setupRoutes(router)

	srv.Router = router
	return srv, nil
}

// setupRoutes defines the router for Server and ties each endpoint to the corresponding method
// from expcontroller.Controller.
func (f *Server) setupRoutes(router *gin.Engine) {
	v1 := router.Group("/v1")

	expGroup := v1.Group("/expressions")
	{
		expGroup.POST("", authmid.AuthMiddleware(), f.expController.Create)
		expGroup.GET("/:id", authmid.AuthMiddleware(), f.expController.Get)
		expGroup.DELETE("/:id", authmid.AuthMiddleware(), f.expController.Delete)
		expGroup.PATCH("", authmid.AuthMiddleware(), f.expController.Update)
		expGroup.GET("", authmid.AuthMiddleware(), f.expController.List)
	}

	evalGroup := v1.Group("/evaluate")
	{
		evalGroup.GET("/:id", authmid.AuthMiddleware(), f.expController.Evaluate)
	}
}

// Start starts the server at the given address.
func (f *Server) Start(address string) error {
	return f.Router.Run(address)
}
