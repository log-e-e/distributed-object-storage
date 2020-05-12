package err_utils

func PanicNonNilError(err error) {
    if err != nil {
        panic(err)
    }
}
