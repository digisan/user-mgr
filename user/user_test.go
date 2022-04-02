package user

import (
	"fmt"
	"reflect"
	"regexp"
	"testing"
	"time"

	gio "github.com/digisan/gotk/io"
	lk "github.com/digisan/logkit"
)

func TestIterTags(t *testing.T) {
	fmt.Println(ListField(User{}, User{}.Core, User{}.Profile, User{}.Admin))
	fmt.Println(ListValidator(User{}, User{}.Core, User{}.Profile, User{}.Admin))
}

func TestFieldValue(t *testing.T) {
	user := &User{
		Core{
			UName:    "unique-user-name",
			Email:    "hello@abc.com",
			Password: "123456789a",
			key:      [16]byte{},
		},
		Profile{
			Name:       "test-name",
			Phone:      "",
			Country:    "",
			City:       "",
			Addr:       "",
			PersonalID: "9876543210",
			Gender:     "",
			DOB:        "",
			Position:   "professor",
			Title:      "",
			Employer:   "",
			Bio:        "",
			AvatarType: "image/png",
			Avatar:     []byte("******"),
		},
		Admin{
			Active:    true,
			SysRole:   "",
			MemLevel:  0,
			MemExpire: time.Time{},
			Regtime:   time.Now(),
			Official:  false,
			Tags:      "",
		},
	}

	fmt.Println(FieldValue(user, "UName"))
	// fmt.Println(user.FieldValue("key")) // panic
	fmt.Println(FieldValue(user, "Avatar"))

	user.AddTags("abc", "def")
	fmt.Println("tags:", user.GetTags())

	user.RmTags("abc")
	fmt.Println("tags:", user.GetTags())
}

func TestUser(t *testing.T) {

	user := &User{
		Core{
			UName:    "unique-user-name",
			Email:    "hello@abc.com",
			Password: "123456789a",
			key:      [16]byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
		},
		Profile{
			Name:       "test-name",
			Phone:      "",
			Country:    "",
			City:       "",
			Addr:       "",
			PersonalID: "987654321011",
			Gender:     "",
			DOB:        "",
			Position:   "professor",
			Title:      "",
			Employer:   "",
			Bio:        "",
			AvatarType: "image/png",
			Avatar:     []byte("******"),
		},
		Admin{
			Active:    true,
			SysRole:   "",
			MemLevel:  3,
			MemExpire: time.Time{},
			Regtime:   time.Now().Truncate(time.Second), // must Truncate, otherwise Unmarshal error
			Official:  false,
			Tags:      "",
		},
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
	fmt.Println("reflect.DeepEqual(*user.COre, *user1.Core) :", reflect.DeepEqual(user.Core, user1.Core))
	fmt.Println("reflect.DeepEqual(*user.Profile, *user1.Profile) :", reflect.DeepEqual(user.Profile, user1.Profile))
	fmt.Println("reflect.DeepEqual(*user.Admin, *user1.Admin) :", reflect.DeepEqual(user.Admin, user1.Admin))
	lk.FailOnErrWhen(!reflect.DeepEqual(*user, *user1), "%v", fmt.Errorf("Marshal-Unmarshal ERROR"))
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

	for _, obj := range []any{Core{}, Profile{}, Admin{}} {
		typ := reflect.TypeOf(obj)
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

}
