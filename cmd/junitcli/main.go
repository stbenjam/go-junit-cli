package main

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/spf13/cobra"

	"github.com/stbenjam/go-junit-cli/pkg/types/junit"
)

func main() {
	var rootCmd = &cobra.Command{Use: "junitcli"}

	var suiteName, fileName, testName, failureOutput, systemOutput string

	var cmdCreate = &cobra.Command{
		Use:   "create",
		Short: "Create a junit.xml with test results",
		Long: `Create or update a junit.xml file with test results:
./junitcli create -s "test suite name" -f junit.xml -t "Test Case 1" --failure-output "Test failed"
./junitcli create -s "test suite name" -f junit.xml -t "Test Case 1" --system-output "Test succeeded"`,
		Run: func(cmd *cobra.Command, args []string) {
			if suiteName == "" || fileName == "" || testName == "" {
				fmt.Println("Missing required arguments")
				os.Exit(1)
			}
			testCase := junit.TestCase{
				Name: testName,
			}

			if failureOutput != "" {
				testCase.FailureOutput = &junit.FailureOutput{
					Message: failureOutput,
				}
			}

			if systemOutput != "" {
				testCase.SystemOut = systemOutput
			}

			file, err := os.Open(fileName)
			defer file.Close()

			var testSuite junit.TestSuite
			if err != nil {
				if os.IsNotExist(err) {
					testSuite = junit.TestSuite{Name: suiteName}
				} else {
					fmt.Println("Error opening file:", err)
					os.Exit(1)
				}
			} else {
				data, err := ioutil.ReadAll(file)
				if err != nil {
					fmt.Println("Error reading file:", err)
					os.Exit(1)
				}
				err = xml.Unmarshal(data, &testSuite)
				if err != nil {
					fmt.Println("Error parsing XML:", err)
					os.Exit(1)
				}
			}

			testSuite.TestCases = append(testSuite.TestCases, &testCase)
			var numTests, numFailed uint
			for _, tc := range testSuite.TestCases {
				numTests++
				if tc.FailureOutput != nil {
					numFailed++
				}
			}
			testSuite.NumTests = numTests
			testSuite.NumFailed = numFailed

			output, err := xml.MarshalIndent(testSuite, "", "  ")
			if err != nil {
				fmt.Println("Error marshaling XML:", err)
				os.Exit(1)
			}

			err = ioutil.WriteFile(fileName, output, 0644)
			if err != nil {
				fmt.Println("Error writing to file:", err)
				os.Exit(1)
			}
		},
	}

	cmdCreate.Flags().StringVarP(&suiteName, "suite", "s", "", "Name of the test suite")
	cmdCreate.Flags().StringVarP(&fileName, "file", "f", "", "Output file name")
	cmdCreate.Flags().StringVarP(&testName, "test", "t", "", "Name of the test case")
	cmdCreate.Flags().StringVar(&failureOutput, "failure-output", "", "Failure output for the test case (implies test case failed)")
	cmdCreate.Flags().StringVar(&systemOutput, "system-output", "", "System output for the test case")

	rootCmd.AddCommand(cmdCreate)
	rootCmd.Execute()
}
