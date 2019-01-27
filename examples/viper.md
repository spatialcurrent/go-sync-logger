# viper

Below is an example for how to initialize a `*gsl.Logger` from configuration provided by [viper](https://github.com/spf13/viper).

```go
func initLogger(v *viper.Viper) *gsl.Logger {

	verbose := v.GetBool("verbose")

	errorDestination := v.GetString("error-destination")
	errorCompression := v.GetString("error-compression")
	errorFormat := v.GetString("error-format")

	errorWriter, err := grw.WriteToResource(errorDestination, errorCompression, true, nil)
	if err != nil {
		fmt.Println(errors.Wrap(err, "error creating error writer"))
		os.Exit(1)
	}

	levels := map[string]int{"error": 0, "fatal": 0}
	writers := []grw.ByteWriteCloser{errorWriter}
	formats := []string{errorFormat}

	if verbose {
		levels["warn"] = 0
	}

	infoDestination := v.GetString("info-destination")
	infoCompression := v.GetString("info-compression")
	infoFormat := v.GetString("info-format")

	if len(infoDestination) > 0 && infoDestination != "/dev/null" && infoDestination != "null" {
		if infoDestination == errorDestination {
			if infoFormat != errorFormat {
				errorWriter.WriteError(fmt.Errorf("info-format ( %s ) and error-format ( %s ) must match when they share a destination", infoFormat, errorFormat)) // #nosec
				errorWriter.Close()                                                                                                                                // #nosec
				os.Exit(1)
			}
			if infoCompression != errorCompression {
				errorWriter.WriteError(fmt.Errorf("info-compression ( %s ) and error-compression ( %s ) must match when they share a destination", infoCompression, errorCompression)) // #nosec
				errorWriter.Close()                                                                                                                                                    // #nosec
				os.Exit(1)
			}

			levels["info"] = 0
			if verbose {
				levels["debug"] = 0
			}
		} else {
			infoWriter, err := grw.WriteToResource(infoDestination, infoCompression, true, nil)
			if err != nil {
				errorWriter.WriteError(errors.Wrap(err, "error creating log writer")) // #nosec
				errorWriter.Close()                                                   // #nosec
				os.Exit(1)
			}

			levels["info"] = 1
			writers = append(writers, infoWriter)
			formats = append(formats, infoFormat)

			if verbose {
				levels["debug"] = 1
			}
		}
	}

	return gsl.NewLogger(levels, writers, formats)
}
```
