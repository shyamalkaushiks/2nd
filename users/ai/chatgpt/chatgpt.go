package chatgpt

import (
	// "fmt"
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
)

/* request
{
	"model":"gpt-3.5-turbo-0125",
	"response_format":{
	   "type":"json_object"
	},
	"messages":[
	   {
		  "role":"user",
		  "content":"Who won the world series in 2020?"
	   }
	]
 }
*/

// ResumeParseRequest :
type ResumeParseRequest struct {
	Model          string         `json:"model"`
	ResponseFormat ResponseFormat `json:"response_format"`
	Messages       []Messages     `json:"messages"`
}

// ResponseFormat :
type ResponseFormat struct {
	Type string `json:"type"`
}

// Messages :
type Messages struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// ResumeParse :
func ResumeParse(content string) (string, error) {
	var result string

	url := "https://api.openai.com/v1/chat/completions"

	var objResumeParseRequest ResumeParseRequest
	objResumeParseRequest.Model = "gpt-3.5-turbo-0125"
	objResumeParseRequest.ResponseFormat.Type = "json_object"

	var messages []Messages
	var message Messages
	message.Role = "user"
	message.Content = content +
		"   this is data of resume pass it in a JSON format like this given example in the most optimized field and in professional summary and in all the description it will be parse with phtml tag as example " +
		"   date format will be mm/dd/yyyy             " +
		"{\"templateData\":{\"templatename\":null,\"firstname\":\"BKDADAD\",\"lastname\":\"ADAD\",\"address\":\"1/ANEWCOLNEY,KALWARROAD,JHOTWARAJAIPUR\",\"postalcode\":\"302012\",\"drivinglicense\":\"adsads\",\"nationality\":\"adsads\",\"placeofbirth\":\"adsadsadsadsd\",\"dateofbirth\":\"2024-02-23\",\"securityclearance\":\"DFKLM'DSKLMF\",\"city\":\"JAIPUR\",\"country\":\"India\",\"email\":\"ADADA\",\"wantedjobtitle\":\"HREXECUTIVR\",\"phone\":\"07976159618\",\"photo\":\"/Users/jyotirmaygaur/depresback/resume_builder/routes/uploads/766594a20d729af8aae35c0c9d9fbd24\",\"education\":[{\"end\":\"2024-02-08\",\"city1\":\"CXVCVXVCXV\",\"start\":\"2024-02-17\",\"degree\":\"CXVCXVCVX\",\"school\":\"BCFVXBVC\",\"edudescription\":\"CXVVCXCVXXCVVVCXCVXV\"}],\"professionalsummary\":\"<p>nullIamahighlymotivatedandexperiencedHRExecutivewithastrongbackgroundinhumanresourcesmanagement</p>\",\"untitledsection\":[],\"employmenthistory\":[{\"jobend\":\"2024-02-16\",\"jobcity\":\"bcvbc\",\"employer\":\"bcbc\",\"jobstart\":\"2024-02-02\",\"jobtitle1\":\"APE\",\"jobdescription\":\"<p>nullIamahighlymotivatedandexperiencedHRExecutivewithastrongbackgroundinhumanresourcesmanagement</p>\"}],\"social\":[{\"link\":\"LINKEDIN\",\"label\":\"LINKEDIN\"}],\"skills\":[\"joyKJMNJO\",\"LJNJKLHN\",\"KL;MJJKLMN\"]}}"
	// "{\"templateData\":{\"templatename\":null,\"firstname\":\"BKDADAD\",\"lastname\":\"ADAD\",\"address\":\"1/ANEWCOLNEY,KALWARROAD,JHOTWARAJAIPUR\",\"postalcode\":\"302012\",\"drivinglicense\":\"adsads\",\"nationality\":\"adsads\",\"placeofbirth\":\"adsadsadsadsd\",\"dateofbirth\":\"2024-02-23\",\"securityclearance\":\"DFKLM'DSKLMF\",\"city\":\"JAIPUR\",\"country\":\"India\",\"email\":\"ADADA\",\"wantedjobtitle\":\"HREXECUTIVR\",\"phone\":\"07976159618\",\"photo\":\"/Users/jyotirmaygaur/depresback/resume_builder/routes/uploads/766594a20d729af8aae35c0c9d9fbd24\",\"education\":[{\"end\":\"2024-02-08\",\"city1\":\"CXVCVXVCXV\",\"start\":\"2024-02-17\",\"degree\":\"CXVCXVCVX\",\"school\":\"BCFVXBVC\",\"edudescription\":\"CXVVCXCVXXCVVVCXCVXV\"}],\"professionalsummary\":\"nullIamahighlymotivatedandexperiencedHRExecutivewithastrongbackgroundinhumanresourcesmanagement\",\"untitledsection\":[],\"employmenthistory\":[{\"jobend\":\"2024-02-16\",\"jobcity\":\"bcvbc\",\"employer\":\"bcbc\",\"jobstart\":\"2024-02-02\",\"jobtitle1\":\"APE\",\"jobdescription\":\"IamadedicatedandexperiencedHRprofessionalwithastrongpassionfordevelopingandimplementingeffectiveHRstrategies\"}],\"social\":[{\"link\":\"LINKEDIN\",\"label\":\"LINKEDIN\"}],\"skills\":[\"joyKJMNJO\",\"LJNJKLHN\",\"KL;MJJKLMN\"]}}"
	// "{\"jobsMyResumeData\":{\"resumeHeadline\":\"JobboardcurrentlylivinginUSA\",\"skillsData\":[\"javascript\",\"css\",\"html\",\"bootstrap\"],\"employmentData\":[{\"jobTitle\":\"JuniorSoftwareDeveloper\",\"company\":\"W3ITEXPERTS\",\"jobStartDate\":\"10-2015\",\"jobEndDate\":\"Present\",\"jobDescription\":\"JuniorSoftwareDeveloper\"}],\"educationData\":[{\"education\":\"Doctorate/PHD\",\"course\":\"SelectCourse\",\"university\":\"SelectUniversity\"}],\"itSkillsData\":[{\"skills\":\"Bootstrap\",\"version\":\"3\",\"lastUsed\":\"2018\",\"experience\":\"1Year5Months\"}],\"projectsData\":[{\"projectTitle\":\"Enterprojecttitle\",\"clientName\":\"EnterClientName\",\"projectStatus\":\"INPROGRESS\",\"workStartedFromYear\":\"2018\",\"workStartedFromMonth\":\"January\",\"workedTillYear\":\"2018\",\"workedTillMonth\":\"January\",\"projectDescription\":\"HeythisistheprojectDescription\",}],\"profileSummary\":\"YourProfileSummaryshouldmentionthehighlightsofyourcareerandeducation,whatyourprofessionalinterestsare,andwhatkindofacareeryouarelookingfor.Writeameaningfulsummaryofmorethan50characters.\"},\"accomplishments\":{\"onlineProfile\":\"AddLinkstoOnlineProfiles\",\"workSample\":\"AddLinkstoyourprojects\",\"whitePaper\":\"AddLinkstoyouronlinepublications\",\"presentation\":\"AddLinkstoyouronlinepresentations\",\"patent\":\"AddDetailsofpatentsyouhavefiled.\",},\"desiredCareerProfile\":{\"industry\":\"ItSoftwareservices\",\"role\":\"WebDesigner\",\"employmentType\":\"FullTime\",\"AvailabilityToJoin\":\"12July\",\"desiredLocation\":\"AddDesiredLocation\",\"functionalArea\":\"design/creative/userExperience\",\"jobType\":\"permanent\",\"desiredShift\":\"addDesiredShift\",\"expectedSalary\":\"1Lakhs\",\"desiredIndustry\":\"adddesiredIndustry\"},\"personalDetails\":{\"dateOfBirth\":\"31July0998\",\"gender\":\"male\",\"maritalStatus\":\"unmarried\",\"passportNumber\":\"+1234567890\",\"differentlyAbled\":\"None\",\"languages\":\"English\",\"addPermanentAddress\":\"AddPermanentAddress\",\"areaPinCode\":\"302010\",\"hometown\":\"Delhi\",\"workPermitOfOtherCountry\":\"USA\"},\"attachResume\":\"\"}"

	messages = append(messages, message)
	objResumeParseRequest.Messages = messages

	out, err := json.Marshal(objResumeParseRequest)
	if err != nil {
		return result, err
	}

	// fmt.Println(string(out))

	requestBody := []byte(string(out))

	request, err := http.NewRequest("POST", url, bytes.NewBuffer(requestBody))
	if err != nil {
		return result, err
	}

	request.Header.Add("Content-Type", "application/json")
	request.Header.Add("Authorization", "")

	respone, err := http.DefaultClient.Do(request)
	if err != nil {
		return result, err
	}
	defer respone.Body.Close()

	body, err := ioutil.ReadAll(respone.Body)
	if err != nil {
		return result, err
	}

	// fmt.Println(string(body))

	result = string(body)

	return result, nil
}
