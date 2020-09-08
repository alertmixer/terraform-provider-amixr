package amixr

import (
	amixr "github.com/alertmixer/amixr-go-client"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"net/http"
)

func stringSetToStringSlice(stringSet *schema.Set) *[]string {
	ret := []string{}
	if stringSet == nil {
		return &ret
	}
	for _, envVal := range stringSet.List() {
		ret = append(ret, envVal.(string))
	}
	return &ret
}

func intSetToIntSlice(intSet *schema.Set) *[]int {
	ret := []int{}
	if intSet == nil {
		return &ret
	}
	for _, envVal := range intSet.List() {
		ret = append(ret, envVal.(int))
	}
	return &ret
}

func handleNonExistentResource(f schema.ReadFunc) schema.ReadFunc {
	return func(d *schema.ResourceData, meta interface{}) error {
		if err := f(d, meta); err != nil {
			// if the error that we receive is an instance of ErrorResponse and
			// the status code is 404, it means we need to re-create resource
			errorResponse, ok := err.(*amixr.ErrorResponse)
			if !ok || errorResponse.Response.StatusCode != http.StatusNotFound {
				return err
			}
			d.SetId("")
			return nil
		}
		return nil
	}
}

func stringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}
