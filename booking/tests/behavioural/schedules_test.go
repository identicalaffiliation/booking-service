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

var _ = ginkgo.Describe("CreateScheduleTest", func() {
	var (
		engine   *echo.Echo
		schedule *mocks.MockSchedulesUsecase
		handler  echo.HandlerFunc
	)

	ginkgo.BeforeEach(func() {
		engine = echo.New()
		schedule = new(mocks.MockSchedulesUsecase)
		handler = controller.CreateSchedule(schedule)
	})

	ginkgo.Context("success case", func() {
		ginkgo.It("will return 201 and create a new schedule", func() {
			body := `{"day":"10-01-2001","start":"15:00","end":"20:00"}`

			roomID := uuid.New()
			r := httptest.NewRequest(http.MethodPost, "/rooms/"+roomID.String()+"/schedule", bytes.NewBufferString(body))
			r.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

			w := httptest.NewRecorder()
			ctx := engine.NewContext(r, w)
			ctx.SetParamNames("roomId")
			ctx.SetParamValues(roomID.String())
			expected := output.NewCreateScheduleOutput(uuid.New(), roomID, "10-01-2001", "15:00", "20:00", time.Now().UTC())

			schedule.EXPECT().CreateSchedule(mock.Anything, mock.Anything).Return(expected, nil).Once()

			err := handler(ctx)
			gomega.Expect(err).To(gomega.BeNil())
			gomega.Expect(http.StatusCreated).To(gomega.Equal(w.Code))
		})
	})
})
