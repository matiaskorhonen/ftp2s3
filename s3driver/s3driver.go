package s3driver

import (
	"bytes"
	"fmt"
	"io"
	"mime"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/awslabs/aws-sdk-go/aws"
	"github.com/awslabs/aws-sdk-go/aws/awsutil"
	"github.com/awslabs/aws-sdk-go/service/s3"
	"github.com/koofr/graval"
)

type S3Driver struct {
	Username               string
	Password               string
	AWSRegion              string
	AWSCredentialsProvider aws.CredentialsProvider
	AWSBucketName          string
}

func (d *S3Driver) s3service() *s3.S3 {
	svc := s3.New(&aws.Config{
		Credentials: d.AWSCredentialsProvider,
		Region:      d.AWSRegion,
	})
	return svc
}

func pathToS3PathPrefix(path string) *string {
	path = strings.TrimPrefix(path, "/")

	if path == "" || strings.HasSuffix(path, "/") {
		return aws.String(path)
	} else {
		p := string(path) + "/"
		return aws.String(p)
	}
}

func (d *S3Driver) s3DirContents(path string, maxKeys int64, marker string) (*s3.ListObjectsOutput, error) {
	svc := d.s3service()

	prefix := pathToS3PathPrefix(path)

	params := &s3.ListObjectsInput{
		Bucket: aws.String(d.AWSBucketName), // Required
		// Delimiter:    aws.String("Delimiter"),
		// EncodingType: aws.String("EncodingType"),
		// Marker:       aws.String("Marker"),
		MaxKeys: aws.Long(maxKeys),
		Prefix:  prefix,
	}

	if marker != "" {
		params.Marker = aws.String(marker)
	}

	resp, err := svc.ListObjects(params)

	if awserr := aws.Error(err); awserr != nil {
		// A service error occurred.
		fmt.Println("Error: ", awserr)
	} else if err != nil {
		// A non-service error occurred.
		panic(err)
	}

	return resp, err
}

func (d *S3Driver) Authenticate(username string, password string) bool {
	return username == d.Username && password == d.Password
}

// Bytes returns the ContentLength for the path if the key exists
func (d *S3Driver) Bytes(path string) int64 {
	svc := d.s3service()

	path = strings.TrimPrefix(path, "/")

	params := &s3.HeadObjectInput{
		Bucket: aws.String(d.AWSBucketName), // Required
		Key:    aws.String(path),            // Required
	}
	resp, err := svc.HeadObject(params)

	if awserr := aws.Error(err); awserr != nil {
		// A service error occurred.
		fmt.Println("Error:", awserr.Code, awserr.Message)
		return -1
	} else if err != nil {
		// A non-service error occurred.
		panic(err)
		return -1
	}

	return *resp.ContentLength
}

// ModifiedTime returns the LastModifiedTime for the path if the key exists
func (d *S3Driver) ModifiedTime(path string) (time.Time, bool) {
	svc := d.s3service()

	path = strings.TrimPrefix(path, "/")

	params := &s3.HeadObjectInput{
		Bucket: aws.String(d.AWSBucketName), // Required
		Key:    aws.String(path),            // Required
	}
	resp, err := svc.HeadObject(params)

	if awserr := aws.Error(err); awserr != nil {
		// A service error occurred.
		fmt.Println("Error:", awserr.Code, awserr.Message)
		return time.Now(), false
	} else if err != nil {
		// A non-service error occurred.
		panic(err)
		return time.Now(), false
	}

	return *resp.LastModified, true
}

// ChangeDir “changes directories” on S3 if there are files under the given path
func (d *S3Driver) ChangeDir(path string) bool {
	resp, err := d.s3DirContents(path, 1, "")

	if err == nil && len(resp.Contents) > 0 {
		return true
	} else {
		return false
	}
}

// DirContents lists “directory” contents on S3
func (d *S3Driver) DirContents(path string) ([]os.FileInfo, bool) {
	moreObjects := true
	objects := make([]*s3.Object, 0)

	var resp *s3.ListObjectsOutput
	var err error
	marker := ""

	for moreObjects {
		resp, err = d.s3DirContents(path, 1000, marker)

		if err == nil {
			for _, obj := range resp.Contents {
				objects = append(objects, obj)
			}

			moreObjects = *resp.IsTruncated

			if moreObjects {
				last := objects[len(objects)-1]
				marker = *last.Key
			}
		}
	}

	prefix := pathToS3PathPrefix(path)
	files := make([]os.FileInfo, 0)
	dirs := make([]string, 0)

	for _, object := range objects {
		p := *object.Key

		p = strings.TrimPrefix(p, *prefix)
		var fi os.FileInfo

		if strings.Contains(p, "/") {

			parts := strings.Split(p, "/")
			dir_part := parts[0]

			// currentDir := strings.TrimSuffix(*prefix, "/") == dir_part

			if dir_part != "" && dir_part != "/" && !stringInSlice(dir_part, dirs) {
				fi = graval.NewDirItem(dir_part)
				files = append(files, fi)

				dirs = append(dirs, dir_part)
			}
		} else {
			fi = graval.NewFileItem(*object.Key, *object.Size, *object.LastModified)
			files = append(files, fi)
		}
	}

	return files, true
}

func (d *S3Driver) DeleteDir(path string) bool {
	// if f, ok := d.Files[path]; ok && f.File.IsDir() {
	// 	haschildren := false
	// 	for p, _ := range d.Files {
	// 		if strings.HasPrefix(p, path+"/") {
	// 			haschildren = true
	// 			break
	// 		}
	// 	}
	//
	// 	if haschildren {
	// 		return false
	// 	}
	//
	// 	delete(d.Files, path)
	//
	// 	return true
	// } else {
	// 	return false
	// }
	return false
}

func (d *S3Driver) DeleteFile(path string) bool {
	// if f, ok := d.Files[path]; ok && !f.File.IsDir() {
	// 	delete(d.Files, path)
	// 	return true
	// } else {
	// 	return false
	// }
	return false
}

func (d *S3Driver) Rename(from_path string, to_path string) bool {
	// if f, from_path_exists := d.Files[from_path]; from_path_exists {
	// 	if _, to_path_exists := d.Files[to_path]; !to_path_exists {
	// 		if _, to_path_parent_exists := d.Files[filepath.Dir(to_path)]; to_path_parent_exists {
	// 			if f.File.IsDir() {
	// 				delete(d.Files, from_path)
	// 				d.Files[to_path] = &MemoryFile{graval.NewDirItem(filepath.Base(to_path)), nil}
	// 				torename := make([]string, 0)
	// 				for p, _ := range d.Files {
	// 					if strings.HasPrefix(p, from_path+"/") {
	// 						torename = append(torename, p)
	// 					}
	// 				}
	// 				for _, p := range torename {
	// 					sf := d.Files[p]
	// 					delete(d.Files, p)
	// 					np := to_path + p[len(from_path):]
	// 					d.Files[np] = sf
	// 				}
	// 			} else {
	// 				delete(d.Files, from_path)
	// 				d.Files[to_path] = &MemoryFile{graval.NewFileItem(filepath.Base(to_path), f.File.Size(), f.File.ModTime()), f.Content}
	// 			}
	// 			return true
	// 		} else {
	// 			return false
	// 		}
	// 	} else {
	// 		return false
	// 	}
	// } else {
	// 	return false
	// }
	return false
}

// MakeDir would normally make a directory but this isn't supported on S3
func (d *S3Driver) MakeDir(path string) bool {
	return false
}

// GetFile returns a reader for the given path on S3
func (d *S3Driver) GetFile(path string, position int64) (io.ReadCloser, bool) {
	svc := d.s3service()

	path = strings.TrimPrefix(path, "/")

	params := &s3.GetObjectInput{
		Bucket: aws.String(d.AWSBucketName), // Required
		Key:    aws.String(path),            // Required
	}
	resp, err := svc.GetObject(params)

	if awserr := aws.Error(err); awserr != nil {
		// A service error occurred.
		fmt.Println("Error:", awserr.Code, awserr.Message)
		return nil, false
	} else if err != nil {
		// A non-service error occurred.
		panic(err)
		return nil, false
	}

	// Pretty-print the response data.
	fmt.Println(awsutil.StringValue(resp))

	return resp.Body, true
}

// PutFile uploads a file to S3
func (d *S3Driver) PutFile(path string, reader io.Reader) bool {
	svc := d.s3service()

	path = strings.TrimPrefix(path, "/")
	fileExt := filepath.Ext(path)

	contentType := mime.TypeByExtension(fileExt)
	if contentType == "" {
		contentType = "application/octet-stream"
	}

	buf := new(bytes.Buffer)
	buf.ReadFrom(reader)

	params := &s3.PutObjectInput{
		Bucket:      aws.String(d.AWSBucketName), // Required
		Key:         aws.String(path),            // Required
		Body:        bytes.NewReader(buf.Bytes()),
		ContentType: aws.String(contentType),
	}
	resp, err := svc.PutObject(params)

	if awserr := aws.Error(err); awserr != nil {
		// A service error occurred.
		fmt.Println("Error:", awserr.Code, awserr.Message)
		return false
	} else if err != nil {
		// A non-service error occurred.
		panic(err)
		return false
	}

	// Pretty-print the response data.
	fmt.Println(awsutil.StringValue(resp))

	return true
}

func stringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}
