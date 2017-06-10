// black-box testing
package host_test

import (
	"net"
	"net/url"
	"testing"

	"github.com/kataras/iris"
	"github.com/kataras/iris/context"
	"github.com/kataras/iris/core/host"
	"github.com/kataras/iris/httptest"
)

func TestProxy(t *testing.T) {
	expectedIndex := "ok /"
	expectedAbout := "ok /about"
	unexpectedRoute := "unexpected"

	// proxySrv := iris.New()
	u, err := url.Parse("https://localhost")
	if err != nil {
		t.Fatalf("%v while parsing url", err)
	}

	// p := host.ProxyHandler(u)
	// transport := &http.Transport{
	// 	TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	// }
	// p.Transport = transport
	// proxySrv.Downgrade(p.ServeHTTP)
	// go proxySrv.Run(iris.Addr(":80"), iris.WithoutBanner, iris.WithoutInterruptHandler)

	go host.NewProxy(":80", u).ListenAndServe()

	app := iris.New()
	app.Get("/", func(ctx context.Context) {
		ctx.WriteString(expectedIndex)
	})

	app.Get("/about", func(ctx context.Context) {
		ctx.WriteString(expectedAbout)
	})

	app.OnErrorCode(iris.StatusNotFound, func(ctx context.Context) {
		ctx.WriteString(unexpectedRoute)
	})

	l, err := net.Listen("tcp", "localhost:443")
	if err != nil {
		t.Fatalf("%v while creating tcp4 listener for new tls local test listener", err)
	}
	// main server
	go app.Run(iris.Listener(httptest.NewLocalTLSListener(l)), iris.WithoutBanner)

	e := httptest.NewInsecure(t, httptest.URL("http://localhost"))
	e.GET("/").Expect().Status(iris.StatusOK).Body().Equal(expectedIndex)
	e.GET("/about").Expect().Status(iris.StatusOK).Body().Equal(expectedAbout)
	e.GET("/notfound").Expect().Status(iris.StatusNotFound).Body().Equal(unexpectedRoute)
}
