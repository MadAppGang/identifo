package config

// Importer has a set of helper functions to import data
// Now it uses to import dummy data for isolated tests
// But could be used to import static data for stateless deployment
// more details in docs
// TODO: implement stateless static data import scheme

// // ImportApps imports apps from file.
// func (s *Server) ImportApps(filename string) error {
// 	data, err := dataFromFile(filename)
// 	if err != nil {
// 		return err
// 	}
// 	return s.AppStorage().ImportJSON(data)
// }

// // ImportUsers imports users from file.
// func (s *Server) ImportUsers(filename string) error {
// 	data, err := dataFromFile(filename)
// 	if err != nil {
// 		return err
// 	}
// 	return s.UserStorage().ImportJSON(data)
// }

// func dataFromFile(filename string) ([]byte, error) {
// 	file, err := os.Open(filename)
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer file.Close()
// 	return ioutil.ReadAll(file)
// }
