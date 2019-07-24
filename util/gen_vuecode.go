/**
 * Created by Wangwei on 2019-07-03 16:56.
 */

package util

import (
	"fmt"
	"strconv"
	"strings"
)

var (
	ApiJsPath          = "%s/output/%s/src/api/api.js"
	BaseVuePagePath    = "%s/output/%s/src/views/base/%ss.vue"
	VueRoutePath       = "%s/output/%s/src/router/routes.js"
	VuePostcssrcjsPath = "%s/output/%s/.postcssrc.js"
	VueBabelrcPath     = "%s/output/%s/.babelrc"
	Space              = "  "
)

// 创建项目
func CreateVueProject() {
	if *NewProject == "" {
		return
	}

	CopyDir(getProjectPath()+"/templates/web", getProjectPath()+"/output/"+getVueProjectName())
	renameProjectName()

	genVueProjectCodes()
	RunProgressBar("前端代码生成:", 50)
}

// 生成项目代码
func genVueProjectCodes() {
	models := make(map[string]interface{})
	ReadJSON(getProjectPath()+"/configs/org.json", &models)

	var apiJsContents, vueRouteContents string
	for modelName := range models {
		if modelName == "desc" {
			modelDescs = models[modelName].(map[string]interface{})
			continue
		}

		fields := models[modelName].(map[string]interface{})
		genVuePageCodes(modelName, fields)

		apiJsContents += genApiJsContent(modelName)
		vueRouteContents += genVueRouteContent(modelName)
	}

	genApiJsCodes(apiJsContents)
	genVueRouteCodes(vueRouteContents)
	genPostcssrcJs()
	genBabelrc()
}

// 生成模块代码
func GenVueModuleCodes() {
	models := make(map[string]interface{})
	ReadJSON(getProjectPath()+"/configs/new_gen_module.json", &models)

	var apiJsContents, vueRouteContents string
	for modelName := range models {
		if modelName == "desc" {
			modelDescs = models[modelName].(map[string]interface{})
			continue
		}

		fields := models[modelName].(map[string]interface{})
		genVuePageCodes(modelName, fields)

		apiJsContents += genApiJsContent(modelName)
		vueRouteContents += genVueRouteContent(modelName)
	}

	genApiJsCodes(apiJsContents)
	genVueRouteCodes(vueRouteContents)
	RunProgressBar("模块前端代码生成:", 100)
}

// 生成page代码
func genVuePageCodes(modelName string, fields map[string]interface{}) {
	modelName = CamelString(modelName)

	templatePath := getTemplatePath("vue_page")
	content := ReadTemplate(templatePath)
	content = formatContent(modelName, content)
	content = formatContentName(modelName, content)

	tableColumnContents, formContents := genVuePageFormContents(fields)
	content = strings.Replace(content, "${formContents}", formContents, -1)
	content = strings.Replace(content, "${tableColumnContents}", tableColumnContents, -1)
	content = strings.Replace(content, "&nbsp;", Space, -1)

	ruleContents := genRuleContents(fields)
	content = strings.Replace(content, "${ruleContents}", ruleContents, -1)

	fileName := GetVueNewFilePath(BaseVuePagePath, getVueProjectName(), SnakeString(modelName))
	GenCodeFile(fileName, content)
}

func genApiJsCodes(apiJsContents string) {
	fileName := GetVueNewFilePath(ApiJsPath, getVueProjectName(), "")
	templatePath := getTemplatePath("api_js")
	content := ReadTemplate(templatePath)
	content = strings.Replace(content, "${api_rows}", apiJsContents, -1)

	GenCodeFile(fileName, content)
}

func genApiJsContent(modelName string) string {
	apisContent1 := "// ${modelNameCn}"
	apisContent2 := `export const find${modelName}s = params => { return axios.get("/v1/admin_api/${snakeModelName}/query", { params: params }).then(res => res.data) }`
	apisContent3 := `export const save${modelName} = params => { return axios.post("/v1/admin_api/${snakeModelName}/save", params).then(res => res.data) }`
	apisContent4 := `export const del${modelName}  = params => { return axios.get("/v1/admin_api/${snakeModelName}/del", {params}).then(res => res.data) }`

	apisContent := fmt.Sprintf("%s\n%s\n%s\n%s\n\n", apisContent1, apisContent2, apisContent3, apisContent4)
	modelName = CamelString(modelName)
	snakeModelName := SnakeString(modelName)
	modelNameCn := modelDescs[snakeModelName].(string)

	apisContent = strings.Replace(apisContent, "${modelNameCn}", modelNameCn, -1)
	apisContent = strings.Replace(apisContent, "${modelName}", modelName, -1)
	apisContent = strings.Replace(apisContent, "${snakeModelName}", snakeModelName, -1)
	return apisContent
}

func genVuePageFormContents(fields map[string]interface{}) (string, string) {
	var sortFormContents, sortTableColumnContents = [30]string{}, [30]string{}
	for fieldName := range fields {
		fieldInfo := fields[fieldName].(string)
		fieldName = CamelString(fieldName)
		if fieldName == SkipFields[0] || fieldName == SkipFields[1] || fieldName == SkipFields[2] {
			continue
		}

		tags := strings.Split(fieldInfo, ",")
		snakeFieldName := SnakeString(fieldName)

		formStart := `<el-form-item label="${fieldNameCn}" prop="${fieldName}">`
		formBody := `<el-input v-model="form.${fieldName}" placeholder="请输入${fieldNameCn}" maxlength="${maxLength}" class="form-item"></el-input>`
		formEnd := `</el-form-item>`
		formContent := fmt.Sprintf("\n\t\t\t\t%s\n\t\t\t\t\t%s\n\t\t\t\t%s", formStart, formBody, formEnd)

		formContent = strings.Replace(formContent, "${fieldName}", snakeFieldName, -1)
		formContent = strings.Replace(formContent, "${fieldNameCn}", strings.TrimSpace(tags[2]), -1)
		formContent = strings.Replace(formContent, "${maxLength}", strings.TrimSpace(tags[1]), -1)

		orderIndex, err := strconv.Atoi(strings.TrimSpace(tags[3]))
		if err != nil {
			panic("排序字段orderIndex不是数字类型！")
		}

		sortFormContents[orderIndex] = formContent

		tableColumnContent := genVuePageTableColumnContent(snakeFieldName, tags[2])
		sortTableColumnContents[orderIndex] = tableColumnContent
	}

	var formContents, tableColumnContents string
	for _, formContent := range sortFormContents {
		if strings.TrimSpace(formContent) == "" {
			continue
		}

		formContents += formContent
	}

	for _, tableColumnContent := range sortTableColumnContents {
		if strings.TrimSpace(tableColumnContent) == "" {
			continue
		}

		tableColumnContents += tableColumnContent
	}

	return tableColumnContents, formContents
}

func genVuePageTableColumnContent(snakeFieldName, fieldNameCn string) string {
	tableColumnContent := `<el-table-column prop="${fieldName}" label="${fieldNameCn}" align="center"></el-table-column>`
	tableColumnContent = strings.Replace(tableColumnContent, "${fieldName}", snakeFieldName, -1)
	tableColumnContent = strings.Replace(tableColumnContent, "${fieldNameCn}", fieldNameCn, -1)
	tableColumnContent = fmt.Sprintf("\n\t\t\t%s", tableColumnContent)

	return tableColumnContent
}

func genRuleContents(fields map[string]interface{}) string {
	var sortRuleContents = [30]string{}
	for fieldName := range fields {
		fieldInfo := fields[fieldName].(string)
		fieldName = CamelString(fieldName)
		snakeFieldName := SnakeString(fieldName)

		if !*NewModule {
			if snakeFieldName == "org_type_id" || snakeFieldName == "account" || snakeFieldName == "password" {
				continue
			}
		}

		tags := strings.Split(fieldInfo, ",")
		orderIndex, err := strconv.Atoi(strings.TrimSpace(tags[3]))
		if err != nil {
			panic("排序字段orderIndex不是数字类型！")
		}

		ruleContent := genRuleContent(snakeFieldName, tags[2])
		sortRuleContents[orderIndex] = ruleContent
	}

	var ruleContents string
	for _, ruleContent := range sortRuleContents {
		if strings.TrimSpace(ruleContent) == "" {
			continue
		}

		ruleContents += fmt.Sprintf("%s\n\t\t", ruleContent)
	}

	ruleContents = ruleContents[0 : len(ruleContents)-3]
	return ruleContents
}

func genRuleContent(snakeFieldName, fieldNameCn string) string {
	ruleContent := `${fieldName}: [{ required: true, message: "${fieldNameCn}不能为空", trigger: "blur" }],`
	ruleContent = strings.Replace(ruleContent, "${fieldName}", snakeFieldName, -1)
	ruleContent = strings.Replace(ruleContent, "${fieldNameCn}", fieldNameCn, -1)

	return ruleContent
}

func genVueRouteCodes(vueRouteContents string) {
	vueRouteContents = vueRouteContents[1 : len(vueRouteContents)-3]

	fileName := GetVueNewFilePath(VueRoutePath, getVueProjectName(), "")
	templatePath := getTemplatePath("vue_routes")
	content := ReadTemplate(templatePath)
	content = strings.Replace(content, "${vueRoutes}", vueRouteContents, -1)

	GenCodeFile(fileName, content)
}

func genVueRouteContent(modelName string) string {
	snakeModelName := SnakeString(modelName)
	modelNameCn := modelDescs[SnakeString(modelName)].(string)

	vueRouteContent := "{ path: '${snakeModelName}s', component: () => import('@/views/base/${snakeModelName}s'), name: '${snakeModelName}s', meta: { title: '${modelNameCn}管理', noCache: true } },"
	vueRouteContent = strings.Replace(vueRouteContent, "${snakeModelName}", snakeModelName, -1)
	vueRouteContent = strings.Replace(vueRouteContent, "${modelNameCn}", modelNameCn, -1)
	vueRouteContent = fmt.Sprintf("\t%s\n\t\t", vueRouteContent)

	return vueRouteContent
}

func genPostcssrcJs() {
	fileName := GetVueNewFilePath(VuePostcssrcjsPath, getVueProjectName(), "")
	content := ReadTemplate(getProjectPath() + "/templates/postcssrc_js.tpl")
	GenCodeFile(fileName, content)
}

func genBabelrc() {
	fileName := GetVueNewFilePath(VueBabelrcPath, getVueProjectName(), "")
	content := ReadTemplate(getProjectPath() + "/templates/babelrc.tpl")
	GenCodeFile(fileName, content)
}

func GetVueNewFilePath(filePath, projectName, modelName string) string {
	if *NewModule {
		projectName = "gencodes/vuecode"
	}

	if strings.TrimSpace(modelName) == "" {
		return fmt.Sprintf(filePath, getProjectPath(), projectName)
	}

	return fmt.Sprintf(filePath, getProjectPath(), projectName, modelName)
}

func getVueProjectName() string {
	return fmt.Sprintf("%s-admin", *NewProject)
}
