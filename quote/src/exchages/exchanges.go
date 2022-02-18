package exchages

type Exchange interface {
	Save(track string, base string)
	GetName() string
}
