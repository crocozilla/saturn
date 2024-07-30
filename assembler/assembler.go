package assembler

var sourceCodePath string

func Run(filePath string) {
	sourceCodePath = filePath
}

func firstStep() {
	scanWords(sourceCodePath, func(word string) {
		
	})
}

func secondStep() {

}
