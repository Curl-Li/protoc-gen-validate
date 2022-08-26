package java

import (
	"fmt"
	"testing"
)

func TestDownloadFormatTool(t *testing.T) {
	t.Log(downloadFormatTool(fmt.Sprintf("google-java-format-%s-all-deps.jar", googleJavaFormatVersion)))
}
