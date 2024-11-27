package main

func main() {
	//Run the initialization sequence

	//todo: Check for brain data backup

	//Check for the existence of the data directory
	err := checkDirectories()
	if err != nil {
		SetupComplete = false
	}

	if !SetupComplete {

	}

}
