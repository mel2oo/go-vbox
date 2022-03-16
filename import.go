package vbox

//ImportOV imports ova or ovf from the given path
func ImportOV(path string) error {
	return Manage().run("import", path)
}
