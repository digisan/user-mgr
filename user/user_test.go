package user

import (
	"fmt"
	"reflect"
	"regexp"
	"testing"

	gio "github.com/digisan/gotk/io"
	lk "github.com/digisan/logkit"
)

func TestUser(t *testing.T) {
	user := &User{
		Active:     "T",
		UName:      "unique-user-name",
		Email:      "hello@abc.com",
		Name:       "test-name",
		Password:   "123456789a",
		Regtime:    "",
		Phone:      "",
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
		Avatar:     []byte("******"),
		key:        "",
	}

	user.AddTags("abc", "def")
	fmt.Println("tags:", user.GetTags())

	user.RmTags("abc")
	fmt.Println("tags:", user.GetTags())

	info, key := user.Marshal()
	fmt.Println("user.key", user.key)

	fmt.Println(user)
	fmt.Println()

	// key[1] = '7'

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
