package repositories_test

import (
	"database/sql/driver"
	"github.com/erikstmartin/go-testdb"
	_ "github.com/erikstmartin/go-testdb"
	"github.com/jinzhu/gorm"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"order/models"
	"order/repositories"
	"testing"
)

func TestRepositories(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Repositories Suite")
}

var _ = Describe("Repositories test with go-testdb", func() {
	var (
		orderRepo repositories.OrderRepository
		orders    models.Orders
		err       error

		userID = 5
	)
	BeforeEach(func() {
		tx, err := gorm.Open("testdb", "")
		Expect(err).To(BeNil())
		orderRepo = repositories.NewOrderRepository(tx)
	})

	Describe("FindAllOrdersByUserID", func() {
		Describe("with no records in the database", func() {
			BeforeEach(func() {
				testdb.SetQueryFunc(func(query string) (driver.Rows, error) {
					columns := []string{"total", "currency_code", "user_id", "restaurant_id", "placed_at"}
					result := ""
					return testdb.RowsFromCSVString(columns, result), nil
				})
			})
			It("returns an empty slice of orders", func() {
				orders, err = orderRepo.FindAllOrdersByUserID(userID)
				Expect(err).To(BeNil())
				Expect(len(orders)).To(Equal(0))
			})
		})

		Describe("when a few records exist", func() {
			BeforeEach(func() {
				testdb.SetQueryFunc(func(query string) (driver.Rows, error) {
					columns := []string{"total", "currency_code", "user_id", "restaurant_id"}
					result := `
		1000,GBP,5,9
		2500,GBP,5,8
		`
					return testdb.RowsFromCSVString(columns, result), nil
				})
			})

			It("returns only the records belonging to the user, in order from latest placed_at first", func() {
				orders, err = orderRepo.FindAllOrdersByUserID(userID)
				Expect(err).To(BeNil())
				Expect(len(orders)).To(Equal(2))
				Expect(orders[0].RestaurantID).To(Equal(9))
				Expect(orders[1].RestaurantID).To(Equal(8))
			})
		})
	})
})

//var _ = Describe("Repositories", func() {
//	var (
//		tx        *gorm.DB
//		orderRepo repositories.OrderRepository
//		orders    models.Orders
//		err       error
//
//		userID = 5
//	)
//
//	BeforeEach(func() {
//		tx = application.ResolveDB().Begin()
//		orderRepo = repositories.NewOrderRepository(tx)
//	})
//
//	Describe("FindAllOrdersByUserID", func() {
//		Describe("with no records in the database", func() {
//			It("returns an empty slice of orders", func() {
//				orders, err = orderRepo.FindAllOrdersByUserID(userID)
//				Expect(err).To(BeNil())
//				Expect(len(orders)).To(Equal(0))
//			})
//		})
//
//		Describe("when a few records exist", func() {
//			BeforeEach(func() {
//				order1 := &models.Order{
//					Total:        1000,
//					CurrencyCode: "GBP",
//					UserID:       userID,
//					RestaurantID: 8,
//					PlacedAt:     time.Now().Add(-72 * time.Hour),
//				}
//				err = tx.Create(order1).Error
//				Expect(err).To(BeNil())
//
//				order2 := &models.Order{
//					Total:        2500,
//					CurrencyCode: "GBP",
//					UserID:       userID,
//					RestaurantID: 9,
//					PlacedAt:     time.Now().Add(-36 * time.Hour),
//				}
//				err = tx.Create(order2).Error
//				Expect(err).To(BeNil())
//
//				order3 := &models.Order{
//					Total:        600,
//					CurrencyCode: "GBP",
//					UserID:       7,
//					RestaurantID: 8,
//					PlacedAt:     time.Now().Add(-24 * time.Hour),
//				}
//				err = tx.Create(order3).Error
//				Expect(err).To(BeNil())
//			})
//
//			It("returns only the records belonging to the user, in order from latest placed_at first", func() {
//				orders, err = orderRepo.FindAllOrdersByUserID(userID)
//				Expect(err).To(BeNil())
//				Expect(len(orders)).To(Equal(2))
//				Expect(orders[0].RestaurantID).To(Equal(9))
//				Expect(orders[1].RestaurantID).To(Equal(8))
//			})
//		})
//	})
//
//	AfterEach(func() {
//		err = tx.Rollback().Error
//		Expect(err).To(BeNil())
//	})
//})
