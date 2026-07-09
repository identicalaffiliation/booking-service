package behavioural

import (
	"bytes"
	"net/http"
	"net/http/httptest"

	"github.com/google/uuid"
	"github.com/identicalaffiliation/booking-service/booking/gen/mocks"
	"github.com/identicalaffiliation/booking-service/booking/internal/controller"
	"github.com/identicalaffiliation/booking-service/booking/internal/domain"
	"github.com/identicalaffiliation/booking-service/booking/internal/dto/output"
	"github.com/labstack/echo/v4"
	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
	"github.com/stretchr/testify/mock"
)

var _ = ginkgo.Describe("CreateBookingTest", func() {
	var (
		engine  *echo.Echo
		booking *mocks.MockBookingsUsecase
		handler echo.HandlerFunc
	)

	ginkgo.BeforeEach(func() {
		engine = echo.New()
		engine.HTTPErrorHandler = controller.HTTPErrorHandler()
		booking = new(mocks.MockBookingsUsecase)
		handler = controller.CreateBooking(booking)
	})

	ginkgo.Context("success", func() {
		ginkgo.It("will return 201 and create booking", func() {
			body := `{"slotId":"` + uuid.New().String() + `"}`
			r := httptest.NewRequest(http.MethodPost, "/bookings", bytes.NewBufferString(body))
			r.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

			w := httptest.NewRecorder()
			ctx := engine.NewContext(r, w)

			expected := output.NewBookingOutput(&domain.Booking{
				ID:     uuid.New(),
				UserID: uuid.New(),
				SlotID: uuid.New(),
			})

			booking.EXPECT().Create(mock.Anything, mock.Anything).Return(expected, nil).Once()

			err := handler(ctx)
			gomega.Expect(err).To(gomega.BeNil())
			gomega.Expect(http.StatusCreated).To(gomega.Equal(w.Code))
		})
	})

	ginkgo.Context("invalid json body", func() {
		ginkgo.It("will return 400 and error: invalid json body", func() {
			body := `invalid json`
			r := httptest.NewRequest(http.MethodPost, "/bookings", bytes.NewBufferString(body))
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

	ginkgo.Context("invalid booking data", func() {
		ginkgo.It("will return 400 and error: invalid booking data", func() {
			body := `{"slotId":"invalid id"}`
			r := httptest.NewRequest(http.MethodPost, "/bookings", bytes.NewBufferString(body))
			r.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

			w := httptest.NewRecorder()
			ctx := engine.NewContext(r, w)

			booking.EXPECT().Create(mock.Anything, mock.Anything).Return(nil, domain.ErrInvalidBookingData).Once()

			err := handler(ctx)
			if err != nil {
				engine.HTTPErrorHandler(err, ctx)
			}

			gomega.Expect(err).To(gomega.Not(gomega.BeNil()))
			gomega.Expect(http.StatusBadRequest).To(gomega.Equal(w.Code))
		})
	})

	ginkgo.Context("invalid user data", func() {
		ginkgo.It("will return 400 and error: invalid user data", func() {
			body := `{"slotId":"o2909329323891"}`
			r := httptest.NewRequest(http.MethodPost, "/bookings", bytes.NewBufferString(body))
			r.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

			w := httptest.NewRecorder()
			ctx := engine.NewContext(r, w)

			booking.EXPECT().Create(mock.Anything, mock.Anything).Return(nil, domain.ErrInvalidUserData).Once()

			err := handler(ctx)
			if err != nil {
				engine.HTTPErrorHandler(err, ctx)
			}

			gomega.Expect(err).To(gomega.Not(gomega.BeNil()))
			gomega.Expect(http.StatusBadRequest).To(gomega.Equal(w.Code))
		})
	})

	ginkgo.Context("slot already booked", func() {
		ginkgo.It("will return 400 and error: slot already booked", func() {
			body := `{"slotId":"o2909329323891"}`
			r := httptest.NewRequest(http.MethodPost, "/bookings", bytes.NewBufferString(body))
			r.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

			w := httptest.NewRecorder()
			ctx := engine.NewContext(r, w)

			booking.EXPECT().Create(mock.Anything, mock.Anything).Return(nil, domain.ErrSlotAlreadyBooked).Once()

			err := handler(ctx)
			if err != nil {
				engine.HTTPErrorHandler(err, ctx)
			}

			gomega.Expect(err).To(gomega.Not(gomega.BeNil()))
			gomega.Expect(http.StatusBadRequest).To(gomega.Equal(w.Code))
		})
	})
})

var _ = ginkgo.Describe("CancelBookingTest", func() {
	var (
		engine  *echo.Echo
		booking *mocks.MockBookingsUsecase
		handler echo.HandlerFunc
	)

	ginkgo.BeforeEach(func() {
		engine = echo.New()
		engine.HTTPErrorHandler = controller.HTTPErrorHandler()
		booking = new(mocks.MockBookingsUsecase)
		handler = controller.CancelBooking(booking)
	})

	ginkgo.Context("success", func() {
		ginkgo.It("will return 200 and cancelled booking", func() {
			r := httptest.NewRequest(http.MethodPost, "/bookings", nil)
			r.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

			w := httptest.NewRecorder()
			ctx := engine.NewContext(r, w)
			ctx.SetParamNames("bookingId")
			ctx.SetParamValues(uuid.New().String())

			booking.EXPECT().Cancel(mock.Anything, mock.Anything).Return(nil).Once()

			err := handler(ctx)
			gomega.Expect(err).To(gomega.BeNil())
			gomega.Expect(http.StatusOK).To(gomega.Equal(w.Code))
		})
	})

	ginkgo.Context("invalid booking id", func() {
		ginkgo.It("will return 400 and error: invalid booking data", func() {
			bookingID := uuid.New()
			r := httptest.NewRequest(http.MethodPost, "/bookings/"+bookingID.String(), nil)
			r.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

			w := httptest.NewRecorder()
			ctx := engine.NewContext(r, w)
			ctx.SetParamNames(controller.BookingIdMuxPattern)
			ctx.SetParamValues(booking.String() + "s$w4q5e512")

			booking.EXPECT().Cancel(mock.Anything, mock.Anything).Return(nil).Once()

			err := handler(ctx)
			if err != nil {
				engine.HTTPErrorHandler(err, ctx)
			}

			gomega.Expect(err).To(gomega.Not(gomega.BeNil()))
			gomega.Expect(http.StatusBadRequest).To(gomega.Equal(w.Code))
		})
	})

	ginkgo.Context("invalid booking data", func() {
		ginkgo.It("will return 400 and error: invalid booking data", func() {
			bookingID := uuid.New()
			r := httptest.NewRequest(http.MethodPost, "/bookings/"+bookingID.String(), nil)
			r.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

			w := httptest.NewRecorder()
			ctx := engine.NewContext(r, w)
			ctx.SetParamNames(controller.BookingIdMuxPattern)
			ctx.SetParamValues(booking.String())

			booking.EXPECT().Cancel(mock.Anything, mock.Anything).Return(domain.ErrInvalidBookingData).Once()

			err := handler(ctx)
			if err != nil {
				engine.HTTPErrorHandler(err, ctx)
			}

			gomega.Expect(err).To(gomega.Not(gomega.BeNil()))
			gomega.Expect(http.StatusBadRequest).To(gomega.Equal(w.Code))
		})
	})
})

var _ = ginkgo.Describe("GetBookingTest", func() {
	var (
		engine  *echo.Echo
		booking *mocks.MockBookingsUsecase
		handler echo.HandlerFunc
	)

	ginkgo.BeforeEach(func() {
		engine = echo.New()
		engine.HTTPErrorHandler = controller.HTTPErrorHandler()
		booking = new(mocks.MockBookingsUsecase)
		handler = controller.GetBooking(booking)
	})

	ginkgo.Context("success", func() {
		ginkgo.It("will return 200 and booking body", func() {
			bookingID := uuid.New()
			r := httptest.NewRequest(http.MethodGet, "/bookings/"+bookingID.String(), nil)
			r.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

			w := httptest.NewRecorder()
			ctx := engine.NewContext(r, w)

			ctx.SetParamNames(controller.BookingIdMuxPattern)
			ctx.SetParamValues(bookingID.String())

			booking.EXPECT().GetMyBooking(mock.Anything, mock.Anything).Return(&output.MyBookingOutput{}, nil).Once()

			err := handler(ctx)
			gomega.Expect(err).To(gomega.BeNil())
			gomega.Expect(http.StatusOK).To(gomega.Equal(w.Code))
		})
	})

	ginkgo.Context("invalid booking data", func() {
		ginkgo.It("will return 400 and error: invalid booking data", func() {
			bookingID := uuid.New()
			r := httptest.NewRequest(http.MethodGet, "/bookings/"+bookingID.String(), nil)
			r.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

			w := httptest.NewRecorder()
			ctx := engine.NewContext(r, w)

			booking.EXPECT().GetMyBooking(mock.Anything, mock.Anything).Return(nil, domain.ErrInvalidBookingData).Once()

			err := handler(ctx)
			if err != nil {
				engine.HTTPErrorHandler(err, ctx)
			}

			gomega.Expect(err).To(gomega.Not(gomega.BeNil()))
			gomega.Expect(http.StatusBadRequest).To(gomega.Equal(w.Code))
		})
	})

	ginkgo.Context("booking not found", func() {
		ginkgo.It("will return 404 and error: booking not found", func() {
			bookingID := uuid.New()
			r := httptest.NewRequest(http.MethodGet, "/bookings/"+bookingID.String(), nil)
			r.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

			w := httptest.NewRecorder()
			ctx := engine.NewContext(r, w)
			ctx.SetParamNames(controller.BookingIdMuxPattern)
			ctx.SetParamValues(bookingID.String())

			booking.EXPECT().GetMyBooking(mock.Anything, mock.Anything).Return(nil, domain.ErrBookingNotFound).Once()

			err := handler(ctx)
			if err != nil {
				engine.HTTPErrorHandler(err, ctx)
			}

			gomega.Expect(err).To(gomega.Not(gomega.BeNil()))
			gomega.Expect(http.StatusNotFound).To(gomega.Equal(w.Code))
		})
	})
})

var _ = ginkgo.Describe("GetBookingsTest", func() {
	var (
		engine  *echo.Echo
		booking *mocks.MockBookingsUsecase
		handler echo.HandlerFunc
	)

	ginkgo.BeforeEach(func() {
		engine = echo.New()
		engine.HTTPErrorHandler = controller.HTTPErrorHandler()
		booking = new(mocks.MockBookingsUsecase)
		handler = controller.GetBookings(booking)
	})

	ginkgo.Context("success", func() {
		ginkgo.It("will return 200 and bookings body", func() {
			r := httptest.NewRequest(http.MethodGet, "/bookings", nil)
			r.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

			w := httptest.NewRecorder()
			ctx := engine.NewContext(r, w)

			booking.EXPECT().GetMyBookings(mock.Anything).Return(&output.MyBookingsOutput{}, nil).Once()

			err := handler(ctx)
			gomega.Expect(err).To(gomega.BeNil())
			gomega.Expect(http.StatusOK).To(gomega.Equal(w.Code))
		})
	})

	ginkgo.Context("invalid user id", func() {
		ginkgo.It("will return 400 and error: invalid user data", func() {
			r := httptest.NewRequest(http.MethodGet, "/bookings", nil)
			r.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

			w := httptest.NewRecorder()
			ctx := engine.NewContext(r, w)

			booking.EXPECT().GetMyBookings(mock.Anything).Return(nil, domain.ErrInvalidUserData).Once()

			err := handler(ctx)
			if err != nil {
				engine.HTTPErrorHandler(err, ctx)
			}
			gomega.Expect(err).To(gomega.Not(gomega.BeNil()))
			gomega.Expect(http.StatusBadRequest).To(gomega.Equal(w.Code))
		})
	})
})
