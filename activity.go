package flogoelasticstream

import (
	"encoding/json"
	"net/http"

	"github.com/TIBCOSoftware/flogo-lib/core/activity"
	"github.com/TIBCOSoftware/flogo-lib/logger"
)

// activityLog is the default logger for the exec Activity
var log = logger.GetLogger("activity-tibco-flogoelasticstream")

type elasticResponse struct {
	ScrollID string `json:"_scroll_id"`
	Hits     struct {
		Hits []map[string]interface{} `json:"hits"`
	} `json:"hits"`
}

type ElasicStreamActivity struct {
	metadata *activity.Metadata
}

// NewActivity creates a new AppActivity
func NewActivity(metadata *activity.Metadata) activity.Activity {
	return &ElasicStreamActivity{metadata: metadata}
}

// Metadata returns the activity's metadata
func (a *ElasicStreamActivity) Metadata() *activity.Metadata {
	return a.metadata
}

func (a *ElasicStreamActivity) Eval(context activity.Context) (done bool, err error) {

	basicAuthUser, _ := context.GetInput("basicAuthUser").(string)
	basicAuthPassword, _ := context.GetInput("basicAuthPassword").(string)
	elasticbaseURL, _ := context.GetInput("elasticbaseURL").(string)
	elasticQuery, _ := context.GetInput("elasticQuery").(string)

	client := &http.Client{}
	// Get first scroll page
	req, err := http.NewRequest("GET", elasticbaseURL+elasticQuery, nil)
	if err != nil {
		panic(err)
	}

	if basicAuthUser != "" {
		req.SetBasicAuth(basicAuthUser, basicAuthPassword)
	}

	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	var allHits []map[string]interface{}
	var er elasticResponse
	// Decode JSON into struct
	err = json.NewDecoder(resp.Body).Decode(&er)
	if err != nil {
		panic(err)
	}
	// Append to allHits
	allHits = append(allHits, er.Hits.Hits...)

	for {
		if er.ScrollID != "" && len(er.Hits.Hits) > 0 {

			// Still more pages to scroll
			req, err := http.NewRequest("GET", elasticbaseURL+"_search/scroll?scroll=1m&scroll_id="+er.ScrollID, nil)
			if err != nil {
				panic(err)
			}
			req.SetBasicAuth(basicAuthUser, basicAuthPassword)
			resp, err := client.Do(req)
			if err != nil {
				panic(err)
			}
			err = json.NewDecoder(resp.Body).Decode(&er)
			if err != nil {
				panic(err)
			}
			// Append to allHits
			allHits = append(allHits, er.Hits.Hits...)
		} else {
			// Done scrolling
			break
		}
	}

	log.Info("allHits length: %v\n", len(allHits))

	context.SetOutput("result", allHits)
	context.SetOutput("hits", len(allHits))

	return true, nil
}
