package controllers

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/gin-gonic/gin"
	"github.com/sajjad3k/contactsmanager/models"
)

func ShowallContacts(c *gin.Context) {
	var resp models.Response
	list, err := models.Getlist()
	if err != nil {
		resp.Data = list
		resp.Status = "Failed"
		resp.Message = "No Contacts Found"
		c.JSON(http.StatusNoContent, resp)
	} else {
		resp.Data = list
		resp.Status = "success"
		resp.Message = "Contacts Found"
		c.JSON(http.StatusOK, resp)
	}

}

func GetcontactbyName(c *gin.Context) {
	var resp models.Response
	var contac models.Contact

	name := c.Params.ByName("name")
	if name == "" {
		c.JSON(http.StatusBadRequest, resp)
	} else {
		list, err := models.Getlist()
		if err != nil {
			resp.Data = list
			resp.Status = "Failed"
			resp.Message = "No Contacts Found"
			c.JSON(http.StatusNoContent, resp)
		} else {

			for _, value := range list {
				if value.Name == name {
					contac = value
					break
				} else {
					continue
				}
			}
			if contac.Number == "" {
				resp.Data = list
				resp.Status = "Failed"
				resp.Message = "Contact not found"
				c.JSON(http.StatusNoContent, resp)
			} else {
				resp.Data = append(resp.Data, contac)
				resp.Status = "success"
				resp.Message = "Contact Found"
				c.JSON(http.StatusOK, resp)
			}

		}
	}

}

func GetcontactbyNumber(c *gin.Context) {
	var resp models.Response
	var contac models.Contact

	number := c.Params.ByName("number")
	if number == "" {
		c.JSON(http.StatusBadRequest, resp)
	} else {
		list, err := models.Getlist()
		if err != nil {
			resp.Data = list
			resp.Status = "Failed"
			resp.Message = "No Contacts Found"
			c.JSON(http.StatusNoContent, resp)
		} else {

			for _, value := range list {
				if value.Number == number {
					contac = value
					break
				} else {
					continue
				}
			}
			if contac.Number == "" {
				resp.Data = list
				resp.Status = "Failed"
				resp.Message = "Contact not found"
				c.JSON(http.StatusNoContent, resp)
			} else {
				resp.Data = append(resp.Data, contac)
				resp.Status = "success"
				resp.Message = "Contact Found"
				c.JSON(http.StatusOK, resp)
			}

		}
	}
}

func CreatenewContact(c *gin.Context) {
	var in models.Contact
	var resp models.Response
	//var list []models.Contact
	c.BindJSON(&in)
	if in.Name == "" || in.Number == "" {
		resp.Status = "failed"
		resp.Message = "add details properly"
		c.JSON(http.StatusBadRequest, resp)
	} else {
		list, err := models.Getlist()
		if err == nil || err != nil {
			list = append(list, in)
			models.Setlist(list)
			resp.Message = "contact added to list"
			resp.Data = append(resp.Data, in)
			resp.Status = "success"
		}
	}
}

func Updatecontact(c *gin.Context) {

	var resp models.Response
	var contac models.Contact

	name := c.Params.ByName("name")

	var flag bool = false

	if name != "" {
		c.BindJSON(&contac)
		p, err := models.Getlist()
		if err != nil {
			resp.Status = "failed"
			resp.Message = "no contacts found to update"
			c.JSON(http.StatusNoContent, resp)
		} else {
			for i, val := range p {
				if val.Name == name {
					p[i].Name = contac.Name
					p[i].Number = contac.Number
					p[i].Email = contac.Email
					flag = true
				}
			}
			if flag == true {
				models.Setlist(p)
				resp.Status = "success"
				resp.Message = "The contact is updated successfully"
				resp.Data = append(resp.Data, contac)
				c.JSON(http.StatusOK, resp)
			}

		}

	} else {
		resp.Status = "Error"
		resp.Message = "The request is not correct"
		c.JSON(http.StatusBadRequest, resp)
	}

}

func DeleteContact(c *gin.Context) {
	var list []models.Contact
	var resp models.Response
	//var ins models.Contact
	//var fla bool
	name := c.Params.ByName("name")
	flag := false
	if name == "" {
		resp.Status = "failed"
		resp.Message = "enter a contact name to delete"
		c.JSON(http.StatusBadRequest, resp)
	} else {

		data, err := models.Getlist()
		if err != nil {
			resp.Status = "failed"
			resp.Message = "No contacts are there to delete"
			c.JSON(http.StatusNoContent, resp)
		} else {
			for _, val := range data {
				if val.Name == name {
					//fla = true
					//ins = val
					resp.Data = append(resp.Data, val)
					flag = true
					continue
				}
				list = append(list, val)
			}
			if flag != false {
				resp.Status = "Success"
				resp.Message = "Deleted the Customer entry"
				models.Setlist(list)
				c.JSON(http.StatusOK, resp)
			} else {
				resp.Status = "failed"
				resp.Message = "contact not found"
				c.JSON(http.StatusNotFound, resp)
			}
		}
	}
}

func Uploadlist(c *gin.Context) {

	var resp models.Response

	list, err := models.Getlist()

	if err != nil {
		resp.Status = "Failed"
		resp.Message = "Empty contact list"
		c.JSON(http.StatusNoContent, resp)
	}

	out, _ := json.MarshalIndent(list, "", " ")

	file, erra := ioutil.TempFile("", "contactslist.json")
	if err != nil {
		log.Fatal(erra)
	}
	defer os.Remove(file.Name())

	if _, errb := file.Write(out); err != nil {
		log.Fatal(errb)
	}

	if errc := file.Close(); err != nil {
		log.Fatal(errc)
	}

	//aws upload part

	aws_access_key_id := "your_access_key_id"
	aws_secret_access_key := "your_secret_access_key"

	const (
		awsregion     = "region"
		AWS_S3_BUCKET = "bucketname"
	)

	ses := func() *session.Session {

		sess, errd := session.NewSession(
			&aws.Config{
				Region: aws.String(awsregion),
				Credentials: credentials.NewStaticCredentials(
					aws_access_key_id,
					aws_secret_access_key,
					"", // a token will be created when the session it's used.
				),
			})

		if errd != nil {
			log.Fatal(errd)
		}
		return sess
	}()

	uploader := s3manager.NewUploader(ses)

	_, err = uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(AWS_S3_BUCKET), // Bucket to be used
		Key:    aws.String("contactlist"), // Name of the file to be saved
		Body:   file,
	})
	if err != nil {
		resp.Status = "failed"
		resp.Message = "contacts upload failed"
	}
	resp.Status = "success"
	resp.Message = "contacts upload successfull"
	c.JSON(http.StatusCreated, resp)
}
