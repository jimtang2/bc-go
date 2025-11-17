package main

import (
	"context"
	"net/http"
)

type HealthDependency interface {
	Health(context.Context) error
}

func HealthHandler(deps ...HealthDependency) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		for _, dep := range deps {
			if err := dep.Health(r.Context()); err != nil {
				http.Error(w, "unhealthy", http.StatusServiceUnavailable)
				return
			}
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok"))
	}
}
