package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path"
	"strconv"

	"github.com/Azure/azure-sdk-for-go/services/cognitiveservices/v1.0/face"
	"github.com/Azure/go-autorest/autorest"
)

func main() {

	// A global context for use in all samples
	faceContext := context.Background()

	// Base url for the Verify and Large Face List examples
	// const imageBaseURL = "https://csdx.blob.core.windows.net/resources/Face/Images/"

	/*
	   Authenticate
	*/
	// Add FACE_SUBSCRIPTION_KEY, FACE_ENDPOINT, and AZURE_SUBSCRIPTION_ID to your environment variables.
	subscriptionKey := os.Getenv("FACE_SUBSCRIPTION_KEY")
	endpoint := os.Getenv("FACE_ENDPOINT")

	// Client used for Detect Faces, Find Similar, and Verify examples.
	client := face.NewClient(endpoint)
	client.Authorizer = autorest.NewCognitiveServicesAuthorizer(subscriptionKey)
	/*
	   END - Authenticate
	*/

	// Detect a face in an image that contains a single face
	singleFaceImageURL := "https://www.biography.com/.image/t_share/MTQ1MzAyNzYzOTgxNTE0NTEz/john-f-kennedy---mini-biography.jpg"
	singleImageURL := face.ImageURL{URL: &singleFaceImageURL}
	singleImageName := path.Base(singleFaceImageURL)
	// Array types chosen for the attributes of Face
	attributes := []face.AttributeType{"age", "emotion", "gender"}
	returnFaceID := true
	returnRecognitionModel := false
	returnFaceLandmarks := false

	// API call to detect faces in single-faced image, using recognition model 3
	// We specify detection model 1 because we are retrieving attributes.
	detectSingleFaces, dErr := client.DetectWithURL(faceContext, singleImageURL, &returnFaceID, &returnFaceLandmarks, attributes, face.Recognition03, &returnRecognitionModel, face.Detection01)
	if dErr != nil {
		log.Fatal(dErr)
	}

	// Dereference *[]DetectedFace, in order to loop through it.
	dFaces := *detectSingleFaces.Value

	fmt.Println("Detected face in (" + singleImageName + ") with ID(s): ")
	fmt.Println(dFaces[0].FaceID)
	fmt.Println()
	// Find/display the age and gender attributes
	for _, dFace := range dFaces {
		fmt.Println("Face attributes:")
		fmt.Printf("  Age: %.0f", *dFace.FaceAttributes.Age)
		fmt.Println("\n  Gender: " + dFace.FaceAttributes.Gender)
	}
	// Get/display the emotion attribute
	emotionStruct := *dFaces[0].FaceAttributes.Emotion
	// Convert struct to a map
	var emotionMap map[string]float64
	result, _ := json.Marshal(emotionStruct)
	json.Unmarshal(result, &emotionMap)
	// Find the emotion with the highest score (confidence level). Range is 0.0 - 1.0.
	var highest float64
	emotion := ""
	dScore := -1.0
	for name, value := range emotionMap {
		if value > highest {
			emotion, dScore = name, value
			highest = value
		}
	}
	fmt.Println("  Emotion: " + emotion + " (score: " + strconv.FormatFloat(dScore, 'f', 3, 64) + ")")

}
