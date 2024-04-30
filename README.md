# JUnit CLI Tool

A tool for easily creating JUnit test cases from scripts. Create or
update a junit.xml file with test results:

```
./junitcli create -s "test suite name" -f junit.xml -t "Test Case 1" --failure-output "Test failed"
```

```
./junitcli create -s "test suite name" -f junit.xml -t "Test Case 1" --system-output "Test succeeded"
```
