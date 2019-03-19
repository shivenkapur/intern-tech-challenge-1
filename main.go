package main

import (
	"bufio"
	"context"
	"fmt"

	"io"

	"flag"

	"github.com/coreos/go-semver/semver"
	"github.com/google/go-github/github"

	"os"

	"sort"

	"strings"
)

// LatestVersions returns a sorted slice with the highest version as its first element and the highest version of the smaller minor versions in a descending order
func LatestVersions(releases []*semver.Version, minVersion *semver.Version) []*semver.Version {

	var versionSlice []*semver.Version
	var tempversions []*semver.Version
	// This is just an example structure of the code, if you implement this interface, the test cases in main_test.go are very easy to run
	for _,release := range releases {
		
		if minVersion.LessThan(*release) {
			tempversions = append(tempversions,release)
		}
	}

	sort.Slice(tempversions, func(i, j int) bool {
		return tempversions[j].LessThan(*tempversions[i])
	})


	fmt.Printf("%s\n",tempversions)

	var(
		minor int64 = -1
		major int64 = -1
	)
	for _,version := range tempversions{
		if(minor != version.Minor || major != version.Major){
			versionSlice = append(versionSlice, version)
			minor = version.Minor
			major = version.Major
		}
	}

	return versionSlice
}

func AllReleases(s1 string, s2 string) ([]*semver.Version){
	client := github.NewClient(nil)
	ctx := context.Background()
	opt := &github.ListOptions{PerPage: 10}
	releases, _, err := client.Repositories.ListReleases(ctx, s1, s2, opt)
	if err != nil {
		fmt.Printf("Error: %s", err)
	}
	allReleases := make([]*semver.Version, len(releases))
	
	for i, release := range releases {
		versionString := *release.TagName
		if versionString[0] == 'v' {
			versionString = versionString[1:]
		}
		allReleases[i] = semver.New(versionString)
	}
	
	return allReleases
}

func getValues(line []byte) (string, string , *semver.Version){
	splitline := strings.Split(string(line) , ",")

	repository := splitline[0]
	minversion := splitline[1]

	s := strings.Split(repository , "/")
	return s[0],s[1], semver.New(minversion)
}

// Here we implement the basics of communicating with github through the library as well as printing the version
// You will need to implement LatestVersions function as well as make this application support the file format outlined in the README
// Please use the format defined by the fmt.Printf line at the bottom, as we will define a passing coding challenge as one that outputs
// the correct information, including this line
func main() {
	// Github
	flag.Parse()
    filename := flag.Arg(0)
	
	file, err := os.Open(filename)

	if err != nil {
		fmt.Printf("failed opening file: %s", err)
	} else{
		reader := bufio.NewReader(file)
		line, _, err := reader.ReadLine()

		for i := 0; ;i++ {
			line, _, err = reader.ReadLine()

			if err == io.EOF{
				break
			}

			s1,s2,minVersion := getValues(line)
			allReleases := AllReleases(s1,s2)
			versionSlice := LatestVersions(allReleases, minVersion)

			fmt.Printf("latest versions of %s/%s: %s\n", s1 , s2 , versionSlice)
			
		}
	}



}
