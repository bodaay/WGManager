package utils

import (
	"bytes"
	"crypto/md5"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"

	// usbdrivedetector "github.com/deepakjois/gousbdrivedetector"
	"github.com/gabriel-vasile/mimetype"
	"github.com/google/uuid"
)

const MyTimeFormatWithTimeZone = "2006-01-02T15-04-05 -0700"
const MyTimeFormatWithoutTimeZone = "2006-01-02T15-04-05"

//FileExists just a simple boolean file exists function
func FileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

// CreateFolderIfNotExists Create folder if not exists, reutrn error if could not create for any reason
func CreateFolderIfNotExists(foldername string) error {
	if _, err := os.Stat(foldername); os.IsNotExist(err) {
		err = os.Mkdir(foldername, os.ModePerm)
		if err != nil {
			return err
		}
	}
	return nil
}

// CreateFolderIfNotExists Create folder if not exists, reutrn error if could not create for any reason
func CreateFolderAllIfNotExists(foldername string) error {
	if _, err := os.Stat(foldername); os.IsNotExist(err) {
		err = os.MkdirAll(foldername, os.ModePerm)
		if err != nil {
			return err
		}
	}
	return nil
}

// GetFileNameWithoutPath Return the file without path
func GetFileNameWithoutPath(filename string) string {
	return filepath.Base(filename)
}

// GetFileMimeType Will return the mime type string of the file
func GetFileMimeType(filename string) (string, error) {

	mime, err := mimetype.DetectFile(filename)
	if err != nil {
		return "", err
	}
	return mime.String(), nil
}

// Getdrives will get list of all drives connected to the machine, its very simple logic
func Getdrives() (r []string) { // got this from: https://stackoverflow.com/questions/23128148/how-can-i-get-a-listing-of-all-drives-on-windows-using-golang
	if runtime.GOOS == "windows" { //TODO: Do other implementation for linux/darwin, you can start by getting some tips from this link: https://unix.stackexchange.com/questions/24182/how-to-get-the-complete-and-exact-list-of-mounted-filesystems-in-linux/264226
		for _, drive := range "ABCDEFGHIJKLMNOPQRSTUVWXYZ" {
			f, err := os.Open(string(drive) + ":\\")
			if err == nil {
				r = append(r, string(drive))
				f.Close()
			}
		}

	}
	return
}

// GetdrivesUSBStorage will get list of all USBG drives connected to the machine
// func GetdrivesUSBStorage() (r []string) {
// 	if drives, err := usbdrivedetector.Detect(); err == nil {
// 		fmt.Printf("%d USB Devices Found\n", len(drives))
// 		for _, d := range drives {
// 			// fmt.Println(d)
// 			r = append(r, string(d))
// 		}
// 	}
// 	return
// }

// GetFileOSState Will stat object of the file
func GetFileOSState(filename string) (os.FileInfo, error) {
	stats, err := os.Stat(filename)
	if err != nil {
		return nil, err
	}

	return stats, nil
}
func CompressDataGZIP(data []byte) ([]byte, error) {
	var buf bytes.Buffer
	err := GzipWrite(&buf, data)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
func CompressFileGZIP(filename string) ([]byte, error) {
	fileBytes, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	var buf bytes.Buffer
	err = GzipWrite(&buf, fileBytes)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
func CompressFileGZIPAndSave(filename string, outputFileName string) error {
	fileBytes, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}
	f, err := os.Create(outputFileName)
	if err != nil {
		return err
	}
	defer f.Close()
	err = GzipWrite(f, fileBytes)
	if err != nil {
		return err
	}
	return nil
}

func DeCompressDataGZIP(data []byte) ([]byte, error) {
	var buf bytes.Buffer
	err := GUnzipWrite(&buf, data)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
func DeCompressFileGZIP(filename string) ([]byte, error) {
	fileBytes, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	var buf bytes.Buffer
	err = GUnzipWrite(&buf, fileBytes)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
func DeCompressFileGZIPAndSave(filename string, outputFileName string) error {
	fileBytes, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}
	f, err := os.Create(outputFileName)
	if err != nil {
		return err
	}
	defer f.Close()
	err = GUnzipWrite(f, fileBytes)
	if err != nil {
		return err
	}
	return nil
}

// HashFileMD5 Hash File with MD5 and return encoded hex string hash
func HashFileMD5(filePath string) (string, error) { //copied from: http://www.mrwaggel.be/post/generate-md5-hash-of-a-file-in-golang/
	//Initialize variable returnMD5String now in case an error has to be returned
	var returnMD5String string

	//Open the passed argument and check for any error
	file, err := os.Open(filePath)
	if err != nil {
		return returnMD5String, err
	}

	//Tell the program to call the following function when the current function returns
	defer file.Close()

	//Open a new hash interface to write to
	hash := md5.New()

	//Copy the file in the hash interface and check for any error
	if _, err := io.Copy(hash, file); err != nil {
		return returnMD5String, err
	}

	//Get the 16 bytes hash
	hashInBytes := hash.Sum(nil)[:16]

	//Convert the bytes to a string
	returnMD5String = hex.EncodeToString(hashInBytes)

	return returnMD5String, nil

}

// HashDataMD5 Hash Data Bytes with MD5 and return encoded hex string hash
func HashDataMD5(data []byte) (string, error) {
	//Initialize variable returnMD5String now in case an error has to be returned
	var returnMD5String string

	//Open a new hash interface to write to
	hash := md5.New()

	_, err := hash.Write(data)
	if err != nil {
		return "", err
	}

	//Get the 16 bytes hash
	hashInBytes := hash.Sum(nil)[:16]

	//Convert the bytes to a string
	returnMD5String = hex.EncodeToString(hashInBytes)

	return returnMD5String, nil

}

//HashFileSha256 Hash File with SHA256 and return encoded hex string hash
func HashFileSha256(filePath string) (string, error) { //copied from: http://www.mrwaggel.be/post/generate-md5-hash-of-a-file-in-golang/
	//Initialize variable returnMD5String now in case an error has to be returned
	var returnSha256 string

	//Open the passed argument and check for any error
	file, err := os.Open(filePath)
	if err != nil {
		return returnSha256, err
	}

	//Tell the program to call the following function when the current function returns
	defer file.Close()

	//Open a new hash interface to write to
	hash := sha256.New()

	//Copy the file in the hash interface and check for any error
	if _, err := io.Copy(hash, file); err != nil {
		return returnSha256, err
	}

	//Get the 32 bytes hash
	hashInBytes := hash.Sum(nil)[:32]

	//Convert the bytes to a string
	returnSha256 = hex.EncodeToString(hashInBytes)

	return returnSha256, nil

}

//HashDataSha256 Hash Data Bytes with SHA256 and return encoded hex string hash
func HashDataSha256(data []byte) (string, error) {
	//Initialize variable returnMD5String now in case an error has to be returned
	var returnSha256 string

	//Open a new hash interface to write to
	hash := sha256.New()
	_, err := hash.Write(data)
	if err != nil {
		return "", nil
	}
	//Get the 32 bytes hash
	hashInBytes := hash.Sum(nil)[:32]

	//Convert the bytes to a string
	returnSha256 = hex.EncodeToString(hashInBytes)

	return returnSha256, nil

}

//GetNewUUID will retrun a new Version 4 UUID string
func GetNewUUID() string {
	return uuid.New().String()
}

//HTTPFileUpload Needs to be further tested, I remember it was working
func HTTPFileUpload(url, file string) (err error) {
	// Prepare a form that you will submit to that URL.
	var b bytes.Buffer
	w := multipart.NewWriter(&b)
	// Add your image file
	f, err := os.Open(file)
	if err != nil {
		log.Println(err)
		return
	}
	defer f.Close()
	fw, err := w.CreateFormFile("file", file)
	if err != nil {
		log.Println(err)
		return
	}
	if _, err = io.Copy(fw, f); err != nil {
		log.Println(err)
		return
	}

	// Don't forget to close the multipart writer.
	// If you don't close it, your request will be missing the terminating boundary.
	w.Close()

	// Now that you have a form, you can submit it to your handler.
	req, err := http.NewRequest("POST", url, &b)
	if err != nil {
		log.Println(err)
		return
	}
	// Don't forget to set the content type, this will contain the boundary.
	req.Header.Set("Content-Type", w.FormDataContentType())

	// Submit the request
	//timeout := time.Duration(5 * time.Second)
	client := &http.Client{}
	fmt.Println(req)
	res, err := client.Do(req)
	if err != nil {
		log.Println("here")
		log.Println(err)
		return
	}

	// Check the response
	if res.StatusCode != http.StatusOK {
		err = fmt.Errorf("bad status: %s", res.Status)
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	fmt.Println(string(body))
	return
}

//GetMeFileListInFolders Get List of of file within the folder specified, pass empty string to spcific extention if you don't want to filter
func GetMeFileListInFolders(folderName string, specificExtension string, IgnoreSubFolders bool, IgnoreRootFolder bool, KeepFullPathForFile bool) []string {
	fileList := []string{}
	//walking through folder, it will just take any file in subdirectory of the folder, I'm supporting subfolders within folder
	err := filepath.Walk(folderName, func(fpath string, f os.FileInfo, err error) error {
		if IgnoreRootFolder {
			if f.IsDir() {
				return nil
			}
			if strings.HasSuffix(fpath, folderName) { //ignore main  folder
				return nil
			}
			if filepath.Dir(fpath) == folderName { //ignore any file in root  folder
				return nil
			}
		}
		//TODO: fix this below, its the cause of the shit // FIXED ALREADY
		//Fixed, Just checking now for stabolity
		if f != nil {
			if IgnoreSubFolders {
				if f.IsDir() { //ignore folders
					return nil
				}
			}
			if f.Mode()&os.ModeSymlink != 0 {
				err = os.Remove(path.Join(fpath, filepath.Dir(f.Name())))
				if err != nil {
					fmt.Println(err)
				}
				return nil
			}
		} else {
			return nil
		}

		if specificExtension != "" {
			if !strings.HasSuffix(strings.ToLower(fpath), strings.ToLower(specificExtension)) {
				return nil
			}
		}
		finalFileName := fpath
		if !KeepFullPathForFile {
			finalFileName = strings.Replace(fpath, folderName, "", 1)
		}
		fileList = append(fileList, finalFileName)
		return nil
	})
	if err != nil {
		log.Fatal(err)
	}
	return fileList
}

//CheckIfAdminOrRoot Cross platform to check if running as Admin or Root
func CheckIfAdminOrRoot() (bool, error) {
	if runtime.GOOS == "windows" { //got it from: https://coolaj86.com/articles/golang-and-windows-and-admins-oh-my/
		// var sid *windows.SID
		// err := windows.AllocateAndInitializeSid(
		// 	&windows.SECURITY_NT_AUTHORITY,
		// 	2,
		// 	windows.SECURITY_BUILTIN_DOMAIN_RID,
		// 	windows.DOMAIN_ALIAS_RID_ADMINS,
		// 	0, 0, 0, 0, 0, 0,
		// 	&sid)
		// if err != nil {
		// 	// log.Fatalf("SID Error: %s", err)
		// 	return false, err
		// }

		// // This appears to cast a null pointer so I'm not sure why this
		// // works, but this guy says it does and it Works for Me™:
		// // https://github.com/golang/go/issues/28804#issuecomment-438838144
		// token := windows.Token(0)

		// memberIsadmin, err := token.IsMember(sid)
		// if err != nil {
		// 	// log.Fatalf("Token Membership Error: %s", err)
		// 	return false, err
		// }
		// // token.IsElevated() //this if you want to check if elevated
		// // log.Println(memberIsadmin)
		// return memberIsadmin, nil
		return false, errors.New("Windows not supported")
	} else { //got it from: https://www.socketloop.com/tutorials/golang-force-your-program-to-run-with-root-permissions
		cmd := exec.Command("id", "-u")
		output, err := cmd.Output()

		if err != nil {
			log.Fatal(err)
		}

		// output has trailing \n
		// need to remove the \n
		// otherwise it will cause error for strconv.Atoi
		// log.Println(output[:len(output)-1])

		// 0 = root, 501 = non-root user
		i, err := strconv.Atoi(string(output[:len(output)-1]))
		if err != nil {
			return false, err
		}
		if i == 0 {
			return true, nil
		} else {
			return false, nil
		}
	}
}

func OpenBrowser(url string) {
	var err error

	switch runtime.GOOS {
	case "linux":
		err = exec.Command("xdg-open", url).Start()
	case "windows":
		err = exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
	case "darwin":
		err = exec.Command("open", url).Start()
	default:
		err = fmt.Errorf("unsupported platform")
	}
	if err != nil {
		log.Fatal(err)
	}

}
