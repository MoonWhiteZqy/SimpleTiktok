package service

type UserService interface {
	UserRegisterSrv(string, string) (int64, error)
	UserLoginSrv(string, string) (int64, error)
	UserBaseInfoSrv(int64) (string, error)
}
