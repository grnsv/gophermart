package handlers_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/grnsv/gophermart/internal/api/handlers"
	"github.com/grnsv/gophermart/internal/api/router"
	"github.com/grnsv/gophermart/internal/logger"
	"github.com/grnsv/gophermart/internal/mocks"
	"github.com/grnsv/gophermart/internal/models"
	"github.com/grnsv/gophermart/internal/services"
	"github.com/grnsv/gophermart/internal/storage"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestHandlers(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Handlers Suite")
}

var _ = Describe("RegisterUser", func() {
	var (
		ctrl *gomock.Controller
		repo *mocks.MockUserRepository
		log  logger.Logger
		h    *handlers.UserHandler
		r    http.Handler
		ts   *httptest.Server
	)

	BeforeEach(func() {
		ctrl = gomock.NewController(GinkgoT())
		repo = mocks.NewMockUserRepository(ctrl)
		log = logger.New()
		h = handlers.NewUserHandler(log, services.NewUserService(repo), services.NewJWTService("", ""))
		r = router.NewRouter(log, h, nil, nil)
		ts = httptest.NewServer(r)
	})

	AfterEach(func() {
		ts.Close()
		ctrl.Finish()
	})

	When("the request body has invalid json", func() {
		It("should return bad request status", func() {
			reqBody := `{"login": "login", "password": "password",}`
			req, err := http.NewRequest("POST", ts.URL+"/api/user/register", strings.NewReader(reqBody))
			Expect(err).ToNot(HaveOccurred())
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Accept-Encoding", "application/json")

			resp, err := http.DefaultClient.Do(req)
			Expect(err).ToNot(HaveOccurred())
			defer resp.Body.Close()
			body := map[string]any{}
			err = json.NewDecoder(resp.Body).Decode(&body)
			Expect(err).ToNot(HaveOccurred())

			Expect(resp.StatusCode).To(Equal(http.StatusBadRequest))
		})
	})

	When("the request body is invalid", func() {
		It("should return bad request status", func() {
			reqBody := `{"login": "", "password": "pass"}`
			req, err := http.NewRequest("POST", ts.URL+"/api/user/register", strings.NewReader(reqBody))
			Expect(err).ToNot(HaveOccurred())
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Accept-Encoding", "application/json")

			resp, err := http.DefaultClient.Do(req)
			Expect(err).ToNot(HaveOccurred())
			defer resp.Body.Close()
			body := map[string]any{}
			err = json.NewDecoder(resp.Body).Decode(&body)
			Expect(err).ToNot(HaveOccurred())

			Expect(resp.StatusCode).To(Equal(http.StatusBadRequest))
		})
	})

	When("the login already exists", func() {
		It("should return conflict status", func() {
			repo.EXPECT().IsLoginExists(gomock.Any(), "login").Return(true, nil)

			reqBody := `{"login": "login", "password": "password"}`
			req, err := http.NewRequest("POST", ts.URL+"/api/user/register", strings.NewReader(reqBody))
			Expect(err).ToNot(HaveOccurred())
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Accept-Encoding", "application/json")

			resp, err := http.DefaultClient.Do(req)
			Expect(err).ToNot(HaveOccurred())
			defer resp.Body.Close()
			body := map[string]any{}
			err = json.NewDecoder(resp.Body).Decode(&body)
			Expect(err).ToNot(HaveOccurred())

			Expect(resp.StatusCode).To(Equal(http.StatusConflict))
		})
	})

	When("the request is valid", func() {
		It("should register user successfully", func() {
			repo.EXPECT().IsLoginExists(gomock.Any(), "login").Return(false, nil)
			repo.EXPECT().CreateUser(gomock.Any(), gomock.Any()).Return(nil)

			reqBody := `{"login": "login", "password": "password"}`
			req, err := http.NewRequest("POST", ts.URL+"/api/user/register", strings.NewReader(reqBody))
			Expect(err).ToNot(HaveOccurred())
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Accept-Encoding", "application/json")

			resp, err := http.DefaultClient.Do(req)
			Expect(err).ToNot(HaveOccurred())
			defer resp.Body.Close()
			body := map[string]any{}
			err = json.NewDecoder(resp.Body).Decode(&body)
			Expect(err).ToNot(HaveOccurred())

			Expect(resp.StatusCode).To(Equal(http.StatusOK))

			var tokenCookie *http.Cookie
			for _, cookie := range resp.Cookies() {
				if cookie.Name == "token" {
					tokenCookie = cookie
					break
				}
			}
			Expect(tokenCookie).NotTo(BeNil())
			Expect(tokenCookie.Value).NotTo(BeEmpty())
		})
	})
})

var _ = Describe("LoginUser", func() {
	var (
		ctrl *gomock.Controller
		repo *mocks.MockUserRepository
		log  logger.Logger
		h    *handlers.UserHandler
		r    http.Handler
		ts   *httptest.Server
	)

	BeforeEach(func() {
		ctrl = gomock.NewController(GinkgoT())
		repo = mocks.NewMockUserRepository(ctrl)
		log = logger.New()
		h = handlers.NewUserHandler(log, services.NewUserService(repo), services.NewJWTService("", ""))
		r = router.NewRouter(log, h, nil, nil)
		ts = httptest.NewServer(r)
	})

	AfterEach(func() {
		ts.Close()
		ctrl.Finish()
	})

	When("the request body has invalid json", func() {
		It("should return bad request status", func() {
			reqBody := `{"login": "login", "password": "password",}`
			req, err := http.NewRequest("POST", ts.URL+"/api/user/login", strings.NewReader(reqBody))
			Expect(err).ToNot(HaveOccurred())
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Accept-Encoding", "application/json")

			resp, err := http.DefaultClient.Do(req)
			Expect(err).ToNot(HaveOccurred())
			defer resp.Body.Close()
			body := map[string]any{}
			err = json.NewDecoder(resp.Body).Decode(&body)
			Expect(err).ToNot(HaveOccurred())

			Expect(resp.StatusCode).To(Equal(http.StatusBadRequest))
		})
	})

	When("the request body is invalid", func() {
		It("should return bad request status", func() {
			reqBody := `{"login": "", "pass": "password"}`
			req, err := http.NewRequest("POST", ts.URL+"/api/user/login", strings.NewReader(reqBody))
			Expect(err).ToNot(HaveOccurred())
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Accept-Encoding", "application/json")

			resp, err := http.DefaultClient.Do(req)
			Expect(err).ToNot(HaveOccurred())
			defer resp.Body.Close()
			body := map[string]any{}
			err = json.NewDecoder(resp.Body).Decode(&body)
			Expect(err).ToNot(HaveOccurred())

			Expect(resp.StatusCode).To(Equal(http.StatusBadRequest))
		})
	})

	When("the user does not exists", func() {
		It("should return unauthorized status", func() {
			repo.EXPECT().FindUserByLogin(gomock.Any(), "login").Return(nil, storage.ErrNotFound)

			reqBody := `{"login": "login", "password": "password"}`
			req, err := http.NewRequest("POST", ts.URL+"/api/user/login", strings.NewReader(reqBody))
			Expect(err).ToNot(HaveOccurred())
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Accept-Encoding", "application/json")

			resp, err := http.DefaultClient.Do(req)
			Expect(err).ToNot(HaveOccurred())
			defer resp.Body.Close()
			body := map[string]any{}
			err = json.NewDecoder(resp.Body).Decode(&body)
			Expect(err).ToNot(HaveOccurred())

			Expect(resp.StatusCode).To(Equal(http.StatusUnauthorized))
		})
	})

	When("the password is wrong", func() {
		It("should return unauthorized status", func() {
			repo.EXPECT().FindUserByLogin(gomock.Any(), "login").Return(&models.User{
				ID:       "00000000-0000-0000-0000-000000000000",
				Login:    "login",
				Password: "$2a$04$00000000000000000000000000000000000000000000000000000",
			}, nil)

			reqBody := `{"login": "login", "password": "password"}`
			req, err := http.NewRequest("POST", ts.URL+"/api/user/login", strings.NewReader(reqBody))
			Expect(err).ToNot(HaveOccurred())
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Accept-Encoding", "application/json")

			resp, err := http.DefaultClient.Do(req)
			Expect(err).ToNot(HaveOccurred())
			defer resp.Body.Close()
			body := map[string]any{}
			err = json.NewDecoder(resp.Body).Decode(&body)
			Expect(err).ToNot(HaveOccurred())

			Expect(resp.StatusCode).To(Equal(http.StatusUnauthorized))
		})
	})

	When("the request is valid", func() {
		It("should login user successfully", func() {
			repo.EXPECT().FindUserByLogin(gomock.Any(), "login").Return(&models.User{
				ID:       "00000000-0000-0000-0000-000000000000",
				Login:    "login",
				Password: "$2a$05$bvIG6Nmid91Mu9RcmmWZfO5HJIMCT8riNW0hEp8f6/FuA2/mHZFpe",
			}, nil)

			reqBody := `{"login": "login", "password": "password"}`
			req, err := http.NewRequest("POST", ts.URL+"/api/user/login", strings.NewReader(reqBody))
			Expect(err).ToNot(HaveOccurred())
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Accept-Encoding", "application/json")

			resp, err := http.DefaultClient.Do(req)
			Expect(err).ToNot(HaveOccurred())
			defer resp.Body.Close()
			body := map[string]any{}
			err = json.NewDecoder(resp.Body).Decode(&body)
			Expect(err).ToNot(HaveOccurred())

			Expect(resp.StatusCode).To(Equal(http.StatusOK))

			var tokenCookie *http.Cookie
			for _, cookie := range resp.Cookies() {
				if cookie.Name == "token" {
					tokenCookie = cookie
					break
				}
			}
			Expect(tokenCookie).NotTo(BeNil())
			Expect(tokenCookie.Value).NotTo(BeEmpty())
		})
	})
})

var _ = Describe("GetOrders", func() {
	var (
		ctrl *gomock.Controller
		repo *mocks.MockOrderRepository
		log  logger.Logger
		h    *handlers.OrderHandler
		jwts services.JWTService
		r    http.Handler
		ts   *httptest.Server
	)

	BeforeEach(func() {
		ctrl = gomock.NewController(GinkgoT())
		repo = mocks.NewMockOrderRepository(ctrl)
		log = logger.New()
		h = handlers.NewOrderHandler(log, services.NewOrderService(repo), services.NewLuhnService())
		jwts = services.NewJWTService("localhost:8080", "secret")
		r = router.NewRouter(log, nil, h, jwts)
		ts = httptest.NewServer(r)
	})

	AfterEach(func() {
		ts.Close()
		ctrl.Finish()
	})

	When("the request does not have token", func() {
		It("should return unauthorized status", func() {
			resp, err := http.Get(ts.URL + "/api/user/orders")
			Expect(err).NotTo(HaveOccurred())
			defer resp.Body.Close()
			body := map[string]any{}
			err = json.NewDecoder(resp.Body).Decode(&body)
			Expect(err).ToNot(HaveOccurred())
			Expect(resp.StatusCode).To(Equal(http.StatusUnauthorized))
		})
	})

	When("the request has invalid token", func() {
		It("should return unauthorized status", func() {
			req, err := http.NewRequest("GET", ts.URL+"/api/user/orders", nil)
			Expect(err).NotTo(HaveOccurred())
			req.AddCookie(&http.Cookie{Name: "token", Value: "invalid"})
			resp, err := http.DefaultClient.Do(req)
			Expect(err).NotTo(HaveOccurred())
			defer resp.Body.Close()
			body := map[string]any{}
			err = json.NewDecoder(resp.Body).Decode(&body)
			Expect(err).ToNot(HaveOccurred())
			Expect(resp.StatusCode).To(Equal(http.StatusUnauthorized))
		})
	})

	When("the request has empty token", func() {
		It("should return unauthorized status", func() {
			req, err := http.NewRequest("GET", ts.URL+"/api/user/orders", nil)
			Expect(err).NotTo(HaveOccurred())
			req.AddCookie(&http.Cookie{Name: "token", Value: ""})
			resp, err := http.DefaultClient.Do(req)
			Expect(err).NotTo(HaveOccurred())
			defer resp.Body.Close()
			body := map[string]any{}
			err = json.NewDecoder(resp.Body).Decode(&body)
			Expect(err).ToNot(HaveOccurred())
			Expect(resp.StatusCode).To(Equal(http.StatusUnauthorized))
		})
	})

	When("the request has empty user ID", func() {
		It("should return unauthorized status", func() {
			req, err := http.NewRequest("GET", ts.URL+"/api/user/orders", nil)
			Expect(err).NotTo(HaveOccurred())
			cookie, err := jwts.BuildCookie("")
			Expect(err).NotTo(HaveOccurred())
			req.AddCookie(cookie)
			resp, err := http.DefaultClient.Do(req)
			Expect(err).NotTo(HaveOccurred())
			defer resp.Body.Close()
			body := map[string]any{}
			err = json.NewDecoder(resp.Body).Decode(&body)
			Expect(err).ToNot(HaveOccurred())
			Expect(resp.StatusCode).To(Equal(http.StatusUnauthorized))
		})
	})
})

var _ = Describe("UploadOrder", func() {
	var (
		ctrl *gomock.Controller
		repo *mocks.MockOrderRepository
		log  logger.Logger
		h    *handlers.OrderHandler
		jwts services.JWTService
		r    http.Handler
		ts   *httptest.Server
	)

	BeforeEach(func() {
		ctrl = gomock.NewController(GinkgoT())
		repo = mocks.NewMockOrderRepository(ctrl)
		log = logger.New()
		h = handlers.NewOrderHandler(log, services.NewOrderService(repo), services.NewLuhnService())
		jwts = services.NewJWTService("localhost:8080", "secret")
		r = router.NewRouter(log, nil, h, jwts)
		ts = httptest.NewServer(r)
	})

	AfterEach(func() {
		ts.Close()
		ctrl.Finish()
	})

	When("the request has empty body", func() {
		It("should return bad request status", func() {
			req, err := http.NewRequest("POST", ts.URL+"/api/user/orders", nil)
			Expect(err).NotTo(HaveOccurred())
			cookie, err := jwts.BuildCookie("00000000-0000-0000-0000-000000000000")
			Expect(err).NotTo(HaveOccurred())
			req.AddCookie(cookie)
			resp, err := http.DefaultClient.Do(req)
			Expect(err).NotTo(HaveOccurred())
			defer resp.Body.Close()
			body := map[string]any{}
			err = json.NewDecoder(resp.Body).Decode(&body)
			Expect(err).ToNot(HaveOccurred())
			Expect(resp.StatusCode).To(Equal(http.StatusUnprocessableEntity))
		})
	})

	When("the request has empty order ID", func() {
		It("should return bad request status", func() {
			reqBody := " "
			req, err := http.NewRequest("POST", ts.URL+"/api/user/orders", strings.NewReader(reqBody))
			Expect(err).ToNot(HaveOccurred())
			cookie, err := jwts.BuildCookie("00000000-0000-0000-0000-000000000000")
			Expect(err).NotTo(HaveOccurred())
			req.AddCookie(cookie)
			resp, err := http.DefaultClient.Do(req)
			Expect(err).NotTo(HaveOccurred())
			defer resp.Body.Close()
			body := map[string]any{}
			err = json.NewDecoder(resp.Body).Decode(&body)
			Expect(err).ToNot(HaveOccurred())
			Expect(resp.StatusCode).To(Equal(http.StatusUnprocessableEntity))
		})
	})

	When("the request has order ID with non-digit characters", func() {
		It("should return bad request status", func() {
			reqBody := "12345678903A"
			req, err := http.NewRequest("POST", ts.URL+"/api/user/orders", strings.NewReader(reqBody))
			Expect(err).ToNot(HaveOccurred())
			cookie, err := jwts.BuildCookie("00000000-0000-0000-0000-000000000000")
			Expect(err).NotTo(HaveOccurred())
			req.AddCookie(cookie)
			resp, err := http.DefaultClient.Do(req)
			Expect(err).NotTo(HaveOccurred())
			defer resp.Body.Close()
			body := map[string]any{}
			err = json.NewDecoder(resp.Body).Decode(&body)
			Expect(err).ToNot(HaveOccurred())
			Expect(resp.StatusCode).To(Equal(http.StatusUnprocessableEntity))
		})
	})

	When("the request has invalid order ID", func() {
		It("should return bad request status", func() {
			reqBody := "12345678904"
			req, err := http.NewRequest("POST", ts.URL+"/api/user/orders", strings.NewReader(reqBody))
			Expect(err).ToNot(HaveOccurred())
			cookie, err := jwts.BuildCookie("00000000-0000-0000-0000-000000000000")
			Expect(err).NotTo(HaveOccurred())
			req.AddCookie(cookie)
			resp, err := http.DefaultClient.Do(req)
			Expect(err).NotTo(HaveOccurred())
			defer resp.Body.Close()
			body := map[string]any{}
			err = json.NewDecoder(resp.Body).Decode(&body)
			Expect(err).ToNot(HaveOccurred())
			Expect(resp.StatusCode).To(Equal(http.StatusUnprocessableEntity))
		})
	})

	When("the order has already been uploaded by current user", func() {
		It("should return OK status", func() {
			repo.EXPECT().FindOrderByID(gomock.Any(), 12345678903).Return(nil, &services.OrderAlreadyExistsError{
				UserID: "00000000-0000-0000-0000-000000000000",
			})

			reqBody := "12345678903"
			req, err := http.NewRequest("POST", ts.URL+"/api/user/orders", strings.NewReader(reqBody))
			Expect(err).ToNot(HaveOccurred())
			cookie, err := jwts.BuildCookie("00000000-0000-0000-0000-000000000000")
			Expect(err).NotTo(HaveOccurred())
			req.AddCookie(cookie)
			resp, err := http.DefaultClient.Do(req)
			Expect(err).NotTo(HaveOccurred())
			defer resp.Body.Close()
			body := map[string]any{}
			err = json.NewDecoder(resp.Body).Decode(&body)
			Expect(err).ToNot(HaveOccurred())
			Expect(resp.StatusCode).To(Equal(http.StatusOK))
		})
	})

	When("the order has already been uploaded by another user", func() {
		It("should return conflict status", func() {
			repo.EXPECT().FindOrderByID(gomock.Any(), 12345678903).Return(nil, &services.OrderAlreadyExistsError{
				UserID: "ffffffff-ffff-ffff-ffff-ffffffffffff",
			})

			reqBody := "12345678903"
			req, err := http.NewRequest("POST", ts.URL+"/api/user/orders", strings.NewReader(reqBody))
			Expect(err).ToNot(HaveOccurred())
			cookie, err := jwts.BuildCookie("00000000-0000-0000-0000-000000000000")
			Expect(err).NotTo(HaveOccurred())
			req.AddCookie(cookie)
			resp, err := http.DefaultClient.Do(req)
			Expect(err).NotTo(HaveOccurred())
			defer resp.Body.Close()
			body := map[string]any{}
			err = json.NewDecoder(resp.Body).Decode(&body)
			Expect(err).ToNot(HaveOccurred())
			Expect(resp.StatusCode).To(Equal(http.StatusConflict))
		})
	})

	When("the order is successfully uploaded", func() {
		It("should return accepted status", func() {
			repo.EXPECT().FindOrderByID(gomock.Any(), 12345678903).Return(nil, storage.ErrNotFound)
			repo.EXPECT().CreateOrder(gomock.Any(), gomock.Any()).Return(nil)

			reqBody := "12345678903"
			req, err := http.NewRequest("POST", ts.URL+"/api/user/orders", strings.NewReader(reqBody))
			Expect(err).ToNot(HaveOccurred())
			cookie, err := jwts.BuildCookie("00000000-0000-0000-0000-000000000000")
			Expect(err).NotTo(HaveOccurred())
			req.AddCookie(cookie)
			resp, err := http.DefaultClient.Do(req)
			Expect(err).NotTo(HaveOccurred())
			defer resp.Body.Close()
			body := map[string]any{}
			err = json.NewDecoder(resp.Body).Decode(&body)
			Expect(err).ToNot(HaveOccurred())
			Expect(resp.StatusCode).To(Equal(http.StatusAccepted))
		})
	})
})
