package telegram

import "gopkg.in/telebot.v3"

const (
	// startStickerID      = "CAACAgIAAxkBAAEM5L5m-H_KjuZObRzIDxb1pO5pmx63GwACqFYAAj5IuErfpmXhx4S2ojYE"
	studentStickerID = "CAACAgIAAxkBAAEM5LZm-H9fKOlKfc0Vr2KRlKsoRyoViAACNV0AAq7DsUovhOnvTx1HNjYE"
	gradesStickerID  = "CAACAgIAAxkBAAEM5Lhm-H99zsd-9qSRrcbpT2yF2Fr1rAACelMAApKPuEqj2t4T4kz55zYE"
	// smStickerID         = "CAACAgIAAxkBAAEM5Lpm-H-qPRbQTb-8mTkGmwzN5JRIGgACVFgAAsPWuUq6Aj9SaHW0iDYE"
	// endStickerID        = "CAACAgIAAxkBAAEM5Lxm-H-8-rYD_839kyWWIMgpzk3s_AAClFgAAhc9uEr4-BYfWKFQhDYE"
	learnMoreStickerID  = "CAACAgIAAxkBAAEM5MBm-IAcFyZWHJ-fRJuuHfvGTPkGUAACOFgAAunOuErUhJ9w1sQcrDYE"
	additionalStickerID = "CAACAgIAAxkBAAEM5MJm-IAtmjYLH7S3xU9EG1tgvyaj0QAC-lIAApNMwUpA5YfcMl1zKzYE"
	// m509StickerID       = "CAACAgIAAxkBAAEM5MRm-IBCcj0HF9deBeXpmZdWb0f2CQACvV4AAvYR-UrhuxHey0mvUDYE"
)

var (
	studentSticker    = telebot.Sticker{File: telebot.File{FileID: studentStickerID}}
	gradesSticker     = telebot.Sticker{File: telebot.File{FileID: gradesStickerID}}
	learnMoreSticker  = telebot.Sticker{File: telebot.File{FileID: learnMoreStickerID}}
	additionalSticker = telebot.Sticker{File: telebot.File{FileID: additionalStickerID}}
)
