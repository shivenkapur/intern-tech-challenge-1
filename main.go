package main

import (
	"bufio"
	"context"
	"fmt"

	"flag"

	"github.com/coreos/go-semver/semver"
	"github.com/google/go-github/github"

	"io"
	"os"

	"sort"
	"strings"
)

// LatestVersions returns a sorted slice with the highest version as its first element and the highest version of the smaller minor versions in a descending order
func LatestVersions(releases []*semver.Version, minVersion *semver.Version) []*semver.Version {

	var versionSlice []*semver.Version

	sort.Slice(releases, func(i, j int) bool {
		return releases[j].LessThan(*releases[i])
	})

	var(
		minor int64 = -1
		major int64 = -1
	)
	for _,release := range releases {
		var isLesserorEqualtoMin bool = minVersion.LessThan(*release) || minVersion.Equal(*release) 
		if isLesserorEqualtoMin {
			if minor != release.Minor || major != release.Major {
				versionSlice = append(versionSlice, release)
				minor = release.Minor
				major = release.Major
			}
		}
	}
	return versionSlice
}

func AllReleases(githubProfileName string, repoName string) ([]*semver.Version){
	client := github.NewClient(nil)
	ctx := context.Background()
	opt := &github.ListOptions{PerPage: 10}
	releases, _, err := client.Repositories.ListReleases(ctx, githubProfileName, repoName, opt)
	if err != nil {
		fmt.Printf("Error getting releases from Github: %s", err)
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

func getValuesFromFile(line []byte) (string, string , *semver.Version){
	splitline := strings.Split(string(line) , ",")

	repository := splitline[0]
	minversion := splitline[1]

	s := strings.Split(repository , "/")
	return s[0],s[1], semver.New(minversion)
}

func main() {
	// get filename entered in command line
	flag.Parse()
    filename := flag.Arg(0)
	
	//open file with the specified file name
	file, err := os.Open(filename)

	if err != nil {
		fmt.Printf("Error while opening file: %s", err)
	} else{
		//initalize reader
		reader := bufio.NewReader(file)
		//read the first lien from the file - repository,min_version
		line, _, err := reader.ReadLine()

		var Iserror bool = false
		//read all subsequent lines and exit if there's an error
		for i := 0; !Iserror ;i++ {
			line, _, err = reader.ReadLine()

			if err == io.EOF{
				Iserror = true
			} else if err == nil{
				//get data from a single line in the file
				githubProfileName,repoName,minVersion := getValuesFromFile(line)
				//get all releases from github repository
				allReleases := AllReleases(githubProfileName,repoName)
				//get sorted max patch versions
				versionSlice := LatestVersions(allReleases, minVersion)

				fmt.Printf("latest versions of %s/%s: %s\n", githubProfileName , repoName , versionSlice)
			} else{
				fmt.Printf("Error while reading file: %s", err)
				Iserror = true
			}
		}
	}



}
