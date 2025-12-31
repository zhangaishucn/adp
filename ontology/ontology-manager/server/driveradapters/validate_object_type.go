package driveradapters

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/dlclark/regexp2"
	libCommon "github.com/kweaver-ai/kweaver-go-lib/common"
	"github.com/kweaver-ai/kweaver-go-lib/rest"

	oerrors "ontology-manager/errors"
	"ontology-manager/interfaces"
)

// 对象类必要创建参数的非空校验。
func ValidateObjectType(ctx context.Context, objectType *interfaces.ObjectType) error {
	// 校验id的合法性
	err := validateID(ctx, objectType.OTID)
	if err != nil {
		return err
	}

	// 校验名称合法性
	// 去掉名称的前后空格
	objectType.OTName = strings.TrimSpace(objectType.OTName)
	err = validateObjectName(ctx, objectType.OTName, interfaces.MODULE_TYPE_OBJECT_TYPE)
	if err != nil {
		return err
	}

	// 若输入了 tags，校验 tags 的合法性
	err = ValidateTags(ctx, objectType.Tags)
	if err != nil {
		return err
	}

	// 去掉tag前后空格以及数组去重
	objectType.Tags = libCommon.TagSliceTransform(objectType.Tags)

	// 校验comment合法性
	err = validateObjectComment(ctx, objectType.Comment)
	if err != nil {
		return err
	}

	// 属性名只包含小写英文字母和数字和下划线(_)和连字符(-)，且不能以下划线开头，不能超过40个字符
	dataPropMap := map[string]*interfaces.DataProperty{}
	for _, prop := range objectType.DataProperties {
		// 校验属性名的合法性,与id的规则一样
		err := ValidatePropertyName(ctx, prop.Name)
		if err != nil {
			return err
		}

		err = ValidateDataProperty(ctx, prop)
		if err != nil {
			return err
		}

		dataPropMap[prop.Name] = prop
	}

	// 校验主键非空
	if len(objectType.PrimaryKeys) == 0 {
		return rest.NewHTTPError(ctx, http.StatusBadRequest,
			oerrors.OntologyManager_ObjectType_NullParameter_PrimaryKeys)
	} else {
		for _, pKey := range objectType.PrimaryKeys {
			if _, ok := dataPropMap[pKey]; !ok {
				return rest.NewHTTPError(ctx, http.StatusBadRequest,
					oerrors.OntologyManager_ObjectType_InvalidParameter).
					WithErrorDetails(fmt.Sprintf("对象类[%s]主键[%s]不存在", objectType.OTName, pKey))
			}
		}
	}

	if objectType.DisplayKey == "" {
		return rest.NewHTTPError(ctx, http.StatusBadRequest,
			oerrors.OntologyManager_ObjectType_NullParameter_DisplayKey)
	} else {
		if _, ok := dataPropMap[objectType.DisplayKey]; !ok {
			return rest.NewHTTPError(ctx, http.StatusBadRequest,
				oerrors.OntologyManager_ObjectType_InvalidParameter).
				WithErrorDetails(fmt.Sprintf("对象类[%s]显示键[%s]不存在", objectType.OTName, objectType.DisplayKey))
		}
	}

	if objectType.IncrementalKey != "" {
		if field, ok := dataPropMap[objectType.IncrementalKey]; !ok {
			return rest.NewHTTPError(ctx, http.StatusBadRequest,
				oerrors.OntologyManager_ObjectType_InvalidParameter).
				WithErrorDetails(fmt.Sprintf("对象类[%s]增量键[%s]不存在", objectType.OTName, objectType.IncrementalKey))
		} else {
			switch field.Type {
			case "integer", "datetime", "timestamp":
			default:
				return rest.NewHTTPError(ctx, http.StatusBadRequest,
					oerrors.OntologyManager_ObjectType_InvalidParameter).
					WithErrorDetails(fmt.Sprintf("不支持的对象类[%s]增量键[%s]类型[%s]", objectType.OTName, field.Name, field.Type))
			}
		}
	}

	// 当逻辑属性是指标模型时，初始化3个请求参数, instant start end
	IfSystemGen := true
	for i, prop := range objectType.LogicProperties {
		// 校验属性名的合法性,与id的规则一样
		err := ValidatePropertyName(ctx, prop.Name)
		if err != nil {
			return err
		}

		// 校验属性类型和绑定的资源是相同的
		if prop.DataSource != nil {
			if prop.Type != prop.DataSource.Type {
				return rest.NewHTTPError(ctx, http.StatusBadRequest, oerrors.OntologyManager_ObjectType_InvalidParameter).
					WithErrorDetails(fmt.Sprintf("对象类[%s]逻辑属性[%s]的数据类型[%s]与其所绑定的数据资源类型[%s]不一致",
						objectType.OTName, prop.Name, prop.Type, prop.DataSource.Type))
			}
		}
		if prop.Type == interfaces.PROPERTY_TYPE_METRIC {
			// 如果不存在指标相关的参数，那么就追加
			paramMap := map[string]interfaces.Parameter{}
			for _, param := range prop.Parameters {
				paramMap[param.Name] = param
			}
			hasInstant := false
			hasStart := false
			hasEnd := false
			hasStep := false
			for pName := range paramMap {
				switch pName {
				case "instant":
					hasInstant = true
				case "start":
					hasStart = true
				case "end":
					hasEnd = true
				case "step":
					hasStep = true
				}
			}
			params := []interfaces.Parameter{}
			if !hasInstant {
				params = append(params, interfaces.Parameter{
					Name:        "instant",
					Type:        "boolean",
					Operation:   "==",
					ValueFrom:   interfaces.VALUE_FROM_INPUT,
					IfSystemGen: &IfSystemGen,
				})
			}
			if !hasStart {
				params = append(params, interfaces.Parameter{
					Name:        "start",
					Type:        "integer",
					Operation:   "==",
					ValueFrom:   interfaces.VALUE_FROM_INPUT,
					IfSystemGen: &IfSystemGen,
				})
			}
			if !hasEnd {
				params = append(params, interfaces.Parameter{
					Name:        "end",
					Type:        "integer",
					Operation:   "==",
					ValueFrom:   interfaces.VALUE_FROM_INPUT,
					IfSystemGen: &IfSystemGen,
				})
			}
			if !hasStep {
				params = append(params, interfaces.Parameter{
					Name:        "step",
					Type:        "string",
					Operation:   "==",
					ValueFrom:   interfaces.VALUE_FROM_INPUT,
					IfSystemGen: &IfSystemGen,
				})
			}
			objectType.LogicProperties[i].Parameters = append(objectType.LogicProperties[i].Parameters, params...)
		}
	}

	return nil
}

func ValidatePropertyName(ctx context.Context, name string) error {
	if name == "" {
		return rest.NewHTTPError(ctx, http.StatusBadRequest, oerrors.OntologyManager_ObjectType_NullParameter_PropertyName)
	}
	//  id，只包含大小写英文字母和数字和下划线(_)和连字符(-)，且不能以下划线开头，不能超过40个字符
	re := regexp2.MustCompile(interfaces.RegexPattern_Property_Name, regexp2.RE2)
	match, err := re.MatchString(name)
	if err != nil || !match {
		errDetails := `The property name can contain only letters, digits and underscores(_),
			it cannot start with underscores and cannot exceed 40 characters`
		return rest.NewHTTPError(ctx, http.StatusBadRequest, oerrors.OntologyManager_ObjectType_InvalidParameter_PropertyName).
			WithErrorDetails(errDetails)
	}
	return nil
}

func ValidateDataProperties(ctx context.Context, propertyNames []string, dataProperties []*interfaces.DataProperty) error {
	if len(propertyNames) != len(dataProperties) {
		httpErr := rest.NewHTTPError(ctx, http.StatusBadRequest, oerrors.OntologyManager_ObjectType_InvalidParameter).
			WithErrorDetails("PropertyNames and DataProperties length not equal")
		return httpErr
	}

	propertyNameMap := map[string]string{}
	for _, propertyName := range propertyNames {
		propertyNameMap[propertyName] = propertyName
	}
	for _, prop := range dataProperties {
		if _, ok := propertyNameMap[prop.Name]; !ok {
			httpErr := rest.NewHTTPError(ctx, http.StatusBadRequest, oerrors.OntologyManager_ObjectType_InvalidParameter).
				WithErrorDetails(fmt.Sprintf("DataProperty %s not in URL", prop.Name))
			return httpErr
		}

		err := ValidateDataProperty(ctx, prop)
		if err != nil {
			return err
		}
	}
	return nil
}

func ValidateDataProperty(ctx context.Context, dataProperty *interfaces.DataProperty) error {
	err := validateID(ctx, dataProperty.Name)
	if err != nil {
		return err
	}

	err = validateObjectName(ctx, dataProperty.DisplayName,
		interfaces.MODULE_TYPE_OBJECT_TYPE)
	if err != nil {
		return err
	}

	err = validateObjectComment(ctx, dataProperty.Comment)
	if err != nil {
		return err
	}

	if dataProperty.IndexConfig != nil {
		err = ValidateIndexConfig(ctx, *dataProperty.IndexConfig)
		if err != nil {
			return err
		}
	}

	return nil
}

func ValidateIndexConfig(ctx context.Context, indexConfig interfaces.IndexConfig) error {
	err := ValidateKeywordConfig(ctx, indexConfig.KeywordConfig)
	if err != nil {
		return err
	}
	err = ValidateFulltextConfig(ctx, indexConfig.FulltextConfig)
	if err != nil {
		return err
	}
	err = ValidateVectorConfig(ctx, indexConfig.VectorConfig)
	if err != nil {
		return err
	}

	return nil
}

func ValidateKeywordConfig(ctx context.Context, keywordConfig interfaces.KeywordConfig) error {
	if !keywordConfig.Enabled {
		return nil
	}
	if keywordConfig.IgnoreAboveLen <= 0 {
		httpErr := rest.NewHTTPError(ctx, http.StatusBadRequest, oerrors.OntologyManager_ObjectType_InvalidParameter).
			WithErrorDetails("KeywordConfig IgnoreAboveLen must be greater than 0")
		return httpErr
	}
	return nil
}

func ValidateFulltextConfig(ctx context.Context, fulltextConfig interfaces.FulltextConfig) error {
	if !fulltextConfig.Enabled {
		return nil
	}
	switch fulltextConfig.Analyzer {
	case "standard", "english", "ik_max_word", "hanlp_standard", "hanlp_index":
	default:
		httpErr := rest.NewHTTPError(ctx, http.StatusBadRequest, oerrors.OntologyManager_ObjectType_InvalidParameter).
			WithErrorDetails("FulltextConfig Analyzer must be standard, english, ik_max_word, hanlp_standard or hanlp_index")
		return httpErr
	}
	return nil
}

func ValidateVectorConfig(ctx context.Context, vectorConfig interfaces.VectorConfig) error {
	if !vectorConfig.Enabled {
		return nil
	}
	if vectorConfig.ModelID == "" {
		httpErr := rest.NewHTTPError(ctx, http.StatusBadRequest, oerrors.OntologyManager_ObjectType_InvalidParameter).
			WithErrorDetails("VectorConfig ModelID must be set")
		return httpErr
	}
	return nil
}
