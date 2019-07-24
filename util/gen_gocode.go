/**
 * Created by Wangwei on 2019-07-02 15:41.
 */

package util

import (
	"flag"
	"fmt"
	"gotools/pkg/gopath"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/gogf/gf/g/util/gconv"
)

var (
	FieldHeaderTpl     = "%s\t%s\t"
	FieldTagTpl        = "`xorm:\"%s(%s) comment('%s')\"`"
	BaseModelPath      = "%s/output/%s/internal/model/base/%s.go"
	BaseServicePath    = "%s/output/%s/internal/service/base/%s.go"
	BaseControllerPath = "%s/output/%s/internal/api/admin/http/v1/%s.go"
	BaseRoutesPath     = "%s/output/%s/internal/api/admin/http/v1/routes.go"
	BaseNewModelPath   = "%s/output/%s/internal/model/model.go"

	SkipFields = []string{"OrgTypeId", "OrgTypeName", "Account"}
	modelDescs = make(map[string]interface{})
)

var NewProject = flag.String("newProject", "", "新建项目名称")
var NewModule = flag.Bool("newModule", false, "新模块生成代码")
var ProjectName = flag.String("projectName", "", "生成新模块代码所需的项目名称")

// 创建项目
func CreateGoProject() {
	if *NewProject == "" {
		return
	}

	CopyDir(getProjectPath()+"/templates/project", getProjectPath()+"/output/"+getGoProjectName())
	renameProjectName()

	genGoProjectCodes()
	RunProgressBar("后端代码生成:", 50)
}

type Org struct {
	Id        int64  `json:"id"`
	Parent_id int64  `json:"parent_id"`
	Name      string `json:"name"`
	Code      string `json:"code"`
}

// 生成项目代码
func genGoProjectCodes() {
	models := make(map[string]interface{})
	ReadJSON(getProjectPath()+"/configs/org.json", &models)

	var routeContents, newModelContents string
	var modelsOrgType []Org
	ReadJSON("../configs/org_type.json", &modelsOrgType)

	for modelName := range models {
		if modelName == "desc" {
			modelDescs = models[modelName].(map[string]interface{})
			continue
		}

		var modelNameId int64 = -1
		var nextName string
		for _, v := range modelsOrgType {
			if v.Code == modelName {
				modelNameId = v.Id
				break
			}
		}
		for _, v := range modelsOrgType {
			if v.Parent_id == modelNameId {
				nextName = v.Code
				break
			}
		}
		fields := models[modelName].(map[string]interface{})
		genModelCodes(modelName, fields)
		genServiceCodes(modelName, nextName)
		genControllerCodes(modelName)

		routeContents += genRouteContent(modelName)
		newModelContents += genNewModelContent(modelName)
	}

	// 生成路由配置代码
	genRouteCodes(routeContents)
	genNewModelCodes(newModelContents)
	genMySQLConfig()
	genCasbinConfig()
}

// 新增加模块生成代码
func GenModuleCodes() {
	if !*NewModule {
		return
	}

	genGoModuleCodes()
	GenVueModuleCodes()
}

func genGoModuleCodes() {
	models := make(map[string]interface{})
	ReadJSON("../configs/new_gen_module.json", &models)

	var modelsOrgType []Org
	ReadJSON("../configs/org_type.json", &modelsOrgType)

	var routeContents string
	for modelName := range models {
		if modelName == "desc" {
			modelDescs = models[modelName].(map[string]interface{})
			continue
		}

		var modelNameId int64 = -1
		var nextName string
		for _, v := range modelsOrgType {
			if v.Code == modelName {
				modelNameId = v.Id
				break
			}
		}
		for _, v := range modelsOrgType {
			if v.Parent_id == modelNameId {
				nextName = v.Code
				break
			}
		}

		fields := models[modelName].(map[string]interface{})
		genModelCodes(modelName, fields)
		genServiceCodes(modelName, nextName)
		genControllerCodes(modelName)

		routeContents += genRouteContent(modelName)

	}

	// 生成路由配置代码
	genRouteCodes(routeContents)
	RunProgressBar("模块后端代码生成:", 100)
}

// 生成model代码
func genModelCodes(modelName string, fields map[string]interface{}) {
	modelName = CamelString(modelName)

	var softJustFields, sortModelFields = [30]string{}, [30]string{}
	for fieldName := range fields {
		fieldInfo := fields[fieldName].(string)
		fieldName = CamelString(fieldName)
		if fieldName == SkipFields[0] || fieldName == SkipFields[1] || fieldName == SkipFields[2] {
			continue
		}

		tags := strings.Split(fieldInfo, ",")
		fieldHeader := fmt.Sprintf(FieldHeaderTpl, fieldName, tags[0])
		fieldTag := fmt.Sprintf(FieldTagTpl, TypeConvert(tags[0]), tags[1], tags[2])
		orderIndex := gconv.Int(tags[3]) - 1

		sortModelFields[orderIndex] = fmt.Sprintf("%s %s\n\t", fieldHeader, fieldTag)
		softJustFields[orderIndex] = fmt.Sprintf("%s\n\t", fieldHeader)
	}

	var justFields, modelFields string
	for _, modelField := range softJustFields {
		modelFields += modelField
	}

	for _, justField := range softJustFields {
		justFields += justField
	}

	modelFields = modelFields[0 : len(modelFields)-2]
	justFields = justFields[0 : len(justFields)-2]

	templatePath := getTemplatePath("model")
	content := ReadTemplate(templatePath)
	content = formatModelContent(modelName, modelFields, justFields, content)

	fileName := GetGoNewFilePath(BaseModelPath, getGoProjectName(), SnakeString(modelName))
	GenCodeFile(fileName, content)
}

func GetGoNewFilePath(filePath, projectName, modelName string) string {
	if *NewModule {
		projectName = "gencodes/gocode"
	}

	if strings.TrimSpace(modelName) == "" {
		return fmt.Sprintf(filePath, getProjectPath(), projectName)
	}

	return fmt.Sprintf(filePath, getProjectPath(), projectName, modelName)
}

// 生成model/model.go中所需的代码
func genNewModelContent(modelName string) string {
	modelName = CamelString(modelName)
	content := fmt.Sprintf("new(base.%s),\n\t\t", modelName)
	return content
}

func genNewModelCodes(newModelContents string) {
	templatePath := getTemplatePath("model_go")

	content := ReadTemplate(templatePath)
	content = formatContent("", content)
	content = strings.Replace(content, "${newModels}", newModelContents, -1)

	fileName := GetGoNewFilePath(BaseNewModelPath, getGoProjectName(), "")
	GenCodeFile(fileName, content)
}

// 生成service代码
func genServiceCodes(modelName, nextName string) {
	modelName = CamelString(modelName)
	nextName = CamelString(nextName)
	templatePath := getTemplatePath("service")

	content := ReadTemplate(templatePath)
	content = formatContent(modelName, content)
	content = formatContentDel(nextName, content)

	fileName := GetGoNewFilePath(BaseServicePath, getGoProjectName(), SnakeString(modelName))
	GenCodeFile(fileName, content)
}

// 生成controller代码
func genControllerCodes(modelName string) {
	modelName = CamelString(modelName)
	templatePath := getTemplatePath("controller")

	content := ReadTemplate(templatePath)
	content = formatContent(modelName, content)

	fileName := GetGoNewFilePath(BaseControllerPath, getGoProjectName(), SnakeString(modelName))
	GenCodeFile(fileName, content)
}

// 生成路由设置代码
func genRouteCodes(routeContent string) {
	// 去除多余的空行
	routeContent = routeContent[0 : len(routeContent)-3]
	templatePath := getTemplatePath("routes")

	content := ReadTemplate(templatePath)
	content = strings.Replace(content, "${project}", getGoProjectName(), -1)
	content = strings.Replace(content, "${SetupController}", routeContent, -1)

	fileName := GetGoNewFilePath(BaseRoutesPath, getGoProjectName(), "")
	GenCodeFile(fileName, content)
}

// 生成单行路由
func genRouteContent(modelName string) string {
	modelName = CamelString(modelName)
	lowerModelName := CamelStringFirstLower(modelName)
	routeContentLine1 := `${lowerModelName}Controller := ${modelName}Controller{router}`
	routeContentLine2 := `${lowerModelName}Controller.Setup()`

	routeContent := fmt.Sprintf("%s\n\t%s\n\n\t", routeContentLine1, routeContentLine2)
	routeContent = strings.Replace(routeContent, "${lowerModelName}", lowerModelName, -1)
	routeContent = strings.Replace(routeContent, "${modelName}", modelName, -1)

	return routeContent
}

func formatModelContent(modelName, modelFields, justFields, content string) string {
	content = strings.Replace(content, "${modelName}", modelName, -1)
	content = strings.Replace(content, "${modelFields}", modelFields, -1)
	content = strings.Replace(content, "${justFields}", justFields, -1)

	return content
}

func formatContent(modelName, content string) string {
	lowerModelName := CamelStringFirstLower(modelName)
	snakeModelName := SnakeString(modelName)

	content = strings.Replace(content, "${project}", getGoProjectName(), -1)
	content = strings.Replace(content, "${modelName}", modelName, -1)
	content = strings.Replace(content, "${lowerModelName}", lowerModelName, -1)
	content = strings.Replace(content, "${snakeModelName}", snakeModelName, -1)

	return content
}

func formatContentName(modelName, content string) string {
	snakeModelName := SnakeString(modelName)

	var modelsOrgType []Org
	ReadJSON("../configs/org_type.json", &modelsOrgType)

	var orgTypeName string
	for _, v := range modelsOrgType {
		if v.Code == snakeModelName {
			orgTypeName = v.Name
			break
		}
	}

	content = strings.Replace(content, "${orgTypeName}", orgTypeName, -1)

	return content
}
func formatContentDel(nextName, content string) string {

	deleteStatement := fmt.Sprintf("	var count int64\ncount, err = DB.Where(\"org_id = ?\", id).Count(&base.%s{})\n"+
		"if err != nil {\nreturn\n}\nif count>0 {\nreturn errors.New(\"先删除下一级用户\")\n}", nextName)

	if nextName != "" {
		content = strings.Replace(content, "${errors}", "\"errors\"", -1)
		content = strings.Replace(content, "${deleteStatement}", deleteStatement, -1)
		return content
	}
	content = strings.Replace(content, "${errors}", "", -1)
	content = strings.Replace(content, "${deleteStatement}", nextName, -1)

	return content
}

// 根据命令行参数获取模板，如果是新建项目，读取的模板不带new_开头，如果是生成新模块，则读取new_开头的模板
// 比如gotools -newProject helloworld，则读取model.tpl, service.tpl等
// 如果是gotools -newModule，则读取new_model.tpl, new_service.tpl等
func getTemplatePath(tplName string) string {
	templatePath := "%s/templates/%s%s.tpl"
	if *NewModule {
		templatePath = fmt.Sprintf(templatePath, getProjectPath(), "new_", tplName)
	} else {
		templatePath = fmt.Sprintf(templatePath, getProjectPath(), "", tplName)
	}

	return templatePath
}

// 批量重命名项目文件中的项目名称
func renameProjectName() {
	projectPath := fmt.Sprintf("%s/output/%s", getProjectPath(), getGoProjectName())

	// 遍历文件夹，获取文件路径
	paths := make([]string, 0)
	filepath.Walk(projectPath, func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			paths = append(paths, path)
		}
		return nil
	})

	// 遍历文件路径，修改文件名
	for _, filepath := range paths {
		data, err := ioutil.ReadFile(filepath)
		if err != nil {
			fmt.Printf("read file data err: %s\n", err)
			return
		}

		content := strings.Replace(string(data), "goadmin", getGoProjectName(), -1)
		GenCodeFile(filepath, string(content))
	}
}

func getGoProjectName() string {
	if *NewProject == "" {
		return fmt.Sprintf("%s-go", *ProjectName)
	}

	return fmt.Sprintf("%s-go", *NewProject)
}

func getProjectPath() string {
	currentDir := gopath.GetCurrentDirectory()
	projectPath := gopath.GetParentDirectory(currentDir)
	return projectPath
}
