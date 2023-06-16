package middleware

import (
	"fmt"

	"github.com/Temctl/E-Notification/util/elog"
)

func PrintZ() {
	fmt.Println()
	fmt.Println(" _____     _   _       _   _  __ _           _   _             ")
	fmt.Println("| ____|   | \\ | | ___ | |_(_)/ _(_) ___ __ _| |_(_) ___  _ __  ")
	fmt.Println("|  _| ____|  \\| |/ _ \\| __| | |_| |/ __/ _` | __| |/ _ \\| '_ \\ ")
	fmt.Println("| |__|____| |\\  | (_) | |_| |  _| | (_| (_| | |_| | (_) | | | |")
	fmt.Println("|_____|   |_| \\_|\\___/ \\__|_|_| |_|\\___\\__,_|\\__|_|\\___/|_| |_|")
	fmt.Println()

	elog.Info().Println()
	elog.Info().Println(" _____     _   _       _   _  __ _           _   _             ")
	elog.Info().Println("| ____|   | \\ | | ___ | |_(_)/ _(_) ___ __ _| |_(_) ___  _ __  ")
	elog.Info().Println("|  _| ____|  \\| |/ _ \\| __| | |_| |/ __/ _` | __| |/ _ \\| '_ \\ ")
	elog.Info().Println("| |__|____| |\\  | (_) | |_| |  _| | (_| (_| | |_| | (_) | | | |")
	elog.Info().Println("|_____|   |_| \\_|\\___/ \\__|_|_| |_|\\___\\__,_|\\__|_|\\___/|_| |_|")
	elog.Info().Println()
}
