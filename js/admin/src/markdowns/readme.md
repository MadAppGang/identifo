how create name file for markdowns

1. filename must be based on url where you want to see documentation
2. if you have nested routes with params you have to add routes in paramRoutes object
3. the name of file for route with param have to start with param name
4. first h1 in md file is excluded and used for button title (ex. #button title)
examples 
    with nested url param + tabs
        url: "/management/applications/c50p9h86n88nubgpm5kg?edit_app_group=tokens"
        filename: "appid_edit_app_group_tokens.md"
        pathtofile:"/management/applications/appid_edit_app_group_tokens.md"
    without param + tabs
        url: "/management?server_group=general"
        filename: "server_group_general.md"
        pathtofile:"/management/server_group_general.md"
    without params and tabs
        url: "/management/applications"
        filename: "index.md"
        pathtofile:"/management/applications/index.md"
