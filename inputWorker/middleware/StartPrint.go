package middleware

import (
	"fmt"

	"github.com/Temctl/E-Notification/util/elog"
)

func PrintZ() {
	logo := ` _____ __  __    _         _   _  ___ _____ ___ _____ ___ ____    _  _____ ___ ___  _   _ 
| ____|  \/  |  / \       | \ | |/ _ \_   _|_ _|  ___|_ _/ ___|  / \|_   _|_ _/ _ \| \ | |
|  _| | |\/| | / _ \ _____|  \| | | | || |  | || |_   | | |     / _ \ | |  | | | | |  \| |
| |___| |  | |/ ___ \_____| |\  | |_| || |  | ||  _|  | | |___ / ___ \| |  | | |_| | |\  |
|_____|_|  |_/_/   \_\    |_| \_|\___/ |_| |___|_|   |___\____/_/   \_\_| |___\___/|_| \_|`
	fmt.Println()
	fmt.Println(logo)
	fmt.Println()

	elog.Info().Println("\n" + logo)
	elog.Info().Println()
}
