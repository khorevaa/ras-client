package rac

import (
	"errors"
)

var ErrUnsupportedWhat = errors.New("telebot: unsupported what argument")

type Clusterable interface {
	ClusterSig() (uuid string, user string, pwd string)
}

type valued interface {
	Values() map[string]string
	Command() string
	RespondParser
}

type RespondParser interface {
	Parse(raw *RawRespond) error
}

type Auth struct {
	User string
	Pwd  string
}

func (a Auth) Sig() (usr string, pwd string) {
	return a.User, a.Pwd
}

type RawRespond struct {
	Status        bool
	raw           []byte
	parsedRespond interface{}
	Error         error
}

func newRawRespond(data []byte, err error) *RawRespond {

	res := &RawRespond{
		raw:    data,
		Error:  err,
		Status: true,
	}

	if err != nil {
		res.Status = false
	}

	return res
}

func extractOptions(how []interface{}) *DoOptions {

	var opts DoOptions

	for _, prop := range how {

		switch opt := prop.(type) {
		case *DoOptions:
			opts = *opt.copy()
		case DoOptions:
			opts = *opt.copy()
		case DoOption:
			opt(&opts)

		default:
			panic("unsupported doOption")
		}
	}

	return &opts
}

func (m *Manager) embedSendOptions(params map[string]string) {
	//if b.parseMode != ModeDefault {
	//	params["parse_mode"] = b.parseMode
	//}
	//
	//if opt == nil {
	//	return
	//}
	//
	//if opt.ReplyTo != nil && opt.ReplyTo.UUID != 0 {
	//	params["reply_to_message_id"] = strconv.Itoa(opt.ReplyTo.UUID)
	//}
	//
	//if opt.DisableWebPagePreview {
	//	params["disable_web_page_preview"] = "true"
	//}
	//
	//if opt.DisableNotification {
	//	params["disable_notification"] = "true"
	//}
	//
	//if opt.ParseMode != ModeDefault {
	//	params["parse_mode"] = opt.ParseMode
	//}
	//
	//if opt.ReplyMarkup != nil {
	//	processButtons(opt.ReplyMarkup.InlineKeyboard)
	//	replyMarkup, _ := json.Marshal(opt.ReplyMarkup)
	//	params["reply_markup"] = string(replyMarkup)
	//}
}
