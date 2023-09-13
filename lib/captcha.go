package lib

import (
	"github.com/mojocn/base64Captcha"
	"go.uber.org/zap"
	"image/color"
	"manuel71sj/go-api-template/constants"
	"time"
)

type Captcha struct {
	*base64Captcha.Captcha
}

type CaptchaStore struct {
	key    string
	redis  *Redis
	logger *zap.SugaredLogger
}

func (s CaptchaStore) getKey(v string) string {
	return s.key + ":" + v
}

func (s CaptchaStore) Set(id string, value string) error {
	err := s.redis.Set(s.getKey(id), value, time.Second*constants.CaptchaExpireTimes)
	if err != nil {
		s.logger.Errorf("captcha - error writing redis :%v", err)
		return err
	}

	return nil
}

func (s CaptchaStore) Get(id string, clear bool) string {
	var (
		key = s.getKey(id)
		val string
	)

	err := s.redis.Get(key, &val)
	if err != nil {
		s.logger.Errorf("captcha - error reading redis :%v", err)
		return ""
	}

	if !clear {
		_, err := s.redis.Delete(key)
		if err != nil {
			s.logger.Errorf("captcha - error deleting item from redis: %v", err)
		}
	}

	return val
}

func (s CaptchaStore) Verify(id, answer string, clear bool) bool {
	v := s.Get(id, clear)
	return v == answer
}

func NewCaptcha(redis Redis, logger Logger, config Config) Captcha {
	ds := base64Captcha.NewDriverString(
		config.Auth.Captcha.Height,
		config.Auth.Captcha.Width,
		config.Auth.Captcha.NoiseCount,
		2,
		4,
		"234567890abcdefghjkmnpqrstuvwxyz",
		&color.RGBA{R: 240, G: 240, B: 246, A: 246},
		base64Captcha.DefaultEmbeddedFonts,
		[]string{"wqy-microhei.ttc"},
	)

	driver := ds.ConvertFonts()
	store := CaptchaStore{
		redis:  &redis,
		key:    constants.CaptchaKeyPrefix,
		logger: logger.Zap.With(zap.String("module", "captcha")),
	}

	return Captcha{Captcha: base64Captcha.NewCaptcha(driver, store)}
}
