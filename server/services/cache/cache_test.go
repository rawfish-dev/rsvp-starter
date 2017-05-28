package cache_test

import (
	"time"

	. "bitbucket.org/rawfish-dev/wedding-rsvp/server/services/cache"
	"bitbucket.org/rawfish-dev/wedding-rsvp/server/testhelpers"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Cache", func() {

	var testCacheService CacheServiceProvider

	BeforeEach(func() {
		testCacheService = testhelpers.NewTestCacheService()
		Expect(testCacheService.Flush()).To(Succeed())
	})

	Context("setting", func() {
		It("should set a value without an error", func() {
			Expect(testCacheService.SetWithExpiry("some_key", "some_value", 1)).To(Succeed())
		})

		It("should not allow setting of blank keys", func() {
			Expect(testCacheService.SetWithExpiry("", "some_value", 1)).ToNot(Succeed())
		})

		It("should allow setting of blank values", func() {
			Expect(testCacheService.SetWithExpiry("some_key", "", 1)).To(Succeed())
		})
	})

	Context("existance", func() {
		It("should check existance for a key that exists correctly", func() {
			err := testCacheService.SetWithExpiry("some_key", "some_value", 1)
			Expect(err).ToNot(HaveOccurred())

			exists, err := testCacheService.Exists("some_key")
			Expect(err).ToNot(HaveOccurred())
			Expect(exists).To(BeTrue())
		})

		It("should check existance for a key that does not exists correctly", func() {
			exists, err := testCacheService.Exists("some_key")
			Expect(err).ToNot(HaveOccurred())
			Expect(exists).To(BeFalse())
		})
	})

	Context("getting", func() {
		It("should get the value of a key that exists successfully", func() {
			err := testCacheService.SetWithExpiry("some_key", "some_value", 1)
			Expect(err).ToNot(HaveOccurred())

			value, err := testCacheService.Get("some_key")
			Expect(err).ToNot(HaveOccurred())
			Expect(value).To(Equal("some_value"))
		})
	})

	Context("deleting", func() {
		It("should a key successfully", func() {
			err := testCacheService.SetWithExpiry("some_key", "some_value", 5)
			Expect(err).ToNot(HaveOccurred())

			err = testCacheService.Delete("some_key")
			Expect(err).ToNot(HaveOccurred())

			value, err := testCacheService.Get("some_key")
			Expect(err).ToNot(HaveOccurred())
			Expect(value).To(Equal(""))
		})
	})

	Context("expiry", func() {
		It("should expire keys according to their expiry times", func() {
			err := testCacheService.SetWithExpiry("some_key", "some_value", 1)
			Expect(err).ToNot(HaveOccurred())

			Eventually(func() string {
				value, err := testCacheService.Get("some_key")
				Expect(err).ToNot(HaveOccurred())
				return value
			}, "2s", "500ms").Should(Equal(""))
		})

		It("should expire keys that were overwritten according to their new expiry times", func() {
			err := testCacheService.SetWithExpiry("some_key", "some_value_1", 1)
			Expect(err).ToNot(HaveOccurred())

			// Set the expiry longer
			err = testCacheService.SetWithExpiry("some_key", "some_value_2", 4)
			Expect(err).ToNot(HaveOccurred())

			// Block for more than 1 second and check some_key is still there
			<-time.After(time.Second * time.Duration(2))

			value, err := testCacheService.Get("some_key")
			Expect(err).ToNot(HaveOccurred())
			Expect(value).To(Equal("some_value_2"))
		})
	})
})
