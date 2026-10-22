package items

// ItemType دسته‌بندی آیتم‌ها را مشخص می‌کند
type ItemType string

const (
	ItemTypeCue    ItemType = "cue"    // چوب بیلیارد
	ItemTypeAvatar ItemType = "avatar" // آواتار
	ItemTypeFrame  ItemType = "frame"  // قاب
	ItemTypeBanner ItemType = "banner" // بنر
)

// InfiniteItem ساختار آیتم‌هایی است که در ستون JSONB دیتابیس ذخیره می‌شوند
type InfiniteItem struct {
	ID   string         `json:"id"`
	Type ItemType       `json:"type"`
	Meta map[string]any `json:"meta"` // شامل افکت‌های خاص مثل {"power": 1.5, "accuracy": 2}
}
