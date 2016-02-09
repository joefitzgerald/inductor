package renderer

import "testing"

func TestSafeComputerNameThatIsTooLong(t *testing.T) {
	invalidComputerName := "/\\*lo<n>g|with?invalidcharsandtoolong"
	actual := SafeComputerName(invalidComputerName)
	if actual != "longwithinvalid" {
		t.Errorf("Expected computer name 'longwithinvalid', but got '%s'", actual)
	}
}

func TestSafeComputerName(t *testing.T) {
	computerName := "valid"
	actual := SafeComputerName(computerName)
	if actual != computerName {
		t.Errorf("Expected computer name '%s', but got '%s'", computerName, actual)
	}
}
