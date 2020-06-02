package main

import(
	"fmt"
	. "github.com/cdipaolo/goml/base"
	"github.com/cdipaolo/goml/text"
	"io/ioutil"
	"math/rand"
	"regexp"
	"strings"
)

const (
	TRAIN_DATA_PATH          = "aclImdb/train/"
	TEST_DATA_PATH           = "aclImdb/test/"
	POSITIVE_TRAIN_DATA_PATH = TRAIN_DATA_PATH + "/pos/"
	NEGATIVE_TRAIN_DATA_PATH = TRAIN_DATA_PATH + "/neg/"
	POSITIVE_TEST_DATA_PATH  = TEST_DATA_PATH + "/pos/"
	NEGATIVE_TEST_DATA_PATH  = TEST_DATA_PATH + "/neg/"
)

func preProcess(text string) string{
	// Find all characters that are not alphabets
	reg := regexp.MustCompile("[^a-zA-Z]+")

	// Replace those chars with spaces
	text = reg.ReplaceAllString(text, " ")

	// Lower case
	text = strings.ToLower(text)

	// Tokenize on whitespaces, while removing excess whitespaces
	tokens := strings.Fields(text)

	// Join the tokens back to string
	return strings.Join(tokens, " ")
}

// Read data from IMDB dataset
func ReadData(dir string) []string{
	fileInfo, err := ioutil.ReadDir(dir)
	if err!= nil{
		panic(err)
	}
	var dataStrings []string
	for _, file := range fileInfo{
		bytes, err := ioutil.ReadFile(dir + file.Name())
		if err!= nil{
			panic(err)
		}
		dataStrings = append(dataStrings, string(bytes))
	}
	return dataStrings
}

// Prepare data in TextDatapoint struct for text classifier
func GetDataForClassifier() ([]TextDatapoint, []TextDatapoint) {
	var TrainData []TextDatapoint
	var TestData []TextDatapoint

	positiveTrainDataStrings := ReadData(POSITIVE_TRAIN_DATA_PATH)
	negativeTrainDataStrings := ReadData(NEGATIVE_TRAIN_DATA_PATH)
	positiveTestDataStrings := ReadData(POSITIVE_TEST_DATA_PATH)
	negativeTestDataStrings := ReadData(NEGATIVE_TEST_DATA_PATH)

	for _, str := range positiveTrainDataStrings {
		d := TextDatapoint{preProcess(str), 1}
		TrainData = append(TrainData, d)
	}
	for _, str := range negativeTrainDataStrings {
		d := TextDatapoint{preProcess(str), 0}
		TrainData = append(TrainData, d)
	}
	for _, str := range positiveTestDataStrings {
		d := TextDatapoint{preProcess(str), 1}
		TestData = append(TrainData, d)
	}
	for _, str := range negativeTestDataStrings {
		d := TextDatapoint{preProcess(str), 0}
		TestData = append(TrainData, d)
	}

	rand.Shuffle(len(TrainData), func(i, j int) {
		TrainData[i], TrainData[j] = TrainData[j], TrainData[i]
	})
	return TrainData, TestData
}

func main() {
	// Get training and test dataset
	TrainData, TestData := GetDataForClassifier()

	// Cross validation dataset
	var CrossValData [][]TextDatapoint

	k := len(TrainData) * 1 / 5
	for idx:=0; idx < len(TrainData); idx += k{
		CrossValData = append(CrossValData, TrainData[idx:idx+k])
	}

	// Use concurrency, create the channels
	stream := make(chan TextDatapoint, 100) // buffered channel
	errors := make(chan error)

	model := text.NewNaiveBayes(stream, 2, OnlyWordsAndNumbers)
	go model.OnlineLearn(errors)  // Defining the model -- making a go routine

	// Train the model
	for _, data := range TrainData{
		stream <- data      // sending data via channel to the stream
	}

	// closing the stream
	close(stream)

	// check for the availability of any error
	for {
		err, _ := <-errors
		if err != nil{
			fmt.Printf("Error Passed: %v", err)
		}else{
			break
		}
	}

	// Now we can predict any statement
	s := "This spacecraft will carry @AstroBehnken and @Astro_Doug to the @Space_Station when it launches atop a " +
		"@SpaceX Falcon 9 rocket on May 27, at 4:33 p.m. ET: https://t.co/yvfOCG4\""

	class := model.Predict(s)
	fmt.Println("Class Predicted", class)
	fmt.Println(model.Probability(s))

	mistakes := 0
	// calculating the mistakes in the training data
	for _, t := range TrainData{
		class := model.Predict(t.X)
		if class != t.Y{
			mistakes += 1
		}
	}

	total := float64(len(TrainData))
	err := float64(mistakes)
	accuracy := (total - err) / total * 100
	fmt.Printf("Train Accuracy is %v%%\n", accuracy)

	// Calculate mistakes in test data
	for _, p := range TestData{
		class := model.Predict(p.X)
		if class != p.Y{
			mistakes += 1
		}
	}
	total = float64(len(TestData))
	err = float64(mistakes)
	accuracy = (total - err) / total * 100
	fmt.Printf("Test accuracy is %v%%\n", accuracy)

}