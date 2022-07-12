package main

import (
	"context"
	"net/http"
	"os"

	"github.com/kelseyhightower/envconfig"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/99designs/gqlgen/graphql"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/jaspeen/sample-users-app/config"
	"github.com/jaspeen/sample-users-app/db"
	"github.com/jaspeen/sample-users-app/graph"
	"github.com/jaspeen/sample-users-app/graph/generated"
	"github.com/jaspeen/sample-users-app/services"
	_ "github.com/lib/pq"
)

func allowAllCorsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Headers", "*")
		next.ServeHTTP(w, r)
	})
}

func main() {
	if len(os.Args) > 1 {
		envconfig.Usage("SMPL", &config.C)
		os.Exit(1)
	}

	config.InitConfig("SMPL")

	// set global log level
	l, err := zerolog.ParseLevel(config.C.LogLevel)
	if err != nil {
		log.Fatal().Err(err)
	}
	zerolog.SetGlobalLevel(l)
	if config.C.PrettyLog {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	}

	db, err := db.InitDB(config.C.DbConnect, l == zerolog.DebugLevel)
	if err != nil {
		log.Fatal().Msgf("Cannot connect to db: %v\n", err)

	}
	c := generated.Config{Resolvers: &graph.Resolver{DB: db}}

	// implement @admin directive to guard admin-only queries and mutations
	c.Directives.Admin = func(ctx context.Context, obj interface{}, next graphql.Resolver) (interface{}, error) {
		user := graph.UserFromContext(ctx)
		if user == nil || !user.Admin {
			return nil, services.Err_UNAUTHORIZED
		}

		return next(ctx)
	}

	srv := handler.NewDefaultServer(generated.NewExecutableSchema(c))

	http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	http.Handle("/query", allowAllCorsMiddleware(graph.AuthMiddleware()(srv)))

	log.Printf("connect to http://%s/ for GraphQL playground", config.C.Listen)
	log.Fatal().Err(http.ListenAndServe(config.C.Listen, nil))
}
