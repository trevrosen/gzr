package comms

import (
	"os/exec"
	"regexp"
	"strings"
	"time"
)

func BuildMetadata() (ImageMetadata, error) {
	meta := ImageMetadata{}
	tags, annotations, err := getTagsAndAnnotations()
	if err != nil {
		return meta, err
	}
	meta.GitTag = tags
	meta.GitAnnotation = annotations

	remote, err := getRemote()
	if err != nil {
		return meta, err
	}
	meta.GitOrigin = remote

	hash, err := getCommitHash()
	if err != nil {
		return meta, err
	}
	meta.GitCommit = hash

	meta.CreatedAt = time.Now().Format(time.RFC3339)

	return meta, nil
}

func getCommitHash() (string, error) {
	hashCmd := exec.Command("git", "rev-parse", "HEAD")

	hash, err := hashCmd.CombinedOutput()
	if err != nil {
		return "", err
	}

	stripped := strings.TrimSpace(string(hash))
	return stripped, nil
}

func getRemote() (string, error) {
	remoteCmd := exec.Command("git", "remote", "-v")
	remote, err := remoteCmd.CombinedOutput()
	if err != nil {
		return "", err
	}

	remotes := strings.Fields(string(remote))
	if len(remotes) == 0 {
		return "", nil
	}
	return remotes[1], nil
}

func getTagsAndAnnotations() ([]string, []string, error) {
	var tags []string
	var annotations []string
	currentTags := exec.Command("git", "tag", "--format", "%(refname:strip=2)~%(contents:subject)", "-l", "-n1", "--points-at", "HEAD")

	tagInfo, err := currentTags.CombinedOutput()
	if err != nil {
		return []string{}, []string{}, err
	}
	regex, _ := regexp.Compile("\n\n")
	tagInfo_ := regex.ReplaceAllString(string(tagInfo), "\n")

	splitTagInfo := strings.Split(tagInfo_, "\n")
	for _, v := range splitTagInfo {
		if v == "" {
			continue
		}
		split := strings.Split(v, "~")
		if len(split) == 2 {
			tag := split[0]
			annotation := split[1]
			tags = append(tags, tag)
			annotations = append(annotations, annotation)
		}
		if len(split) == 1 {
			tag := split[0]
			tags = append(tags, tag)
		}
	}
	return tags, annotations, nil
}
