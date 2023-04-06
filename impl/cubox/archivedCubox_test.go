package cubox

import (
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("ArchivedCubox", func() {
	var _ = Describe("convertTime", func() {
		const (
			rightTimeStr = "2022-07-31T01:26:33:603+08:00"
			emptyTimeStr = ""
		)
		var (
			timeStr      string
			expectedTime time.Time
		)
		JustBeforeEach(func() {
			actualTime, err := convertTime(timeStr)
			Expect(err).To(BeNil())
			Expect(actualTime).To(Equal(expectedTime))
		})

		When("timeStr is right", func() {
			BeforeEach(func() {
				timeStr = rightTimeStr
				expectedTime = time.Date(2022, 7, 31, 1, 26, 33, 603000000, time.Local)
			})
			It("should return right time", func() {
			})
		})

		When("timeStr is empty", func() {
			BeforeEach(func() {
				timeStr = emptyTimeStr
				expectedTime = time.Time{}
			})
			It("should return empty time", func() {
			})
		})
	})
})
