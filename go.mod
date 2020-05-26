module github.com/open-kingfisher/king-kf

go 1.14

require (
	github.com/DeanThompson/ginpprof v0.0.0-20190408063150-3be636683586
	github.com/dgrijalva/jwt-go v3.2.0+incompatible
	github.com/gin-gonic/gin v1.6.2
	github.com/open-kingfisher/king-utils v0.0.0-20200526073742-a2d1d095795b
	golang.org/x/crypto v0.0.0-20190820162420-60c769a6c586
	golang.org/x/sync v0.0.0-20200317015054-43a5402ce75a
)

replace (
	k8s.io/api => k8s.io/api v0.17.3
	k8s.io/apimachinery => k8s.io/apimachinery v0.17.3
	k8s.io/client-go => k8s.io/client-go v0.17.3
	k8s.io/metrics => k8s.io/metrics v0.17.3
)
