package callbacks

import "github.com/mviner000/eyymi/src/germ"

type BeforeCreateInterface interface {
	BeforeCreate(*germ.DB) error
}

type AfterCreateInterface interface {
	AfterCreate(*germ.DB) error
}

type BeforeUpdateInterface interface {
	BeforeUpdate(*germ.DB) error
}

type AfterUpdateInterface interface {
	AfterUpdate(*germ.DB) error
}

type BeforeSaveInterface interface {
	BeforeSave(*germ.DB) error
}

type AfterSaveInterface interface {
	AfterSave(*germ.DB) error
}

type BeforeDeleteInterface interface {
	BeforeDelete(*germ.DB) error
}

type AfterDeleteInterface interface {
	AfterDelete(*germ.DB) error
}

type AfterFindInterface interface {
	AfterFind(*germ.DB) error
}
