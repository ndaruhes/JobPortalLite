package validators

import (
	lang "job-portal-lite/shared/validators/lang"
	"reflect"
	"strings"

	enTranslations "github.com/go-playground/validator/v10/translations/en"

	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	"golang.org/x/exp/slices"
)

func ValidateStruct(payload interface{}) map[string]string {
	validate := validator.New()
	english := en.New()
	uni := ut.New(english, english)
	trans, _ := uni.GetTranslator("en")

	registerTagName(validate)
	registerTranslation(validate, lang.CustomEnValidatorTranslation, trans)

	_ = enTranslations.RegisterDefaultTranslations(validate, trans)

	validate.RegisterValidation("num_boolean", validateNumberBoolean)
	validate.RegisterValidation("phone_code", validatePhoneCode)

	err := validate.Struct(payload)
	if err != nil {
		validationErrors := err.(validator.ValidationErrors)
		errors := make(map[string]string)

		for _, fieldError := range validationErrors {
			field := strings.ToLower(fieldError.Field())
			errors[field] = fieldError.Translate(trans)
		}

		return errors
	}

	return nil
}

func registerTagName(validate *validator.Validate) {
	validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}

		if name == "" {
			name = strings.SplitN(fld.Tag.Get("form"), ",", 2)[0]
			if name == "-" {
				return ""
			}
		}
		return name
	})
}

func registerTranslation(validate *validator.Validate, rules map[string]string, trans ut.Translator) {
	for key, value := range rules {
		validate.RegisterTranslation(key, trans, func(ut ut.Translator) error {
			return ut.Add(key, value, true)
		}, func(ut ut.Translator, fe validator.FieldError) string {
			name := strings.TrimSuffix(fe.Field(), "_id")
			name = strings.ReplaceAll(name, "_", " ")
			t, _ := ut.T(fe.Tag(), name)
			return t
		})
	}
}

func validateNumberBoolean(fl validator.FieldLevel) bool {
	value := fl.Field().Interface()
	return value == 0 || value == 1
}

func validatePhoneCode(fl validator.FieldLevel) bool {
	value := fl.Field().Interface()
	if value.(string) == "" {
		return true
	}
	codes := []string{"+93", "+355", "+213", "+1-684", "+376", "+244", "+1-264", "+672", "+1-268", "+54", "+374", "+297", "+61", "+43", "+994", "+1-242", "+973", "+880", "+1-246", "+375", "+32", "+501", "+229", "+1-441", "+975", "+591", "+387", "+267", "+55", "+246", "+1-284", "+673", "+359", "+226", "+257", "+855", "+237", "+1", "+238", "+1-345", "+236", "+235", "+56", "+86", "+61", "+57", "+269", "+682", "+506", "+385", "+53", "+599", "+357", "+420", "+243", "+45", "+253", "+1-767", "+1-809", "+1-829", "+1-849", "+670", "+593", "+20", "+503", "+240", "+291", "+372", "+251", "+500", "+298", "+679", "+358", "+33", "+689", "+241", "+220", "+995", "+49", "+233", "+350", "+30", "+299", "+1-473", "+1-671", "+502", "+44-1481", "+224", "+245", "+592", "+509", "+504", "+852", "+36", "+354", "+91", "+62", "+98", "+964", "+353", "+44-1624", "+972", "+39", "+225", "+1-876", "+81", "+44-1534", "+962", "+7", "+254", "+686", "+383", "+965", "+996", "+856", "+371", "+961", "+266", "+231", "+218", "+423", "+370", "+352", "+853", "+389", "+261", "+265", "+60", "+960", "+223", "+356", "+692", "+222", "+230", "+262", "+52", "+691", "+373", "+377", "+976", "+382", "+1-664", "+212", "+258", "+95", "+264", "+674", "+977", "+31", "+599", "+687", "+64", "+505", "+227", "+234", "+683", "+850", "+1-670", "+47", "+968", "+92", "+680", "+970", "+507", "+675", "+595", "+51", "+63", "+64", "+48", "+351", "+1-787, 1-939", "+974", "+242", "+262", "+40", "+7", "+250", "+590", "+290", "+1-869", "+1-758", "+590", "+508", "+1-784", "+685", "+378", "+239", "+966", "+221", "+381", "+248", "+232", "+65", "+1-721", "+421", "+386", "+677", "+252", "+27", "+82", "+211", "+34", "+94", "+249", "+597", "+47", "+268", "+46", "+41", "+963", "+886", "+992", "+255", "+66", "+228", "+690", "+676", "+1-868", "+216", "+90", "+993", "+1-649", "+688", "+1-340", "+256", "+380", "+971", "+44", "+1", "+598", "+998", "+678", "+379", "+58", "+84", "+681", "+212", "+967", "+260", "+263"}
	return slices.Contains(codes, value.(string))
}
