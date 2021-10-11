// use this router mapper to exclude route params for useMarkdown hook
// You have to exclude nasted routes for correct markdown file name
export const paramRoutes = {
  newApp: '/management/applications/new',
  editApp: '/management/applications/:appid',
  newUser: '/management/users/new',
  editUser: '/management/users/:userid',
};
