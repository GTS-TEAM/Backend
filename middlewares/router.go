package middlewares

type SecurityRouter struct {
	Admin []string
	User  []string
}

func GetSecurityRouters() SecurityRouter {
	return SecurityRouter{
		Admin: []string{
			"/api/product",
			"/api/category",
			"/api/reviews",
			"/api/metadata",
			"/api/stock",
		},
		User: []string{},
	}
}
