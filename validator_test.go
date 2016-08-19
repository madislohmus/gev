package main

import (
	"testing"
)

var (
	validEmails = []string{
		`first.last@iana.org`,
		`(comment)first.last@iana.org`,
		`()first.last@iana.org`,
		`(comment)first.last(comment)@iana.org`,
		`first.last(comment)@iana.org`,
		`first.last()@iana.org`,
		`1234567890123456789012345678901234567890123456789012345678901234@iana.org`,
		`"first\"last"@iana.org`,
		`"first@last"@iana.org`,
		`prettyandsimple@example.com`,
		`very.common@example.com`,
		`disposable.style.email.with+symbol@example.com`,
		`other.email-with-dash@example.com`,
		`x@example.com`,
		`"much.more unusual"@example.com`,
		`"very.unusual.@.unusual.com"@example.com`,
		`"very.(),:;<>[]\".VERY.\"very@\\ \"very\".unusual"@strange.example.com`,
		`example-indeed@strange-example.com`,
		`admin@mailserver1`,
		"#!$%&'*+-/=?^_`{}|~@example.org",
		"\"()<>[]:,;@\\\"!#$%&'-/=?^_`{}| ~.a\"@example.org",
		`" "@example.org`,
		`example@localhost`,
		`example@s.solutions`,
		`user@com`,
		`user@localserver`,
		`user@[IPv6:2001:db8::1]`,
		`user@[172.0.0.1]`,
		`user@[IPv6:2001:db8::1](comment)`,
		`user@(comment)[IPv6:2001:db8::1](comment)`,
		`user@(comment)[IPv6:2001:db8::1]`,
		`“email”@example.com`,
		`_______@example.com`,
	}

	invalidEmails = []string{
		`Abc.example.com`,
		`A@b@c@example.com`,
		`a"b(c)d,e:f;g<h>i[j\k]l@example.com`,
		`just"not"right@example.com`,
		`this is"not\allowed@example.com`,
		`this\ still\"not\\allowed@example.com`,
		`john..doe@example.com`,
		`john.doe@example..com`,
		` first.last@iana.org`,
		`first.last@iana.org `,
		`first.last.@iana.org`,
		`.first.last@iana.org`,
		`user@(comment)[2001:db8::1]`,
		`user@(comment)[666.666.666.666]`,
		`user@[666.666.666.666]`,
		`user@[666.666.666.666](comment)`,
		`(comment)@sada.com`,
		`local@(comment)`,
		`user@[IPv6:2001:db8::1]example.com`,
		`user@[127.0.0.1](coimment)example.com`,
	}
)

func TestValidEmails(t *testing.T) {
	for _, email := range validEmails {
		if !IsValid(email) {
			t.Errorf("Email %s should be valid", email)
		}
	}
}

func TestInvalidEmails(t *testing.T) {
	for _, email := range invalidEmails {
		if IsValid(email) {
			t.Errorf("Email %s should be invalid", email)
		}
	}
}
