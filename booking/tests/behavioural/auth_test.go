package behavioural

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"time"

	"github.com/google/uuid"
	"github.com/identicalaffiliation/booking-service/booking/gen/mocks"
	"github.com/identicalaffiliation/booking-service/booking/internal/controller"
	_ "github.com/identicalaffiliation/booking-service/booking/internal/controller"
	"github.com/identicalaffiliation/booking-service/booking/internal/domain"
	"github.com/identicalaffiliation/booking-service/booking/internal/dto/output"
	"github.com/labstack/echo/v4"
	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
	"github.com/stretchr/testify/mock"
)

var _ = ginkgo.Describe("RegistrationTest", func() {
	var (
		engine  *echo.Echo
		auth    *mocks.MockAuthUsecase
		handler echo.HandlerFunc
	)

	ginkgo.BeforeEach(func() {
		engine = echo.New()
		auth = new(mocks.MockAuthUsecase)
		handler = controller.Registration(auth)

		engine.HTTPErrorHandler = controller.HTTPErrorHandler()
	})

	ginkgo.Context("success case", func() {
		ginkgo.It("returns 201", func() {
			body := `{"nickname":"test user","password":"test password","role":"client"}`
			r := httptest.NewRequest(http.MethodPost, "/registration", bytes.NewBufferString(body))
			r.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

			w := httptest.NewRecorder()
			ctx := engine.NewContext(r, w)

			expected := output.NewUserOutput(uuid.New(), "test user", domain.Client, time.Now().UTC())

			auth.EXPECT().Registration(mock.Anything, mock.Anything).Return(expected, nil)

			err := handler(ctx)
			gomega.Expect(err).To(gomega.BeNil())
			gomega.Expect(http.StatusCreated).To(gomega.Equal(w.Code))
		})
	})

	ginkgo.Context("invalid json body", func() {
		ginkgo.It("returns 400 and error: invalid json body", func() {
			body := `{"nickname":"test user","password":"test password","role":"client}`
			r := httptest.NewRequest(http.MethodPost, "/registration", bytes.NewBufferString(body))
			r.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

			w := httptest.NewRecorder()
			ctx := engine.NewContext(r, w)

			err := handler(ctx)
			if err != nil {
				engine.HTTPErrorHandler(err, ctx)
			}
			gomega.Expect(err).To(gomega.Not(gomega.BeNil()))
			gomega.Expect(http.StatusBadRequest).To(gomega.Equal(w.Code))
		})
	})

	ginkgo.Context("invalid user data", func() {
		ginkgo.It("returns 400 and error: invalid user data", func() {
			body := `{"nickname":"test user","password":"1","role":"client"}`
			r := httptest.NewRequest(http.MethodPost, "/registration", bytes.NewBufferString(body))
			r.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

			w := httptest.NewRecorder()
			ctx := engine.NewContext(r, w)

			auth.EXPECT().Registration(mock.Anything, mock.Anything).Return(nil, domain.ErrInvalidUserData)

			err := handler(ctx)
			if err != nil {
				engine.HTTPErrorHandler(err, ctx)
			}
			gomega.Expect(err).To(gomega.Not(gomega.BeNil()))
			gomega.Expect(http.StatusBadRequest).To(gomega.Equal(w.Code))
		})
	})

	ginkgo.Context("user already exists", func() {
		ginkgo.It("returns 400 and error: user already exists with nickname", func() {
			body := `{"nickname":"test user","password":"test password","role":"client"}`
			r := httptest.NewRequest(http.MethodPost, "/registration", bytes.NewBufferString(body))
			r.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

			w := httptest.NewRecorder()
			ctx := engine.NewContext(r, w)

			auth.EXPECT().Registration(mock.Anything, mock.Anything).Return(nil, domain.ErrUserAlreadyExists)

			err := handler(ctx)
			if err != nil {
				engine.HTTPErrorHandler(err, ctx)
			}
			gomega.Expect(err).To(gomega.Not(gomega.BeNil()))
			gomega.Expect(http.StatusBadRequest).To(gomega.Equal(w.Code))
		})
	})
})

var _ = ginkgo.Describe("LoginTest", func() {
	var (
		engine  *echo.Echo
		auth    *mocks.MockAuthUsecase
		handler echo.HandlerFunc
	)

	ginkgo.BeforeEach(func() {
		engine = echo.New()
		auth = new(mocks.MockAuthUsecase)
		handler = controller.Login(auth)

		engine.HTTPErrorHandler = controller.HTTPErrorHandler()
	})

	ginkgo.Context("success case", func() {
		ginkgo.It("returns 201", func() {
			body := `{"nickname":"test user","password":"test password"}`
			r := httptest.NewRequest(http.MethodPost, "/login", bytes.NewBufferString(body))
			r.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

			w := httptest.NewRecorder()
			ctx := engine.NewContext(r, w)

			expected := &output.LoginOutput{
				Tokens: output.Tokens{
					AccessToken:  "some token",
					RefreshToken: "some token",
				},
			}

			auth.EXPECT().Login(mock.Anything, mock.Anything).Return(expected, nil)

			err := handler(ctx)
			gomega.Expect(err).To(gomega.BeNil())
			gomega.Expect(http.StatusOK).To(gomega.Equal(w.Code))
		})
	})

	ginkgo.Context("invalid json body", func() {
		ginkgo.It("returns 400 and error: invalid json body", func() {
			body := `{"nickname":"test user","password":"test pass}`
			r := httptest.NewRequest(http.MethodPost, "/login", bytes.NewBufferString(body))
			r.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

			w := httptest.NewRecorder()
			ctx := engine.NewContext(r, w)

			err := handler(ctx)
			if err != nil {
				engine.HTTPErrorHandler(err, ctx)
			}
			gomega.Expect(err).To(gomega.Not(gomega.BeNil()))
			gomega.Expect(http.StatusBadRequest).To(gomega.Equal(w.Code))
		})
	})

	ginkgo.Context("invalid user data", func() {
		ginkgo.It("returns 400 and error: invalid user data", func() {
			body := `{"nickname":"test user","password":"test pass","role":"client"}`
			r := httptest.NewRequest(http.MethodPost, "/login", bytes.NewBufferString(body))
			r.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

			w := httptest.NewRecorder()
			ctx := engine.NewContext(r, w)

			auth.EXPECT().Login(mock.Anything, mock.Anything).Return(nil, domain.ErrInvalidUserData)

			err := handler(ctx)
			if err != nil {
				engine.HTTPErrorHandler(err, ctx)
			}
			gomega.Expect(err).To(gomega.Not(gomega.BeNil()))
			gomega.Expect(http.StatusBadRequest).To(gomega.Equal(w.Code))
		})
	})
})

var _ = ginkgo.Describe("RefreshTest", func() {
	var (
		engine  *echo.Echo
		auth    *mocks.MockAuthUsecase
		handler echo.HandlerFunc
	)

	ginkgo.BeforeEach(func() {
		engine = echo.New()
		auth = new(mocks.MockAuthUsecase)
		handler = controller.Refresh(auth)

		engine.HTTPErrorHandler = controller.HTTPErrorHandler()
	})

	ginkgo.Context("success case", func() {
		ginkgo.It("returns 200", func() {
			body := `{"refreshToken":"some token"}`
			r := httptest.NewRequest(http.MethodPost, "/refresh", bytes.NewBufferString(body))
			r.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

			w := httptest.NewRecorder()
			ctx := engine.NewContext(r, w)

			expected := &output.LoginOutput{
				Tokens: output.Tokens{
					AccessToken:  "some token",
					RefreshToken: "some token",
				},
			}

			auth.EXPECT().Refresh(mock.Anything, mock.Anything).Return(expected, nil)

			err := handler(ctx)
			gomega.Expect(err).To(gomega.BeNil())
			gomega.Expect(http.StatusOK).To(gomega.Equal(w.Code))
		})
	})

	ginkgo.Context("invalid json body", func() {
		ginkgo.It("returns 400 and error: invalid json body", func() {
			body := `{dasdas`
			r := httptest.NewRequest(http.MethodPost, "/refresh", bytes.NewBufferString(body))
			r.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

			w := httptest.NewRecorder()
			ctx := engine.NewContext(r, w)

			err := handler(ctx)
			if err != nil {
				engine.HTTPErrorHandler(err, ctx)
			}
			gomega.Expect(err).To(gomega.Not(gomega.BeNil()))
			gomega.Expect(http.StatusBadRequest).To(gomega.Equal(w.Code))
		})
	})

	ginkgo.Context("invalid refresh token data", func() {
		ginkgo.It("returns 400 and error: invalid refresh token data", func() {
			body := `{"refreshToken":"some token""}`
			r := httptest.NewRequest(http.MethodPost, "/refresh", bytes.NewBufferString(body))
			r.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

			w := httptest.NewRecorder()
			ctx := engine.NewContext(r, w)

			auth.EXPECT().Refresh(mock.Anything, mock.Anything).Return(nil, domain.ErrInvalidRefreshTokenData)

			err := handler(ctx)
			if err != nil {
				engine.HTTPErrorHandler(err, ctx)
			}
			gomega.Expect(err).To(gomega.Not(gomega.BeNil()))
			gomega.Expect(http.StatusBadRequest).To(gomega.Equal(w.Code))
		})

		ginkgo.Context("token not found", func() {
			ginkgo.It("returns 404 and error: token not found", func() {
				body := `{"refreshToken":"some token"}`
				r := httptest.NewRequest(http.MethodPost, "/refresh", bytes.NewBufferString(body))
				r.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

				w := httptest.NewRecorder()
				ctx := engine.NewContext(r, w)

				auth.EXPECT().Refresh(mock.Anything, mock.Anything).Return(nil, domain.ErrTokenNotFound)

				err := handler(ctx)
				if err != nil {
					engine.HTTPErrorHandler(err, ctx)
				}
				gomega.Expect(err).To(gomega.Not(gomega.BeNil()))
				gomega.Expect(http.StatusNotFound).To(gomega.Equal(w.Code))
			})
		})
	})
})
