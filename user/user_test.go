package user

import (
	"fmt"
	"reflect"
	"regexp"
	"testing"

	gio "github.com/digisan/gotk/io"
	lk "github.com/digisan/logkit"
)

func TestUserFieldValue(t *testing.T) {
	user := &User{
		Active:     "T",
		UName:      "unique-user-name",
		Email:      "hello@abc.com",
		Name:       "test-name",
		Password:   "123456789a",
		Regtime:    "",
		Phone:      "",
		Country:    "",
		City:       "",
		Addr:       "",
		SysRole:    "",
		MemLevel:   "",
		MemExpire:  "",
		NationalID: "9876543210",
		Gender:     "",
		Position:   "professor",
		Title:      "",
		Employer:   "",
		Tags:       "",
		AvatarType: "image/png",
		key:        [16]byte{},
		Avatar:     []byte("******"),
	}

	fmt.Println(user.FieldValue("UName"))
	// fmt.Println(user.FieldValue("key")) // panic
	fmt.Println(user.FieldValue("Avatar"))

	user.AddTags("abc", "def")
	fmt.Println("tags:", user.GetTags())

	user.RmTags("abc")
	fmt.Println("tags:", user.GetTags())
}

func TestUser(t *testing.T) {
	user := &User{
		Active:     "T",
		UName:      "unique-user-name",
		Email:      "hello@abc.com",
		Name:       "test-name",
		Password:   "123456789a",
		Regtime:    "",
		Phone:      "",
		Country:    "",
		City:       "",
		Addr:       "",
		SysRole:    "",
		MemLevel:   "",
		MemExpire:  "",
		NationalID: "9876543210",
		Gender:     "",
		Position:   "professor",
		Title:      "",
		Employer:   "",
		Tags:       "",
		AvatarType: "image/png",
		key:        [16]byte{1, 2, 3, 4, 5, 6, 7, 8, 9},
		Avatar:     []byte("******"),
	}

	// ava := make([]byte, 0, 20000000)
	// for i := 0; i < 20000000; i++ {
	// 	ava = append(ava, byte(i%100))
	// }
	// user.Avatar = ava

	info, key := user.Marshal()
	fmt.Println("user.key", user.key)

	fmt.Println()

	fmt.Println(user)
	fmt.Println()

	// key[1] = 2

	user1 := &User{}
	user1.Unmarshal(info, key)
	fmt.Println(user1)

	fmt.Println("user == user1 :", user == user1)
	fmt.Println("reflect.DeepEqual(*user, *user1) :", reflect.DeepEqual(*user, *user1))
	lk.FailOnErrWhen(!reflect.DeepEqual(*user, *user1), "%v", fmt.Errorf("Marshal-Unmarshal ERROR"))
}

func TestIterTags(t *testing.T) {
	fmt.Println(ListUserField())
	fmt.Println(ListUserValidator())
}

// *** Auto Generate Validate Field Tag Const *** //
func TestMakeUserFieldTag(t *testing.T) {
	const pkg = "valfield"
	const file = pkg + "/auto-tags.go"
	gio.MustAppendFile(file, []byte("package "+pkg+"\n"), true)
	gio.MustAppendFile(file, []byte("const ("), true)
	defer gio.MustAppendFile(file, []byte("\n)"), false)

	const TAG = "validate"
	// r := regexp.MustCompile(`((required),?)|((email),?)`) // exclude default validator tags
	r := regexp.MustCompile(`(required),?`)

	typ := reflect.TypeOf(User{})
	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)
		tag := field.Tag.Get(TAG)
		// fmt.Printf("%d. %v (%v), tag: '%v'\n", i+1, field.Name, field.Type.Name(), tag)
		tag = r.ReplaceAllString(tag, "")
		if len(tag) > 0 {
			line := fmt.Sprintf("\t%s = \"%s\"", field.Name, tag)
			gio.MustAppendFile(file, []byte(line), true)
		}
	}
}
