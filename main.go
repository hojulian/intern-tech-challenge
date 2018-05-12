// Intern-tech-challenge
// For Lalamove
// @author Julian Ho
// 2018
package main

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/coreos/go-semver/semver"
	"github.com/google/go-github/github"
)

// LatestVersions returns a sorted slice with the highest version
// as its first element and the highest version of the smaller minor
// versions in a descending order
// @return nothing if the minVersion doesn't exist
func LatestVersions(releases []*semver.Version, minVersion *semver.Version) []*semver.Version {
	var versionSlice []*semver.Version
	// This is just an example structure of the code,
	// if you implement this interface, the test cases
	// in main_test.go are very easy to run

	// Sort before filtering
	sort.Sort(byVersion(releases))
	// Filtering
	// Compare each element in releases with minVersion
	// and append the element that is highest in its MINOR
	for _, r := range releases {
		if r.Compare(*minVersion) >= 0 {
			// Append an element first(Latest version)
			if versionSlice == nil {
				versionSlice = append(versionSlice, r)
			} else {
				// Compare the MINOR of the last appended element
				cm := r.Slice()
				tcm := versionSlice[len(versionSlice)-1].Slice()
				if cm[1] != tcm[1] {
					versionSlice = append(versionSlice, r)
				}
			}
		} else {
			// No need to check releases smaller than minVersion
			break
		}
	}
	// Check if minVersion actually exists
	if versionSlice == nil {
		fmt.Printf("%s doesn't exist\n", minVersion)
	}
	return versionSlice
}

// Implementing sort interface
type byVersion []*semver.Version

func (v byVersion) Len() int           { return len(v) }
func (v byVersion) Swap(i, j int)      { v[i], v[j] = v[j], v[i] }
func (v byVersion) Less(i, j int) bool { return v[j].LessThan(*v[i]) }

// Simple error handle
func check(e error) {
	if e != nil {
		panic(fmt.Sprintln(e.Error()))
	}
}

// Github searching
// Takes in the input from file
// and print it out
func search(s []string) {
	// Github
	client := github.NewClient(nil)
	ctx := context.Background()
	opt := &github.ListOptions{PerPage: 10}
	// Input data
	ver := semver.Version{}
	if err := ver.Set(s[1]); err != nil {
		fmt.Printf("%s has an invalid version format: %s\n", s[0], err.Error())
		return
	}
	tr, v := s[0], semver.New(s[1])
	// Checking data before searching
	if v == nil {
		fmt.Printf("%s has an invalid version format\n", tr)
		return
	}
	if !strings.ContainsAny(tr, "/") {
		fmt.Printf("%s has an invalid repo name\n", tr)
		return
	}
	// Searching(Existing code)
	r := strings.Split(tr, "/")
	releases, _, err := client.Repositories.ListReleases(ctx, r[0], r[1], opt)
	if err != nil {
		fmt.Printf("%s is invalid or doesn't exist: %s\n", tr, err.Error())
		return
	}
	allReleases := make([]*semver.Version, len(releases))
	for i, release := range releases {
		versionString := *release.TagName
		if versionString[0] == 'v' {
			versionString = versionString[1:]
		}
		allReleases[i] = semver.New(versionString)
	}
	versionSlice := LatestVersions(allReleases, v)
	fmt.Printf("latest versions of %s: %s\n", tr, versionSlice)
}

// Here we implement the basics of communicating with github through
// the library as well as printing the version
// You will need to implement LatestVersions function as well as make
// this application support the file format outlined in the README
// Please use the format defined by the fmt.Printf line at the bottom,
// as we will define a passing coding challenge as one that outputs
// the correct information, including this line
//
// Takes in one argument as file path
// e.g. <code> go run main.go mock_data.txt </code>
func main() {
	// Opening file
	file, err := os.Open(os.Args[1])
	check(err)
	defer file.Close()
	// Scanning file
	scanner := bufio.NewScanner(file)
	// Split on comma
	for scanner.Scan() {
		if !strings.ContainsAny(scanner.Text(), ",") {
			fmt.Printf("%s has an invalid format\n", scanner.Text())
			continue
		}
		result := strings.Split(scanner.Text(), ",")
		search([]string{result[0], result[1]})
	}
}
