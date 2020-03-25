package err_utils

import "log"

func PanicNonNilError(err error) {
    if err != nil {
        panic(err)
    }
}

func PrintNonNilError(err error, errInfo string) {
    if err != nil {
        log.Println(errInfo)
    }
}
