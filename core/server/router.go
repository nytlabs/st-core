package server

import (
	"net/http"

	"github.com/gorilla/mux"
)

func (s *Server) NewRouter() *mux.Router {
	type Route struct {
		Name        string
		Pattern     string
		Method      string
		HandlerFunc http.HandlerFunc
	}

	routes := []Route{
		Route{
			"UpdateSocket",
			"/updates",
			"GET",
			s.UpdateSocketHandler,
		},
		Route{
			"GroupIndex",
			"/groups",
			"GET",
			s.GroupIndexHandler,
		},
		Route{
			"GroupCreate",
			"/groups",
			"POST",
			s.GroupCreateHandler,
		},
		Route{
			"GroupDelete",
			"/groups/{id}",
			"DELETE",
			s.GroupDeleteHandler,
		},
		Route{
			"BlockIndex",
			"/blocks",
			"GET",
			s.BlockIndexHandler,
		},
		Route{
			"BlockCreate",
			"/blocks",
			"POST",
			s.BlockCreateHandler,
		},
		Route{
			"BlockDelete",
			"/blocks/{id}",
			"DELETE",
			s.BlockDeleteHandler,
		},
		Route{
			"BlockModifyName",
			"/blocks/{id}/name",
			"PUT",
			s.BlockModifyNameHandler,
		},
		Route{
			"BlockModifyRoute",
			"/blocks/{id}/routes/{index}",
			"PUT",
			s.BlockModifyRouteHandler,
		},
		Route{
			"BlockModifyGroup",
			"/blocks/{id}/group",
			"PUT",
			s.BlockModifyGroupHandler,
		},
		Route{
			"ConnectionIndex",
			"/connections",
			"GET",
			s.ConnectionIndexHandler,
		},
		Route{
			"ConnectionCreate",
			"/connections",
			"POST",
			s.ConnectionCreateHandler,
		},
		Route{
			"ConnectionDelete",
			"/connections/{id}",
			"DELETE",
			s.ConnectionDeleteHandler,
		},
		Route{
			"SourceCreate",
			"/sources",
			"POST",
			s.SourceCreateHandler,
		},
		Route{
			"SourceIndex",
			"/sources",
			"GET",
			s.SourceIndexHandler,
		},
		Route{
			"SourceModify",
			"/sources/{id}",
			"/PUT",
			s.SourceModifyHandler,
		},
		Route{
			"Source",
			"/sources/{id}",
			"GET",
			s.SourceHandler,
		},
		Route{
			"Source",
			"/sources/{id}",
			"DELETE",
			s.SourceDeleteHandler,
		},
		Route{
			"Library",
			"/library",
			"GET",
			s.LibraryHandler,
		},
	}
	router := mux.NewRouter().StrictSlash(true)
	for _, route := range routes {
		var handler http.Handler

		handler = route.HandlerFunc
		handler = Logger(handler, route.Name)

		router.
			Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(handler)
	}

	router.PathPrefix("/").Handler(http.FileServer(http.Dir("./static/")))

	return router

}
