package domain

type ProfileRepo interface {
	Save(profile Profile) (Profile, error)
}
