package utils

import (
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"github.com/rs/zerolog/log"
	"net/url"
	"regexp"
)

// RegisterCustomBindings - Bindings for form validation.
func RegisterCustomBindings() {
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		// Optional field, but require a six-digit hexadecimal color without the preceding #
		_ = v.RegisterValidation("hex", func(fl validator.FieldLevel) bool {
			val := fl.Field().String()
			if val == "" {
				return true
			}
			matched, _ := regexp.MatchString("^[\\da-fA-f]{6}$", val)
			return matched
		})

		// Optional field, but require a URL
		_ = v.RegisterValidation("optionalurl", func(fl validator.FieldLevel) bool {
			val := fl.Field().String()
			if val == "" {
				return true
			}
			_, err := url.ParseRequestURI(val)
			return err == nil
		})

		_ = v.RegisterValidation("regexpattern", func(fl validator.FieldLevel) bool {
			f := fl.Field().Interface().([]SpyglassNotificationRegexPattern)
			for _, p := range f {
				if !(p.Field == "game" || p.Field == "title") {
					return false
				}
				if p.Pattern == "" {
					return false
				}
				_, err := CompileWithTimeout(p)
				if err != nil {
					log.Debug().Err(err).Msg("RegExp failed compilation in <10ms")
					return false
				}
			}
			return true
		})

		// Require at least one digit or letter proceeded by 0-24 digits, letters, or underscores
		// https://discuss.dev.twitch.tv/t/twitch-channel-name-regex/3855
		_ = v.RegisterValidation("twitchusername", func(fl validator.FieldLevel) bool {
			matched, _ := regexp.MatchString("^[\\da-zA-Z]\\w{0,24}$", fl.Field().String())
			return matched
		})
	}
}
