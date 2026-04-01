package fixtures

var (
	ValidEmails = []string{
		"test@gmail.com",
		"y@yandex.ru",
		"user@mail.ru",
	}

	InvalidEmails = []string{
		"41235612873",
		"user.com",
		"invalid@",
	}

	ValidPasswords = []string{
		"qwerty123",
		"printsf.f_mtA",
		"ersmkruwbrlnh123",
	}

	InvalidPasswords = []string{
		"rthjekpofwpoifjliuwekfmwjefwuikfmwel;f,wiejri",
		"1",
		"q",
		"wdf",
	}
)
