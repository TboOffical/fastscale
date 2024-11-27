package main

/*
All application variables can be found here along with a description of what they are used for.
*/

//Mappings: maps a string value to the corresponding value in the brain
//used for terminal queries and saving the information from the brain into a backup so the orchestrator can bounce back if it crashes or is stopped

var Mappings = map[string]interface{}{
	"setup_completed": SetupComplete,
}

// SetupComplete is a boolean that is used to determine if the setup process has been completed
var SetupComplete bool