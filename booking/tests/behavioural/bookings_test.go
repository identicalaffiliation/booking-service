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
})

var _ = ginkgo.Describe("CancelBookingTest", func() {
	var (
		engine  *echo.Echo
		booking *mocks.MockBookingsUsecase
		handler echo.HandlerFunc
	)

	ginkgo.BeforeEach(func() {
		engine = echo.New()
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
})
