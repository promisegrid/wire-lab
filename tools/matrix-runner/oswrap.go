package main

import (
	"os"
)

// Small wrappers keep filesystem use explicit and easy to replace in narrow
// tests without hiding the fact that matrix-runner is a file-oriented tool.
var (
	osCreate   = os.Create
	osMkdirAll = os.MkdirAll
	osOpen     = os.Open
	osReadDir  = os.ReadDir
	osStat     = os.Stat
)
