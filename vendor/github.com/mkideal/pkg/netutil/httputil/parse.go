package httputil

import (
	"net/http"
	"strconv"

	"github.com/mkideal/pkg/option"
)

func getArgument(r *http.Request, key string, required ...bool) (value string, err error) {
	if option.Bool(false, required...) {
		if r.Form == nil {
			r.ParseMultipartForm(32 << 20)
		}
		if vs := r.Form[key]; len(vs) > 0 {
			value = vs[0]
		} else {
			err = ErrMissingRequiredArgument
		}
	} else {
		value = r.FormValue(key)
	}
	return
}

func parseInt64(r *http.Request, key string, required ...bool) (int64, error) {
	if value, err := getArgument(r, key, required...); err != nil {
		return 0, err
	} else {
		if value == "" {
			return 0, nil
		}
		return strconv.ParseInt(value, 0, 64)
	}
}

func parseUint64(r *http.Request, key string, required ...bool) (uint64, error) {
	if value, err := getArgument(r, key, required...); err != nil {
		return 0, err
	} else {
		if value == "" {
			return 0, nil
		}
		return strconv.ParseUint(value, 0, 64)
	}
}

func parseFloat64(r *http.Request, key string, required ...bool) (float64, error) {
	if value, err := getArgument(r, key, required...); err != nil {
		return 0, err
	} else {
		if value == "" {
			return 0, nil
		}
		return strconv.ParseFloat(value, 64)
	}
}

func ParseInt(r *http.Request, key string, required ...bool) (int, error) {
	i, err := parseInt64(r, key, required...)
	if err != nil {
		return 0, err
	}
	return int(i), nil
}

func ParseInt8(r *http.Request, key string, required ...bool) (int8, error) {
	i, err := parseInt64(r, key, required...)
	if err != nil {
		return 0, err
	}
	return int8(i), nil
}

func ParseInt16(r *http.Request, key string, required ...bool) (int16, error) {
	i, err := parseInt64(r, key, required...)
	if err != nil {
		return 0, err
	}
	return int16(i), nil
}

func ParseInt32(r *http.Request, key string, required ...bool) (int32, error) {
	i, err := parseInt64(r, key, required...)
	if err != nil {
		return 0, err
	}
	return int32(i), nil
}

func ParseInt64(r *http.Request, key string, required ...bool) (int64, error) {
	return parseInt64(r, key, required...)
}

func ParseUint(r *http.Request, key string, required ...bool) (uint, error) {
	i, err := parseUint64(r, key, required...)
	if err != nil {
		return 0, err
	}
	return uint(i), nil
}

func ParseUint8(r *http.Request, key string, required ...bool) (uint8, error) {
	i, err := parseUint64(r, key, required...)
	if err != nil {
		return 0, err
	}
	return uint8(i), nil
}

func ParseUint16(r *http.Request, key string, required ...bool) (uint16, error) {
	i, err := parseUint64(r, key, required...)
	if err != nil {
		return 0, err
	}
	return uint16(i), nil
}

func ParseUint32(r *http.Request, key string, required ...bool) (uint32, error) {
	i, err := parseUint64(r, key, required...)
	if err != nil {
		return 0, err
	}
	return uint32(i), nil
}

func ParseUint64(r *http.Request, key string, required ...bool) (uint64, error) {
	return parseUint64(r, key, required...)
}

func ParseFloat32(r *http.Request, key string, required ...bool) (float32, error) {
	f, err := parseFloat64(r, key, required...)
	if err != nil {
		return 0, err
	}
	return float32(f), err
}

func ParseFloat64(r *http.Request, key string, required ...bool) (float64, error) {
	return parseFloat64(r, key, required...)
}

func ParseBool(r *http.Request, key string, required ...bool) (bool, error) {
	if value, err := getArgument(r, key, required...); err != nil {
		return false, err
	} else {
		if value == "" {
			return false, nil
		}
		return strconv.ParseBool(value)
	}
}

func ParseString(r *http.Request, key string, required ...bool) (string, error) {
	if value, err := getArgument(r, key, required...); err != nil {
		return "", err
	} else {
		return value, nil
	}
}
