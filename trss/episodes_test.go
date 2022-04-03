package transmissionrss_test

import (
	"database/sql"
	"database/sql/driver"
	"regexp"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	trss "github.com/iben12/transmission-rss-go/trss"
)

type AnyTime struct{}

// Match satisfies sqlmock.Argument interface
func (a AnyTime) Match(v driver.Value) bool {
	_, ok := v.(time.Time)
	return ok
}

var _ = Describe("Episodes", func() {
	var episodes *trss.Episodes
	var sqlMock sqlmock.Sqlmock

	BeforeEach(func() {
		var db *sql.DB
		var err error

		db, sqlMock, err = sqlmock.New()
		Expect(err).NotTo(HaveOccurred())

		gdb, err := gorm.Open(mysql.New(mysql.Config{Conn: db, SkipInitializeWithVersion: true}), &gorm.Config{
			Logger: logger.Default.LogMode(logger.Silent),
		})
		Expect(err).NotTo(HaveOccurred())

		episodes = &trss.Episodes{Db: gdb}
	})
	AfterEach(func() {
		err := sqlMock.ExpectationsWereMet() // make sure all expectations were met
		Expect(err).NotTo(HaveOccurred())
	})

	Context("All()", func() {
		It("empty", func() {
			var sqlSelectAll = "SELECT * FROM `episodes` WHERE `episodes`.`deleted_at` IS NULL"

			sqlMock.ExpectQuery(regexp.QuoteMeta(sqlSelectAll)).
				WillReturnRows(sqlmock.NewRows(nil))

			result, err := episodes.All()
			Expect(err).NotTo(HaveOccurred())
			Expect(result).To(BeEmpty())
		})

		It("returns items", func() {
			episode := &trss.Episode{
				Model:     gorm.Model{ID: 1},
				ShowTitle: "Show Title",
				Title:     "Episode Title",
			}

			var sqlSelectAll = "SELECT * FROM `episodes` WHERE `episodes`.`deleted_at` IS NULL"

			sqlMock.ExpectQuery(regexp.QuoteMeta(sqlSelectAll)).
				WillReturnRows(
					sqlmock.NewRows([]string{"id", "show_title", "title"}).
						AddRow(1, episode.ShowTitle, episode.Title))

			result, err := episodes.All()
			Expect(err).NotTo(HaveOccurred())
			Expect(result[0]).To(BeEquivalentTo(*episode))
		})
	})

	Context("FindEpisode()", func() {
		It("empty", func() {
			episode := &trss.Episode{
				Model: gorm.Model{ID: 1},
			}

			var sqlSelectAll = "SELECT * FROM `episodes` WHERE `episodes`.`id` = ? AND `episodes`.`deleted_at` IS NULL ORDER BY `episodes`.`id` LIMIT 1"

			sqlMock.ExpectQuery(regexp.QuoteMeta(sqlSelectAll)).WithArgs(episode.Model.ID).
				WillReturnRows(sqlmock.NewRows(nil))

			_, err := episodes.FindEpisode(episode)
			Expect(err).To(MatchError(gorm.ErrRecordNotFound))
		})

		It("returns found item", func() {
			episode := &trss.Episode{
				Model: gorm.Model{ID: 1},
			}

			var sqlSelectAll = "SELECT * FROM `episodes` WHERE `episodes`.`id` = ? AND `episodes`.`deleted_at` IS NULL ORDER BY `episodes`.`id` LIMIT 1"

			sqlMock.ExpectQuery(regexp.QuoteMeta(sqlSelectAll)).WithArgs(episode.Model.ID).
				WillReturnRows(sqlmock.NewRows([]string{"id", "show_title", "title"}).AddRow(1, episode.ShowTitle, episode.Title))

			result, err := episodes.FindEpisode(episode)
			Expect(err).NotTo(HaveOccurred())
			Expect(result).To(BeEquivalentTo(*episode))
		})
	})

	Context("AddEpisode()", func() {
		It("saves item", func() {
			episode := &trss.Episode{
				ShowId:    "1",
				EpisodeId: "2",
				ShowTitle: "Show Title",
				Title:     "Episode Title",
				Link:      "url",
			}

			var sqlSelectAll = "INSERT INTO `episodes` (`created_at`,`updated_at`,`deleted_at`,`show_id`,`episode_id`,`show_title`,`title`,`link`) VALUES (?,?,?,?,?,?,?,?)"

			sqlMock.ExpectBegin()
			sqlMock.ExpectExec(regexp.QuoteMeta(sqlSelectAll)).
				WithArgs(
					AnyTime{},
					AnyTime{},
					nil,
					episode.ShowId,
					episode.EpisodeId,
					episode.ShowTitle,
					episode.Title,
					episode.Link).
				WillReturnResult(sqlmock.NewResult(1, 1))
			sqlMock.ExpectCommit()

			err := episodes.AddEpisode(episode)

			Expect(err).NotTo(HaveOccurred())
			Expect(episode.ID).To(Equal(uint(1)))
		})
	})
})
