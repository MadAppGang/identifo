package mongo

import mgo "gopkg.in/mgo.v2"

//NewDB creates new database connection
func NewDB(conn string, db string) (*DB, error) {
	m := DB{}
	session, err := mgo.Dial(conn)
	if err != nil {
		return nil, err
	}
	m.S = session
	m.DB = db
	//TODO: ensure all indexes
	return &m, nil
}

//DB is database connection structure
type DB struct {
	S  *mgo.Session
	DB string
}

//Session creates new database connection session
//don't dorget to close it, defer is commonly used way to do it
func (db *DB) Session(col string) *Session {
	s := Session{}
	s.S = db.S.Clone()
	s.C = s.S.DB(db.DB).C(col)
	return &s
}

//Close closes database connection
func (db *DB) Close() {
	db.S.Close()
}

//Session implements one single session connection to database
type Session struct {
	C *mgo.Collection
	S *mgo.Session
}

//Close closes connection to current session
func (s *Session) Close() {
	s.S.Close()
}
