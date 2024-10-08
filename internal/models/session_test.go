package models_test

import (
	"time"

	. "github.com/jimmykarily/quizmaker/internal/models"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Session", func() {
	var session Session
	BeforeEach(func() {
		session = Session{}
	})

	Describe("#HasExpiredQuestions", func() {
		When("there are expired questions", func() {
			BeforeEach(func() {
				session.Questions = []Question{
					{
						Text:           "started and expired question",
						StartedAt:      time.Now().Add(-10 * time.Second),
						AllowedSeconds: 5,
					},
					{
						Text:           "not started question",
						AllowedSeconds: 5,
					},
					{
						Text:           "started but not expired question",
						StartedAt:      time.Now(),
						AllowedSeconds: 1000000,
					},
				}
			})

			It("returns true", func() {
				Expect(session.HasExpiredQuestions()).To(BeTrue())
			})
		})

		// An answered question is never considered "expired"
		When("there are no expired questions", func() {
			BeforeEach(func() {
				session.Questions = []Question{
					{
						Text:           "started, anwswered (and 'expired') question",
						StartedAt:      time.Now().Add(-10),
						AllowedSeconds: 5,
						UserAnswer:     2,
					},
					{
						Text:           "not started question",
						AllowedSeconds: 5,
					},
					{
						Text:           "started but not expired question",
						StartedAt:      time.Now(),
						AllowedSeconds: 1000000,
					},
				}
			})

			It("returns false", func() {
				Expect(session.HasExpiredQuestions()).To(BeFalse())
			})
		})
	})

	Describe("#CurrentQuestion", func() {
		When("there is an expired question", func() {
			BeforeEach(func() {
				session.Questions = []Question{
					{
						Text:           "started and expired question",
						StartedAt:      time.Now().Add(-10 * time.Second),
						AllowedSeconds: 5,
						Index:          1,
					},
					{
						Text:           "not started question",
						AllowedSeconds: 5,
						Index:          2,
					},
				}
			})
			When("there is a started question", func() {
				BeforeEach(func() {
					session.Questions = append(session.Questions, Question{
						Text:           "started but not expired question",
						StartedAt:      time.Now(),
						AllowedSeconds: 1000000,
						Index:          3,
					})
				})
				It("returns the started question regardless of Index", func() {
					q, err := session.CurrentQuestion()
					Expect(err).ToNot(HaveOccurred())
					Expect(q.Text).To(Equal("started but not expired question"))
					Expect(q.Index).To(Equal(3))
				})
			})

			It("returns the next question", func() {
				q, err := session.CurrentQuestion()
				Expect(err).ToNot(HaveOccurred())
				Expect(q.Text).To(Equal("not started question"))
				Expect(q.Index).To(Equal(2))
			})
		})
		When("there are no expired questions", func() {
			When("there is a started question", func() {
				BeforeEach(func() {
					session.Questions = []Question{
						{
							Text:           "answered and expired question",
							StartedAt:      time.Now().Add(-10 * time.Second),
							AllowedSeconds: 5,
							UserAnswer:     2,
							Index:          1,
						},
						{
							Text:           "not answered, not expired, not started question",
							AllowedSeconds: 10,
							Index:          2,
						},
						{
							Text:           "not answered, not expired, started question",
							StartedAt:      time.Now(),
							AllowedSeconds: 5000,
							Index:          3,
						},
					}
				})

				It("returns the started question (regardless of Index)", func() {
					q, err := session.CurrentQuestion()
					Expect(err).ToNot(HaveOccurred())
					Expect(q.Text).To(Equal("not answered, not expired, started question"))
				})
			})

			When("there is no started question", func() {
				BeforeEach(func() {
					session.Questions = []Question{
						{
							Text:           "answered and expired question",
							StartedAt:      time.Now().Add(-10 * time.Second),
							AllowedSeconds: 5,
							UserAnswer:     2,
							Index:          1,
						},
						{
							Text:           "not answered, not expired, not started question (3)",
							AllowedSeconds: 10,
							Index:          3,
						},
						{
							Text:           "not answered, not expired, not started question (2)",
							AllowedSeconds: 10,
							Index:          2,
						},
					}
				})

				It("returns the first not answered question based on Index", func() {
					q, err := session.CurrentQuestion()
					Expect(err).ToNot(HaveOccurred())
					Expect(q.Text).To(Equal("not answered, not expired, not started question (2)"))
				})
			})

			When("there is no unanswered question left", func() {
				BeforeEach(func() {
					session.Questions = []Question{
						{
							Text:           "answered and expired question",
							StartedAt:      time.Now().Add(-10 * time.Second),
							AllowedSeconds: 5,
							UserAnswer:     2,
							Index:          1,
						},
						{
							Text:           "answered, not expired, not started question (3)",
							AllowedSeconds: 10,
							Index:          3,
							StartedAt:      time.Now().Add(-10 * time.Second),
							UserAnswer:     2,
						},
						{
							Text:           "answered, not expired, started question (2)",
							AllowedSeconds: 10,
							Index:          2,
							UserAnswer:     1,
						},
					}
				})

				It("returns no error and an empty question", func() {
					q, err := session.CurrentQuestion()
					Expect(err).ToNot(HaveOccurred())
					Expect(q.ID).To(BeZero())
					Expect(q.Index).To(BeZero())
					Expect(q.Text).To(BeZero())
				})
			})
		})
	})

	Describe("#UpdateCacheColumns", func() {
		BeforeEach(func() {
			session.Questions = []Question{
				{
					Text:           "started, expired, not answered",
					StartedAt:      time.Now().Add(-10 * time.Second),
					AllowedSeconds: 5,
					UserAnswer:     0,
					RightAnswer:    2,
				},
				{
					Text:           "not started",
					AllowedSeconds: 5,
				},
				{
					Text:           "started, not answered, not expired",
					StartedAt:      time.Now(),
					AllowedSeconds: 1000000,
				},
				{
					Text:           "correctly answered",
					StartedAt:      time.Now().Add(-10 * time.Second),
					AllowedSeconds: 1,
					UserAnswer:     2,
					RightAnswer:    2,
				},
				{
					Text:           "wrongly answered",
					StartedAt:      time.Now().Add(-10 * time.Second),
					AllowedSeconds: 1,
					UserAnswer:     1,
					RightAnswer:    2,
				},
			}
		})

		It("updates the Score field", func() {
			Expect(session.Score).To(BeZero())
			Expect(session.Complete).To(BeFalse())
			session.UpdateCacheColumns()
			Expect(session.Score).To(Equal(33))
			Expect(session.Complete).To(BeFalse())
		})

		It("updated the Complete field", func() {
			for i := range session.Questions {
				session.Questions[i].StartedAt = time.Now().Add(-10 * time.Second)
				session.Questions[i].UserAnswer = 2
			}
			Expect(session.Complete).To(BeFalse())
			session.UpdateCacheColumns()
			Expect(session.Complete).To(BeTrue())
		})
	})
})
