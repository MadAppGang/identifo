package adminpanel

// Setup all routes for admin panel router.
func (apr *Router) initRoutes() {
	if apr.router == nil {
		return
	}

	handlers := apr.staticFilesStorage.AdminPanelHandlers()

	apr.router.PathPrefix("/config.json").Handler(handlers.ConfigHandler).Methods("GET")
	apr.router.PathPrefix("/src/").Handler(handlers.SrcHandler).Methods("GET")
	apr.router.PathPrefix("/management").Handler(handlers.ManagementHandler).Methods("GET")
	apr.router.PathPrefix("/").Handler(handlers.BuildHandler).Methods("GET")
}
