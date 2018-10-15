package controller

type SessionDao struct {
	Context *Context
}

func NewSession(context *Context) *SessionDao {
	var session = &SessionDao{}
	session.Context = context

	return session
}

// func (this *SessionDao) Create(session *model.Session) *model.Session {
// 	timeUUID, _ := uuid.NewV4()
// 	session.SessionId = string(timeUUID.String())
// 	this.Context.DB.Create(session)

// 	return session
// }
