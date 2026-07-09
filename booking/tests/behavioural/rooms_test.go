package behavioural

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"time"

	"github.com/google/uuid"
	"github.com/identicalaffiliation/booking-service/booking/gen/mocks"
	"github.com/identicalaffiliation/booking-service/booking/internal/controller"
	"github.com/identicalaffiliation/booking-service/booking/internal/dto/output"
	"github.com/labstack/echo/v4"
	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
	"github.com/stretchr/testify/mock"
)

var _ = ginkgo.Describe("CreateRoomTest", func() {
	var (
		engine  *echo.Echo
		rooms   *mocks.MockRoomsUsecase
		handler echo.HandlerFunc
	)

	ginkgo.BeforeEach(func() {
		engine = echo.New()
		rooms = new(mocks.MockRoomsUsecase)
		handler = controller.CreateRoom(rooms)
	})

	ginkgo.Context("success case", func() {
		ginkgo.It("created a new room and returned 201", func() {
			body := `{"name":"test room","capacity":5}`
			r := httptest.NewRequest(http.MethodPost, "/rooms", bytes.NewBufferString(body))
			r.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

			w := httptest.NewRecorder()
			ctx := engine.NewContext(r, w)

			expected := output.NewCreateRoomOutput(uuid.New(), "test room", 5, time.Now().UTC())

			rooms.EXPECT().CreateRoom(mock.Anything, mock.Anything).Return(expected, nil)

			err := handler(ctx)
			gomega.Expect(err).To(gomega.BeNil())
			gomega.Expect(http.StatusCreated).To(gomega.Equal(w.Code))
		})
	})
})

var _ = ginkgo.Describe("GetRoomTest", func() {
	var (
		engine  *echo.Echo
		rooms   *mocks.MockRoomsUsecase
		handler echo.HandlerFunc
	)

	ginkgo.BeforeEach(func() {
		engine = echo.New()
		rooms = new(mocks.MockRoomsUsecase)
		handler = controller.GetRoom(rooms)
	})

	ginkgo.Context("success case", func() {
		ginkgo.It("return a room and returned 200", func() {
			r := httptest.NewRequest(http.MethodGet, "/rooms/"+uuid.New().String(), nil)
			r.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

			w := httptest.NewRecorder()
			ctx := engine.NewContext(r, w)

			expected := output.NewRoomOutput(uuid.New(), "test room", 5, time.Now().UTC())

			rooms.EXPECT().GetRoom(mock.Anything, mock.Anything).Return(expected, nil)

			err := handler(ctx)
			gomega.Expect(err).To(gomega.BeNil())
			gomega.Expect(http.StatusOK).To(gomega.Equal(w.Code))
		})
	})
})

var _ = ginkgo.Describe("GetRoomsTest", func() {
	var (
		engine  *echo.Echo
		rooms   *mocks.MockRoomsUsecase
		handler echo.HandlerFunc
	)

	ginkgo.BeforeEach(func() {
		engine = echo.New()
		rooms = new(mocks.MockRoomsUsecase)
		handler = controller.GetRooms(rooms)
	})

	ginkgo.Context("success case", func() {
		ginkgo.It("return rooms and returned 200", func() {
			r := httptest.NewRequest(http.MethodGet, "/rooms", nil)
			r.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

			w := httptest.NewRecorder()
			ctx := engine.NewContext(r, w)

			expected := &output.RoomsOutput{Rooms: []*output.RoomOutput{output.NewRoomOutput(uuid.New(), "test room", 5, time.Now().UTC())}}

			rooms.EXPECT().GetRooms(mock.Anything).Return(expected, nil)

			err := handler(ctx)
			gomega.Expect(err).To(gomega.BeNil())
			gomega.Expect(http.StatusOK).To(gomega.Equal(w.Code))
		})
	})
})

var _ = ginkgo.Describe("GetRoomTest", func() {
	var (
		engine  *echo.Echo
		rooms   *mocks.MockRoomsUsecase
		handler echo.HandlerFunc
	)

	ginkgo.BeforeEach(func() {
		engine = echo.New()
		rooms = new(mocks.MockRoomsUsecase)
		handler = controller.DeleteRoom(rooms)
	})

	ginkgo.Context("success case", func() {
		ginkgo.It("delete a room and returned 200", func() {
			r := httptest.NewRequest(http.MethodGet, "/rooms/"+uuid.New().String(), nil)
			r.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

			w := httptest.NewRecorder()
			ctx := engine.NewContext(r, w)

			rooms.EXPECT().DeleteRoom(mock.Anything, mock.Anything).Return(nil)

			err := handler(ctx)
			gomega.Expect(err).To(gomega.BeNil())
			gomega.Expect(http.StatusOK).To(gomega.Equal(w.Code))
		})
	})
})
