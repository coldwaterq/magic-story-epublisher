# magic-story-epublisher
 Have calibre installed, thi swill add the ebook-convert binary to your path.

grab the epub, azw3, html, or mobi format of each story.

run `magic-story-epublisher.exe` while in the git repo directory to download all the books if you don't trust my versions

build with `go build` while in the git repo directory

Update the following code block in main.go to get new versions. Older versions don't seem to be in the same format.
    stories_to_get := []string{
		"Strixhaven: School of Mages",
		"Kaldheim",
		"Zendikar Rising",
	}

Copy one of the blocks and paste it bellow it to make all the books an aditional file format.
    // .mobi
    cmd = exec.Command("ebook-convert", "html/"+file.Name(), "mobi/"+strings.ReplaceAll(file.Name(), ".html", ".mobi"))
    stdoutStderr, err = cmd.CombinedOutput()
    if err != nil {
        log.Println(err)
    }
    fmt.Printf("%s\n", stdoutStderr)

Change the comment to be clear, change the .mobi at the end of the first line to adjust the output, and create a new folder and change "mobi/" to be the folder created.