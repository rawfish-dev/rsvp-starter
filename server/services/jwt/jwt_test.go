package jwt_test

import (
	"time"

	"github.com/rawfish-dev/rsvp-starter/server/config"
	"github.com/rawfish-dev/rsvp-starter/server/interfaces"
	. "github.com/rawfish-dev/rsvp-starter/server/services/jwt"

	"github.com/Sirupsen/logrus"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"golang.org/x/net/context"
)

var _ = Describe("Jwt", func() {

	var testJWTService interfaces.JWTServiceProvider
	var jwtConfig config.JWTConfig

	BeforeEach(func() {
		ctxlogger := logrus.New()
		ctx := context.Background()
		ctx = context.WithValue(ctx, "logger", ctxlogger)

		jwtConfig = config.JWTConfig{
			HMACSecret:  "some-secret-hmac",
			TokenIssuer: "rsvp-starter-test",
		}

		testJWTService = NewService(ctx, jwtConfig)
	})

	Context("creation", func() {

		It("should create a jwt token given expiry time and additional claims", func() {
			additionalClaims := map[string]string{
				"userID": "123123",
			}

			token, err := testJWTService.GenerateAuthToken(additionalClaims, 10)
			Expect(err).ToNot(HaveOccurred())
			Expect(token).ToNot(BeEmpty())
		})

		It("should return an error if the expiry time is invalid", func() {
			additionalClaims := map[string]string{
				"userID": "123123",
			}

			token, err := testJWTService.GenerateAuthToken(additionalClaims, 0)
			Expect(err).To(HaveOccurred())
			Expect(token).To(BeEmpty())
		})
	})

	Context("parsing", func() {

		var token string

		BeforeEach(func() {
			additionalClaims := map[string]string{
				"userID": "123123",
			}

			token, _ = testJWTService.GenerateAuthToken(additionalClaims, 10)
			Expect(token).ToNot(BeEmpty())
		})

		It("should return an error if the JWT is malformed", func() {
			token = "abc" + token

			claims, err := testJWTService.ParseToken(token)
			Expect(err).To(HaveOccurred())
			Expect(claims).To(BeNil())
		})

		XIt("should return an error if the JWT is unverifiable", func() {

		})

		XIt("should return an error if the JWT issuer did not match", func() {

		})

		It("should return an error if the JWT has an invalid signature", func() {
			token += "abc"

			claims, err := testJWTService.ParseToken(token)
			Expect(err).To(HaveOccurred())
			Expect(claims).To(BeNil())
		})

		It("should return an error if the JWT is expired", func() {
			additionalClaims := map[string]string{
				"userID": "123123",
			}

			token, _ = testJWTService.GenerateAuthToken(additionalClaims, 1)
			Expect(token).ToNot(BeEmpty())

			time.Sleep(2 * time.Second)

			claims, err := testJWTService.ParseToken(token)
			Expect(err).To(HaveOccurred())
			Expect(claims).To(BeNil())
		})

		It("should not return an error if the JWT is valid", func() {
			claims, err := testJWTService.ParseToken(token)
			Expect(err).ToNot(HaveOccurred())
			Expect(claims["userID"].(string)).To(Equal("123123"))
			Expect(claims["iss"].(string)).To(Equal("rsvp-starter-test"))
			Expect(claims["exp"]).ToNot(BeNil())
			Expect(claims["iat"]).ToNot(BeNil())
		})
	})

	Context("validation", func() {

		var token string

		BeforeEach(func() {
			additionalClaims := map[string]string{
				"userID": "123123",
			}

			token, _ = testJWTService.GenerateAuthToken(additionalClaims, 10)
			Expect(token).ToNot(BeEmpty())
		})

		It("should return an error if the JWT is malformed", func() {

			token = "abc" + token

			valid := testJWTService.IsAuthTokenValid(token)
			Expect(valid).To(BeFalse())
		})

		XIt("should return an error if the JWT is unverifiable", func() {

		})

		XIt("should return an error if the JWT issuer did not match", func() {

		})

		It("should return an error if the JWT has an invalid signature", func() {
			token += "abc"

			valid := testJWTService.IsAuthTokenValid(token)
			Expect(valid).To(BeFalse())
		})

		It("should return an error if the JWT is expired", func() {
			additionalClaims := map[string]string{
				"userID": "123123",
			}

			token, _ = testJWTService.GenerateAuthToken(additionalClaims, 1)
			Expect(token).ToNot(BeEmpty())

			Eventually(func() bool {
				return testJWTService.IsAuthTokenValid(token)
			}, 3, 0.2).Should(BeFalse())
		})

		It("should not return an error if the JWT is valid", func() {
			valid := testJWTService.IsAuthTokenValid(token)
			Expect(valid).To(BeTrue())
		})
	})
})
