package main_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"

	"testing"

	"github.com/labstack/echo"
	"github.com/sighupio/permission-manager/pkg/server"
	"github.com/stretchr/testify/assert"
	"k8s.io/client-go/kubernetes/fake"
)

func AreEqualJSON(s1, s2 string) (bool, error) {
	var o1 interface{}
	var o2 interface{}

	var err error
	err = json.Unmarshal([]byte(s1), &o1)
	if err != nil {
		return false, fmt.Errorf("Error mashalling string 1 :: %s", err.Error())
	}
	err = json.Unmarshal([]byte(s2), &o2)
	if err != nil {
		return false, fmt.Errorf("Error mashalling string 2 :: %s", err.Error())
	}

	return reflect.DeepEqual(o1, o2), nil
}

func TestMain(t *testing.T) {
	kc := fake.NewSimpleClientset()
	s := server.New(kc)

	createRolebindingJSON := `{
		"roleName":"template-namespaced-resources___developer","generated_for_user":"montana","namespace":"yellow","roleKind":"ClusterRole","subjects":[{"kind":"User","name":"montana","apiGroup":"rbac.authorization.k8s.io"}],"rolebindingName":"montana___template-namespaced-resources___developer___yellow"
		}`
	req := httptest.NewRequest(echo.POST, "/create-rolebinding", strings.NewReader(createRolebindingJSON))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	res := httptest.NewRecorder()
	c := s.NewContext(req, res)
	appContext := &server.AppContext{Context: c, Kubeclient: kc}

	if assert.NoError(t, server.CreateRolebinding(appContext)) {
		assert.Equal(t, http.StatusOK, res.Code)
		assert.JSONEq(t, `{"ok":true}`, res.Body.String())
	}

	req = httptest.NewRequest(echo.GET, "/list-users", nil)
	res = httptest.NewRecorder()
	c = s.NewContext(req, res)
	appContext = &server.AppContext{Context: c, Kubeclient: kc}
	if assert.NoError(t, server.ListRbac(appContext)) {
		assert.Equal(t, http.StatusOK, res.Code)
		assert.JSONEq(t, `{"clusterRoles":null,"clusterRoleBindings":null,"roles":null,"roleBindings":[{"metadata":{"name":"montana___tem plate-namespaced-resources___developer___yellow","namespace":"yellow","c reationTimestamp":null,"labels":{"generated_for_user":"montana"}},"subje cts":[{"kind":"User","apiGroup":"rbac.authorization.k8s.io","name":"mont ana"}],"roleRef":{"apiGroup":"rbac.authorization.k8s.io","kind":"Cluster Role","name":"template-namespaced-resources___developer"}}]}`, res.Body.String())
	}
}
