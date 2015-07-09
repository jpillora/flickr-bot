package flickr

type User struct {
	c        *Client
	ID       string `json:"nsid"`
	Username string
	Fullname string
	Token    string
	Perms    string
}

func (u *User) userArgs(args Args) Args {
	if args == nil {
		args = Args{}
	}
	if _, ok := args["user_id"]; !ok {
		args["user_id"] = u.ID
	}
	args["auth_token"] = u.Token
	return args
}

func (u *User) Test(method string, args Args) {
	u.c.Test(method, u.userArgs(args))
}

func (u *User) Do(method string, args Args) ([]byte, error) {
	return u.c.Do(method, u.userArgs(args))
}
