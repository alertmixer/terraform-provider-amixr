package amixr

import (
	"github.com/alertmixer/amixr-go-client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
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

func listOfSetsToStringSlice(listSet []interface{}) *[][]string {
	ret := [][]string{}

	if listSet == nil {
		return &ret
	}

	for _, envVal := range listSet {
		res := []string{}
		for _, envVal2 := range envVal.(*schema.Set).List() {
			res = append(res, envVal2.(string))
		}
		ret = append(ret, res)
	}
	return &ret
}

func flattenRouteSlack(in *amixr.SlackRoute) []map[string]interface{} {
	slack := make([]map[string]interface{}, 0, 1)

	out := make(map[string]interface{})

	if in.ChannelId != nil {
		out["channel_id"] = in.ChannelId
		slack = append(slack, out)
	}
	return slack
}

func expandRouteSlack(in []interface{}) *amixr.SlackRoute {
	slackRoute := amixr.SlackRoute{}

	for _, r := range in {
		inputMap := r.(map[string]interface{})
		if inputMap["channel_id"] != "" {
			channelId := inputMap["channel_id"].(string)
			slackRoute.ChannelId = &channelId
		}
	}

	return &slackRoute

}
