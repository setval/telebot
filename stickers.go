package telebot

import (
	"encoding/json"
	"strconv"
)

// StickerSet represents a sticker set.
type StickerSet struct {
	Name          string        `json:"name"`
	Title         string        `json:"title"`
	Animated      bool          `json:"is_animated"`
	Stickers      []Sticker     `json:"stickers"`
	Thumbnail     *Photo        `json:"thumb"`
	PNG           *File         `json:"png_sticker"`
	TGS           *File         `json:"tgs_sticker"`
	Emojis        string        `json:"emojis"`
	ContainsMasks bool          `json:"contains_masks"`
	MaskPosition  *MaskPosition `json:"mask_position"`
}

// MaskPosition describes the position on faces where
// a mask should be placed by default.
type MaskPosition struct {
	Feature MaskFeature `json:"point"`
	XShift  float32     `json:"x_shift"`
	YShift  float32     `json:"y_shift"`
	Scale   float32     `json:"scale"`
}

// UploadSticker uploads a PNG file with a sticker for later use.
func (b *Bot) UploadSticker(to Recipient, png *File) (*File, error) {
	files := map[string]File{
		"png_sticker": *png,
	}
	params := map[string]string{
		"user_id": to.Recipient(),
	}

	data, err := b.sendFiles("uploadStickerFile", files, params)
	if err != nil {
		return nil, err
	}

	var resp struct {
		Result File
	}
	if err := json.Unmarshal(data, &resp); err != nil {
		return nil, wrapError(err)
	}
	return &resp.Result, nil
}

// StickerSet returns a sticker set on success.
func (b *Bot) StickerSet(name string) (*StickerSet, error) {
	data, err := b.Raw("getStickerSet", map[string]string{"name": name})
	if err != nil {
		return nil, err
	}

	var resp struct {
		Result *StickerSet
	}
	if err := json.Unmarshal(data, &resp); err != nil {
		return nil, wrapError(err)
	}
	return resp.Result, nil
}

// CreateStickerSet creates a new sticker set.
func (b *Bot) CreateStickerSet(to Recipient, s StickerSet) error {
	files := make(map[string]File)
	if s.PNG != nil {
		files["png_sticker"] = *s.PNG
	}
	if s.TGS != nil {
		files["tgs_sticker"] = *s.TGS
	}

	params := map[string]string{
		"user_id":        to.Recipient(),
		"name":           s.Name,
		"title":          s.Title,
		"emojis":         s.Emojis,
		"contains_masks": strconv.FormatBool(s.ContainsMasks),
	}

	if s.MaskPosition != nil {
		data, _ := json.Marshal(&s.MaskPosition)
		params["mask_position"] = string(data)
	}

	_, err := b.sendFiles("createNewStickerSet", files, params)
	return err
}

// AddSticker adds a new sticker to the existing sticker set.
func (b *Bot) AddSticker(to Recipient, s StickerSet) error {
	files := make(map[string]File)
	if s.PNG != nil {
		files["png_sticker"] = *s.PNG
	} else if s.TGS != nil {
		files["tgs_sticker"] = *s.TGS
	}

	params := map[string]string{
		"user_id": to.Recipient(),
		"name":    s.Name,
		"emojis":  s.Emojis,
	}

	if s.MaskPosition != nil {
		data, _ := json.Marshal(&s.MaskPosition)
		params["mask_position"] = string(data)
	}

	_, err := b.sendFiles("addStickerToSet", files, params)
	return err
}

// SetStickerPosition moves a sticker in set to a specific position.
func (b *Bot) SetStickerPosition(sticker string, position int) error {
	params := map[string]string{
		"sticker":  sticker,
		"position": strconv.Itoa(position),
	}

	_, err := b.Raw("setStickerPositionInSet", params)
	return err
}

// DeleteSticker deletes a sticker from a set created by the bot.
func (b *Bot) DeleteSticker(sticker string) error {
	_, err := b.Raw("deleteStickerFromSet", map[string]string{"sticker": sticker})
	return err

}

// SetStickerSetThumb sets a thumbnail of the sticker set.
// Animated thumbnails can be set for animated sticker sets only.
//
// Thumbnail must be a PNG image, up to 128 kilobytes in size
// and have width and height exactly 100px, or a TGS animation
// up to 32 kilobytes in size.
//
// Animated sticker set thumbnail can't be uploaded via HTTP URL.
//
func (b *Bot) SetStickerSetThumb(to Recipient, s StickerSet) error {
	files := make(map[string]File)
	if s.PNG != nil {
		files["thumb"] = *s.PNG
	} else if s.TGS != nil {
		files["thumb"] = *s.TGS
	}

	params := map[string]string{
		"name":    s.Name,
		"user_id": to.Recipient(),
	}

	_, err := b.sendFiles("setStickerSetThumb", files, params)
	return err
}
