package html

// Setup all routes.
func (ar *Router) initRoutes() {
	if ar.Router == nil {
		panic("Empty HTML router")
	}

	// If serve new web static files than just set web handlers and return
	if ar.staticFilesStorageSettings.ServeNewWeb {
		appHandler := ar.staticFilesStorage.WebHandlers()
		ar.Router.PathPrefix(`/`).Handler(appHandler.AppHandler).Methods("GET")

		return
	}
}
